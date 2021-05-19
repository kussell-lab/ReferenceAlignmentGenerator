import glob
import os

pileupList = glob.glob('*.pileup.fasta')
file = open('APS207_completepiles', 'w+')
i = 0
for pileup in pileupList:
    size = os.stat(pileup).st_size
    if size != 0:
        i = i + 1
        pilestr = pileup.split('.')
        file.write(pilestr[0]+"\n")
file.close()
print(str(i))