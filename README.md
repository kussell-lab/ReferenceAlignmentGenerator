# ReferenceAlignmentGenerator
This program uses [SMALT](http://www.sanger.ac.uk/science/tools/smalt-0) to align reads to the reference genome, and then 
uses [SAMtools](https://github.com/samtools/samtools) to generate a consensus genome for each read file. 
The alignments of genes are extracted, alignments with gaps removed, and alignments for core and accessory genes are outputted.
There are several additional features that have been added, as described in "Usage".

# Installation

For basic installation (making reference-guided alignments), clone this github repository to your workstation, and
install the program RefAligner via pip:
* `pip install ~/go/src/github.com/apsteinberg/ReferenceAlignmentGenerator/RefAlignGenerate`

The following dependencies need to be installed:
* [SMALT](http://www.sanger.ac.uk/science/tools/smalt-0)
* [SAMtools](https://github.com/samtools/samtools)
*  [Parallel](https://www.gnu.org/software/parallel/)
* [sra-tools](https://github.com/ncbi/sra-tools/wiki/01.-Downloading-SRA-Toolkit)

As well as several in-house developed programs found in this repo:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/GenomicConsensus`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/FilterGaps`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/CollectGeneAlignments`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/SplitGenome`

We have tested this on Mac OS 11.4 and on the NYU HPC (which uses Slurm) using Python 3.8 and 3.9, 
Go version 1.15.7 and 1.16, SMALT 0.7.6, SAMtools 1.11, sra-tools 2.10.9, parallel 20201022

# Basic Usage

After installation, enter `RefAligner --help` into your terminal to see inputs. For an example, see the RUN_ME.sh
file in the "hinfluenzae_example" subdirectory. For large jobs (>500 genomes), 
we recommend doing each of the steps that the RefAligner does separately, because there may be an issue with
one step that creates downstream problems. See the following file for the step-by-step process:

`~/go/src/github.com/apsteinberg/ReferenceAlignmentGenerator/RefAlignGenerate/RefAlignGenerate/RefAligner/cli.py`

# Additional Features

(1) To split XMFA files into XMFA files which only include genes which are in a certain percentage of strains,
(e.g., 50-70% of strains) use the GeneBins program, which can be installed via:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/GeneBins`

After installation, enter `GeneBins --help` into your terminal to see inputs/outputs.

(2) To measure the frequency of gaps in an XMFA file generated using the RefAligner, use MeasureGaps, which can be installed
via:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/MeasureGaps`

After installation, enter `MeasureGaps --help` into your terminal for inputs/outputs.

# Example

For an example of how to run RefAligner, see the directory "hinfluenzae_example" in this repo.
To run the example, run the RUN_ME.sh script. Downloads and read mapping are the time limiting step. To run the full 
example takes ~2-3 hours on a standard PC. To skip the download step and run everything else, use the flag 
`--skipdownloads=True`. This should then only take ~5 minutes to run.
Example outputs are given in the folder as well.
