package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	dbhost     = "sa.homework"
	dbport     = 32715
	dbuser     = "admin"
	dbpassword = "admin"
	dbname     = "db"
)

var (
	db  *sql.DB
	ctx context.Context
)

func main() {
	pgdbInfo := fmt.Sprintf("host:=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", pgdbInfo)
	if err != nil {
		panic(err)
	}
	text, err := fmt.Printf("%s Succesfull connected to DB", time.Now().Format("02.01.2006 15:04:05"))
	if err != nil {
		return
	}
	fmt.Println(text)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	text, err = fmt.Printf("%s Succesfull pinged to DB", time.Now().Format("02.01.2006 15:04:05"))
	if err != nil {
		return
	}
	fmt.Println(text)

	text, err = fmt.Printf("%s Listening on port :8000", time.Now().Format("02.01.2006 15:04:05"))
	if err != nil {
		return
	}
	fmt.Println(text)
	http.HandleFunc("/health", handleHealth())
	http.HandleFunc("/ping", handlePing())
	http.HandleFunc("/time", handleTime())
	http.HandleFunc("/weather", handleWeather())
	http.HandleFunc("/rates", handleRates())
	http.HandleFunc("/createTask", handleCreate())
	http.HandleFunc("/deleteTask", handleDelete())
	http.HandleFunc("/updateTask", handleUpdate())
	http.HandleFunc("/viewTasks", handleView())

	http.ListenAndServe(":8000", nil)
}

func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Status struct {
			Status string `json:"status"`
		}
		status := &Status{Status: "OK"}
		body, err := json.MarshalIndent(status, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}

func handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func handleWeather() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type WeatherCity struct {
			City           string  `json:"city"`
			Day            float64 `json:"day"`
			Night          float64 `json:"night"`
			RainPercentage string  `json:"rainPercentage"`
		}
		type Weather struct {
			WeatherCity []WeatherCity
		}
		weatherAlmaty := &WeatherCity{City: "Almaty", Day: 18.8, Night: 4.6, RainPercentage: "32"}
		weatherAstana := &WeatherCity{City: "Astana", Day: 11.4, Night: 0.3, RainPercentage: "11"}
		weather := []WeatherCity{}
		data := Weather{weather}
		data.WeatherCity = append(data.WeatherCity, *weatherAlmaty)
		data.WeatherCity = append(data.WeatherCity, *weatherAstana)
		body, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}

func handleTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, time.Now().Format("15:04 02.01.2006"))
	}
}

func handleRates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Rate struct {
			Rate1 string  `json:"buyRate"`
			Rate2 string  `json:"sellRate"`
			Buy   float64 `json:"buy"`
			Sell  float64 `json:"sell"`
		}
		type Rates struct {
			Rate []Rate
		}
		kztUsd := &Rate{Rate1: "KZT", Rate2: "USD", Buy: 433.1, Sell: 434.7}
		kztEur := &Rate{Rate1: "KZT", Rate2: "EUR", Buy: 514.7, Sell: 517.1}
		kztRub := &Rate{Rate1: "KZT", Rate2: "RUB", Buy: 5.61, Sell: 5.67}
		kztKgs := &Rate{Rate1: "KZT", Rate2: "KGS", Buy: 4.85, Sell: 5.25}
		kztGbp := &Rate{Rate1: "KZT", Rate2: "GBP", Buy: 593.0, Sell: 603.0}
		rate := []Rate{}
		data := Rates{rate}
		data.Rate = append(data.Rate, *kztUsd)
		data.Rate = append(data.Rate, *kztEur)
		data.Rate = append(data.Rate, *kztRub)
		data.Rate = append(data.Rate, *kztKgs)
		data.Rate = append(data.Rate, *kztGbp)
		body, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.URL.Query().Get("task")
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Fatal(err)
		}
		sqlInsert := `insert into kbtu.tasks values ($2,$3)`
		_, err = db.Exec(sqlInsert, task, time.Now())
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Values inserted")
	}
}

func handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Fatal(err)
		}
		type Task struct {
			Id          int       `json:"id" db:"id"`
			Task        string    `json:"task" db:"task"`
			CorrectDate time.Time `json:"correctDate" db:"correct_date"`
		}
		type Tasks struct {
			Task []Task
		}
		var tasks []Tasks
		sqlSelect := `select * from kbtu.tasks`
		res, err := db.Query(sqlSelect)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		defer res.Close()
		for res.Next() {
			if err = res.Scan(&tasks); err != nil {
				log.Fatal(err)
			}
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		body, err := json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Fatal(err)
		}
		sqlDelete := `delete from kbtu.tasks where id = $1`
		_, err = db.Exec(sqlDelete, id)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Values deleted")
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		task := r.URL.Query().Get("task")
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Fatal(err)
		}
		sqlUpdate := `update kbtu.tasks set task = $2, correct_date = now() where id = $1`
		_, err = db.Exec(sqlUpdate, id, task)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Values updated")
	}
}
