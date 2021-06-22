# ReferenceAlignmentGenerator
This program uses [SMALT](http://www.sanger.ac.uk/science/tools/smalt-0) to align reads to the reference genome, and then uses [SAMtools](https://github.com/samtools/samtools) to generate a consensus genome for each read file. The alignments of genes are extracted.
There are several additional features that have been added, as described in "Usage".
# Installation
This program requires [SMLAT](http://www.sanger.ac.uk/science/tools/smalt-0), [SAMtools](https://github.com/samtools/samtools), and [Parallel](https://www.gnu.org/software/parallel/). It also needs two in-host developed programs:

* `go get -u github.com/kussell-lab/go-misc/cmd/GenomicConsensus`
* `go get -u github.com/kussell-lab/go-misc/cmd/CollectGeneAlignments`

# Usage

(1) For generating an XMFA file using reference-guided alignment, you can clone this github repository to your workstation, then
install the program ReferenceAlignmentGenerate via pip:


`ReferenceAlignmentGenerate <accession list file> <working directory> <the reference genome> <the gff file> <the output file>`
* `<accession list file>` is a file containing read accessions that can be downloaded from NCBI [SRA](https://www.ncbi.nlm.nih.gov/sra) database;
* `<working director>` is the working space; 
* `<the reference genome>` is the reference genome for read mapping;
* `<the gff file>` is the corresponding GFF file of the reference genome;
* `<the output file>` contains the output alignment results.
