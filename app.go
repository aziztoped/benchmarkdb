package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"math/rand"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var DBtype, SqlTruncate, SqlInsert, SqlRead, SqlReadUsingWhere string

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func init() {
	var err error
	var dbStr, dbDriver string

	DBtype = "mysql" //change this value & adjust connection below

	//init connection
	switch DBtype {
	case "postgre":
		dbDriver = "postgres"
		dbStr = "postgres://akangaziz@localhost:5432/products?sslmode=disable" //for postgreSQL
		SqlTruncate = "truncate table product_shippings;"
		SqlInsert = "insert into product_shippings(shop_id, shipping_ids) values($1,'|1|2|3|');"
		SqlReadUsingWhere = "SELECT id, shop_id, shipping_ids FROM product_shippings WHERE shipping_ids LIKE '|2|' ORDER BY id DESC"
		SqlRead = "SELECT id, shop_id, shipping_ids FROM product_shippings  ORDER BY id DESC"
	case "mysql":
		dbDriver = "mysql"
		dbStr = "mysql://root:@localhost:3306/products?parseTime=true&loc=Local" //for postgreSQL
		SqlTruncate = "truncate table product_shippings;"
		SqlInsert = "insert into product_shippings(shop_id, shipping_ids) values($1,'|1|2|3|');"
		SqlReadUsingWhere = "SELECT id, shop_id, shipping_ids FROM product_shippings WHERE shipping_ids LIKE '|2|' ORDER BY id DESC"
		SqlRead = "SELECT id, shop_id, shipping_ids FROM product_shippings  ORDER BY id DESC"
	case "cockroachDB":
		dbDriver = "postgres"
		dbStr = "postgresql://root@localhost:26257?sslmode=disable" //for cockroachDB
		SqlTruncate = "truncate table products.product_shippings;"
		SqlInsert = "insert into products.product_shippings(shop_id, shipping_ids) values($1,'|1|2|3|');"
		SqlReadUsingWhere = "SELECT id, shop_id, shipping_ids FROM products.product_shippings WHERE product_shippings.shipping_ids LIKE '|2|' ORDER BY id DESC"
		SqlRead = "SELECT id, shop_id, shipping_ids FROM products.product_shippings  ORDER BY id DESC"
	}

	DB, err = sql.Open(dbDriver, dbStr)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
}

func main() {
	doTest(50)
}

func doTest(loop int) {
	i := 1
	for i <= loop {
		fmt.Printf("loop i = %d\n", i)
		//set n
		n := 1000
		fmt.Printf("set iteration n to %d\n", n)

		//truncate
		truncate()

		//test write
		timeWrite := time.Now()
		testWrite(n)
		tWrite := time.Now().Sub(timeWrite).Seconds() / 1000
		fmt.Printf("test write done in %f second\n", tWrite)

		//test read with where condition
		timeRead := time.Now()
		testRead(false)
		tRead1 := time.Now().Sub(timeRead).Seconds() / 1000
		fmt.Printf("test read using where done in %f second\n", tRead1)

		//test read without where condition
		timeRead = time.Now()
		testRead(true)
		tRead2 := time.Now().Sub(timeRead).Seconds() / 1000
		fmt.Printf("test read without where condition done in %f second\n", tRead2)

		//write result to csv
		tWriteString := strconv.FormatFloat(tWrite, 'f', -1, 64)
		tRead1String := strconv.FormatFloat(tRead1, 'f', -1, 64)
		tRead2String := strconv.FormatFloat(tRead2, 'f', -1, 64)

		writeCSV([][]string{{tWriteString, tRead1String, tRead2String}})
		fmt.Printf("====================\n")
		i++
	}
}

func truncate() {
	if _, err := DB.Exec(SqlTruncate); err != nil {
		log.Fatalf("fatal truncate DB: %s", err)
	}
}

func testWrite(n int) {
	i := 1
	for i <= n {
		shopID := random(100, 100000)
		if _, err := DB.Exec(SqlInsert, shopID); err != nil {
			log.Fatalf("fatal test write: %s", err)
		}
		i++
	}
}

func testRead(usingWhere bool) {
	var q string

	if usingWhere {
		q = SqlReadUsingWhere
	} else {
		q = SqlRead
	}
	rows, err := DB.Query(q)
	if err != nil {
		log.Fatalf("fatal test read: %s", err)
	}
	defer rows.Close()
	var id, balance int
	var shippingIds string
	for rows.Next() {
		if err := rows.Scan(&id, &balance, &shippingIds); err != nil {
			log.Fatalf("fatal scan read rows: %s", err)
		}
		//fmt.Printf("%d %d %s\n", id, balance, shippingIds)
	}
}

func writeCSV(data [][]string) {
	file, err := os.OpenFile("result.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}
