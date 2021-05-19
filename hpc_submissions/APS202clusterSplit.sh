#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=4
#SBATCH --time=4:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=SP_clustersplit
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=SP_clustersplit_slurm%j.out


##INPUTS
out=/scratch/aps376/recombo/APS202_SP_Archive
cluster_dict=/scratch/aps376/recombo/APS169_SP_Archive/cluster_list
msa_dir=/scratch/aps376/recombo/APS180_SP_Archive/widerbins

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
for bin in "0-20" "20-40" "40-60" "60-80" "80-100";
do
  mkdir -p ${out}/bin_${bin}
  clusterSplit ${msa_dir}/bin${bin}/MSA_${bin} ${out}/bin_${bin} ${cluster_dict} --num-cpu=4
done

