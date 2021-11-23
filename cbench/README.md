# Benchmark for C library

To provide a performance comparison, the `Hyperscan`, `re2` and `pcre2` performance testing tools are provided here.

## Test suite

The basic test suite is consistent with the golang implementation.

| Index | Level | Pattern |
|-------|-------|---------|
| 0 | Easy0 | ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 1 |Easy0i | (?i)ABCDEFGHIJklmnopqrstuvwxyz$ |
| 2 |Easy1 | A[AB]B[BC]C[CD]D[DE]E[EF]F[FG]G[GH]H[HI]I[IJ]J$ |
| 3 |Medium | [XYZ]ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 4 |Hard | [ -~]*ABCDEFGHIJKLMNOPQRSTUVWXYZ$ |
| 5 |Hard1 | ABCD\|CDEF\|EFGH\|GHIJ\|IJKL\|KLMN\|MNOP\|OPQR\|QRST\|STUV\|UVWX\|WXYZ |


## Build & Run

```
$ mkdir build && cd build
$ cmake .. && make
$ ./scan_test
2021-11-23T19:38:31+08:00
Running ./scan_test
Run on (8 X 2800 MHz CPU s)
CPU Caches:
  L1 Data 32 KiB (x4)
  L1 Instruction 32 KiB (x4)
  L2 Unified 256 KiB (x4)
  L3 Unified 6144 KiB (x1)
Load Average: 5.54, 5.54, 5.60
----------------------------------------------------------------------------------------------
Benchmark                                    Time             CPU   Iterations UserCounters...
----------------------------------------------------------------------------------------------
BM_BlockScan/regex:0/size:16              13.1 ns         12.9 ns     54734965 bytes_per_second=1.15367G/s
BM_BlockScan/regex:1/size:16              13.3 ns         13.1 ns     53372778 bytes_per_second=1.13868G/s
BM_BlockScan/regex:2/size:16              13.2 ns         12.9 ns     48543353 bytes_per_second=1.15306G/s
BM_BlockScan/regex:3/size:16              13.1 ns         12.9 ns     54102098 bytes_per_second=1.15388G/s
BM_BlockScan/regex:4/size:16              15.8 ns         14.7 ns     54422615 bytes_per_second=1037.07M/s
BM_BlockScan/regex:5/size:16              43.8 ns         43.3 ns     16206968 bytes_per_second=352.318M/s
BM_BlockScan/regex:0/size:32              58.9 ns         51.2 ns     18245131 bytes_per_second=596.407M/s
BM_BlockScan/regex:1/size:32              40.5 ns         40.0 ns     14648427 bytes_per_second=763.552M/s
BM_BlockScan/regex:2/size:32              49.9 ns         47.8 ns     17368908 bytes_per_second=638.658M/s
BM_BlockScan/regex:3/size:32              58.3 ns         52.3 ns     10597390 bytes_per_second=583.062M/s
BM_BlockScan/regex:4/size:32              52.0 ns         50.7 ns     10000000 bytes_per_second=602.185M/s
BM_BlockScan/regex:5/size:32              44.4 ns         44.1 ns     15308075 bytes_per_second=692.547M/s
BM_BlockScan/regex:0/size:1024             134 ns          132 ns      5367810 bytes_per_second=7.20062G/s
BM_BlockScan/regex:1/size:1024             153 ns          152 ns      4598033 bytes_per_second=6.27993G/s
BM_BlockScan/regex:2/size:1024             220 ns          217 ns      3212955 bytes_per_second=4.38958G/s
BM_BlockScan/regex:3/size:1024             149 ns          148 ns      4248424 bytes_per_second=6.45684G/s
BM_BlockScan/regex:4/size:1024             145 ns          143 ns      5057621 bytes_per_second=6.67505G/s
BM_BlockScan/regex:5/size:1024             220 ns          215 ns      3352201 bytes_per_second=4.43066G/s
BM_BlockScan/regex:0/size:32768           2300 ns         2252 ns       316330 bytes_per_second=13.5497G/s
BM_BlockScan/regex:1/size:32768           2519 ns         2479 ns       288425 bytes_per_second=12.3098G/s
BM_BlockScan/regex:2/size:32768           5445 ns         5392 ns       130692 bytes_per_second=5.66013G/s
BM_BlockScan/regex:3/size:32768           2918 ns         2523 ns       295193 bytes_per_second=12.097G/s
BM_BlockScan/regex:4/size:32768           2936 ns         2809 ns       218587 bytes_per_second=10.8633G/s
BM_BlockScan/regex:5/size:32768           5462 ns         5406 ns       130978 bytes_per_second=5.64524G/s
BM_BlockScan/regex:0/size:1048576        70009 ns        69330 ns        10585 bytes_per_second=14.0858G/s
BM_BlockScan/regex:1/size:1048576        79814 ns        78804 ns         9106 bytes_per_second=12.3923G/s
BM_BlockScan/regex:2/size:1048576       188209 ns       183665 ns         4011 bytes_per_second=5.31708G/s
BM_BlockScan/regex:3/size:1048576        68653 ns        68100 ns        10031 bytes_per_second=14.3401G/s
BM_BlockScan/regex:4/size:1048576        68838 ns        68030 ns        10294 bytes_per_second=14.3548G/s
BM_BlockScan/regex:5/size:1048576       172155 ns       170730 ns         4130 bytes_per_second=5.71991G/s
BM_BlockScan/regex:0/size:33554432     3092394 ns      3050563 ns          238 bytes_per_second=10.244G/s
BM_BlockScan/regex:1/size:33554432     3642805 ns      3540807 ns          212 bytes_per_second=8.82567G/s
BM_BlockScan/regex:2/size:33554432     6650852 ns      6563688 ns          109 bytes_per_second=4.76104G/s
BM_BlockScan/regex:3/size:33554432     4832506 ns      4276776 ns          183 bytes_per_second=7.30691G/s
BM_BlockScan/regex:4/size:33554432     3555352 ns      3464279 ns          208 bytes_per_second=9.02064G/s
BM_BlockScan/regex:5/size:33554432     6482280 ns      6296726 ns          113 bytes_per_second=4.9629G/s
BM_StreamScan/regex:0/size:16              189 ns          185 ns      3712012 bytes_per_second=82.6134M/s
BM_StreamScan/regex:1/size:16              195 ns          190 ns      3682156 bytes_per_second=80.2166M/s
BM_StreamScan/regex:2/size:16              221 ns          214 ns      2883839 bytes_per_second=71.1413M/s
BM_StreamScan/regex:3/size:16              212 ns          206 ns      3649064 bytes_per_second=74.1362M/s
BM_StreamScan/regex:4/size:16              254 ns          225 ns      3856834 bytes_per_second=67.8782M/s
BM_StreamScan/regex:5/size:16              230 ns          200 ns      3652168 bytes_per_second=76.2538M/s
BM_StreamScan/regex:0/size:32              195 ns          194 ns      3433443 bytes_per_second=157.424M/s
BM_StreamScan/regex:1/size:32              201 ns          200 ns      3637649 bytes_per_second=152.886M/s
BM_StreamScan/regex:2/size:32              229 ns          228 ns      3054661 bytes_per_second=133.858M/s
BM_StreamScan/regex:3/size:32              220 ns          218 ns      3405961 bytes_per_second=139.907M/s
BM_StreamScan/regex:4/size:32              202 ns          199 ns      2722464 bytes_per_second=153.116M/s
BM_StreamScan/regex:5/size:32              169 ns          166 ns      4138534 bytes_per_second=184.039M/s
BM_StreamScan/regex:0/size:1024            360 ns          315 ns      2288352 bytes_per_second=3.02681G/s
BM_StreamScan/regex:1/size:1024            506 ns          400 ns      2497271 bytes_per_second=2.38429G/s
BM_StreamScan/regex:2/size:1024            327 ns          325 ns      2038196 bytes_per_second=2.93577G/s
BM_StreamScan/regex:3/size:1024            283 ns          282 ns      2320486 bytes_per_second=3.38686G/s
BM_StreamScan/regex:4/size:1024            458 ns          383 ns      2607746 bytes_per_second=2.48936G/s
BM_StreamScan/regex:5/size:1024            352 ns          340 ns      1981291 bytes_per_second=2.802G/s
BM_StreamScan/regex:0/size:32768          4143 ns         3764 ns       229939 bytes_per_second=8.10675G/s
BM_StreamScan/regex:1/size:32768          4157 ns         3837 ns       202278 bytes_per_second=7.95361G/s
BM_StreamScan/regex:2/size:32768          4723 ns         4493 ns       139387 bytes_per_second=6.7918G/s
BM_StreamScan/regex:3/size:32768          4303 ns         3883 ns       207454 bytes_per_second=7.85936G/s
BM_StreamScan/regex:4/size:32768          3154 ns         3113 ns       229737 bytes_per_second=9.80196G/s
BM_StreamScan/regex:5/size:32768          6835 ns         6713 ns        96918 bytes_per_second=4.54637G/s
BM_StreamScan/regex:0/size:1048576      138374 ns       123307 ns         6753 bytes_per_second=7.91976G/s
BM_StreamScan/regex:1/size:1048576      103879 ns       102705 ns         5209 bytes_per_second=9.50843G/s
BM_StreamScan/regex:2/size:1048576      159993 ns       148346 ns         5602 bytes_per_second=6.58302G/s
BM_StreamScan/regex:3/size:1048576      101154 ns       100183 ns         6007 bytes_per_second=9.74777G/s
BM_StreamScan/regex:4/size:1048576       91726 ns        91219 ns         7758 bytes_per_second=10.7057G/s
BM_StreamScan/regex:5/size:1048576      203027 ns       201990 ns         3291 bytes_per_second=4.8347G/s
BM_StreamScan/regex:0/size:33554432    5027585 ns      4639724 ns          196 bytes_per_second=6.73531G/s
BM_StreamScan/regex:1/size:33554432    3503289 ns      3488969 ns          193 bytes_per_second=8.9568G/s
BM_StreamScan/regex:2/size:33554432    6545575 ns      6015150 ns          100 bytes_per_second=5.19522G/s
BM_StreamScan/regex:3/size:33554432    3424809 ns      3374344 ns          160 bytes_per_second=9.26106G/s
BM_StreamScan/regex:4/size:33554432    5505244 ns      5026410 ns          100 bytes_per_second=6.21716G/s
BM_StreamScan/regex:5/size:33554432    9952991 ns      9852947 ns           76 bytes_per_second=3.17164G/s
BM_RE2Match/regex:0/size:16             80.5 ns         79.7 ns      8790105 bytes_per_second=191.385M/s
BM_RE2Match/regex:1/size:16             86.7 ns         86.1 ns      7886702 bytes_per_second=177.25M/s
BM_RE2Match/regex:2/size:16             93.2 ns         92.6 ns      8668731 bytes_per_second=164.702M/s
BM_RE2Match/regex:3/size:16              120 ns          120 ns      6156498 bytes_per_second=127.461M/s
BM_RE2Match/regex:4/size:16              113 ns          111 ns      6346098 bytes_per_second=137.132M/s
BM_RE2Match/regex:5/size:16              125 ns          122 ns      5273189 bytes_per_second=125.377M/s
BM_RE2Match/regex:0/size:32             94.7 ns         93.9 ns      7302697 bytes_per_second=325.005M/s
BM_RE2Match/regex:1/size:32              110 ns          109 ns      5576401 bytes_per_second=279.493M/s
BM_RE2Match/regex:2/size:32             95.8 ns         93.9 ns      8632597 bytes_per_second=325.01M/s
BM_RE2Match/regex:3/size:32              165 ns          159 ns      4230093 bytes_per_second=192.132M/s
BM_RE2Match/regex:4/size:32              133 ns          130 ns      5603452 bytes_per_second=234.627M/s
BM_RE2Match/regex:5/size:32              129 ns          126 ns      5725878 bytes_per_second=241.591M/s
BM_RE2Match/regex:0/size:1024            198 ns          197 ns      3524921 bytes_per_second=4.83336G/s
BM_RE2Match/regex:1/size:1024            674 ns          670 ns      1033088 bytes_per_second=1.42407G/s
BM_RE2Match/regex:2/size:1024            223 ns          220 ns      3162541 bytes_per_second=4.34093G/s
BM_RE2Match/regex:3/size:1024           1883 ns         1820 ns       409817 bytes_per_second=536.602M/s
BM_RE2Match/regex:4/size:1024           1679 ns         1670 ns       412527 bytes_per_second=584.801M/s
BM_RE2Match/regex:5/size:1024           1695 ns         1685 ns       414182 bytes_per_second=579.535M/s
BM_RE2Match/regex:0/size:32768          3075 ns         3056 ns       228863 bytes_per_second=9.98527G/s
BM_RE2Match/regex:1/size:32768         19041 ns        18907 ns        37146 bytes_per_second=1.61409G/s
BM_RE2Match/regex:2/size:32768          3744 ns         3721 ns       186880 bytes_per_second=8.20214G/s
BM_RE2Match/regex:3/size:32768         50233 ns        49893 ns        13444 bytes_per_second=626.345M/s
BM_RE2Match/regex:4/size:32768         50097 ns        49802 ns        13859 bytes_per_second=627.488M/s
BM_RE2Match/regex:5/size:32768         49975 ns        49668 ns        13302 bytes_per_second=629.179M/s
BM_RE2Match/regex:0/size:1048576      207667 ns       206329 ns         3386 bytes_per_second=4.73304G/s
BM_RE2Match/regex:1/size:1048576      707228 ns       647207 ns         1181 bytes_per_second=1.50889G/s
BM_RE2Match/regex:2/size:1048576      246559 ns       244550 ns         2858 bytes_per_second=3.9933G/s
BM_RE2Match/regex:3/size:1048576     1685822 ns      1629115 ns          443 bytes_per_second=613.83M/s
BM_RE2Match/regex:4/size:1048576     1597964 ns      1588290 ns          442 bytes_per_second=629.608M/s
BM_RE2Match/regex:5/size:1048576     1602594 ns      1594127 ns          442 bytes_per_second=627.303M/s
BM_RE2Match/regex:0/size:33554432    7518029 ns      7474000 ns           93 bytes_per_second=4.18116G/s
BM_RE2Match/regex:1/size:33554432   20942655 ns     20293714 ns           35 bytes_per_second=1.53989G/s
BM_RE2Match/regex:2/size:33554432    9424043 ns      9116707 ns           82 bytes_per_second=3.42777G/s
BM_RE2Match/regex:3/size:33554432   56217759 ns     55287692 ns           13 bytes_per_second=578.791M/s
BM_RE2Match/regex:4/size:33554432   53830618 ns     53255000 ns           12 bytes_per_second=600.883M/s
BM_RE2Match/regex:5/size:33554432   52235840 ns     51905231 ns           13 bytes_per_second=616.508M/s
BM_PCRE2Match/regex:0/size:16             66.8 ns         60.8 ns     13875124 bytes_per_second=250.976M/s
BM_PCRE2Match/regex:1/size:16             88.7 ns         79.4 ns      9992862 bytes_per_second=192.165M/s
BM_PCRE2Match/regex:2/size:16             59.0 ns         57.3 ns      9481626 bytes_per_second=266.186M/s
BM_PCRE2Match/regex:3/size:16             59.0 ns         58.4 ns     11561266 bytes_per_second=261.062M/s
BM_PCRE2Match/regex:4/size:16             46.6 ns         46.1 ns     14931518 bytes_per_second=331.169M/s
BM_PCRE2Match/regex:5/size:16              310 ns          302 ns      2447321 bytes_per_second=50.4791M/s
BM_PCRE2Match/regex:0/size:32             51.6 ns         50.6 ns     13541485 bytes_per_second=602.697M/s
BM_PCRE2Match/regex:1/size:32             68.3 ns         67.7 ns      9596271 bytes_per_second=450.962M/s
BM_PCRE2Match/regex:2/size:32             48.5 ns         48.0 ns     14152965 bytes_per_second=635.407M/s
BM_PCRE2Match/regex:3/size:32             77.4 ns         75.7 ns      9662236 bytes_per_second=403.292M/s
BM_PCRE2Match/regex:4/size:32             54.5 ns         53.3 ns     13442409 bytes_per_second=572.369M/s
BM_PCRE2Match/regex:5/size:32              754 ns          741 ns       945346 bytes_per_second=41.1732M/s
BM_PCRE2Match/regex:0/size:1024            645 ns          639 ns      1031520 bytes_per_second=1.49177G/s
BM_PCRE2Match/regex:1/size:1024           1562 ns         1543 ns       463073 bytes_per_second=632.905M/s
BM_PCRE2Match/regex:2/size:1024            823 ns          817 ns       834148 bytes_per_second=1.16752G/s
BM_PCRE2Match/regex:3/size:1024           2741 ns         2529 ns       352872 bytes_per_second=386.077M/s
BM_PCRE2Match/regex:4/size:1024         360123 ns       356793 ns         1933 bytes_per_second=2.73706M/s
BM_PCRE2Match/regex:5/size:1024          27062 ns        26861 ns        25701 bytes_per_second=36.3555M/s
BM_PCRE2Match/regex:0/size:32768         23063 ns        22210 ns        36458 bytes_per_second=1.37403G/s
BM_PCRE2Match/regex:1/size:32768         48069 ns        44611 ns        14660 bytes_per_second=700.493M/s
BM_PCRE2Match/regex:2/size:32768         19899 ns        19697 ns        35921 bytes_per_second=1.54936G/s
BM_PCRE2Match/regex:3/size:32768        111223 ns       101239 ns         8664 bytes_per_second=308.674M/s
BM_PCRE2Match/regex:4/size:32768      16666791 ns     16480976 ns           42 bytes_per_second=1.89613M/s
BM_PCRE2Match/regex:5/size:32768        944913 ns       936398 ns          714 bytes_per_second=33.3726M/s
BM_PCRE2Match/regex:0/size:1048576      761132 ns       754457 ns          929 bytes_per_second=1.29439G/s
BM_PCRE2Match/regex:1/size:1048576     1620822 ns      1601384 ns          445 bytes_per_second=624.46M/s
BM_PCRE2Match/regex:2/size:1048576      825227 ns       815228 ns          907 bytes_per_second=1.1979G/s
BM_PCRE2Match/regex:3/size:1048576     2911238 ns      2837960 ns          248 bytes_per_second=352.366M/s
BM_PCRE2Match/regex:4/size:1048576   560549647 ns    555134000 ns            1 bytes_per_second=1.80137M/s
BM_PCRE2Match/regex:5/size:1048576    30511217 ns     30194417 ns           24 bytes_per_second=33.1187M/s
BM_PCRE2Match/regex:0/size:33554432   21465550 ns     21243485 ns           33 bytes_per_second=1.47104G/s
BM_PCRE2Match/regex:1/size:33554432   51115303 ns     49525833 ns           12 bytes_per_second=646.127M/s
BM_PCRE2Match/regex:2/size:33554432   24387834 ns     23899156 ns           32 bytes_per_second=1.30758G/s
BM_PCRE2Match/regex:3/size:33554432   86278970 ns     85179500 ns            8 bytes_per_second=375.677M/s
BM_PCRE2Match/regex:4/size:33554432 19615792401 ns   18874714000 ns            1 bytes_per_second=1.69539M/s
BM_PCRE2Match/regex:5/size:33554432 1035454278 ns   1014120000 ns            1 bytes_per_second=31.5545M/s
```
