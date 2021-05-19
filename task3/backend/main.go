package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
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
	pgdbInfo := fmt.Sprintf("host=sa.homework port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", pgdbInfo)
	if err != nil {
		log.Fatal(err)
	}
	text, err := fmt.Printf("%s Succesfull connected to DB", time.Now().Format("02.01.2006 15:04:05"))
	if err != nil {
		return
	}
	fmt.Println(text)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
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
	http.HandleFunc("/createTask", handleCreate(*db))
	http.HandleFunc("/deleteTask", handleDelete(*db))
	http.HandleFunc("/updateTask", handleUpdate(*db))
	http.HandleFunc("/viewTasks", handleView(*db))

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

func handleCreate(db sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.URL.Query().Get("task")

		_, err := db.Exec("insert into kbtu.tasks (task) values ($1)", task)
		if err != nil {
			log.Printf(err.Error())
		}
		fmt.Fprintf(w, "Values inserted")
	}
}

func handleView(db sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Fatal(err)
		}*/
		type Task struct {
			Id          int       `json:"id" db:"id"`
			Task        string    `json:"task" db:"task"`
			CorrectDate time.Time `json:"correctDate" db:"correct_date"`
		}

		var tasks []Task
		//sqlSelect := `select * from kbtu.tasks`
		res, err := db.Query("select id, task, correct_date from kbtu.tasks")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Close()
		for res.Next() {
			var t Task
			if err = res.Scan(&t.Id, &t.Task, &t.CorrectDate); err != nil {
				log.Fatal(err)
			}
			tasks = append(tasks, t)
		}
		body, err := json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}

func handleDelete(db sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		_, err := db.Exec("delete from kbtu.tasks where id = $1", id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Values deleted")
	}
}

func handleUpdate(db sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := strings.TrimSpace(query.Get("id"))
		task := strings.TrimSpace(query.Get("task"))
		update := fmt.Sprintf("update kbtu.tasks set task = '%s', correct_date = now() where id = %s;", task, id)
		_, err := db.Exec(update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Values updated")
	}
}
