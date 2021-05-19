#!/usr/bin/env python3
import os
import numpy as np
import glob
outdir = "/Users/asherpreskasteinberg/Desktop/code/recombo/APS194_EC_analysis"
genomedir = os.path.join(outdir, "local/genomes")
genomes = glob.glob(genomedir+"/*/*.fna")
outfile = os.path.join(outdir, "unaligned_genomes.fasta")
with open(outfile, "w+") as out:
    for genome in genomes:
        with open(genome, "r") as g:
            one_genome = g.readlines()
        for line in one_genome:
            out.write(line)



