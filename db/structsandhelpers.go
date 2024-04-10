package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Polls struct {
	status          int
	timestampOfPoll *time.Time
}

type Timings struct {
	dayOfWeek      int
	startTimeLocal []int
	endTimeLocal   []int
}

var db, _ = sql.Open("sqlite3", "./db/data.db")

var Stores []string

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func StrtoTime1(s string) *time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return &t
}

func StrtoTime2(timestamp string) (int, int, int, int) {
	t, _ := time.Parse("15:04:05", timestamp)
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	n := t.Nanosecond()
	return h, m, s, n
}
