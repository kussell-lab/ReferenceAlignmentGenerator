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

    # open FTP connection.
    # ftp = ftplib.FTP("ftp.ncbi.nlm.nih.gov")
    # ftp.login("anonymous", "mingzhi.lin@nyu.edu")
    # ftp.cwd("sra/sra-instant/reads/ByRun/sra")

    print("Fetching SRA from NCBI via sra toolkit")
    for accession in tqdm(accession_list):
        # read path: SRR/SRR000/SRR000001/SRR000001.sra
        os.system('fasterq-dump '+str(accession)+' -O '+working_dir+' -t '+'$BEEGFS')
       # sra_file_path = "%s/%s/%s/%s.sra" % (accession[:3], accession[:6], accession, accession)
        #sra_file_path = "%s/%s.sra" % (ncbi, accession) #for cluster
        #local_file_path = os.path.join(working_dir, accession + ".sra")
        #i think basically this then moves the sra file onto the working path

       # with open(local_file_path, 'wb') as writer:
            #ftp.retrbinary('RETR %s' % ftp_file_path, writer.write)
    #ftp.close()
    print("Completed downloading %d read files." % len(accession_list))

if __name__ == "__main__":
    main()
