package main

//this program splits a master MSA file into core and flexible genomes
//written by Asher Preska Steinberg (apsteinberg@nyu.edu)
import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/kussell-lab/biogo/seq"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	app := kingpin.New("SplitGenome", "splits a master MSA file of all strains into core and flexible genomes")
	app.Version("v20210112")
	alnFile := app.Arg("master_MSA", "multi-sequence alignment file for all genes").Required().String()
	sampleFile := app.Arg("strain list", "list of all strains").Required().String()
	cutoff := app.Arg("core-cutoff", "Percentage above which to be considered a core gene (0 to 100)").Required().Int()
	outdir := app.Arg("outdir", "output directory for the core/flex MSA gene percentage csv").Required().String()
	ncpu := app.Flag("num-cpu", "Number of CPUs (default: using all available cores)").Default("0").Int()
	numSplitters := app.Flag("threads", "Number of alignments to process at a time (default: 8)").Default("8").Int()

	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *ncpu == 0 {
		*ncpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(*ncpu)
	//alnFile := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/1224_properheader"
	////strain list
	//sampleFile := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/strain_list"
	//outdir := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/threshold99"
	//numSplitters := 4
	//cutoff := 99
	//timer
	start := time.Now()

	//make the outdir and core and flexible MSAs
	makeCFMSA(*outdir)
	//prepare the gene percentage out csv
	makeGeneCSV(*cutoff, *outdir)
	//set the threshold
	threshold := float64(*cutoff) / 100
	samples := readSamples(*sampleFile)
	//get the total number of sequences
	totSeqs := len(samples)
	done := make(chan struct{})
	//read in alignments
	alignments, errc := readAlignments(done, *alnFile)

	//start a fixed number of goroutines to read alignments and split into core/flex
	c := make(chan result)
	var wg sync.WaitGroup
	for i := 0; i < *numSplitters; i++ {
		wg.Add(1)
		go splitter(done, alignments, c, totSeqs, threshold, i, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	//end of pipeline; write files
	for gene := range c {
		writeMSA(gene, *outdir)
		getGenePercentage(gene, *cutoff, *outdir)
	}
	if err := <-errc; err != nil { // HLerrc
		panic(err)
	}
	//add the number of core and flex to the bottom of the spreadsheet

	duration := time.Since(start)
	fmt.Println("Time to split into core and flex:", duration)
}

// Alignment is an array of multiple sequences with same length.
type Alignment struct {
	ID        string
	num       int
	Sequences []seq.Sequence
}

// A result is a single gene alignment belonging to the core or flexible genome
type result struct {
	Alignment Alignment
	genome    string  //"CORE" or "FLEX"
	frac      float64 //fraction of strains that have the gene
}

// readAlignments reads sequence alignment from a extended Multi-FASTA file,
// and return a channel of alignment, which is a list of seq.Sequence
func readAlignments(done <-chan struct{}, file string) (<-chan Alignment, <-chan error) {
	alignments := make(chan Alignment)
	errc := make(chan error, 1)
	go func() {
		defer close(alignments)

		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		xmfaReader := seq.NewXMFAReader(f)
		numAln := 0
		for {
			alignment, err := xmfaReader.Read()
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			if len(alignment) > 0 {
				numAln++
				alnID := strings.Split(alignment[0].Id, " ")[0]
				select {
				case alignments <- Alignment{alnID, numAln, alignment}:
					fmt.Printf("\rRead %d alignments.", numAln)
					fmt.Printf("\r alignment ID: %s", alnID)
				case <-done:
					fmt.Printf(" Total alignments %d\n", numAln)
				}
			}
		}
		errc <- err
	}()
	return alignments, errc
}

// splitter reads gene alignments from the master MSA, figures out if the gene is core/flex
// then sends these processed results on alnChan until either the master MSA or done channel is closed.
func splitter(done <-chan struct{}, alignments <-chan Alignment, genes chan<- result, totSeqs int, threshold float64, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Printf("Worker %d starting\n", id)
	for aln := range alignments { // HLpaths
		//get the fraction of sequences which have the gene
		var frac float64
		//define core/flex string
		var genome string
		//count number of strains with the gene; the strain needs to have at least one full codon
		//to say the gene is present
		var count int
		for _, s := range aln.Sequences {
			NumFullCodons := extractFullCodons(s)
			if NumFullCodons > 0 {
				count++
			}
		}
		frac = float64(count) / float64(totSeqs)
		//is it core or flex
		if frac > threshold {
			genome = "CORE"
		} else {
			genome = "FLEX"
		}
		gene := result{aln, genome, frac}
		//writeAln(aln, outdir)
		select {
		//case c <- aln.num:
		case genes <- gene:
		//	writeAln(aln, outdir)
		case <-done:
			return
		}
	}
	//fmt.Printf("Worker %d done\n", id)

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

//check for errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//writeMSA write the gene to the correct MSA (core or flex)
func writeMSA(c result, outdir string) {
	MSAname := "MSA_" + c.genome
	MSA := filepath.Join(outdir, MSAname)
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile(MSA, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	aln := c.Alignment
	for _, s := range aln.Sequences {
		f.WriteString(">" + s.Id + "\n")
		f.Write(s.Seq)
		f.WriteString("\n")
	}
	f.WriteString("=\n")
}

func getGenePercentage(c result, cutoff int, outdir string) {
	threshold := strconv.Itoa(cutoff)
	name := "gene_percentages_" + threshold + "%_cutoff.csv"
	path := filepath.Join(outdir, name)
	w, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	defer w.Close()
	if err != nil {
		panic(err)
	}
	csvwriter := csv.NewWriter(w)
	defer csvwriter.Flush()
	//get percentages and write a line
	p := fmt.Sprintf("%f", c.frac)
	aln := c.Alignment
	genePercent := []string{aln.ID, p, c.genome}
	csvwriter.Write(genePercent)
	//w.Close()
}

//makeGeneCSV initiates the gene percentage CSV
func makeGeneCSV(cutoff int, outdir string) {
	//prepare the gene percentage out csv
	t := strconv.Itoa(cutoff)
	name := "gene_percentages_" + t + "%_cutoff.csv"
	path := filepath.Join(outdir, name)
	w, err := os.Create(path)
	check(err)
	defer w.Close()
	csvwriter := csv.NewWriter(w)
	defer csvwriter.Flush()
	header := []string{"gene", "fraction of strains", "genome"}
	err = csvwriter.Write(header)
	check(err)

	return
}

//makeCFMSA makes the outdir and initializes the MSA files for core and flexible genomes
func makeCFMSA(outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}
	MSA := filepath.Join(outdir, "MSA_CORE")
	f, err := os.Create(MSA)
	check(err)
	f.Close()
	MSA = filepath.Join(outdir, "MSA_FLEX")
	f, err = os.Create(MSA)
	check(err)
	f.Close()
}

// extractFullCodons returns the number of full codons
//there needs to be at least 1 full codon for us to say the strain "has the gene"
func extractFullCodons(s seq.Sequence) (NumFullCodons int) {
	var codons []Codon
	for i := 0; i+3 <= len(s.Seq); i += 3 {
		c := s.Seq[i:(i + 3)]
		//check for gaps
		containsGap := false
		for k := 0; k < 3; k++ {
			if c[k] == '-' || c[k] == 'N' {
				containsGap = true
				break
			}
		}
		if containsGap {
			continue
		} else {
			codons = append(codons, c)
		}

	}
	NumFullCodons = len(codons)
	return
}

// Codon is a byte list of length 3
type Codon []byte

// readAlignments reads sequence alignment from a extended Multi-FASTA file,
// and return a channel of alignment, which is a list of seq.Sequence
//func readAlignments(file string) (alnChan chan Alignment) {
//	alnChan = make(chan Alignment)
//	read := func() {
//		defer close(alnChan)
//
//		f, err := os.Open(file)
//		if err != nil {
//			panic(err)
//		}
//		defer f.Close()
//		xmfaReader := seq.NewXMFAReader(f)
//		numAln := 0
//		for {
//			alignment, err := xmfaReader.Read()
//			if err != nil {
//				if err != io.EOF {
//					panic(err)
//				}
//				break
//			}
//			if len(alignment) > 0 {
//				numAln++
//				alnID := strings.Split(alignment[0].Id, " ")[0]
//				alnChan <- Alignment{ID: alnID, Sequences: alignment}
//				fmt.Printf("\rRead %d alignments.", numAln)
//				fmt.Printf("\r alignment ID: %s", alnID)
//			}
//		}
//		fmt.Printf(" Total alignments %d\n", numAln)
//	}
//	go read()
//	return
//}
