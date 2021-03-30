##make a list of pileups to be completed
import os
import numpy as np

dir = '/Users/asherpreskasteinberg/Desktop/code/recombo/APS171_NG_analysis/'
outdir = dir+'piles_tbc1/'
complete = set()
piles = open(dir + 'APS171_completepiles', 'r')
for _, pile in enumerate(piles):
    pile = str.rstrip(pile)
    complete.add(pile)

all_fastqs = set()

dwnlds = open(dir + 'sra_list', 'r')
#dwnlds = open(dir + 'APS154_blankpiles', 'r')
for _, sra in enumerate(dwnlds):
    sra = str.rstrip(sra)
    all_fastqs.add(sra)

print("all SRA: %d" % len(all_fastqs))
tbc = all_fastqs.difference(complete)
print("to be completed: %d" % len(tbc))
#tbc = list(all_fastqs)
tbc = list(tbc)
split_tbc = np.array_split(tbc, 200)

if not os.path.exists(outdir):
    os.makedirs(outdir)
for i in np.arange(0, len(split_tbc)):
    tbc_i = split_tbc[i]
    with open(outdir+'piles_TBC_'+str(i), 'w+') as pilesupTBC:
        for sra in tbc_i:
            if sra != "":
                pilesupTBC.write(sra+'\n')