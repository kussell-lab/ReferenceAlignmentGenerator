package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kussell-lab/biogo/feat/gff"
	"github.com/kussell-lab/biogo/seq"
)

// Splitter is a interface for splitting genome sequences into multiple parts.
type Splitter interface {
	Split(s *seq.Sequence) []*seq.Sequence
}

// GeneSplitter splites a genome into multiple genes.
type GeneSplitter struct {
	features []*gff.Record
}

// NewGeneSplitter returns a new GeneSplitter
func NewGeneSplitter(features []*gff.Record) *GeneSplitter {
	cdsFeatures := []*gff.Record{}
	for _, feat := range features {
		if feat.Feature == "CDS" {
			cdsFeatures = append(cdsFeatures, feat)
		}
	}
	return &GeneSplitter{features: cdsFeatures}
}

// Split splits a genomic sequence into multiple genes.
func (spliter *GeneSplitter) Split(s *seq.Sequence) []*seq.Sequence {
	return split(s, spliter.features)
}

func split(s *seq.Sequence, features []*gff.Record) []*seq.Sequence {
	subSeqs := []*seq.Sequence{}
	for _, rec := range features {
		if s.Id != rec.SeqName {
			//            continue
		}
		start := rec.Start - 1
		if start < len(s.Seq) {
			end := rec.End
			if end > len(s.Seq) {
				end = len(s.Seq)
			}

			if end > len(s.Seq) {
				end = len(s.Seq)
			}
			ss := &seq.Sequence{}
			ss.Id = getGffID(rec)
			ss.Seq = s.Seq[start:end]

			// add gaps into the end if the sequence is of short.
			length := rec.End - rec.Start + 1 // desired length.
			for len(ss.Seq) < length {
				ss.Seq = append(ss.Seq, '-')
			}

			if rec.Strand == gff.ReverseStrand {
				ss.Seq = seq.Reverse(ss.Seq)
				if rec.Feature == "CDS" {
					ss.Seq = seq.Complement(ss.Seq)
				}
			}

			subSeqs = append(subSeqs, ss)
		}
	}

	return subSeqs
}

func getGffID(rec *gff.Record) string {
	attrs := strings.Split(rec.Attribute, ";")
	for _, att := range attrs {
		keyvalue := strings.Split(att, "=")
		key := keyvalue[0]
		value := keyvalue[1]
		if key == "ID" {
			strand := "+"
			if rec.Strand == gff.ReverseStrand {
				strand = "-"
			}
			return fmt.Sprintf("%s|%s %d%s%d", rec.SeqName, value, rec.Start, strand, rec.End)
		}
	}
	return fmt.Sprintf("%s|start:%d_end:%d", rec.SeqName, rec.Start, rec.End)
}

// UpstreamSplitter split a genome sequence into mutiple upstream region.
type UpstreamSplitter struct {
	UpstreamLen int
	features    []*gff.Record
}

// NewUpstreamSplitter return a new UpstreamSplitter
func NewUpstreamSplitter(features []*gff.Record, upstreamLen int) *UpstreamSplitter {
	us := &UpstreamSplitter{UpstreamLen: upstreamLen}

	genomeFeatures := [][]*gff.Record{}
	current := ""
	currentFeatures := []*gff.Record{}
	for _, feat := range features {
		if current == "" {
			current = feat.SeqName
		}
		if current != feat.SeqName {
			genomeFeatures = append(genomeFeatures, currentFeatures)
			current = feat.SeqName
			currentFeatures = []*gff.Record{}
		}
		currentFeatures = append(currentFeatures, feat)
	}
	genomeFeatures = append(genomeFeatures, currentFeatures)

	upstreamRegions := []*gff.Record{}
	for _, genomeFeats := range genomeFeatures {
		upstreamRegions = append(upstreamRegions, extractUpstreamRegions(genomeFeats, upstreamLen)...)
	}

	us.features = upstreamRegions
	return us
}

// Split splits a genomic sequence into multiple genes.
func (spliter *UpstreamSplitter) Split(s *seq.Sequence) []*seq.Sequence {
	return split(s, spliter.features)
}

// ByStart is a sorting function to sort gff records by start position.
type ByStart []*gff.Record

func (s ByStart) Len() int {
	return len(s)
}
func (s ByStart) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByStart) Less(i, j int) bool {
	return s[i].Start < s[j].Start
}

// extractUpstreamRegions return upstream regions for each gene.
func extractUpstreamRegions(features []*gff.Record, upstreamLen int) []*gff.Record {
	cdsGffs := []*gff.Record{}
	for _, feat := range features {
		if feat.Feature == "CDS" {
			cdsGffs = append(cdsGffs, feat)
		}
	}

	sort.Sort(ByStart(cdsGffs))

	regions := []*gff.Record{}
	for i, feat := range cdsGffs {
		start := 0
		end := -1
		if feat.Strand == gff.ForwardStrand {
			start = feat.Start - upstreamLen
			end = feat.Start - 1
			if start < 0 {
				start = 0
			}

			previousEnd := 1
			if i > 0 {
				previousEnd = cdsGffs[i-1].End + 1
			}
			if start < previousEnd {
				start = previousEnd
			}
		} else {
			start = feat.End + 1
			end = feat.End + upstreamLen
			if i < len(cdsGffs)-1 {
				nextStart := cdsGffs[i+1].Start
				if nextStart < end {
					end = nextStart - 2
				}
			}
		}

		if end > start {
			rec := gff.Record{}
			rec.SeqName = feat.SeqName
			rec.Source = feat.Source
			rec.Attribute = feat.Attribute
			rec.End = end
			rec.Start = start + 1
			rec.Strand = feat.Strand
			rec.Feature = "upstream"
			rec.Frame = feat.Frame
			rec.Score = feat.Score
			regions = append(regions, &rec)
		}
	}

	return regions
}
