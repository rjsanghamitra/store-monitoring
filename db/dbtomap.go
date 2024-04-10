package db

import (
	"time"
)

func PollstoMap() map[string][]*Polls {
	data := make(map[string][]*Polls)
	row, err := db.Query("SELECT * FROM polls")
	CheckError(err)
	defer row.Close()
	for row.Next() {
		var a string
		var b int
		var c string
		row.Scan(&a, &b, &c)
		storeId := a
		Stores = append(Stores, storeId)
		timestampOfPoll := StrtoTime1(c)
		status := b
		data[storeId] = append(data[storeId], &Polls{
			status:          status,
			timestampOfPoll: timestampOfPoll,
		})
	}
	return data
}

func StoreDatatoMap() map[string][]Timings {
	data := make(map[string][]Timings)
	row, err := db.Query("SELECT * FROM store_data")
	CheckError(err)
	defer row.Close()

	for row.Next() {
		var storeId string
		var day int
		var temp1, temp2 string
		row.Scan(&storeId, &day, &temp1, &temp2)
		starth, startm, starts, startn := StrtoTime2(temp1)
		endh, endm, ends, endn := StrtoTime2(temp2)
		data[storeId] = append(data[storeId], Timings{
			dayOfWeek:      day,
			startTimeLocal: []int{starth, startm, starts, startn},
			endTimeLocal:   []int{endh, endm, ends, endn},
		})
	}
	return data
}

func GetTimeZone() map[string]*time.Location {
	data := make(map[string]*time.Location)
	row, err := db.Query("SELECT * FROM timezone")
	CheckError(err)
	for row.Next() {
		var a, b string
		row.Scan(&a, &b)
		x, _ := time.LoadLocation(b)
		data[a] = x
	}
	return data
}
