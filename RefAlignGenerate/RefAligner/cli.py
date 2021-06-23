#!/usr/bin/env python3
import time
import argparse
import os
"""
this commandline program which executes each step of the ReferenceAlignmentGenerator
"""
def main():
    parser = argparse.ArgumentParser(description="Downloads raw reads from NCBI, aligns them to a\
                                                    reference genome, filters out alignments with >2% gaps,\
                                                 and generates XMFAs for core and accessory genes")
    parser.add_argument("accession_list", help="a list of read accessions for SRA files to be downloaded from the NCBI SRA database")
    parser.add_argument("working_dir", help="the working space")
    parser.add_argument("ref_fasta", help="FASTA file for reference genome")
    parser.add_argument("ref_gff", help="GFF file for reference genome for extracting CDS regions")
    parser.add_argument("out_prefix", type=str, help="prefix for master XMFA")
    parser.add_argument("--threshold", type=int, default=95, help="threshold percentage above which you're considered a core gene (Default: 95)")
    ##define commandline args as variables
    args = parser.parse_args()
    acc_list = args.accession_list
    wrkdir = args.working_dir
    fasta = args.ref_fasta
    gff = args.ref_gff
    msa = args.out_prefix
    t = args.threshold

    start_time = time.time()

    #Step 1: Fetch SRA files, convert to FASTQ, and map to reference genome
    print("Fetching raw reads and map to reference genome ...")
    os.system("ConvertMap.sh %s %s %s" %(acc_list, wrkdir, fasta))

    #Step 2: Extract CDS regions and output alignments into the master XMFA
    print("extracting CDS regions ...")
    geneinputs = (acc_list, gff, wrkdir, msa)
    os.system("CollectGeneAlignments %s %s %s_MASTER --appendix '.pileup.fasta' --progress" % geneinputs)
    #Step 3: Filter gapped sequences
    print("filter gapped sequences ...")
    outdir = os.curdir
    os.system("FilterGaps %s_MASTER %s" % (msa, outdir))
    #Step 4: Split into XMFA files for core and accessory genes
    print("splitting into XMFA files for core and flexible genes ...")
    #name of gapfiltered msa
    gapfiltered = msa + "_MASTER_GAPFILTERED"
    splitinput = (gapfiltered, acc_list, t, outdir)
    os.system("splitGenome %s %s %s %s" % splitinput)
    print("Done with making clusters, time to boogie")
    print("Total run time: %s minutes" % str((time.time() - start_time)/60))

if __name__ == "__main__":
    main()