#!/bin/bash
#SBATCH --job-name=sra-downloads
#SBATCH --cpus-per-task=4
#SBATCH --mem=32000
#SBATCH --time=2:00:00
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu


##INPUTS
#The directory you want create serotype folders in. You need lots of space.
WRKD=/scratch/aps376
#SRC is the directory your SRA accession files and genome files are within.
SRC=/home/aps376/salmonella

echo "Loading modules."
module load git/gnu/2.16.2
module load go/1.10.2 #try go/1.13.6
module load python3/intel/3.6.3 ##do 3.7.3!
module load parallel/20171022
module load prokka/1.12
module load muscle/intel/3.8.31
#module load sra-tools/intel/2.9.6 #try 2.9.6
module load sra-tools/2.10.5
module load samtools/intel/1.6
module load smalt/intel/0.7.6
alias roary='singularity exec /beegfs/work/public/singularity/roary-20181203.simg roary'

##Making the AssemblyAlignmentGenerator and ReferenceAlignmentGenerator in path
echo "Making everything in path."
#mcorr
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin

#ReferenceAlignmentGenerator
export PATH=$PATH:~/opt/AssemblyAlignmentGenerator/
export PATH=$PATH:~/opt/ReferenceAlignmentGenerator

mkdir ${WRKD}/2cladetest

cd ${WRKD}/2cladetest

for sero in Kentucky_1 Kentucky_2 Kentucky_both
do

  ReferenceAlignmentGenerate ${SRC}/SRA_files/sra_accession_$sero ${WRKD}/2cladetest ${SRC}/Reference/GCF_000006945.2_ASM694v2_genomic.fna ${SRC}/Reference/GCF_000006945.2_ASM694v2_genomic.gff ${WRKD}/2cladetest/REFGEN_$sero
done
