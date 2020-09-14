import csv

import pandas as pd
"""
This script grabs fastq.gz files from ENA
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

    #accession_list_file = '/Users/asherpreskasteinberg/Desktop/code/recombo/APS138_cjejuni/filereport_read_run_PRJEB31119_tsv.txt'
    #working_dir = '/Users/asherpreskasteinberg/Desktop/fetchsra_test/'


    dat = pd.read_csv(accession_list_file, sep = '\t')

    fastqlist = dat['fastq_ftp']
    accessions = dat['run_accession']
    accession_list = []
    i = 0
    for fastq in tqdm(fastqlist):
        if ';' not in fastq:
            i = i + 1
            continue
        else:
            fastqs = fastq.split(";")
            sra = fastq.split("/")
            dir = os.getcwd()
            if os.path.isfile(dir + '/'+sra[5] + '_1.fastq'):
                i = i + 1
                continue
            os.system("ssh dtn")
            os.system("wget -nv "+ fastqs[0])
            os.system("wget -nv " + fastqs[1])
            os.system("logout")
            ##unzip the files
            os.system("gunzip "+ sra[5]+'_1.fastq')
            os.system("gunzip " + sra[5]+'_2.fastq')
            accession_list.append(sra[5])
            i = i + 1

    with open('ENA_names', 'w') as f:
        i = 0
        for sra in accession_list:
            if i < len(accession_list)-1:
                f.write("%s\n" %sra)
            else:
                f.write("%s" %sra)
            i = i + 1

if __name__ == "__main__":
    main()
