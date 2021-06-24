#!/bin/bash
# This script runs the haemophilus influenzae example
# Created by Asher Preska Steinberg (apsteinberg@nyu.edu)

##INPUTS
refdir=$HOME/go/src/github.com/apsteinberg/ReferenceAlignmentGenerator/hinfluenzae_example
wrkdir=${refdir}/hinfluenzae
sra_list=${refdir}/sra_list.txt
fasta=${refdir}/reference/GCF_000767075.1_ASM76707v1_genomic.fna
gff=${refdir}/reference/GCF_000767075.1_ASM76707v1_genomic.gff
out=MSA_HI

mkdir -p ${wrkdir}
RefAligner ${sra_list} ${wrkdir} ${fasta} ${gff} ${out}
# to run without downloads, comment out line 14 and uncomment line 16 ...
#RefAligner ${sra_list} ${wrkdir} ${fasta} ${gff} ${out} --skipdownloads=True