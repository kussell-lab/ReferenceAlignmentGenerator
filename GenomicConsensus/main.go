// GenomicConsensus converts a samtools pileup file to the genomic sequence.
// Created by Mingzhi Lin (mingzhi9@gmail.com).
package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kussell-lab/biogo/pileup"
	"github.com/kussell-lab/biogo/seq"
)

func main() {
	var pileFile = kingpin.Flag("pileup_file", "pileup file (default reading from stdin).").Default("").String()
	var baseQual = kingpin.Flag("base_qual", "base quality").Default("30").Int8()
	var outFile = kingpin.Flag("out_file", "output sequence file in FASTA format (default writing to stdout).").Default("").String()
	var fnaFile = kingpin.Flag("fna_file", "reference sequences file").Required().String()
	kingpin.Parse()

	ss := readFastas(*fnaFile)
	seqLens := getSequenceLengths(ss)
	snpChan := readPileup(*pileFile)
	baseChan := pileupSupport(snpChan, *baseQual)
	seqChan := generateSequence(baseChan, seqLens)
	writeFasta(seqChan, *outFile)
}

// openFile is a helper to open a file.
func openFile(fileName string) (f *os.File) {
	var err error
	if fileName == "" {
		f = os.Stdin
		return
	}

	f, err = os.Open(fileName)
	if err != nil {
		log.Fatalf("Error when openning file %s: %v", fileName, err)
	}
	return
}

// createFile is a helper to open a file.
func createFile(fileName string) (f *os.File) {
	var err error
	if fileName == "" {
		f = os.Stdout
		return
	}

	f, err = os.Create(fileName)
	if err != nil {
		log.Fatalf("Error when openning file %s: %v", fileName, err)
	}
	return
}

// readPileup returns a channel of *pileup.SNP.
func readPileup(fileName string) chan *pileup.SNP {
	c := make(chan *pileup.SNP)
	go func() {
		defer close(c)

		f := openFile(fileName)
		defer f.Close()

		r := pileup.NewReader(f)
		for {
			snp, err := r.Read()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("Error when parsing pileup file: %v", err)
				}
				break
			}
			c <- snp
		}
	}()
	return c
}

// Base is a record store the position, ref, alt base and its supports.
type Base struct {
	Ref     string
	Pos     int
	RefBase byte
	AltBase byte
}

// pileupSupport filter and return SNP calls from a pileup source.
func pileupSupport(c chan *pileup.SNP, baseQual int8) chan Base {
	baseChan := make(chan Base)

	go func() {
		defer close(baseChan)
		letters := []byte{'a', 'A', 'c', 'C', 'g', 'G', 't', 'T'}

		for s := range c {
			cc := make([]int, len(letters))
			for i, b := range s.Bases {
				if int8(s.Quals[i]) >= baseQual {
					index := bytes.IndexByte(letters, b)
					if index >= 0 {
						cc[index]++
					}
				}
			}

			dp := 0
			for _, c := range cc {
				dp += c
			}

			var baseCalled byte = '-'
			for i := 0; i < 4; i++ {
				n := cc[i*2] + cc[i*2+1]
				if float64(n)/float64(dp) >= 0.75 && cc[i*2] >= 2 && cc[i*2+1] >= 2 {
					baseCalled = byte(letters[i*2+1])
				}
			}

			if baseCalled != '-' {
				refBase := bytes.ToUpper([]byte{s.RefBase})[0]
				altBase := bytes.ToUpper([]byte{baseCalled})[0]
				b := Base{
					Ref:     s.Reference,
					Pos:     s.Position,
					RefBase: refBase,
					AltBase: altBase,
				}
				baseChan <- b
			}
		}
	}()
	return baseChan
}

// generateSequence generates a sequence from the SNP chan.
func generateSequence(baseChan chan Base, seqLens map[string]int) chan seq.Sequence {
	seqChan := make(chan seq.Sequence)
	go func() {
		defer close(seqChan)
		currentRef := ""
		ss := []byte{}
		for b := range baseChan {
			ref := b.Ref
			if ref != currentRef {
				if len(ss) > 0 {
					l := seqLens[currentRef]
					for len(ss) < l {
						ss = append(ss, '-')
					}
					s := seq.Sequence{}
					s.Id = currentRef
					s.Name = currentRef
					s.Seq = ss
					seqChan <- s
				}
				currentRef = ref
				ss = []byte{}
			}
			pos := b.Pos
			for len(ss) < pos-1 {
				ss = append(ss, '-')
			}
			ss = append(ss, b.AltBase)
		}

		if len(ss) > 0 {
			l := seqLens[currentRef]
			for len(ss) < l {
				ss = append(ss, '-')
			}
			s := seq.Sequence{}
			s.Id = currentRef
			s.Name = currentRef
			s.Seq = ss
			seqChan <- s
		}
	}()
	return seqChan
}

// writeFasta write genomes into a fasta file.
func writeFasta(seqChan chan seq.Sequence, file string) {
	f := createFile(file)
	defer f.Close()

	for s := range seqChan {
		f.WriteString(">" + s.Id + "\n")
		bases := s.Seq
		numOfLines := len(bases)/80 + 1
		for i := 0; i < numOfLines; i++ {
			start := i * 80
			end := (i + 1) * 80
			if end > len(bases) {
				end = len(bases)
			}

			if start < len(bases) {
				f.WriteString(string(bases[start:end]) + "\n")
			}
		}
	}
}

// readFastas returns a genome sequence.
func readFastas(file string) []*seq.Sequence {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fastaReader := seq.NewFastaReader(f)
	s, err := fastaReader.ReadAll()
	if err != nil {
		log.Fatalf("Error when reading fasta file: %v", err)
	}
	return s
}

// getSequenceLengths returns a map of sequence name to its length.
func getSequenceLengths(ss []*seq.Sequence) map[string]int {
	m := make(map[string]int)
	for _, s := range ss {
		m[s.Id] = len(s.Seq)
	}
	return m
}
