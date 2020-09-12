#!/bin/bash
# This script maps reads to the reference genome to generate the concensus genomic sequence.
# Created by Mingzhi Lin (mingzhi9@gmail.com).
# It requires samtools (https://github.com/samtools/samtools) and smalt (http://www.sanger.ac.uk/science/tools/smalt-0).
# Inputs:
#   (1) the accession list file;
#   (2) the working directory;
#   (3) the reference genome.
# Output:
#   the consensus genomic sequences in FASTA format.

function map2reference {
    READ='ERR3144502'
    WORKING_DIR='/scratch/aps376/APS138_test'
    REFERENCE='GCF_000009085.1_ASM908v1_genomic.fna'
    LABEL=pileup
    smalt map -n 1 ${REFERENCE} ${WORKING_DIR}/${READ}.fastq | \
        samtools sort -T /tmp -o ${WORKING_DIR}/${READ}.sorted.bam
        samtools mpileup -f ${REFERENCE} ${WORKING_DIR}/${READ}.sorted.bam | \
        GenomicConsensus --fna_file ${REFERENCE} > ${WORKING_DIR}/${READ}.${LABEL}.fasta 
        rm ${WORKING_DIR}/${READ}.sorted.bam
}

export -f map2reference

accession_list_file=$1
working_dir=$2
reference=$3
smalt index ${reference} ${reference}
samtools faidx ${reference}
parallel map2reference {} ${working_dir} ${reference} :::: ${accession_list_file}
