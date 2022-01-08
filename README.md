# SIMPLE API

Simple golang app using std lib net/http package to serve two endpoints REST API

### Endpoints GetAllUsers

```bash
curl localhost:8080/users/all
```
### Endpoints CreateUser

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"first_name": "linuxize","last_name": "linuxize","phone": "53555555","job_title": "DEV", "email": "linuxize@example.com"}' \
    localhost:8080/users/create
```

### Bench

```bash

# simple command two use 
# The values can be change. 
ab -n 100 -c 10 -g out.data http://localhost:8080/users/all > ab.txt


This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient).....done


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /users/all
Document Length:        253554 bytes

Concurrency Level:      10
Time taken for tests:   57.831 seconds
Complete requests:      100
Failed requests:        0
Total transferred:      25365100 bytes
HTML transferred:       25355400 bytes
Requests per second:    1.73 [#/sec] (mean)
Time per request:       5783.054 [ms] (mean)
Time per request:       578.305 [ms] (mean, across all concurrent requests)
Transfer rate:          428.33 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       0
Processing:  2858 5393 1317.4   5252    9309
Waiting:     2858 5393 1317.4   5252    9309
Total:       2858 5393 1317.4   5252    9309

Percentage of the requests served within a certain time (ms)
  50%   5252
  66%   5690
  75%   6154
  80%   6467
  90%   7428
  95%   7706
  98%   8363
  99%   9309
 100%   9309 (longest request)

```