#!/bin/bash
# This script converts .SRA file to .fastq using fastq-dump in sra-toolkits (https://github.com/ncbi/sra-tools/wiki).
# Created by Mingzhi Lin (mingzhi9@gmail.com).
# Inputs: 
#   (1) the accession list file;
#   (2) the working directory containing the SRA files.
# Output:
#   the converted fastq files will be saved in the working directory.

accession_file=$1
working_dir=$2

function convertSRA2Fastq {
    accession=$1
    working_dir=$2
    sra_file=${working_dir}/${accession}.sra
    fastq-dump -O ${working_dir} --split-3 ${sra_file}
}
export -f convertSRA2Fastq

parallel convertSRA2Fastq {} ${working_dir} :::: ${accession_file}
