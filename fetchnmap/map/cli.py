#!user/bin/env python3

import argparse
from tqdm import tqdm
import os
from . import functions as fn
"""
commandline program for fetching SRAs from NCBI and aligning them to a reference genome
"""

def main():
    parser = argparse.ArgumentParser(description="Fetchs SRA files from NCBI and aligns them to a reference genome;\
                                                  output is a FASTA file (suffix is '.pileup.fasta'. Additional requirements:\
                                                    fasterq-dump, SMALT, SAMtools, GenomicConsensus (see README.MD)")
    parser.add_argument("accession_list", help="a list of read accessions of which SRA files can be downloaded from NCBI SRA database")
    parser.add_argument("working_dir", help="the working space and output directory")
    parser.add_argument("reference", help="the reference genome (FASTA format) that will be used for read mapping")
    parser.add_argument("--tmp", help="Can specify the directory for temp files created by fasterq-dump \
                                        (may speed up downloads)")
    ##define commandline args as variables
    args = parser.parse_args()
    accession_list_file = args.accession_list
    wrkdir = args.working_dir
    reference = args.reference
    tmp = args.tmp

    #build index of k-mer words for the reference genome
    os.system("smalt index %s %s" % (reference, reference))
    #index and extract FASTA file
    os.system("samtools faidx %s" % reference)

    #read the list of accessions
    accession_list = []
    with open(accession_list_file, 'r') as reader:
        for line in reader:
            accession_list.append(line.rstrip())
    ##download from NCBI
    dwnlds = 0
    maps = 0
    if tmp:
        tmpdir = ' -t ' + tmp
    else:
        tmpdir= ''
    for accession in tqdm(accession_list):
        os.system('fasterq-dump ' + str(accession) + ' -O ' + wrkdir + tmpdir)
        read1 = os.path.join(wrkdir, accession + "_1.fastq")
        read2 = os.path.join(wrkdir, accession + "_2.fastq")
        if os.path.isfile(read1) and os.path.isfile(read2):
            j = fn.map2reference(accession, wrkdir, reference)
            dwnlds = dwnlds + 1
            maps = maps + j
    print("%d successful SRA downloads" % dwnlds)
    print("%d reads mapped to reference genome" % maps)






if __name__ == "__main__":
    main()