package main

import (
	"database/sql"
	"fmt" // For test
	"log"
	rand "math/rand"
	"os"
	"runtime/trace"

	_ "github.com/go-sql-driver/mysql" // go get github.com/go-sql-driver/mysql
)

/*
type uffw struct {
    buben, gusli	string
    digits		int
    hash_anybody	float64
}
*/
//
func main() {
	// Debuger: go tool trace//
	traceDebug, err := os.Create("trace.out")
	checkError(err)
	defer traceDebug.Close()

	err = trace.Start(traceDebug)
	checkError(err)
	defer trace.Stop()

	tablExist()

	dbIsert()
}

// Инcтрукция по драйверу https://github.com/go-sql-driver/mysql/
func dbConnect() (db *sql.DB) {
	const (
		dbDriver string = "mysql"
		dbUser   string = "tests"
		dbPassw  string = "A47Y%k[@"
		dbName   string = "treq"
		dbHost   string = "localhost"
	)

	//db, err := sql.Open(dbDriver, dbUser+":"+dbPassw+"@tcp("+dbHost+":3306)/"+dbName) // DSN  username:password@protocol(address)/dbname?param=value
	//db, err := sql.Open("mysql", "user:password@tcp(address:3306)/db_name?multiStatements=false&charset=utf8mb4,utf8")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassw+"@tcp("+dbHost+":3306)/"+dbName+"?multiStatements=false&collation=utf8mb4_unicode_ci")
	checkError(err)
	log.Println("Have connect")
	return db
}

func tablExist() {
	db := dbConnect()
	defer db.Close()
	// check exists table
	//row, err := db.Query("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'bzdyk') AND (TABLE_NAME = 'uffw');")
	//row, err := db.Query("SHOW TABLES LIKE 'uffw';")
	var tableName = string("uffw")
	var tLimit = string("LIMIT 1;")
	var newTable string = `(
		identy INT UNSIGNED NOT NULL AUTO_INCREMENT,
		buben TEXT,
		gusli TEXT,
		hash_anybody VARCHAR(72) NULL,
		digits INT NULL,
		PRIMARY KEY (identy)
		) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET utf8mb4 COLLATE = utf8mb4_unicode_ci;`

	_, err := db.Query("SELECT 1 FROM" + tableName + " " + tLimit)
	if err != nil {
		var tableStatements = string("create table" + tableName + newTable)

		_, err := db.Exec(tableStatements)
		checkError(err)
		fmt.Println("Table was create")
	} else {
		fmt.Println("Table is exists")
	}

	/*
		columns, err := row.Columns()
		checkError(err)
		count := len(columns)

		fmt.Printf("%T,  %v\n", row, row)
		fmt.Println(row)

		fmt.Printf("%T,  %v\n", columns, columns)
		fmt.Println(columns)

		fmt.Printf("%T,  %v\n", count, count)
		fmt.Println(count)
	*/
}

func dbIsert() {
	db := dbConnect()
	defer db.Close()
	// check exists table
	//rows, err := db.Query("SHOW TABLES LIKE 'uffw';")
	var buben = string(randData())
	gusli := randData()
	digits := rand.Intn(7715)
	hash_anybody := (rand.Float64() * 774) + 255

	fmt.Printf("%T,  %v\n", buben, buben)
	fmt.Printf("%T,  %v\n", gusli, gusli)

	var tableName = string("uffw")
	var insertString string = "(buben, gusli, hash_anybody, digits) VALUES(?, ?, ?, ?)"
	insT, err := db.Prepare("INSERT INTO " + tableName + insertString)
	//insT, err := db.Prepare("INSERT INTO uffw (buben, gusli, hash_anybody, digits) VALUES(?, ?, ?, ?)")
	checkError(err)
	res, err := insT.Exec(buben, gusli, hash_anybody, digits)
	checkError(err)
	id, _ := res.LastInsertId()
	fmt.Printf("Inserted row: %d", id)
}

/*
//
func dbSelect () {
    db := dbConnect()
    rows, err := db.Query("SELECT * FROM uffw")
    checkError(err)

    defer db.Close()

    for rows.Next() {
    	var identy int
    	var buben string
    	var gusli string
    	var hash_anybody float64
    	var digits int
    	err = rows.Scan(&identy, &buben, &gusli, &hash_anybody, &digits)
    	checkError(err)
    	fmt.Printf("%T,  %v\n", identy, identy)
    	fmt.Printf("%T,  %v\n", buben, buben)
    	fmt.Printf("%T,  %v\n", gusli, gusli)
    	fmt.Printf("%T,  %v\n", hash_anybody, hash_anybody)
    }
}
*/

func randData() string {
	var alpha = "muabcdQK947Vefghijkmnpq4ei4st ]]]'}"
	q := rand.Intn(112)
	fmt.Println(q)
	//w := rand.Float64()
	//fmt.Println(w)
	//fmt.Println((rand.Float64() * 'q') + 255*'w'/'q')
	buf := make([]byte, q)
	for i := 0; i < q; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	//e:= string(buf)
	//fmt.Println(e)
	//fmt.Println(string(buf))
	return string(buf)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func getError(err string) {
	fmt.Println("Error is: ", err)
	os.Exit(1)
}
