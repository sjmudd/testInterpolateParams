package main

/*

See: https://bugs.mysql.com/bug.php?id=101630

Test with

$ ./testInterpolateParams interpolate 2>&1 | less
$ ./testInterpolateParams 2>&1 | less

Tests as of 17/11/2020 show the go driver works properly.

An error is triggered with > 64k rows without the interpolateParams
setting and works fine with it set.

Requires a mysql server running on 127.0.0.1:3306 with user root
and password root.

*/

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	size := 100 * 1000
	dsn := "root:root@tcp(127.0.0.1:3306)/"

	if len(os.Args) > 1 {
		if os.Args[1] == "interpolate" {
			log.Print("setting interpolateParams=true")
			dsn += "?interpolateParams=true"
		}
	}

	log.Print("dsn:", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}
	log.Print("Open succeeded")
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	// do a simple query to check things are workingg
	sql := "SELECT 1"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatalf("simple query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		log.Printf("Got back: %v", id)
	}

	// build the very large list
	if len(os.Args) > 2 {
		size, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("strconvAtoi of arg 2 failed: %v", err)
		}
		log.Printf("size set to %v", size)
	}
	fixed := []string{}
	args := []interface{}{}
	for i := 0; i < size; i++ {
		fixed = append(fixed, "?")
		args = append(args, 1)
	}
	sql = "SELECT " + strings.Join(fixed, ",")

	_, err = db.Query(sql, args...)
	if err != nil {
		log.Fatalf("Query failed: err: %q, sql: %q, args: %+v", err, sql, args)
	}
	log.Printf("query worked with size: %v", size)
}
