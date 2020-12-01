### fetchnmap
This program downloads SRA files from NCBI using fasterq-dump, uses SMALT to align to reference genome, and SAMtools to 
generate a consensus genome for each read file. It then "cleans up" by deleting ".fastq" files from the working directory
which can build up and take up prohibitive amounts of storage space.

## Installation
To install the package, take the following steps:\
(1) Pull the fetchnmap_v2 directory from github using git pull or copy it to your computer\
(2) Use the command "pip install /path/fetchnmap_v2" to install the package (where "path" is the path to the directory)\
(3) Install the following requirements: SMALT, SAMtools, GNU Parallel, and fasterq-dump. Also install the following in-house
program using the command:
- go get -u github.com/kussell-lab/go-misc/cmd/GenomicConsensus

## Usage

the following commandline tools can then be used to download and align sequences:\
(1)  fetchnzip:
This downloads SRA files from NCBI and stores them as .gz files to conserve space. For help type "fetchnzip --help" into commandline.\
Basic input is as follows:\

fetchnzip --tmp=TMP accession_list working_dir

positional arguments:\
accession_list: list of read accessions of SRA files (to be downloaded from NCBI)
working_dir: the working space and output directory

optional arguments:\
--tmp: can specify directory for temporary files created by fasterq-dump (could offer a speedup if this is an SSD, e.g.)

(2) mapnclean:\
Maps the downloaded fastq files to the reference genome, then removes fastqs to conserve storage space.\
Basic input is as follows:\

mapnclean accession_list working_dir reference

positional arguments:\
accession_list: see fetchnzip description above\
working_dir: see fetchnzip description above\
reference: reference genome as a .fna file