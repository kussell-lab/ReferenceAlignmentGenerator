# ReferenceAlignmentGenerator
This program uses [SMALT](http://www.sanger.ac.uk/science/tools/smalt-0) to align reads to the reference genome, and then 
uses [SAMtools](https://github.com/samtools/samtools) to generate a consensus genome for each read file. 
The alignments of genes are extracted, alignments with gaps removed, and alignments for core and accessory genes are outputted.
There are several additional features that have been added, as described in "Usage".

# Installation/Usage

(1) For generating an XMFA file using reference-guided alignment, you can clone this github repository to your workstation, then
install the program RefAligner via pip:
* `pip install ~/go/src/github.com/apsteinberg/ReferenceAlignmentGenerator/RefAlignGenerate`

The following dependencies need to be installed:
* [SMALT](http://www.sanger.ac.uk/science/tools/smalt-0)
* [SAMtools](https://github.com/samtools/samtools)
*  [Parallel](https://www.gnu.org/software/parallel/)

As well as several in-house developed programs found in this repo:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/GenomicConsensus`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/FilterGaps`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/CollectGeneAlignments`
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/SplitGenome`

After installation, enter `RefAligner --help` into your terminal to see inputs. For an example, see the RUN_ME.sh
file in the "hinfluenzae_example" subdirectory. For large jobs (>500 genomes), 
we recommend doing each of the steps that the RefAligner does separately, because there may be an issue with
one step that creates downstream problems. See the following file for the step-by-step process:

`~/go/src/github.com/apsteinberg/ReferenceAlignmentGenerator/RefAlignGenerate/RefAlignGenerate/RefAligner/cli.py`


(2) To split XMFA files into XMFA files which only include genes which are in a certain percentage of strains,
(e.g., 50-70% of strains) use the GeneBins program, which can be installed via:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/GeneBins`

After installation, enter `GeneBins --help` into your terminal to see inputs/outputs.

(3) To measure the frequency of gaps in an XMFA file generated using the RefAligner, use MeasureGaps, which can be installed
via:
* `go get -u github.com/kussell-lab/ReferenceAlignmentGenerator/MeasureGaps`

After installation, enter `MeasureGaps --help` into your terminal for inputs/outputs.
