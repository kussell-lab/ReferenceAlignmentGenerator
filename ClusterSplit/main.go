package main

//this program splits MSA into MSA files for each sequence cluster
//written by Asher Preska Steinberg (apsteinberg@nyu.edu)
import (
	"bufio"
	"fmt"
	"github.com/kussell-lab/biogo/seq"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	app := kingpin.New("ClusterSplit", "split MSA into MSA files for each sequence cluster; option to run on both CORE and FLEX MSA files")
	app.Version("v20210113")
	alnFile := app.Arg("MSA", "multi-sequence alignment file").Required().String()
	outdir := app.Arg("outdir", "output directory for cluster MSA files").Required().String()
	clusterdict := app.Arg("cluster_dict", "hash table from makeCluster.py as a text file, relating cluster # to sequence name").Required().String()
	ncpu := app.Flag("num-cpu", "Number of CPUs (default: using all available cores)").Default("0").Int()
	secondAln := app.Flag("FLEX_MSA", "multi-sequence alignment file for flexible genome").Default("None").String()
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *ncpu == 0 {
		*ncpu = runtime.NumCPU()
	}

	runtime.GOMAXPROCS(*ncpu)

	//for local testing
	//alnFile := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/1224_properheader"
	//secondAln := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/testxmfa"
	//outdir := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/threshold99"
	//clusterdict := "/Volumes/aps_timemachine/recombo/APS160_splitGenome/clusterlist"
	//timer
	start := time.Now()
	//get the MSA file name
	var alnFiles []string
	if *secondAln != "None" {
		alnFiles = []string{*alnFile, *secondAln}
	} else {
		alnFiles = []string{*alnFile}
	}
	for _, alnFile := range alnFiles {
		terms := strings.Split(alnFile, "/")
		MSAname := terms[len(terms)-1]
		//convert clusterdict to a golang map and a slice of strings
		clustermap, clusters := makeClusterMap(*clusterdict)
		//make cluster folders and blank MSA files for writing
		initClusterMSAs(*outdir, MSAname, clusters)

		//define alignment channel
		var alnChan chan Alignment
		alnChan = make(chan Alignment)
		go func() {
			defer close(alnChan)
			count := 0
			c := readAlignments(alnFile)
			for a := range c {
				alnChan <- a
				count++
			}
		}()

		for aln := range alnChan {
			alnMap := AssembleAlignments(aln, clustermap)
			//same as above, but for the CFsplit = False case
			var wg sync.WaitGroup
			wg.Add(len(clusters))
			for i := 0; i < len(clusters); i++ {
				go func(i int) {
					defer wg.Done()
					//define the cluster ID
					ID := clusters[i]
					//get the cluster alignment
					cAln, found := alnMap[ID]
					if found {
						cluster := cAlignment{clusterID: ID, Sequences: cAln}
						WriteClusterMSA(cluster, *outdir, MSAname)
					}
				}(i)
			}
			wg.Wait()

		}
	}

	duration := time.Since(start)
	fmt.Println("Time to split into clusters:", duration)
}

//makeClusterMap takes the clusterdict and returns both a map relating sequence to cluster ID
//and returns a vector of cluster IDs
func makeClusterMap(clusterdict string) (clustermap map[string]string, clusters []string) {
	f, err := os.Open(clusterdict)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	clustermap = make(map[string]string)
	var intclusters []int
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		ln := strings.Split(line, ",")
		seq := ln[0]
		cluster := strings.Split(ln[1], "\n")
		clusternum, _ := strconv.Atoi(cluster[0])
		clustermap[seq] = cluster[0]
		intclusters = append(intclusters, clusternum)
		if err == io.EOF {
			break
		}
	}
	clusters = unique(intclusters)

	return
}

//unique returns a unique slice of strings from a slice of integers
func unique(intSlice []int) []string {
	keys := make(map[int]bool)
	list := []int{}
	strlist := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	sort.Ints(list)
	for _, num := range list {
		strnum := strconv.Itoa(num)
		strlist = append(strlist, strnum)
	}
	return strlist
}

//initClusterMSAs makes cluster folders and blank MSAs for writing
func initClusterMSAs(outdir string, MSAname string, clusters []string) {
	//make the cluster folders and MSAs
	for _, cluster := range clusters {
		clusterfolder := "cluster" + cluster
		clusterpath := filepath.Join(outdir, clusterfolder)
		if _, err := os.Stat(clusterpath); os.IsNotExist(err) {
			os.Mkdir(clusterpath, 0755)
		}
		MSA := filepath.Join(clusterpath, MSAname+"_cluster"+cluster)
		f, err := os.Create(MSA)
		check(err)
		f.Close()
	}
}

//check for errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Alignment is an array of mutliple sequences with same length.
type Alignment struct {
	ID        string
	Sequences []seq.Sequence
}

// readAlignments reads sequence alignment from a extended Multi-FASTA file,
// and return a channel of alignment, which is a list of seq.Sequence
func readAlignments(file string) (alnChan chan Alignment) {
	alnChan = make(chan Alignment)
	read := func() {
		defer close(alnChan)

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
				alnChan <- Alignment{ID: alnID, Sequences: alignment}
				fmt.Printf("\rRead %d alignments.", numAln)
				fmt.Printf("\r alignment ID: %s", alnID)
			}
		}
		fmt.Printf(" Total alignments %d\n", numAln)
	}
	go read()
	return
}

//AssembleAlignments maps strain sequences to their respective clusters
func AssembleAlignments(aln Alignment, clustermap map[string]string) (alnMap map[string][]seq.Sequence) {
	Sequences := aln.Sequences
	alnMap = make(map[string][]seq.Sequence)
	//double check for duplicate sequences ....
	duplicates := make(map[string]bool)
	for _, s := range Sequences {
		_, _, strain := getNames(s.Id)
		_, found := duplicates[strain]
		if found {
			continue
		} else {
			duplicates[strain] = true
			cluster, found := clustermap[strain]
			if found {
				alnMap[cluster] = append(alnMap[cluster], s)
			}
		}
	}
	return
}

func getNames(s string) (geneName, genomePos, genomeName string) {
	terms := strings.Split(s, " ")
	geneName = terms[0]
	genomePos = terms[1]
	genomeName = terms[2]
	return
}

//cAlignment is a channel structure to pump a gene alignment for a cluster into
type cAlignment struct {
	clusterID string         //ID for cluster
	Sequences []seq.Sequence //the gene alignment for the cluster
}

//WriteClusterMSA write the MSA file for the cluster
func WriteClusterMSA(c cAlignment, outdir string, MSAname string) {
	cluster := "cluster" + c.clusterID
	clusterpath := filepath.Join(outdir, cluster)
	MSA := filepath.Join(clusterpath, MSAname+"_"+cluster)
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile(MSA, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	for _, s := range c.Sequences {
		f.WriteString(">" + s.Id + "\n")
		f.Write(s.Seq)
		f.WriteString("\n")
	}
	f.WriteString("=\n")
	f.Close()
}
