package main

// author: Asher Preska Steinberg
import (
	"bufio"
	"fmt"
	"github.com/kussell-lab/biogo/seq"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

func main() {
	app := kingpin.New("measureGaps", "measure fraction of MSA that's been aligned to the reference genome")
	app.Version("v20210302")
	alnFile := app.Arg("MSA", "multi-sequence alignment file").Required().String()
	ncpu := app.Flag("num-cpu", "Number of CPUs (default: using all available cores)").Default("0").Int()
	numDigesters := app.Flag("threads", "Number of alignments to process at a time (default: 8)").Default("8").Int()

	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *ncpu == 0 {
		*ncpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(*ncpu)

	//alnFile := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/1224_properheader"
	//numDigesters := 20

	//start := time.Now()
	//done channel will close when we're done reading alignments
	done := make(chan struct{})
	//read in alignments
	alignments, errc := readAlignments(done, *alnFile)

	//start a fixed number of goroutines to read alignments and split into core/flex
	c := make(chan result)
	var wg sync.WaitGroup
	for i := 0; i < *numDigesters; i++ {
		wg.Add(1)
		go MeasureGaps(done, alignments, c, i, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	//end of pipeline; write files
	//total gap length
	var totGaps float64
	//total sequence length
	var totLength float64
	//sum of aligned fractions
	var sumFrac float64
	//number of genes
	var numGenes float64
	//initializing counters ...
	totGaps = 0
	totLength = 0
	sumFrac = 0
	for gene := range c {
		totGaps = totGaps + gene.numGaps
		totLength = totLength + gene.seqLength
		sumFrac = sumFrac + gene.frac
		numGenes++
	}
	if err := <-errc; err != nil { // HLerrc
		panic(err)
	}
	//total aligned fraction
	fracAligned := (totLength - totGaps) / totLength
	//average aligned fraction for each gene
	avgAligned := sumFrac / numGenes
	fmt.Printf("Statistics for MSA file: %s\n", *alnFile)
	fmt.Print("-------------------------------\n")
	fmt.Printf("Total aligned fraction: %f\n", fracAligned)
	fmt.Printf("Aligned fraction per gene: %f\n", avgAligned)
	fmt.Print("-------------------------------\n")
	//duration := time.Since(start)
	//fmt.Println("Time to measure sequence gaps:", duration)
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
					// fmt.Printf("\rRead %d alignments.", numAln)
					// fmt.Printf("\r alignment ID: %s", alnID)
				case <-done:
					// fmt.Printf(" Total alignments %d\n", numAln)
				}
			}
		}
		errc <- err
	}()
	return alignments, errc
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
	frac      float64 //fraction of gene that's aligned
	seqLength float64 //total gene sequence length
	numGaps   float64 //number of gaps in sequence length
}

// MeasureGaps reads gene alignments from the master MSA, figures out what portion of the gene is aligned
// then sends these processed results on alnChan until either the master MSA or done channel is closed.
func MeasureGaps(done <-chan struct{}, alignments <-chan Alignment, genes chan<- result, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Printf("Worker %d starting\n", id)
	for aln := range alignments { // HLpaths
		//total number of gaps for a gene sequence
		var totGaps float64
		totGaps = 0
		//total length of the gene alignment
		var totLength float64

		for _, s := range aln.Sequences {
			NumGaps := countGaps(s)
			seqLength := float64(len(s.Seq))
			//don't take sequences from strains which don't have the gene
			//(depending on how the alignment was done, these show up as sequences with all gaps)
			if NumGaps == seqLength {
				continue
			} else {
				totGaps = NumGaps + totGaps
				totLength = totLength + seqLength
			}

		}
		//fraction aligned for that gene
		frac := (totLength - totGaps) / totLength
		gene := result{aln, frac, totLength, totGaps}
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

// countGaps counts the number of gaps in a gene sequence
func countGaps(s seq.Sequence) (NumGaps float64) {
	for i := 0; i < len(s.Seq); i++ {
		b := s.Seq[i]
		if b == '-' || b == 'N' {
			NumGaps++
		}
	}
	return
}

// Codon is a byte list of length 3
type Codon []byte
