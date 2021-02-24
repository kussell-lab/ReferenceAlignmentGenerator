#!/bin/bash
# This script maps reads to the reference genome to generate the consensus genomic sequence.
# Created by Asher Preska Steinberg (apsteinberg@nyu.edu) based on scripts by Mingzhi Lin (mingzhi9@gmail.com)
# It requires samtools (https://github.com/samtools/samtools) and smalt (http://www.sanger.ac.uk/science/tools/smalt-0).
# Inputs:
#   (1) the accession list file;
#   (2) the working directory;
#   (3) the reference genome.
#   (4) number of threads given to mapping/compression from sam to bam
# Output:
#   the consensus genomic sequences in FASTA format.

function map2reference {
    READ=$1
    WORKING_DIR=$2
    REFERENCE=$3
    LABEL=pileup
    sra_file=${WORKING_DIR}/${READ}/${READ}.sra
    prefetch ${READ}
    fastq-dump -O ${WORKING_DIR} --split-3 ${sra_file}
    smalt map -n 1 ${REFERENCE} ${WORKING_DIR}/${READ}_1.fastq ${WORKING_DIR}/${READ}_2.fastq | \
        samtools sort -T /tmp -o ${WORKING_DIR}/${READ}.sorted.bam
        samtools mpileup -f ${REFERENCE} ${WORKING_DIR}/${READ}.sorted.bam | \
        GenomicConsensus --fna_file ${REFERENCE} > ${WORKING_DIR}/${READ}.${LABEL}.fasta
        rm ${WORKING_DIR}/${READ}.sorted.bam
        rm ${WORKING_DIR}/${READ}_1.fastq
        rm ${WORKING_DIR}/${READ}_2.fastq
        rm -r ${WORKING_DIR}/${READ}
}

export -f map2reference

accession_list_file=$1
working_dir=$2
reference=$3
smalt index ${reference} ${reference}
samtools faidx ${reference}
parallel map2reference {} ${working_dir} ${reference} :::: ${accession_list_file}