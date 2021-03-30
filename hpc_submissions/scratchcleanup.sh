#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=1
#SBATCH --time=12:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=APS169splitgenome
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=scratchcleanup_slurm%j.out

rm -r $SCRATCH/recombo/APS168_SARS-CoV-2

