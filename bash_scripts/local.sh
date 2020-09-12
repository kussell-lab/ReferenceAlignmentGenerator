##INPUTS
#The directory you want create serotype folders in. You need lots of space.
#SRC is the directory your SRA accession files and genome files are within.
SRC=/Users/asherpreskasteinberg/Desktop/code/recombo/salmonella

mkdir /Users/asherpreskasteinberg/Desktop/code/recombo/2cladestest

WRKD=/Users/asherpreskasteinberg/Desktop/code/recombo/2cladestest


#for sero in typhimurium_both typhimurium_1 typhimurium_2
for sero in typhimurium_1
do
  ReferenceAlignmentGenerate ${SRC}/SRA_files/sra_accession_$sero ${WRKD} ${SRC}/Reference/GCF_000006945.2_ASM694v2_genomic.fna ${SRC}/Reference/GCF_000006945.2_ASM694v2_genomic.gff ${WRKD}/REFGEN_$sero
done