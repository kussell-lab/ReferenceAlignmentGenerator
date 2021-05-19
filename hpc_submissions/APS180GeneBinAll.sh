#!/bin/bash
DATE=210330
ARCHIVE=/scratch/aps376/recombo/APS169_SP_Archive
OUTDIR=/scratch/aps376/recombo/APS180_SP_Archive
MSA=${ARCHIVE}/MSA_SP_MASTER_GAPFILTERED
list=${ARCHIVE}/strain_list
SUBMITDIR=/scratch/aps376/recombo/APS180genebins/${DATE}_submissions
SLURMDIR=/scratch/aps376/recombo/APS180genebins/${DATE}_slurm

mkdir -p ${SUBMITDIR}
mkdir -p ${SLURMDIR}

for i in {0..8}
do
  min=$(expr $i \* 10)
  j=$(expr $i + 1)
  max=$(expr $j \* 10)
  echo "submitting ${min}-${max}"
  jobfile=${SUBMITDIR}/bin_${min}-${max}.sh
  echo "#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=16
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

GeneBins ${MSA} ${list} ${min} ${max} ${OUTDIR}/bin${min}-${max} --threads=16 --num-cpu=16" > $jobfile
    sbatch "$jobfile"
    echo "I'm taking a 1 second break"
    sleep 1 #pause the script for a second so we don't break the cluster with our magic
done


