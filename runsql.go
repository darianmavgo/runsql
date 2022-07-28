package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// #run a script to clean for next run.
// func clean_run(script)
type Config struct {
	loadtablepath string
}

func main() {

	dbPath, err := filepath.Abs(os.Args[1]) // passing in settings.db
	if err != nil {
		log.Fatal("Failed db string ", os.Args[1], err)
	}
	log.Print("db ", dbPath)

	db := sqlx.MustConnect("sqlite3", dbPath)
	defer db.Close()

	cfgqry := `select loadtablepath from config;`
	c := Config{}
	db.Get(&c, cfgqry)
	// now I have the config so time to run backtest against the right file.

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

	loadDbPath, err := filepath.Abs(`/Users/darianhickman/Documents/wc_study/history_large.db`)
	// loadDbPath, err := filepath.Abs(c.loadtablepath)
	if err != nil {
		log.Fatalln(`Couldnot create path`, err, c.loadtablepath)
	}
	log.Println("loadDbPath ", loadDbPath)

	target_db := sqlx.MustConnect("sqlite3", loadDbPath)
	for idx, file := range files {
		log.Print(orders[idx], " ")
		if len(file) > 80 {
			println(idx, file[len(file)-80:])
		} else {
			println(idx, file)
		}

		sqlScript, err := ioutil.ReadFile(file)
		// cmd := `sqlite3 ` + dbPath + ` ".read ` + file + `"`
		// result, err := db.Exec(cmd)
		// result, err := exec.Command(cmd).Output()
		if err != nil {
			log.Println("Fail: ", file, err)
		}
		if result, err := target_db.Exec(string(sqlScript)); err != nil {
			log.Println(string(sqlScript[:80]))
			log.Println(err)

		} else {
			log.Println(result.LastInsertId())
		}
	}

}
