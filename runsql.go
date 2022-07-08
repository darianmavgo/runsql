package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// dbPath := "file:/Users/darianhickman/Documents/wc_study/history.db?cache=shared"
	dbPath := "file:/Users/darianhickman/Documents/wc_study/history.db"
	db := sqlx.MustConnect("sqlite3", dbPath)
	defer db.Close()
	// db.SetMaxOpenConns(1)

	// qry := `select list_folder||'/'|| sql_file as script from run_list where 'order' > 0 ORDER BY 'order'; `
	qry := `select list_folder||'/'|| sql_file as script from run_list 
	where run_order > 0 
	ORDER BY run_order;`

	scripts, err := db.Query(qry)
	if err != nil {
		println("script qry failed ", qry)
	}
	defer scripts.Close()
	var files []string
	var file string

	// Giving up on executing sql scripts while reading the query of scripts.
	for scripts.Next() {
		err = scripts.Scan(&file)
		if err != nil {
			println(err)
		} else {
			files = append(files, file)
		}
	}

	for _, file := range files {

		println(file[len(file)-80:])
		sqlScript, err := ioutil.ReadFile(file)
		if err != nil {
			println("reading script file failed\n", file)
		}
		if _, err := db.Exec(string(sqlScript)); err != nil {
			fmt.Println(string(sqlScript[:80]))
			fmt.Println(err)

		} else {
			fmt.Println("Success ", string(sqlScript[:80]))
		}
	}

}
