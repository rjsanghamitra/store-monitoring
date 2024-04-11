package main

import (
	"net/http"
	"reports/db"
)

func main() {
	// fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/trigger_report", db.HandleTrigger)
	http.HandleFunc("/get_report", db.HandleGet)
	http.HandleFunc("/download", db.HandleDownload)
	http.ListenAndServe(":8080", nil)
}
