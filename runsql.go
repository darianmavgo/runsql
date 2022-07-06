package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "/Users/darianhickman/Documents/wc_study/history.db"

	engine := sqlx.MustConnect("sqlite3", dbPath)

	// qry := `select list_folder||'/'|| sql_file as script from run_list where 'order' > 0 ORDER BY 'order'; `
	qry := `select list_folder||'/'|| sql_file as script from run_list 
	where run_order > 0 
	ORDER BY run_order;`

	scripts, err := engine.Query(qry)
	if err != nil {
		println("queary failed ", qry)
	}
	var file string
	for scripts.Next() {
		err = scripts.Scan(&file)
		if err != nil {
			println(err)
		}
		println(file[len(file)-80:])
		sqlScript, err := ioutil.ReadFile(file)
		if err != nil {
			println("reading script file failed\n", file)
		}
		if _, err := engine.Exec(string(sqlScript)); err != nil {
			fmt.Println(string(sqlScript[:80]))
			fmt.Println(err)

		} else {
			fmt.Println("Success ", string(sqlScript[:80]))
		}
	}

}
