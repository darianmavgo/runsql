package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// dbPath := "file:/Users/darianhickman/Documents/wc_study/history.db?cache=shared"
	// dbPath := "file:/Users/darianhickman/Documents/wc_study/history.db"
	dbPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal("Failed db string ", os.Args[1], err)
	}
	log.Print("db ", dbPath)

	db := sqlx.MustConnect("sqlite3", dbPath)
	defer db.Close()
	// db.SetMaxOpenConns(1)

	// qry := `select list_folder||'/'|| sql_file as script from run_list where 'order' > 0 ORDER BY 'order'; `
	qry := `select list_folder||'/'|| sql_file as script, run_order from run_list 
	where run_order > 0 
	ORDER BY run_order;`

	scripts, err := db.Query(qry)
	if err != nil {
		println("script qry failed ", qry)
	}
	defer scripts.Close()
	var files []string
	var orders []int
	var file string
	var run_order int
	// Giving up on executing sql scripts while reading the query of scripts.
	for scripts.Next() {
		err = scripts.Scan(&file, &run_order)
		if err != nil {
			println(err)
		} else {
			files = append(files, file)
			orders = append(orders, run_order)
		}
	}

	for idx, file := range files {
		log.Print(orders[idx], " ")
		if len(file) > 80 {
			println(idx, file[len(file)-80:])
		} else {
			println(idx, file)
		}

		sqlScript, err := ioutil.ReadFile(file)
		if err != nil {
			println("reading script file failed\n", file)
		}
		if result, err := db.Exec(string(sqlScript)); err != nil {
			fmt.Println(string(sqlScript[:80]))
			fmt.Println(err)

		} else {
			spew.Dump(result)
		}
	}

}
