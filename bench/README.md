# Benchmark

To provide a performance comparison, the `Hyperscan`, `regexp`, `re2` and `pcre2` performance testing tools are provided here.

The testing tools are divided into two versions, `cpp` and `golang`. The `cpp` version providing baseline of each library and the `golang` providing performance comparison.

## Test suite

| Index | Level | Pattern |
|-------|-------|---------|
| 0 | Easy0 | ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 1 |Easy0i | (?i)ABCDEFGHIJklmnopqrstuvwxyz$ |
| 2 |Easy1 | A[AB]B[BC]C[CD]D[DE]E[EF]F[FG]G[GH]H[HI]I[IJ]J$ |
| 3 |Medium | [XYZ]ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 4 |Hard | [ -~]*ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 5 |Hard1 | ABCD\|CDEF\|EFGH\|GHIJ\|IJKL\|KLMN\|MNOP\|OPQR\|QRST\|STUV\|UVWX\|WXYZ |


## Golang benchmarks

The golang performance testing support `hyperscan` and `regexp` package.

```sh
$ cd go && go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/flier/gohs/bench/go
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkHyperscanBlockScan/Easy0/16-8         	 1000000	      1079 ns/op	  14.83 MB/s
BenchmarkHyperscanBlockScan/Easy0/32-8         	 1449081	       823.3 ns/op	  38.87 MB/s
BenchmarkHyperscanBlockScan/Easy0/1K-8         	 1000000	      1192 ns/op	 859.32 MB/s
BenchmarkHyperscanBlockScan/Easy0/32K-8        	  339740	      3402 ns/op	9633.36 MB/s
BenchmarkHyperscanBlockScan/Easy0/1M-8         	   16390	     71762 ns/op	14611.95 MB/s
BenchmarkHyperscanBlockScan/Easy0/32M-8        	     380	   3029810 ns/op	11074.77 MB/s
BenchmarkHyperscanBlockScan/Easy0i/16-8        	 1476424	       820.8 ns/op	  19.49 MB/s
BenchmarkHyperscanBlockScan/Easy0i/32-8        	 1487445	       880.8 ns/op	  36.33 MB/s
BenchmarkHyperscanBlockScan/Easy0i/1K-8        	 1000000	      1032 ns/op	 992.20 MB/s
BenchmarkHyperscanBlockScan/Easy0i/32K-8       	  293540	      4369 ns/op	7500.92 MB/s
BenchmarkHyperscanBlockScan/Easy0i/1M-8        	   11388	     90519 ns/op	11584.08 MB/s
BenchmarkHyperscanBlockScan/Easy0i/32M-8       	     361	   3256664 ns/op	10303.31 MB/s
BenchmarkHyperscanBlockScan/Easy1/16-8         	 1611226	       743.4 ns/op	  21.52 MB/s
BenchmarkHyperscanBlockScan/Easy1/32-8         	 1527103	       776.3 ns/op	  41.22 MB/s
BenchmarkHyperscanBlockScan/Easy1/1K-8         	 1235462	       973.8 ns/op	1051.50 MB/s
BenchmarkHyperscanBlockScan/Easy1/32K-8        	  189573	      6250 ns/op	5242.92 MB/s
BenchmarkHyperscanBlockScan/Easy1/1M-8         	    6976	    174281 ns/op	6016.58 MB/s
BenchmarkHyperscanBlockScan/Easy1/32M-8        	     194	   6162829 ns/op	5444.65 MB/s
BenchmarkHyperscanBlockScan/Medium/16-8        	 1574637	       741.2 ns/op	  21.59 MB/s
BenchmarkHyperscanBlockScan/Medium/32-8        	 1547606	       776.0 ns/op	  41.24 MB/s
BenchmarkHyperscanBlockScan/Medium/1K-8        	 1324886	       903.7 ns/op	1133.09 MB/s
BenchmarkHyperscanBlockScan/Medium/32K-8       	  373332	      3019 ns/op	10854.57 MB/s
BenchmarkHyperscanBlockScan/Medium/1M-8        	   17019	     69800 ns/op	15022.59 MB/s
BenchmarkHyperscanBlockScan/Medium/32M-8       	     392	   3045703 ns/op	11016.97 MB/s
BenchmarkHyperscanBlockScan/Hard/16-8          	 1609479	       740.4 ns/op	  21.61 MB/s
BenchmarkHyperscanBlockScan/Hard/32-8          	 1503058	       795.7 ns/op	  40.21 MB/s
BenchmarkHyperscanBlockScan/Hard/1K-8          	 1239666	       935.7 ns/op	1094.39 MB/s
BenchmarkHyperscanBlockScan/Hard/32K-8         	  360604	      3259 ns/op	10055.91 MB/s
BenchmarkHyperscanBlockScan/Hard/1M-8          	   16080	     70097 ns/op	14959.01 MB/s
BenchmarkHyperscanBlockScan/Hard/32M-8         	     394	   3200098 ns/op	10485.44 MB/s
BenchmarkHyperscanBlockScan/Hard1/16-8         	 1438864	       786.2 ns/op	  20.35 MB/s
BenchmarkHyperscanBlockScan/Hard1/32-8         	 1448848	       828.0 ns/op	  38.65 MB/s
BenchmarkHyperscanBlockScan/Hard1/1K-8         	 1000000	      1029 ns/op	 995.04 MB/s
BenchmarkHyperscanBlockScan/Hard1/32K-8        	  159253	      6441 ns/op	5087.15 MB/s
BenchmarkHyperscanBlockScan/Hard1/1M-8         	    6280	    191017 ns/op	5489.45 MB/s
BenchmarkHyperscanBlockScan/Hard1/32M-8        	     163	   6870650 ns/op	4883.73 MB/s
BenchmarkHyperscanStreamScan/Easy0/16-8        	  624003	      2091 ns/op	   7.65 MB/s
BenchmarkHyperscanStreamScan/Easy0/32-8        	  629965	      2129 ns/op	  15.03 MB/s
BenchmarkHyperscanStreamScan/Easy0/1K-8        	  609904	      2337 ns/op	 438.19 MB/s
BenchmarkHyperscanStreamScan/Easy0/32K-8       	  116311	     10648 ns/op	3077.48 MB/s
BenchmarkHyperscanStreamScan/Easy0/1M-8        	    3904	    315493 ns/op	3323.62 MB/s
BenchmarkHyperscanStreamScan/Easy0/32M-8       	     100	  10809844 ns/op	3104.06 MB/s
BenchmarkHyperscanStreamScan/Easy0i/16-8       	  436660	      2331 ns/op	   6.86 MB/s
BenchmarkHyperscanStreamScan/Easy0i/32-8       	  611118	      1984 ns/op	  16.13 MB/s
BenchmarkHyperscanStreamScan/Easy0i/1K-8       	  512800	      2165 ns/op	 472.99 MB/s
BenchmarkHyperscanStreamScan/Easy0i/32K-8      	  105780	     10371 ns/op	3159.51 MB/s
BenchmarkHyperscanStreamScan/Easy0i/1M-8       	    3740	    360452 ns/op	2909.06 MB/s
BenchmarkHyperscanStreamScan/Easy0i/32M-8      	     106	  12012712 ns/op	2793.24 MB/s
BenchmarkHyperscanStreamScan/Easy1/16-8        	  667262	      1978 ns/op	   8.09 MB/s
BenchmarkHyperscanStreamScan/Easy1/32-8        	  623806	      1990 ns/op	  16.08 MB/s
BenchmarkHyperscanStreamScan/Easy1/1K-8        	  537493	      2187 ns/op	 468.20 MB/s
BenchmarkHyperscanStreamScan/Easy1/32K-8       	  105459	     12185 ns/op	2689.23 MB/s
BenchmarkHyperscanStreamScan/Easy1/1M-8        	    3510	    350149 ns/op	2994.66 MB/s
BenchmarkHyperscanStreamScan/Easy1/32M-8       	      98	  11259620 ns/op	2980.07 MB/s
BenchmarkHyperscanStreamScan/Medium/16-8       	  535310	      1920 ns/op	   8.33 MB/s
BenchmarkHyperscanStreamScan/Medium/32-8       	  550485	      2025 ns/op	  15.80 MB/s
BenchmarkHyperscanStreamScan/Medium/1K-8       	  621860	      2031 ns/op	 504.11 MB/s
BenchmarkHyperscanStreamScan/Medium/32K-8      	  115676	     11011 ns/op	2975.84 MB/s
BenchmarkHyperscanStreamScan/Medium/1M-8       	    3949	    301671 ns/op	3475.89 MB/s
BenchmarkHyperscanStreamScan/Medium/32M-8      	     100	  10302615 ns/op	3256.89 MB/s
BenchmarkHyperscanStreamScan/Hard/16-8         	  642331	      1884 ns/op	   8.49 MB/s
BenchmarkHyperscanStreamScan/Hard/32-8         	  657028	      1881 ns/op	  17.01 MB/s
BenchmarkHyperscanStreamScan/Hard/1K-8         	  604287	      1952 ns/op	 524.72 MB/s
BenchmarkHyperscanStreamScan/Hard/32K-8        	  116262	      9926 ns/op	3301.33 MB/s
BenchmarkHyperscanStreamScan/Hard/1M-8         	    3474	    288090 ns/op	3639.75 MB/s
BenchmarkHyperscanStreamScan/Hard/32M-8        	     122	   9664560 ns/op	3471.90 MB/s
BenchmarkHyperscanStreamScan/Hard1/16-8        	  674084	      1809 ns/op	   8.84 MB/s
BenchmarkHyperscanStreamScan/Hard1/32-8        	  649875	      1815 ns/op	  17.63 MB/s
BenchmarkHyperscanStreamScan/Hard1/1K-8        	  574249	      1996 ns/op	 512.91 MB/s
BenchmarkHyperscanStreamScan/Hard1/32K-8       	   88580	     13241 ns/op	2474.79 MB/s
BenchmarkHyperscanStreamScan/Hard1/1M-8        	    2907	    397630 ns/op	2637.06 MB/s
BenchmarkHyperscanStreamScan/Hard1/32M-8       	      74	  15599489 ns/op	2151.00 MB/s
BenchmarkRegexpMatch/Easy0/16-8                	277881128	         4.329 ns/op	3696.18 MB/s
BenchmarkRegexpMatch/Easy0/32-8                	25172512	        47.98 ns/op	 666.98 MB/s
BenchmarkRegexpMatch/Easy0/1K-8                	 4673600	       257.7 ns/op	3973.92 MB/s
BenchmarkRegexpMatch/Easy0/32K-8               	  267102	      4418 ns/op	7417.06 MB/s
BenchmarkRegexpMatch/Easy0/1M-8                	    4774	    248210 ns/op	4224.55 MB/s
BenchmarkRegexpMatch/Easy0/32M-8               	     138	   9835445 ns/op	3411.58 MB/s
BenchmarkRegexpMatch/Easy0i/16-8               	272014188	         4.546 ns/op	3519.30 MB/s
BenchmarkRegexpMatch/Easy0i/32-8               	 1253239	       839.2 ns/op	  38.13 MB/s
BenchmarkRegexpMatch/Easy0i/1K-8               	   48343	     24527 ns/op	  41.75 MB/s
BenchmarkRegexpMatch/Easy0i/32K-8              	    1208	   1008136 ns/op	  32.50 MB/s
BenchmarkRegexpMatch/Easy0i/1M-8               	      36	  37949778 ns/op	  27.63 MB/s
BenchmarkRegexpMatch/Easy0i/32M-8              	       1	1035730738 ns/op	  32.40 MB/s
BenchmarkRegexpMatch/Easy1/16-8                	264849393	         4.477 ns/op	3574.18 MB/s
BenchmarkRegexpMatch/Easy1/32-8                	18934544	        53.00 ns/op	 603.79 MB/s
BenchmarkRegexpMatch/Easy1/1K-8                	 1546407	       826.3 ns/op	1239.19 MB/s
BenchmarkRegexpMatch/Easy1/32K-8               	   33408	     35859 ns/op	 913.80 MB/s
BenchmarkRegexpMatch/Easy1/1M-8                	    1014	   1178962 ns/op	 889.41 MB/s
BenchmarkRegexpMatch/Easy1/32M-8               	      28	  46157110 ns/op	 726.96 MB/s
BenchmarkRegexpMatch/Medium/16-8               	246332023	         4.952 ns/op	3230.88 MB/s
BenchmarkRegexpMatch/Medium/32-8               	 1437381	       830.0 ns/op	  38.56 MB/s
BenchmarkRegexpMatch/Medium/1K-8               	   45045	     27213 ns/op	  37.63 MB/s
BenchmarkRegexpMatch/Medium/32K-8              	     987	   1167692 ns/op	  28.06 MB/s
BenchmarkRegexpMatch/Medium/1M-8               	      32	  35372326 ns/op	  29.64 MB/s
BenchmarkRegexpMatch/Medium/32M-8              	       1	1159234100 ns/op	  28.95 MB/s
BenchmarkRegexpMatch/Hard/16-8                 	269736577	         4.506 ns/op	3551.06 MB/s
BenchmarkRegexpMatch/Hard/32-8                 	  889212	      1315 ns/op	  24.33 MB/s
BenchmarkRegexpMatch/Hard/1K-8                 	   32110	     40720 ns/op	  25.15 MB/s
BenchmarkRegexpMatch/Hard/32K-8                	     796	   1615068 ns/op	  20.29 MB/s
BenchmarkRegexpMatch/Hard/1M-8                 	      24	  50312315 ns/op	  20.84 MB/s
BenchmarkRegexpMatch/Hard/32M-8                	       1	1597974950 ns/op	  21.00 MB/s
BenchmarkRegexpMatch/Hard1/16-8                	  300464	      3573 ns/op	   4.48 MB/s
BenchmarkRegexpMatch/Hard1/32-8                	  180313	      6662 ns/op	   4.80 MB/s
BenchmarkRegexpMatch/Hard1/1K-8                	    5784	    206433 ns/op	   4.96 MB/s
BenchmarkRegexpMatch/Hard1/32K-8               	     153	   9048851 ns/op	   3.62 MB/s
BenchmarkRegexpMatch/Hard1/1M-8                	       5	 240529097 ns/op	   4.36 MB/s
BenchmarkRegexpMatch/Hard1/32M-8               	       1	8055742592 ns/op	   4.17 MB/s
PASS
ok  	github.com/flier/gohs/bench/go	171.939s
```

## C++ benchmarks

```sh
$ cd cpp && mkdir build && cd build
$ cmake .. && make
$ ./scan_test
2021-11-24T11:05:49+08:00
Running ./scan_test
Run on (8 X 2800 MHz CPU s)
CPU Caches:
  L1 Data 32 KiB (x4)
  L1 Instruction 32 KiB (x4)
  L2 Unified 256 KiB (x4)
  L3 Unified 6144 KiB (x1)
Load Average: 4.48, 4.21, 3.82
-------------------------------------------------------------------------------------------------
Benchmark                                       Time             CPU   Iterations UserCounters...
-------------------------------------------------------------------------------------------------
BM_HS_BlockScan/regex:0/size:16              13.8 ns         13.4 ns     53895489 bytes_per_second=1.11038G/s
BM_HS_BlockScan/regex:1/size:16              14.1 ns         13.9 ns     48056789 bytes_per_second=1099.35M/s
BM_HS_BlockScan/regex:2/size:16              14.2 ns         14.1 ns     52080621 bytes_per_second=1085.4M/s
BM_HS_BlockScan/regex:3/size:16              13.9 ns         13.8 ns     47940280 bytes_per_second=1106.69M/s
BM_HS_BlockScan/regex:4/size:16              14.2 ns         14.1 ns     52540719 bytes_per_second=1081.79M/s
BM_HS_BlockScan/regex:5/size:16              50.7 ns         48.7 ns     13575626 bytes_per_second=313.536M/s
BM_HS_BlockScan/regex:0/size:32              40.3 ns         40.0 ns     16780173 bytes_per_second=762.113M/s
BM_HS_BlockScan/regex:1/size:32              39.2 ns         38.9 ns     16946444 bytes_per_second=784.018M/s
BM_HS_BlockScan/regex:2/size:32              43.2 ns         43.0 ns     17058525 bytes_per_second=710.424M/s
BM_HS_BlockScan/regex:3/size:32              40.6 ns         40.2 ns     16723846 bytes_per_second=759.925M/s
BM_HS_BlockScan/regex:4/size:32              37.4 ns         37.1 ns     18784080 bytes_per_second=821.478M/s
BM_HS_BlockScan/regex:5/size:32              44.6 ns         44.3 ns     15687997 bytes_per_second=689.273M/s
BM_HS_BlockScan/regex:0/size:1024             134 ns          133 ns      5233567 bytes_per_second=7.18029G/s
BM_HS_BlockScan/regex:1/size:1024             150 ns          149 ns      4715011 bytes_per_second=6.40565G/s
BM_HS_BlockScan/regex:2/size:1024             211 ns          210 ns      3347665 bytes_per_second=4.5518G/s
BM_HS_BlockScan/regex:3/size:1024             143 ns          142 ns      4799616 bytes_per_second=6.73836G/s
BM_HS_BlockScan/regex:4/size:1024             133 ns          132 ns      5278955 bytes_per_second=7.21186G/s
BM_HS_BlockScan/regex:5/size:1024             205 ns          204 ns      3425495 bytes_per_second=4.68064G/s
BM_HS_BlockScan/regex:0/size:32768           2142 ns         2126 ns       327950 bytes_per_second=14.3518G/s
BM_HS_BlockScan/regex:1/size:32768           2703 ns         2414 ns       290296 bytes_per_second=12.6425G/s
BM_HS_BlockScan/regex:2/size:32768           5265 ns         5225 ns       134401 bytes_per_second=5.84089G/s
BM_HS_BlockScan/regex:3/size:32768           2149 ns         2132 ns       328655 bytes_per_second=14.3115G/s
BM_HS_BlockScan/regex:4/size:32768           2138 ns         2123 ns       318570 bytes_per_second=14.3717G/s
BM_HS_BlockScan/regex:5/size:32768           5431 ns         5380 ns       131401 bytes_per_second=5.67285G/s
BM_HS_BlockScan/regex:0/size:1048576        67697 ns        67094 ns        10493 bytes_per_second=14.5551G/s
BM_HS_BlockScan/regex:1/size:1048576        77143 ns        75189 ns         9560 bytes_per_second=12.9881G/s
BM_HS_BlockScan/regex:2/size:1048576       168406 ns       167027 ns         4213 bytes_per_second=5.84673G/s
BM_HS_BlockScan/regex:3/size:1048576        76277 ns        71355 ns        10642 bytes_per_second=13.6859G/s
BM_HS_BlockScan/regex:4/size:1048576       100923 ns        96562 ns         9801 bytes_per_second=10.1133G/s
BM_HS_BlockScan/regex:5/size:1048576       175106 ns       172563 ns         3458 bytes_per_second=5.65915G/s
BM_HS_BlockScan/regex:0/size:33554432     3231445 ns      3168914 ns          175 bytes_per_second=9.86142G/s
BM_HS_BlockScan/regex:1/size:33554432     3327615 ns      3267303 ns          211 bytes_per_second=9.56446G/s
BM_HS_BlockScan/regex:2/size:33554432     6138092 ns      6054447 ns          123 bytes_per_second=5.1615G/s
BM_HS_BlockScan/regex:3/size:33554432     2873969 ns      2856500 ns          226 bytes_per_second=10.94G/s
BM_HS_BlockScan/regex:4/size:33554432     3271825 ns      3120980 ns          245 bytes_per_second=10.0129G/s
BM_HS_BlockScan/regex:5/size:33554432     6146788 ns      6070000 ns          116 bytes_per_second=5.14827G/s
BM_HS_StreamScan/regex:0/size:16              178 ns          175 ns      4044746 bytes_per_second=86.9931M/s
BM_HS_StreamScan/regex:1/size:16              177 ns          175 ns      3968501 bytes_per_second=87.31M/s
BM_HS_StreamScan/regex:2/size:16              190 ns          186 ns      3768628 bytes_per_second=81.8494M/s
BM_HS_StreamScan/regex:3/size:16              173 ns          172 ns      4022896 bytes_per_second=88.8386M/s
BM_HS_StreamScan/regex:4/size:16              174 ns          173 ns      4060019 bytes_per_second=88.3356M/s
BM_HS_StreamScan/regex:5/size:16              153 ns          151 ns      4594141 bytes_per_second=100.842M/s
BM_HS_StreamScan/regex:0/size:32              189 ns          187 ns      3756796 bytes_per_second=162.874M/s
BM_HS_StreamScan/regex:1/size:32              190 ns          189 ns      3743996 bytes_per_second=161.879M/s
BM_HS_StreamScan/regex:2/size:32              218 ns          216 ns      3225925 bytes_per_second=141.432M/s
BM_HS_StreamScan/regex:3/size:32              201 ns          199 ns      3532570 bytes_per_second=153.088M/s
BM_HS_StreamScan/regex:4/size:32              188 ns          187 ns      3666284 bytes_per_second=163.531M/s
BM_HS_StreamScan/regex:5/size:32              156 ns          155 ns      4454315 bytes_per_second=197.107M/s
BM_HS_StreamScan/regex:0/size:1024            254 ns          252 ns      2761232 bytes_per_second=3.79139G/s
BM_HS_StreamScan/regex:1/size:1024            275 ns          273 ns      2602221 bytes_per_second=3.49582G/s
BM_HS_StreamScan/regex:2/size:1024            303 ns          301 ns      2290149 bytes_per_second=3.17297G/s
BM_HS_StreamScan/regex:3/size:1024            265 ns          263 ns      2633708 bytes_per_second=3.63114G/s
BM_HS_StreamScan/regex:4/size:1024            254 ns          252 ns      2748245 bytes_per_second=3.79172G/s
BM_HS_StreamScan/regex:5/size:1024            319 ns          317 ns      2226994 bytes_per_second=3.00715G/s
BM_HS_StreamScan/regex:0/size:32768          3108 ns         2904 ns       246494 bytes_per_second=10.5104G/s
BM_HS_StreamScan/regex:1/size:32768          3129 ns         3083 ns       229862 bytes_per_second=9.89757G/s
BM_HS_StreamScan/regex:2/size:32768          3686 ns         3654 ns       190657 bytes_per_second=8.35126G/s
BM_HS_StreamScan/regex:3/size:32768          2842 ns         2819 ns       251109 bytes_per_second=10.8269G/s
BM_HS_StreamScan/regex:4/size:32768          2876 ns         2851 ns       251446 bytes_per_second=10.7042G/s
BM_HS_StreamScan/regex:5/size:32768          6082 ns         6034 ns       114226 bytes_per_second=5.05767G/s
BM_HS_StreamScan/regex:0/size:1048576       85409 ns        84686 ns         8232 bytes_per_second=11.5316G/s
BM_HS_StreamScan/regex:1/size:1048576       97327 ns        95888 ns         7639 bytes_per_second=10.1845G/s
BM_HS_StreamScan/regex:2/size:1048576      109066 ns       107604 ns         6749 bytes_per_second=9.07551G/s
BM_HS_StreamScan/regex:3/size:1048576       87596 ns        86786 ns         8026 bytes_per_second=11.2525G/s
BM_HS_StreamScan/regex:4/size:1048576       87201 ns        86453 ns         8092 bytes_per_second=11.2959G/s
BM_HS_StreamScan/regex:5/size:1048576      196266 ns       194618 ns         3609 bytes_per_second=5.01783G/s
BM_HS_StreamScan/regex:0/size:33554432    3220054 ns      3192269 ns          212 bytes_per_second=9.78928G/s
BM_HS_StreamScan/regex:1/size:33554432    3408965 ns      3380164 ns          207 bytes_per_second=9.24511G/s
BM_HS_StreamScan/regex:2/size:33554432    4276215 ns      4241557 ns          167 bytes_per_second=7.36758G/s
BM_HS_StreamScan/regex:3/size:33554432    3384055 ns      3319072 ns          222 bytes_per_second=9.41528G/s
BM_HS_StreamScan/regex:4/size:33554432    4062913 ns      3908705 ns          166 bytes_per_second=7.99498G/s
BM_HS_StreamScan/regex:5/size:33554432   13788005 ns     10931746 ns           63 bytes_per_second=2.85865G/s
BM_RE2_Match/regex:0/size:16                 85.0 ns         83.1 ns      8783377 bytes_per_second=183.534M/s
BM_RE2_Match/regex:1/size:16                 90.9 ns         89.7 ns      7899162 bytes_per_second=170.156M/s
BM_RE2_Match/regex:2/size:16                 81.5 ns         80.6 ns      8205822 bytes_per_second=189.248M/s
BM_RE2_Match/regex:3/size:16                  104 ns          102 ns      7243602 bytes_per_second=149.233M/s
BM_RE2_Match/regex:4/size:16                  109 ns          106 ns      6903184 bytes_per_second=144.456M/s
BM_RE2_Match/regex:5/size:16                  152 ns          120 ns      6840747 bytes_per_second=127.362M/s
BM_RE2_Match/regex:0/size:32                  125 ns          113 ns      8113026 bytes_per_second=270.581M/s
BM_RE2_Match/regex:1/size:32                  126 ns          115 ns      7086311 bytes_per_second=264.988M/s
BM_RE2_Match/regex:2/size:32                  101 ns         95.6 ns      6919389 bytes_per_second=319.304M/s
BM_RE2_Match/regex:3/size:32                  132 ns          128 ns      5503707 bytes_per_second=238.463M/s
BM_RE2_Match/regex:4/size:32                  160 ns          154 ns      4186928 bytes_per_second=198.531M/s
BM_RE2_Match/regex:5/size:32                  213 ns          179 ns      3128576 bytes_per_second=170.191M/s
BM_RE2_Match/regex:0/size:1024                265 ns          258 ns      2843887 bytes_per_second=3.70037G/s
BM_RE2_Match/regex:1/size:1024               1006 ns          929 ns       774371 bytes_per_second=1051.65M/s
BM_RE2_Match/regex:2/size:1024                571 ns          341 ns      2151119 bytes_per_second=2.79619G/s
BM_RE2_Match/regex:3/size:1024               9241 ns         3203 ns       215940 bytes_per_second=304.9M/s
BM_RE2_Match/regex:4/size:1024              12832 ns         3771 ns       196416 bytes_per_second=258.961M/s
BM_RE2_Match/regex:5/size:1024              11484 ns         3472 ns       192504 bytes_per_second=281.274M/s
BM_RE2_Match/regex:0/size:32768              5930 ns         4008 ns       162166 bytes_per_second=7.61508G/s
BM_RE2_Match/regex:1/size:32768             26884 ns        22840 ns        29320 bytes_per_second=1.33616G/s
BM_RE2_Match/regex:2/size:32768              4312 ns         4174 ns       162000 bytes_per_second=7.31131G/s
BM_RE2_Match/regex:3/size:32768             65004 ns        60592 ns        12719 bytes_per_second=515.748M/s
BM_RE2_Match/regex:4/size:32768             78039 ns        69632 ns         9670 bytes_per_second=448.791M/s
BM_RE2_Match/regex:5/size:32768            139792 ns        81426 ns         7737 bytes_per_second=383.782M/s
BM_RE2_Match/regex:0/size:1048576          345010 ns       288741 ns         2282 bytes_per_second=3.38214G/s
BM_RE2_Match/regex:1/size:1048576         1163356 ns       791127 ns          795 bytes_per_second=1.23439G/s
BM_RE2_Match/regex:2/size:1048576          573456 ns       307397 ns         1967 bytes_per_second=3.17688G/s
BM_RE2_Match/regex:3/size:1048576         2007048 ns      1881170 ns          365 bytes_per_second=531.584M/s
BM_RE2_Match/regex:4/size:1048576         4623146 ns      2758967 ns          397 bytes_per_second=362.454M/s
BM_RE2_Match/regex:5/size:1048576         6291603 ns      2995235 ns          238 bytes_per_second=333.864M/s
BM_RE2_Match/regex:0/size:33554432       12292832 ns      9057543 ns           81 bytes_per_second=3.45016G/s
BM_RE2_Match/regex:1/size:33554432       22297519 ns     21933029 ns           35 bytes_per_second=1.42479G/s
BM_RE2_Match/regex:2/size:33554432        8928809 ns      8713319 ns           69 bytes_per_second=3.58646G/s
BM_RE2_Match/regex:3/size:33554432       52855948 ns     52358385 ns           13 bytes_per_second=611.172M/s
BM_RE2_Match/regex:4/size:33554432       65909571 ns     55910167 ns           12 bytes_per_second=572.347M/s
BM_RE2_Match/regex:5/size:33554432       52798528 ns     52308846 ns           13 bytes_per_second=611.751M/s
BM_PCRE2_Match/regex:0/size:16               50.7 ns         50.2 ns     10000000 bytes_per_second=303.798M/s
BM_PCRE2_Match/regex:1/size:16               56.2 ns         55.7 ns     12394207 bytes_per_second=274.033M/s
BM_PCRE2_Match/regex:2/size:16               50.5 ns         50.1 ns     10000000 bytes_per_second=304.868M/s
BM_PCRE2_Match/regex:3/size:16               61.7 ns         61.0 ns     12079170 bytes_per_second=250.015M/s
BM_PCRE2_Match/regex:4/size:16               49.0 ns         48.2 ns     15172106 bytes_per_second=316.656M/s
BM_PCRE2_Match/regex:5/size:16                327 ns          321 ns      2381406 bytes_per_second=47.567M/s
BM_PCRE2_Match/regex:0/size:32               54.8 ns         53.7 ns     10950846 bytes_per_second=568.05M/s
BM_PCRE2_Match/regex:1/size:32               75.1 ns         73.5 ns     10052561 bytes_per_second=415.27M/s
BM_PCRE2_Match/regex:2/size:32               59.7 ns         58.8 ns     12316788 bytes_per_second=518.661M/s
BM_PCRE2_Match/regex:3/size:32               72.7 ns         72.1 ns      9453838 bytes_per_second=423.234M/s
BM_PCRE2_Match/regex:4/size:32               52.2 ns         51.8 ns     13532585 bytes_per_second=589.708M/s
BM_PCRE2_Match/regex:5/size:32                729 ns          722 ns       963404 bytes_per_second=42.264M/s
BM_PCRE2_Match/regex:0/size:1024              653 ns          642 ns      1065725 bytes_per_second=1.48485G/s
BM_PCRE2_Match/regex:1/size:1024             1514 ns         1500 ns       468921 bytes_per_second=651.015M/s
BM_PCRE2_Match/regex:2/size:1024              821 ns          813 ns       848176 bytes_per_second=1.1725G/s
BM_PCRE2_Match/regex:3/size:1024             1983 ns         1964 ns       355673 bytes_per_second=497.25M/s
BM_PCRE2_Match/regex:4/size:1024           356774 ns       353173 ns         1976 bytes_per_second=2.76511M/s
BM_PCRE2_Match/regex:5/size:1024            27190 ns        26907 ns        25903 bytes_per_second=36.2942M/s
BM_PCRE2_Match/regex:0/size:32768           19096 ns        18911 ns        36775 bytes_per_second=1.61378G/s
BM_PCRE2_Match/regex:1/size:32768           43862 ns        43452 ns        16055 bytes_per_second=719.187M/s
BM_PCRE2_Match/regex:2/size:32768           20136 ns        19882 ns        35196 bytes_per_second=1.53497G/s
BM_PCRE2_Match/regex:3/size:32768           82007 ns        81227 ns         8663 bytes_per_second=384.723M/s
BM_PCRE2_Match/regex:4/size:32768        16029648 ns     15878250 ns           44 bytes_per_second=1.9681M/s
BM_PCRE2_Match/regex:5/size:32768          908637 ns       899426 ns          772 bytes_per_second=34.7444M/s
BM_PCRE2_Match/regex:0/size:1048576        742967 ns       736052 ns          955 bytes_per_second=1.32676G/s
BM_PCRE2_Match/regex:1/size:1048576       1539122 ns      1523551 ns          452 bytes_per_second=656.361M/s
BM_PCRE2_Match/regex:2/size:1048576        758259 ns       751683 ns          932 bytes_per_second=1.29917G/s
BM_PCRE2_Match/regex:3/size:1048576       2765804 ns      2738792 ns          255 bytes_per_second=365.124M/s
BM_PCRE2_Match/regex:4/size:1048576     539847874 ns    534908000 ns            1 bytes_per_second=1.86948M/s
BM_PCRE2_Match/regex:5/size:1048576      29264190 ns     28985708 ns           24 bytes_per_second=34.4998M/s
BM_PCRE2_Match/regex:0/size:33554432     21107502 ns     20900559 ns           34 bytes_per_second=1.49518G/s
BM_PCRE2_Match/regex:1/size:33554432     53189948 ns     50076313 ns           16 bytes_per_second=639.025M/s
BM_PCRE2_Match/regex:2/size:33554432     21585020 ns     21385938 ns           32 bytes_per_second=1.46124G/s
BM_PCRE2_Match/regex:3/size:33554432     85048097 ns     83850571 ns            7 bytes_per_second=381.631M/s
BM_PCRE2_Match/regex:4/size:33554432   18672117602 ns   18025074000 ns            1 bytes_per_second=1.7753M/s
BM_PCRE2_Match/regex:5/size:33554432    944520308 ns    933643000 ns            1 bytes_per_second=34.2743M/s
