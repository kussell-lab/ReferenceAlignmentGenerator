# ReferenceAlignmentGenerator
This program uses [SMLAT](http://www.sanger.ac.uk/science/tools/smalt-0) to align reads to the reference genome, and then uses [SAMtools](https://github.com/samtools/samtools) to generate a consensus genome for each read file. The aligned genes are extracted.

# Installation
This program requires [SMLAT](http://www.sanger.ac.uk/science/tools/smalt-0), [SAMtools](https://github.com/samtools/samtools), and [Parallel](https://www.gnu.org/software/parallel/). It also needs two in-host developed programs:
* `go get -u github.com/kussell-lab/go-misc/cmd/GenomicConsensus`
* `go get -u github.com/kussell-lab/go-misc/cmd/CollectGeneAlignments`

# Usage
`ReferenceAlignmentGenerate <accession list file> <working directory> <the reference genome> <the gff file> <the output file>`
* `<accession list file>` is a file containing read accessions that can be downloaded from NCBI [SRA](https://www.ncbi.nlm.nih.gov/sra) database;
* `<working director>` is the working space; 
* `<the reference genome>` is the reference genome for read mapping;
* `<the gff file>` is the corresponding GFF file of the reference genome;
* `<the output file>` contains the output alignment results.
