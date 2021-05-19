#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=1
#SBATCH --time=12:00:00
#SBATCH --mem=2GB
#SBATCH --job-name=APS187EC
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS187EC_scratchcleanup_slurm%j.out

rm -r $SCRATCH/recombo/APS187_EC_Archive