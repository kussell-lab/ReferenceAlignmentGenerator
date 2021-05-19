#!/bin/bash
##first write for one replicate
##then write for multiple

##make job directory
DATE=0409
PROJDIR=/scratch/aps376/recombo
JOBDIR=${PROJDIR}/APS184map2ref
WRKDIR=${PROJDIR}/APS184_cjejuni
REF=${PROJDIR}/APS184_CJ_Archive/reference/GCF_000009085.1_ASM908v1_genomic.fna
LISTS=$JOBDIR/piles_tbc
SUBMITDIR=${DATE}_submissions
SLURMDIR=${DATE}_slurm

mkdir -p $SUBMITDIR
mkdir -p $SLURMDIR
mkdir -p $WRKDIR

##will change to 0 to 9 once confirmed that it werks
#for line in {0};
for line in 0; do
#for line in {1..199}; do #for line in 0
  echo "submitting list ${line}"
  jobfile=$SUBMITDIR/APS184map2ref_${line}.sh

  echo "#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=2
#SBATCH --time=24:00:00
#SBATCH --mem=2GB
#SBATCH --job-name=CJmap_${line}
#SBATCH --mail-type=END,FAIL
#SBATCH --mail-user=aps376@nyu.edu
#SBATCH --output=${SLURMDIR}/slurm%j_map2ref_${line}.out

echo \"Loading modules.\"
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
cd ${PROJDIR}
source venv/bin/activate;
export OMP_NUM_THREADS=\$SLURM_CPUS_PER_TASK;

#put go on path
export PATH=\$PATH:\$HOME/go/bin:\$HOME/.local/bin

##set perl language variable; this will give you fewer warnings
export LC_ALL=C

echo \"let's rock\"
cd ${WRKDIR}
bash ${JOBDIR}/ConvertMap.sh ${LISTS}/piles_TBC_${line} $WRKDIR $REF" >$jobfile
  sbatch "$jobfile"
  echo "I'm taking a 2 second break"
  sleep 2 #pause the script for a second so we don't break the cluster with our magic
done
