package main

import (
	"fmt"
	"github.com/kussell-lab/biogo/seq"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	app := kingpin.New("ReduceCoreGenome", "Make the core genome look like the flexible genome by\t"+
		"taking the flexible genome MSA and replacing the alignments with core genes")
	app.Version("v20210302")
	coreAlnFile := app.Arg("core-MSA", "multi-sequence alignment file for the core genome").Required().String()
	flexAlnFile := app.Arg("flex-MSA", "multi-sequence alignment file for the flexible genome").Required().String()
	outdir := app.Arg("outdir", "output directory for the 'flexible' genome sampled from core genes").Required().String()
	ncpu := app.Flag("num-cpu", "Number of CPUs (default: using all available cores)").Default("0").Int()
	numDigesters := app.Flag("threads", "Number of alignments to process at a time (default: 20)").Default("20").Int()

	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *ncpu == 0 {
		*ncpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(*ncpu)

	//coreAlnFile := "/Volumes/aps_timemachine/recombo/APS162_ResampleCore/threshold99/MSA_CORE"
	//flexAlnFile := "/Volumes/aps_timemachine/recombo/APS162_ResampleCore/threshold99/MSA_FLEX"
	//outdir := "/Volumes/aps_timemachine/recombo/APS162_ResampleCore/threshold99"
	//numDigesters := 4
	//timer
	start := time.Now()
	makeResampledMSA(*outdir)
	done := make(chan struct{})
	//read through the flexible genome
	flexAln, errc := readAlignments(done, *flexAlnFile)
	//get the number of core alignments
	lenCore := getNumberOfAlignments(*coreAlnFile)
	//start a fixed number of goroutines to read alignments and split into core/flex
	c := make(chan Alignment)
	var wg sync.WaitGroup
	for i := 0; i < *numDigesters; i++ {
		wg.Add(1)
		go resampleCoreAln(done, flexAln, *coreAlnFile, lenCore, c, i, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	//end of pipeline; write files
	for rgene := range c {
		//fmt.Print(gene.ID)
		//mutex.Lock()
		writeMSA(rgene, *outdir)
		//mutex.Unlock()
	}
	if err := <-errc; err != nil { // HLerrc
		panic(err)
	}
	//add the number of core and flex to the bottom of the spreadsheet

	duration := time.Since(start)
	fmt.Println("Time to make a reduced core genome blahhhhh:", duration)
}

//makeResampledMSA makes the outdir and initializes the MSA files for core and flexible genomes
func makeResampledMSA(outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}
	MSA := filepath.Join(outdir, "MSA_REDUCED_CORE")
	f, err := os.Create(MSA)
	check(err)
	f.Close()
}

//check for errors
func check(e error) {
	if e != nil {
		panic(e)
	}
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
				case alignments <- Alignment{alnID, alignment}:
					//fmt.Printf("\rRead %d alignments.", numAln)
					//fmt.Printf("\r alignment ID: %s", alnID)
				case <-done:
					fmt.Printf(" Total alignments %d\n", numAln)
				}
			}
		}
		//f.Close()
		errc <- err
	}()
	return alignments, errc
}

// Alignment is an array of multiple sequences with same length.
type Alignment struct {
	ID        string
	Sequences []seq.Sequence
}

//getCoreAln grabs a random core gene to be re-sampled to look like a flexble gene
func getCoreAln(file string, lenCore int) (coreAln chan Alignment) {
	coreAln = make(chan Alignment)
	go func() {
		defer close(coreAln)
		//generate a random number from 0 to (# of core genes - 1)
		randNum := rand.Intn(lenCore)
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		xmfaReader := seq.NewXMFAReader(f)
		//start the count
		count := 0
		for {
			alignment, err := xmfaReader.Read()
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			if count == randNum {
				header := strings.Split(alignment[0].Id, " ")
				coreAlnID := header[0]
				coreAln <- Alignment{ID: coreAlnID, Sequences: alignment}
				break
			}
			count++
		}
	}()

	return
}

// mustOpen is a helper function to open a file.
// and panic if error occurs.
func mustOpen(file string) (f *os.File) {
	var err error
	f, err = os.Open(file)
	if err != nil {
		panic(err)
	}
	return
}

// resampleCoreAln reads flexible gene alignments, makes a map of strains with that gene,
// then grabs a core gene and grabs all the gene alignments for the strains in the map
// then sends this resampled alignment on genes until the flexible MSA or done channel is closed
func resampleCoreAln(done <-chan struct{}, alignments <-chan Alignment, coreAlnFile string, lenCore int, rgenes chan<- Alignment, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Printf("Worker %d starting\n", id)
	for aln := range alignments { // HLpaths
		//make a map of strains with the flexible gene
		strainMap := make(map[string]bool)
		for _, s := range aln.Sequences {
			strain := getStrainName(s)
			strainMap[strain] = true
		}
		//grab a random core gene
		//mutex.Lock()
		coreAlnChan := getCoreAln(coreAlnFile, lenCore)
		//mutex.Unlock()
		coreAln := <-coreAlnChan
		fmt.Printf("\r Resampling from: %s\n", coreAln.ID)
		// ... and only grab gene alignments in strainMap
		var resampledSeqs []seq.Sequence
		//stop the loop if you've found all the strains
		numFlex := len(strainMap)
		count := 0
		for _, s := range coreAln.Sequences {
			strain := getStrainName(s)
			if _, found := strainMap[strain]; found {
				resampledSeqs = append(resampledSeqs, s)
				count++
				if count == numFlex {
					break
				}
			}
		}
		//send out the resampled gene to the world
		rgene := Alignment{aln.ID, resampledSeqs}
		//writeAln(aln, outdir)
		select {
		//case c <- aln.num:
		case rgenes <- rgene:
		//	writeAln(aln, outdir)
		case <-done:
			return
		}
	}
	//fmt.Printf("Worker %d done\n", id)

}

//getStrainName gets the strain name from a sequence header
func getStrainName(s seq.Sequence) (strain string) {
	header := s.Id
	terms := strings.Split(header, " ")
	strain = terms[len(terms)-1]
	return
}

// getNumberOfAlignments return total number of alignments in a xmfa file.
func getNumberOfAlignments(file string) (count int) {
	c := readXMFA(file)
	for a := range c {
		if len(a) >= 2 {
			count++
		}
	}
	return
}

// readXMFA reads a xmfa format file and returns a channel of []seq.Sequence.
func readXMFA(file string) chan []seq.Sequence {
	c := make(chan []seq.Sequence)
	go func() {
		defer close(c)

		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rd := seq.NewXMFAReader(f)
		for {
			a, err := rd.Read()
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			//temporary change from 2 to 1
			if len(a) >= 1 {
				c <- a
			}
		}
		//f.Close()
	}()
	return c
}

//writeMSA write the gene to the correct MSA (core or flex)
func writeMSA(c Alignment, outdir string) {
	MSAname := "MSA_REDUCED_CORE"
	MSA := filepath.Join(outdir, MSAname)
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile(MSA, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, s := range c.Sequences {
		f.WriteString(">" + s.Id + "\n")
		f.Write(s.Seq)
		f.WriteString("\n")
	}
	f.WriteString("=\n")
	//f.Close()
}
