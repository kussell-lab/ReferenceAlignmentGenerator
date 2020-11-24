#!user/bin/env python3

import os
"""
commandline program for fetching SRAs from NCBI and aligning them to a reference genome
"""

def map2reference(accession, wrkdir, reference):
    j = 0
    read1 = os.path.join(wrkdir, accession+"_1.fastq")
    read2 = os.path.join(wrkdir, accession + "_2.fastq")
    ##sort the alignment file
    bam = os.path.join(wrkdir, accession + ".sorted.bam")
    os.system("smalt map -n 1 %s %s %s | samtools sort -T /temp -o %s" % (reference, read1, read2, bam))
    ##multi-way pileup to consensus genome
    pile = os.path.join(wrkdir, accession+".pileup.fasta")
    os.system("samtools mpileup -f %s %s | GenomicConsensus --fna_file %s > %s" % (reference, bam, reference, pile))
    os.system("rm " + bam)
    if os.path.isfile(pile):
        size = os.stat(pile).st_size
        if size != 0:
            os.system("rm " + read1)
            os.system("rm " + read2)
            j = 1
    return j

