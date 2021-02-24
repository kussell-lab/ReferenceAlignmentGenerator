// This program collects gene alignments from multiple FASTA files.
// Created by Mingzhi Lin (mingzhi9@gmail.com).
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"runtime"

	"io/ioutil"

	"github.com/alecthomas/kingpin"
	"github.com/boltdb/bolt"
	"github.com/cheggaaa/pb"
	"github.com/kussell-lab/biogo/feat/gff"
	"github.com/kussell-lab/biogo/seq"
)

// BucketName is the name of bucket.
const BucketName = "col"

var (
	// ShowProgress indicates that progress can be showed.
	ShowProgress bool
)

// Genome stores genome meta and sequence data.
type Genome struct {
	Sample string
	Seq    *seq.Sequence
}

func main() {
	app := kingpin.New("CollectGeneAlignments", "collect gene alignments.")
	app.Version("0.4.1")
	sampleFile := app.Arg("list_file", "file contains a list of genomes").String()
	gffFile := app.Arg("gff_file", "gff file").String()
	dataDir := app.Arg("data_dir", "genome sequence folder").String()
	outFile := app.Arg("out_file", "output file").String()
	bufSize := app.Flag("buf_size", "gene load buf size (memory usage), default 1000 genes").Default("1000").Int()
	numGenome := app.Flag("num_genome", "number of samples per db").Default("1000").Int()
	appendix := app.Flag("appendix", "appendix (default .fasta)").Default(".fasta").String()
	tmpDir := app.Flag("temp_dir", "temp dir").Default(".").String()
	showProgress := app.Flag("progress", "showing progress").Default("false").Bool()
	ncpu := app.Flag("num_cpu", "number of threads").Default("0").Int()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	ShowProgress = *showProgress
	if *ncpu <= 0 || *ncpu > runtime.NumCPU() {
		*ncpu = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*ncpu)

	// prepare gene splitter.
	gffs := readGff(*gffFile)
	splitter := NewGeneSplitter(gffs)

	samples := readSamples(*sampleFile)

	sampleChan := make(chan string)
	go func() {
		defer close(sampleChan)
		var pbar *pb.ProgressBar
		if ShowProgress {
			pbar = pb.StartNew(len(samples))
			defer pbar.FinishPrint("Finish collecting alignments")
		}
		for _, sample := range samples {
			sampleChan <- sample
			if ShowProgress {
				pbar.Increment()
			}
		}
	}()

	jobChan := make(chan *LoadingJob)
	done := make(chan bool)
	for i := 0; i < *ncpu; i++ {
		go func() {
			for {
				dbFile, err := ioutil.TempFile(*tmpDir, fmt.Sprintf("boltdb_%d", i))
				if err != nil {
					panic(err)
				}
				job := &LoadingJob{
					dbFile:   dbFile.Name(),
					dataDir:  *dataDir,
					appendix: *appendix,
					bufSize:  *bufSize,
					splitter: splitter,
				}
				job.Load(sampleChan, *numGenome)
				jobChan <- job
				if len(job.samples) < *numGenome {
					break
				}
			}
			done <- true
		}()
	}

	go func() {
		defer close(jobChan)
		for i := 0; i < *ncpu; i++ {
			<-done
		}
	}()

	geneSet := make(map[string]bool)

	var jobs []*LoadingJob
	for job := range jobChan {
		for gene := range job.geneSet {
			geneSet[gene] = true
		}
		jobs = append(jobs, job)
		defer os.Remove(job.dbFile)
	}

	geneIDChan := make(chan string)
	go func() {
		defer close(geneIDChan)
		for id := range geneSet {
			geneIDChan <- id
		}
	}()

	w, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	wb := bufio.NewWriter(w)
	defer wb.Flush()
	// pull gene alignments, and write to file.
	var pbar *pb.ProgressBar
	if ShowProgress {
		pbar = pb.StartNew(len(geneSet))
		defer pbar.FinishPrint("Finish collecting alignments")
	}
	for id := range geneIDChan {
		alignChan := pullGenes(id, jobs, *ncpu)
		numGenes := 0
		alnLength := -1
		for aln := range alignChan {
			for _, gene := range aln {
				if alnLength == -1 {
					if len(gene.Seq) > 0 {
						alnLength = len(gene.Seq)
					}
				}
				if len(gene.Seq) == alnLength && len(strings.TrimSpace(gene.Id)) > 0 {
					wb.WriteString(">" + gene.Id + "\n")
					wb.Write(gene.Seq)
					wb.WriteString("\n")
					numGenes++
				}
			}
		}
		if numGenes > 0 {
			wb.WriteString("=\n")
		}
		if ShowProgress {
			pbar.Increment()
		}
	}
}

func pullGenes(id string, jobList []*LoadingJob, ncpu int) chan []seq.Sequence {
	jobChan := make(chan *LoadingJob)
	go func() {
		defer close(jobChan)
		for _, job := range jobList {
			jobChan <- job
		}
	}()

	resChan := make(chan []seq.Sequence, ncpu)
	done := make(chan bool)
	for i := 0; i < ncpu; i++ {
		go func() {
			for job := range jobChan {
				db := job.db
				var geneIDs []string
				for j := 0; j < len(job.samples); j++ {
					geneID := fmt.Sprintf("%s_%d", id, j)
					geneIDs = append(geneIDs, geneID)
				}
				genes := getGene(db, BucketName, geneIDs)
				for j, gene := range genes {
					gene.Id = id + " " + job.samples[j]
					genes[j] = gene
				}
				resChan <- genes
			}
			done <- true
		}()
	}

	go func() {
		defer close(resChan)
		for i := 0; i < ncpu; i++ {
			<-done
		}
	}()
	return resChan
}

func splitSamples(samples []string, perDB int) (subSamples [][]string) {
	numSubSamples := len(samples)/perDB + 1
	for i := 0; i < numSubSamples; i++ {
		start := i * perDB
		end := (i + 1) * perDB
		if end > len(samples) {
			end = len(samples)
		}
		subSamples = append(subSamples, samples[start:end])
	}

	return
}

// LoadingJob load genomes into a db.
type LoadingJob struct {
	samples  []string
	db       *bolt.DB
	dbFile   string
	dataDir  string
	appendix string
	bufSize  int
	splitter Splitter
	geneSet  map[string]bool
}

// Load loads genomes into a db.
func (job *LoadingJob) Load(sampleChan chan string, numGenome int) {

	// create bolt db.
	job.db = createDB(job.dbFile)
	createBucket(job.db, BucketName)

	// geneCount stores number of genes of the same id.
	job.geneSet = make(map[string]bool)

	// splitting genes and push them into a channel.
	geneChan := make(chan *seq.Sequence)
	go func() {
		defer close(geneChan)
		for sample := range sampleChan {
			genomes := readGenome(sample, job.dataDir, job.appendix)
			for _, genome := range genomes {
				genes := job.splitter.Split(genome)
				for _, gene := range genes {
					job.geneSet[gene.Id] = true
					gene.Id = fmt.Sprintf("%s_%d", gene.Id, len(job.samples))
					geneChan <- gene
				}
			}
			job.samples = append(job.samples, sample)

			if len(job.samples) > numGenome {
				break
			}
		}
	}()

	// push genes into bolt db.
	putGene(job.db, BucketName, geneChan, job.bufSize)
	return
}

// readSamples return a list of samples from a sample file.
func readSamples(filename string) (samples []string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error when reading file %s:%v", filename, err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				log.Fatalf("Error when reading file %s: %v", filename, err)
			}
			break
		}
		samples = append(samples, strings.TrimSpace(line))
	}
	return
}

func readGenome(sample, workDir, appendix string) (genomes []*seq.Sequence) {
	filename := filepath.Join(workDir, sample+appendix)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error when opening file %s: %v", filename, err)
	}
	defer f.Close()

	rd := seq.NewFastaReader(f)
	genomes, err = rd.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	return genomes
}

// readGff return gene features.
func readGff(filename string) []*gff.Record {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error when reading file %s:%v", filename, err)
	}
	defer f.Close()

	r := gff.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error when reading file %s:%v", filename, err)
	}
	return records
}

// createDB creates a bolt db.
func createDB(dbFile string) *bolt.DB {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// createBucket creates a bucket.
func createBucket(db *bolt.DB, bucketName string) {
	fn := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	}

	err := db.Update(fn)
	if err != nil {
		log.Fatal(err)
	}
}

func putGene(db *bolt.DB, bucketName string, geneChan chan *seq.Sequence, bufSize int) {
	buf := make([]*seq.Sequence, bufSize)
	count := 0
	for gene := range geneChan {
		if count >= len(buf) {
			loadGenes(db, bucketName, buf)
			count = 0
		}
		buf[count] = gene
		count++
	}

	loadGenes(db, bucketName, buf[:count])
}

func loadGenes(db *bolt.DB, bucketName string, genes []*seq.Sequence) {
	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for _, gene := range genes {
			err := b.Put([]byte(gene.Id), gene.Seq)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := db.Update(fn)
	if err != nil {
		log.Fatal(err)
	}
}

func getGene(db *bolt.DB, bucketName string, ids []string) (genes []seq.Sequence) {
	fn := func(tx *bolt.Tx) error {
		for _, id := range ids {
			b := tx.Bucket([]byte(bucketName))
			v := b.Get([]byte(id))
			var gene seq.Sequence
			gene.Id = id
			gene.Seq = v
			genes = append(genes, gene)
		}

		return nil
	}

	db.View(fn)

	return
}
