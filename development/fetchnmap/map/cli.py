#!user/bin/env python3
import uuid
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

    # #build index of k-mer words for the reference genome
    # os.system("smalt index %s %s" % (reference, reference))
    # #index and extract FASTA file
    # os.system("samtools faidx %s" % reference)

    #read the list of accessions
    accession_list = []
    with open(accession_list_file, 'r') as reader:
        for line in reader:
            accession_list.append(line.rstrip())
    ##download from NCBI
    dwnlds = 0
    maps = 0
    if args.tmp is not None:
        tmpdir = ' -t ' + tmp
        print("set temp dir to: %s" % tmp)
    else:
        tmpdir= ''

    scriptdir = os.path.dirname(os.path.abspath(__file__))
    mapRead2Ref = os.path.join(scriptdir, "mapnclean")
    count = 0
    totalcount = 0
    numaccessions = len(accession_list)
    #print(numaccessions)

    ##temp file for mapread2ref
    ##generate a unique tag for this file
    tag = str(uuid.uuid4())
    reads = os.path.join(wrkdir, "reads_"+tag)
    ##open the file
    f = open(reads, "a+")
    for accession in tqdm(accession_list):
        totalcount = totalcount+1
        os.system('fasterq-dump ' + str(accession) + ' -O ' + wrkdir + tmpdir)
        read1 = os.path.join(wrkdir, accession + "_1.fastq")
        read2 = os.path.join(wrkdir, accession + "_2.fastq")
        if os.path.isfile(read1) and os.path.isfile(read2):
            count = count + 1
            f.write(accession+"\n")
            if count == 10 or totalcount == numaccessions:
                f.close()
                os.system("bash %s %s %s %s" % (mapRead2Ref, reads, wrkdir, reference))
                ##reset count
                #reads = []
                os.remove(reads)
                count = 0
                ##generate another reads file if you're not done
                if totalcount != numaccessions:
                    ##generate a unique tag for this file
                    tag = str(uuid.uuid4())
                    reads = os.path.join(wrkdir, "reads_" + tag)
                    f = open(reads, "a+")
            dwnlds = dwnlds + 1
    print("%d downloads and maps" % dwnlds)
    #print("%d reads mapped to reference genome" % maps)






if __name__ == "__main__":
    main()