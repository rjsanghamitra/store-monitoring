package main

import (
	"reports/db"
)

func main() {
	db.GenerateReport()
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "hello")
	// })
	// http.ListenAndServe(":8080", nil)
}
