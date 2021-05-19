#!/bin/bash
#SBATCH --job-name=NMgenes
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=16
#SBATCH --mem=4GB
#SBATCH --time=24:00:00
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS183collectgenes_slurm%j.out


##INPUTS
wrkd=${SCRATCH}/recombo/APS183_nmeningitidis
archive=${SCRATCH}/recombo/APS183_NM_Archive
fasta=${archive}/reference/GCF_008330805.1_ASM833080v1_genomic.fna
gff=${archive}/reference/GCF_008330805.1_ASM833080v1_genomic.gff
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

CollectGeneAlignments ${sra_list} ${gff} ${wrkd} ${archive}/MSA_NM_MASTER --appendix ".pileup.fasta" --progress