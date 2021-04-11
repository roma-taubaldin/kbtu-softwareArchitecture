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