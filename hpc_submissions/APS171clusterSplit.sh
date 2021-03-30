#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=4
#SBATCH --time=4:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=NGclustersplit
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=APS171clustersplit_slurm%j.out


##INPUTS
ARCHIVE=/scratch/aps376/recombo/APS171_NG_Archive
cluster_dict=${ARCHIVE}/cluster_list

mkdir -p ${OUTDIR}

echo "Loading modules."
module load go/1.15.7

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
cd ${ARCHIVE}/corethreshold95
clusterSplit MSA_CORE ${ARCHIVE}/corethreshold95 ${cluster_dict} --FLEX_MSA=MSA_FLEX --num-cpu=4

