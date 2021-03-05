##make comma separated lists for the NCBI run selector ....
### this is so i can get location tags
import os
import numpy as np

dir = '/Users/asherpreskasteinberg/Desktop/code/recombo/APS168_sars-cov-2/'
outdir = dir+'runselector/'
complete = set()
piles = open(dir + 'APS168_completepiles', 'r')
for _, pile in enumerate(piles):
    pile = str.rstrip(pile)
    complete.add(pile)

tbc = complete
print("total SRAs: %d" % len(tbc))
#tbc = list(all_fastqs)
tbc = list(tbc)
split_tbc = np.array_split(tbc, 10)

if not os.path.exists(outdir):
    os.makedirs(outdir)
for i in np.arange(0, len(split_tbc)):
    tbc_i = split_tbc[i]
    with open(outdir+'runselector_'+str(i), 'w+') as pilesupTBC:
        for sra in tbc_i:
            if sra != "":
                pilesupTBC.write(sra+',\n')