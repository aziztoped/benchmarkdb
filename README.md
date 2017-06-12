
Benchmark DB  

## Introduction 

i created this simple repo to test performance for comparing postgreDB, mysql & cockroachDB


## Getting started

### Pre Requirements
1. mysql, postgrelsql & cockroachDB installed in your local machine
2. golang v1.8 or latest
3. db & table created before running this benchmark app

for cockroachDB
```
create database products;
create table products.product_shippings(id SERIAL, shop_id INT , shipping_ids string);
```

for postgreSQL
```
CREATE TABLE "product_shippings" (
  "id" serial NOT NULL,
  "shop_id" integer NOT NULL DEFAULT '0',
  "shipping_ids" character varying(150) NOT NULL
);
```

for mysql 
```
CREATE TABLE `product_shippings` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `shop_id` int NOT NULL,
  `shipping_ids` varchar(150) NOT NULL
);
```

### Run

```
$ git clone https://github.com/aziztoped/benchmarkdb
$ cd benchmarkdb
$ go build
$ ./benchmarkdb
    loop i = 1
    set iteration n to 1000
    test write done in 0.001664 second
    test read using where done in 0.000007 second
    test read without where condition done in 0.000002 second
    ====================
    loop i = 2
    set iteration n to 1000
    test write done in 0.001410 second
    test read using where done in 0.000004 second
    test read without where condition done in 0.000001 second
    ====================
    loop i = 3
    set iteration n to 1000
    test write done in 0.001475 second
    test read using where done in 0.000003 second
    test read without where condition done in 0.000001 second
    ====================
    loop i = 4
    set iteration n to 1000
    test write done in 0.001469 second
    test read using where done in 0.000004 second
    test read without where condition done in 0.000001 second

```
### Output Benchmark Result
output file is located at result.csv

### Result 


#### cockroachDB benchmark result

cockroachDB use 3 nodes in same machine <br>

<img src="results/result-cockroach.png">
<br><br>
Average write: 0.003860905	<br>
Average read using where query: 0.000604311	<br>
Average read without where query: 0.000570346<br><br><br>

#### postgreSQL benchmark result
<img src="results/result-postgresql.png">
<br><br>
Average write: 0.001467372	<br>
Average read using where query: 0.000003697182	<br>
Average read without where query: 0.000000808635 <br><br><br>

####  mysql benchmark result<br>
@TO-DO: no time to convert the query, will DO next day.. 

 

