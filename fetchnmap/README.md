### fetchnmap
This program downloads SRA files from NCBI using fasterq-dump, uses SMALT to align to reference genome, and SAMtools to 
generate a consensus genome for each read file. It then "cleans up" by deleting ".fastq" files from the working directory
which can build up and take up prohibitive amounts of storage space.

## Installation & Usage
Like ReferenceAlignmentGenerator, this program requires SMALT, SAMtools, Parallel, and fasterq-dump. It also requires an
in-house developed program:

- go get -u github.com/kussell-lab/go-misc/cmd/GenomicConsensus

After the above requirements have been installed