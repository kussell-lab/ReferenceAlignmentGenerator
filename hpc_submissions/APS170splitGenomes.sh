#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=16
#SBATCH --time=1:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=APS170splitgenome
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS170splitgenome_slurm%j.out


##INPUTS
ARCHIVE=/scratch/aps376/recombo/APS170_SA_Archive
MSA=${ARCHIVE}/MSA_SA_MASTER_GAPFILTERED
list=${ARCHIVE}/strain_list

mkdir -p ${OUTDIR}

echo "Loading modules."
module load go/1.15.7
module load singularity/3.6.4

##aliases for singularity

##Making the AssemblyAlignmentGenerator and ReferenceAlignmentGenerator in path
echo "Making everything in path."
#mcorr
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin

#ReferenceAlignmentGenerator
export PATH=$PATH:~/opt/AssemblyAlignmentGenerator/
export PATH=$PATH:~/opt/ReferenceAlignmentGenerator

##set perl language variable; this will give you fewer warnings
export LC_ALL=C

echo "let's rock"
mkdir -p ${ARCHIVE}/corethreshold95
splitGenome ${MSA} ${list} 95 ${ARCHIVE}/corethreshold95 --threads=16 --num-cpu=16
