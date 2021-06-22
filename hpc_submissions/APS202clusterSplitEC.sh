#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=4
#SBATCH --time=4:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=EC_clustersplit
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=EC_clustersplit_slurm%j.out


##INPUTS
cluster_dict=/scratch/aps376/recombo/APS197_EC_Archive/cluster_list
msa_dir=/scratch/aps376/recombo/APS202_EC_Archive
out=${msa_dir}

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
for bin in "0-15" "15-55" "55-95";
do
  mkdir -p ${out}/bin_${bin}
  clusterSplit ${msa_dir}/bin_${bin}/MSA_${bin} ${out}/bin_${bin} ${cluster_dict} --num-cpu=4
done

