"""
This script downloads read files (SRA format) from NCBI FTP.
Created by Mingzhi Lin (mingzhi9@gmail.com).
Modified by Asher Preska Steinberg (asherps@gmail.com)
to accommodate updates to NCBI SRA database as of 200727
This version just relies on fasterq from sra tool-kit and converts to fastq in one step
Inputs: 
    (1) a file containing a list of accessions;
    (2) working directory.
"""
import sys
import os
import ftplib
import shutil as shu
from tqdm import tqdm
def main():
    """ main function """
    accession_list_file = sys.argv[1]
    working_dir = sys.argv[2]

    #accession_list_file = '/Users/asherpreskasteinberg/Desktop/code/recombo/salmonella/SRA_files/sra_accession_Kentucky_1'
    #working_dir = '/Users/asherpreskasteinberg/Desktop/fetchsra_test/'
    ncbi='/home/aps376/ncbi/public/sra'

    # read the list of accessions.
    accession_list = []
    with open(accession_list_file, 'r') as reader:
        for line in reader:
            accession_list.append(line.rstrip())

    print("Fetching SRA from NCBI via sra toolkit")
    for accession in tqdm(accession_list):
        # read path: SRR/SRR000/SRR000001/SRR000001.sra
        #os.system('fasterq-dump '+str(accession)+' -O '+working_dir+' -t '+'$SCRATCH'+' -p')
        ##cluster -- trying a speedup with -e to increase number of threads, ... we will see how it works
        os.system('fasterq-dump ' + str(accession) + ' -O ' + working_dir + ' -t ' + '$BEEGFS'+ ' -e 10')
        ##local
        #os.system('fasterq-dump ' + str(accession) + ' -O ' + working_dir + ' -t ' + working_dir)

    print("Completed downloading %d read files." % len(accession_list))

if __name__ == "__main__":
    main()
