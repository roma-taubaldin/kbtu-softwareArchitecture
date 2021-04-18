package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	text, err := fmt.Printf("%s Listening on port :8000",time.Now().Format("02.01.2006 15:04:05"))
	if err != nil {
		return
	}
	fmt.Println(text)
	http.HandleFunc("/health", handleHealth())
	http.HandleFunc("/ping", handlePing())
	http.HandleFunc("/time", handleTime())
	http.HandleFunc("/weather", handleWeather())
	http.HandleFunc("/rates", handleRates())
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
			City string `json:"city"`
			Day  float64 `json:"day"`
			Night float64 `json:"night"`
			RainPercentage string `json:"rainPercentage"`
		}
		type Weather struct {
			WeatherCity []WeatherCity
		}
		weatherAlmaty := &WeatherCity{City: "Almaty", Day: 18.8, Night: 4.6, RainPercentage: "32"}
		weatherAstana := &WeatherCity{City: "Astana", Day: 11.4, Night: 0.3, RainPercentage: "11"}
		weather := []WeatherCity{}
		data := Weather{weather}
		data.WeatherCity = append(data.WeatherCity,*weatherAlmaty)
		data.WeatherCity = append(data.WeatherCity,*weatherAstana)
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
			Rate1 string `json:"buyRate"`
			Rate2  string `json:"sellRate"`
			Buy float64 `json:"buy"`
			Sell float64 `json:"sell"`
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
		data.Rate = append(data.Rate,*kztUsd)
		data.Rate = append(data.Rate,*kztEur)
		data.Rate = append(data.Rate,*kztRub)
		data.Rate = append(data.Rate,*kztKgs)
		data.Rate = append(data.Rate,*kztGbp)
		body, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return
		}
		fmt.Fprintf(w, string(body))
	}
}
