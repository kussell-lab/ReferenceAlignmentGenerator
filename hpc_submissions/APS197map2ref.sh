#!/bin/bash
##first write for one replicate
##then write for multiple

##make job directory
DATE=0420
SPECIES=EC
PROJDIR=/scratch/aps376/recombo
JOBDIR=${PROJDIR}/APS197map2ref
WRKDIR=${PROJDIR}/APS197_ecoli
REF=${PROJDIR}/APS197_EC_Archive/reference/GCF_000005845.2_ASM584v2_genomic.fna
LISTS=$JOBDIR/piles_tbc
SUBMITDIR=${DATE}_submissions
SLURMDIR=${DATE}_slurm

mkdir -p $SUBMITDIR
mkdir -p $SLURMDIR
mkdir -p $WRKDIR

for line in 0; do
#for line in {1..199}; do #for line in 0
  echo "submitting list ${line}"
  jobfile=$SUBMITDIR/APS197map2ref_${line}.sh

  echo "#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=2
#SBATCH --time=24:00:00
#SBATCH --mem=2GB
#SBATCH --job-name=${SPECIES}map_${line}
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
bash ${JOBDIR}/ConvertMapHPC.sh ${LISTS}/piles_TBC_${line} $WRKDIR $REF" >$jobfile
  sbatch "$jobfile"
  echo "I'm taking a 2 second break"
  sleep 2 #pause the script for a second so we don't break the cluster with our magic
done
