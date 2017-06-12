
Benchmark DB  

## Introduction 

i created this simple repo to test performance for comparing postgreDB, mysql & cockroachDB


## Getting started

### Pre Requirements
1. mysql, postgrelsql & cockroachDB installed in your local machine
2. golang v1.8 or latest

### Run

```
$ git clone <git-repo-url>
$ cd benchmarkdb
$ go build
$ ./benchmarkdb
    set iteration n to 1000
    test write done in 0.003378 second
    test read using where done in 0.000534 second
    test read without where condition done in 0.000459 second

```
### Output Benchmark Result
output file is located at result.csv

### Result 

cockroachDB benchmark result
<img src="results/cockroach-result.png">
<br><br>


