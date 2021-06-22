#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=8
#SBATCH --time=6:00:00
#SBATCH --mem=16GB
#SBATCH --job-name=SAgalfilter
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS185_filter_slurm%j.out


##INPUTS
projdir=/scratch/aps376/recombo
archive=${projdir}/APS185_SAgal_Archive
msa=${archive}/MSA_SAgal_MASTER

echo "Loading modules."
# Load modules
module purge
module load samtools/intel/1.11
module load sra-tools/2.10.9
module load parallel/20201022
module load python/intel/3.8.6
module load smalt/intel/0.7.6
module load go/1.15.7
module load bowtie2/2.4.2
module load bedtools/intel/2.29.2

#activate virtual environment
projdir=/scratch/aps376/recombo
cd ${projdir}
source venv/bin/activate;
export OMP_NUM_THREADS=$SLURM_CPUS_PER_TASK;

#put go on path
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin

##set perl language variable; this will give you fewer warnings
export LC_ALL=C


##MSA stands for multi sequence alignment in the below
  #the '$1' command tells it to grab the argument of pipe_dream

echo "let's rock"

FilterGaps ${msa} ${archive} --threads=20