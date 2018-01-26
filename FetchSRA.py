"""
This script download read files (SRA format) from NCBI FTP.
Created by Mingzhi Lin (mingzhi9@gmail.com).
It requires 
    (1) a file containing a list of accessions;
    (2) working directory.
"""
import sys
import os
import ftplib
from tqdm import tqdm
def main():
    """ main function """
    accession_list_file = sys.argv[1]
    working_dir = sys.argv[2]

    # read the list of accessions.
    accession_list = []
    with open(accession_list_file, 'r') as reader:
        for line in reader:
            accession_list.append(line.rstrip())

    # open FTP connection.
    ftp = ftplib.FTP("ftp.ncbi.nlm.nih.gov")
    ftp.login("anonymous", "mingzhi.lin@nyu.edu")
    ftp.cwd("sra/sra-instant/reads/ByRun/sra")

    print("Fetching SRA from NCBI ftp...")
    for accession in tqdm(accession_list):
        # read path: SRR/SRR000/SRR000001/SRR000001.sra
        ftp_file_path = "%s/%s/%s/%s.sra" % (accession[:3], accession[:6], accession, accession)
        local_file_path = os.path.join(working_dir, accession + ".sra")
        with open(local_file_path, 'wb') as writer:
            ftp.retrbinary('RETR %s' % ftp_file_path, writer.write)
    ftp.close()
    print("Completed downloading %d read files." % len(accession_list))

if __name__ == "__main__":
    main()
