package db

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var fileReady bool = false

func HandleTrigger(w http.ResponseWriter, r *http.Request) {
	go FileCreation()
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Home Page</title>
		</head>
		<body>
			<a href="/get_report">Click here to get the report</a>
		</body>
		</html>
	`)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	if !fileReady {
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<body>
			<p>The File is Downloading in the Home Directory</p>
		</body>
		</html>
	`)
	} else {
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<body>
			<p>Report is generated and stored in the Home Directory</p>
		</body>
		</html>
	`)
	}
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	home, _ := os.UserHomeDir()
	file, err := os.Open(home + "/" + ReportName + ".csv")
	CheckError(err)
	defer file.Close()

	w.Header().Set("Content-type", "text/csv")

	_, err = io.Copy(w, file)
	CheckError(err)
}

func FileCreation() {
	GenerateReport()
	fileReady = true
}
