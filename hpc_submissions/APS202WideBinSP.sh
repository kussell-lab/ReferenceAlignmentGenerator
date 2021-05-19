#!/bin/bash
DATE=210517
ARCHIVE=/scratch/aps376/recombo/APS169_SP_Archive
OUTDIR=/scratch/aps376/recombo/APS202_SP_Archive/pangenome_bins
MSA=${ARCHIVE}/MSA_SP_MASTER_GAPFILTERED
list=${ARCHIVE}/strain_list
SUBMITDIR=/scratch/aps376/recombo/APS202genebins/${DATE}_submissions
SLURMDIR=/scratch/aps376/recombo/APS202genebins/${DATE}_slurm

mkdir -p ${SUBMITDIR}
mkdir -p ${SLURMDIR}

bins=(0 15 55 95)
for i in {0..2}
do
  min=${bins[$i]}
  max=${bins[$i+1]}
  echo "submitting ${min}-${max}"
  jobfile=${SUBMITDIR}/bin_${min}-${max}.sh
  echo "#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=8
#SBATCH --time=4:00:00
#SBATCH --mem=4GB
#SBATCH --job-name=bin${min}-${max}
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=${SLURMDIR}/slurm%j_bin_${min}-${max}.out


##INPUTS

echo \"Loading modules.\"
module load go/1.15.7
module load singularity/3.6.4

##aliases for singularity

##Making the AssemblyAlignmentGenerator and ReferenceAlignmentGenerator in path
echo \"Making everything in path.\"
#mcorr
export PATH=\$PATH:\$HOME/go/bin:\$HOME/.local/bin

#ReferenceAlignmentGenerator
export PATH=\$PATH:~/opt/AssemblyAlignmentGenerator/
export PATH=\$PATH:~/opt/ReferenceAlignmentGenerator

##set perl language variable; this will give you fewer warnings
export LC_ALL=C

echo \"let's rock\"
mkdir -p ${OUTDIR}

GeneBins ${MSA} ${list} ${min} ${max} ${OUTDIR}/bin${min}-${max} --threads=16 --num-cpu=8" > $jobfile
    sbatch "$jobfile"
    echo "I'm taking a 1 second break"
    sleep 1 #pause the script for a second so we don't break the cluster with our magic
done


