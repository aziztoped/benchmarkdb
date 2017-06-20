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

	DBtype = "cockroach" //change this value & adjust connection below

	//init connection
	switch DBtype {
	case "postgre":
		dbDriver = "postgres"
		dbStr = "postgres://akangaziz@localhost:5432/testing?sslmode=disable" //for postgreSQL
		SqlTruncate = "truncate table shop_campaigns;"
		SqlInsert = "insert into shop_campaigns(campaign_id, shop_id, fg_status, read_status, read_time) values(123123123, $1, 1,1,$2);"
		SqlReadUsingWhere = "SELECT campaign_id, fg_status, read_status, read_time FROM shop_campaigns WHERE shop_id = $1"
	case "cockroach":
		dbDriver = "postgres"
		dbStr = "postgresql://root@localhost:26257?sslmode=disable" //for cockroachDB
		SqlTruncate = "truncate table testing.shop_campaigns;"
		SqlInsert = "insert into testing.shop_campaigns(campaign_id, shop_id, fg_status, read_status, read_time) values(123123123, $1, 1,1,$2);"
		SqlReadUsingWhere = "SELECT campaign_id, fg_status, read_status, read_time FROM testing.shop_campaigns WHERE shop_id = $1"
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
		n := 1000000
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

		// //test read without where condition
		// timeRead = time.Now()
		// testRead(true)
		// tRead2 := time.Now().Sub(timeRead).Seconds() / 1000
		// fmt.Printf("test read without where condition done in %f second\n", tRead2)

		//write result to csv
		tWriteString := strconv.FormatFloat(tWrite, 'f', -1, 64)
		tRead1String := strconv.FormatFloat(tRead1, 'f', -1, 64)
		// tRead2String := strconv.FormatFloat(tRead2, 'f', -1, 64)

		writeCSV([][]string{{tWriteString, tRead1String}})
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
		now := time.Now().Format("2006-01-02 15:04:05")
		if _, err := DB.Exec(SqlInsert, shopID, now); err != nil {
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
