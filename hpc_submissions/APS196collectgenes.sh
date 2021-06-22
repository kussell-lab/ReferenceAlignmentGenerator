#!/bin/bash
#SBATCH --job-name=SSgenes
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=16
#SBATCH --mem=4GB
#SBATCH --time=24:00:00
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS196collectgenes_slurm%j.out


##INPUTS
wrkd=${SCRATCH}/recombo/APS196_ssonnei
archive=${SCRATCH}/recombo/APS196_SS_Archive
fasta=${archive}/reference/GCF_002247485.1_ASM224748v1_genomic.fna
gff=${archive}/reference/GCF_002247485.1_ASM224748v1_genomic.gff
sra_list=${archive}/strain_list

echo "Loading modules."
module purge
module load go/1.15.7

##Making the AssemblyAlignmentGenerator and ReferenceAlignmentGenerator in path
echo "Making everything in path."
#mcorr
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin

#ReferenceAlignmentGenerator
export PATH=$PATH:~/opt/AssemblyAlignmentGenerator/
export PATH=$PATH:~/opt/ReferenceAlignmentGenerator

##set perl language variable; this will give you fewer warnings
export LC_ALL=C


##MSA stands for multi sequence alignment in the below
  #the '$1' command tells it to grab the argument of pipe_dream

echo "let's rock"

CollectGeneAlignments ${sra_list} ${gff} ${wrkd} ${archive}/MSA_SS_MASTER --appendix ".pileup.fasta" --progress