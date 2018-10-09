# proxy

In a Dell Precision T1500: 

```bash
docker run -p 7777:80 -d nginx
PROXY_SERV=":7778" PROXY_DEST="http://localhost:7777/" go run main.go 2> /dev/null
docker run --rm --read-only -v `pwd`:`pwd` -w `pwd` jordi/ab -k -c 128 -n 521000 http://172.17.0.1:7778/                                                                                                                                         
This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 172.17.0.1 (be patient)
Completed 52100 requests
Completed 104200 requests
Completed 156300 requests
Completed 208400 requests
Completed 260500 requests
Completed 312600 requests
Completed 364700 requests
Completed 416800 requests
Completed 468900 requests
Completed 521000 requests
Finished 521000 requests


Server Software:        nginx/1.15.3
Server Hostname:        172.17.0.1
Server Port:            7778

Document Path:          /
Document Length:        612 bytes

Concurrency Level:      128
Time taken for tests:   30.459 seconds
Complete requests:      521000
Failed requests:        0
Keep-Alive requests:    521000
Total transferred:      442850000 bytes
HTML transferred:       318852000 bytes
Requests per second:    17105.13 [#/sec] (mean)
Time per request:       7.483 [ms] (mean)
Time per request:       0.058 [ms] (mean, across all concurrent requests)
Transfer rate:          14198.59 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       4
Processing:     0    7   2.0      7      33
Waiting:        0    7   2.0      7      33
Total:          0    7   2.0      7      33
```

```bash
netstat -nat | grep TIME | wc -l
5381
```

In a MBP mid 2012:

```bash
$ docker run --rm --read-only -v `pwd`:`pwd` -w `pwd` jordi/ab -k -c 32 -n 52100 http://192.168.1.38:7778/
This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 192.168.1.38 (be patient)
Completed 5210 requests
Completed 10420 requests
Completed 15630 requests
Completed 20840 requests
Completed 26050 requests
Completed 31260 requests
Completed 36470 requests
Completed 41680 requests
Completed 46890 requests
Completed 52100 requests
Finished 52100 requests


Server Software:        nginx/1.15.5
Server Hostname:        192.168.1.38
Server Port:            7778

Document Path:          /
Document Length:        612 bytes

Concurrency Level:      32
Time taken for tests:   62.874 seconds
Complete requests:      52100
Failed requests:        0
Keep-Alive requests:    52100
Total transferred:      44285000 bytes
HTML transferred:       31885200 bytes
Requests per second:    828.64 [#/sec] (mean)
Time per request:       38.618 [ms] (mean)
Time per request:       1.207 [ms] (mean, across all concurrent requests)
Transfer rate:          687.84 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0      12
Processing:     4   38  16.2     37     291
Waiting:        1   38  16.1     36     291
Total:          4   39  16.2     37     292

Percentage of the requests served within a certain time (ms)
  50%     37
  66%     43
  75%     47
  80%     50
  90%     58
  95%     67
  98%     77
  99%     88
 100%    292 (longest request)
```