lendisto: length distribution for records in common bioinformatics files
========================================================================

`lendisto` is an API and command line app to calculate the length distribution
of the records in common bioinformatics files. Currently supported are FASTA,
FASTQ and BED files. The API tries to be safe for concurrent use.

An example of how to use the provided command line apps.

Given the following `file1.bed` and `file2.bed`

```bash
$> cat file1.bed
chr1	0	2
chr2	0	7
chr2	10	15

$> cat file2.bed
chrX	0	2
chrX	10	15
chrX	10	15
chrY	0	7
chrY	0	7
chrY	0	7
chrY	0	7
```

the command

```bash
bed-len-distro file1.bed file2.bed
```

prints

```
len	count	density
0	0	0
1	0	0
2	2	0.2
3	0	0
4	0	0
5	3	0.3
6	0	0
7	5	0.5
```

