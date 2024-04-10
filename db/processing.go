package db

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

var pollsData map[string][]*Polls = PollstoMap()
var storeTimings map[string][]Timings = StoreDatatoMap()
var storeTimezone map[string]*time.Location = GetTimeZone()

func FindStartAndEndtime(storeId string, dayNumber int, timestamp time.Time, timezone *time.Location) (time.Time, time.Time) {
	var startTime, endTime time.Time
	for _, timings := range storeTimings[storeId] { // O(1) since this will run only once
		if timings.dayOfWeek == dayNumber {
			startTime = time.Date(timestamp.Year(),
				timestamp.Month(),
				timestamp.Day(),
				timings.startTimeLocal[0],
				timings.startTimeLocal[1],
				timings.startTimeLocal[2],
				timings.startTimeLocal[3], timezone)
			endTime = time.Date(timestamp.Year(),
				timestamp.Month(),
				timestamp.Day(),
				timings.endTimeLocal[0],
				timings.endTimeLocal[1],
				timings.endTimeLocal[2],
				timings.endTimeLocal[3], timezone)
			break
		}
	}
	return startTime, endTime
}

func GenerateReport() {

	reportName := strconv.FormatInt(time.Now().UTC().UnixNano(), 16) // returns the string representation of the first arg in base 16
	reportFile, err := os.Create(reportName + ".csv")
	CheckError(err)
	csvWriter := csv.NewWriter(reportFile)
	headings := []string{"store ID", "uptime  Last Hour", "uptime Last Day Final", "uptime Last Week Final", "downtime Last Hour Final",
		"downtime Last Day Final", "downtime Last Week Final"}
	csvWriter.Write(headings)
	csvWriter.Flush()

	for _, storeId := range Stores {
		referenceTime := pollsData[storeId][0]

		timezone, ok := storeTimezone[storeId]
		if !ok {
			temp, err := time.LoadLocation("America/Chicago")
			CheckError(err)
			timezone = temp
		}
		*referenceTime.timestampOfPoll = referenceTime.timestampOfPoll.In(timezone)

		startTime, endTime := FindStartAndEndtime(storeId, int(referenceTime.timestampOfPoll.Weekday()), *referenceTime.timestampOfPoll, timezone)
		if referenceTime.timestampOfPoll.After(endTime) {
			referenceTime.timestampOfPoll = &endTime
		} else if referenceTime.timestampOfPoll.Before(startTime) {
			referenceTime.timestampOfPoll = &startTime
		}
		lastHourTime := referenceTime.timestampOfPoll.Add(-time.Hour).In(timezone)

		n := len(pollsData[storeId])
		var timeToCheckTill, timeToCheckFrom time.Time
		timeToCheckFrom = *referenceTime.timestampOfPoll

		// calculating uptime and downtime in the last hour
		uptimeLastHour, downtimeLastHour := 0, 0
		ind := 0 // this variable is used so that uptime/downtime in the last day can be added to the uptime/downtime in the last hour.
		for i := 0; i < n-1; i++ {

			*pollsData[storeId][i].timestampOfPoll = pollsData[storeId][i].timestampOfPoll.In(timezone)

			// checking if the timestamp lies within the business hours....to work on this more. corner case missed.
			if pollsData[storeId][i].timestampOfPoll.Before(startTime) && pollsData[storeId][i].timestampOfPoll.After(endTime) {
				timeToCheckFrom = endTime // since the polls are in descending order, we'll be starting again from the endtime
				continue
			}

			if pollsData[storeId][i].timestampOfPoll.Before(lastHourTime) { // we are checking 'before' because it is sorted in decreasing order.
				break
			}

			if pollsData[storeId][i].status == 1 {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					uptimeLastHour += int(timeToCheckFrom.Sub(timeToCheckTill))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastHourTime) {
					timeToCheckTill = lastHourTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				uptimeLastHour += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			} else {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					downtimeLastHour += int((timeToCheckFrom.Sub(timeToCheckTill)))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastHourTime) {
					timeToCheckTill = lastHourTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				downtimeLastHour += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			}
			ind++
		}

		// calculating uptime and downtime in the last day
		uptimeLastDay, downtimeLastDay := uptimeLastHour, downtimeLastHour
		lastDayTime := referenceTime.timestampOfPoll.Add(-time.Hour * 23)
		for i := ind; i < n-1; i++ {
			// we have to calculate startime and endtime everyday since they can be different for each day for some stores.
			startTime, endTime := FindStartAndEndtime(storeId, int(pollsData[storeId][i].timestampOfPoll.Weekday()), *pollsData[storeId][i].timestampOfPoll, timezone)

			*pollsData[storeId][i].timestampOfPoll = pollsData[storeId][i].timestampOfPoll.In(timezone)

			// checking if the timestamp lies within the business hours
			if pollsData[storeId][i].timestampOfPoll.Before(startTime) && pollsData[storeId][i].timestampOfPoll.After(endTime) {
				timeToCheckFrom = endTime // since the polls are in descending order, we'll be starting again from the endtime
				continue
			}

			if pollsData[storeId][i].timestampOfPoll.Before(lastDayTime) { // we are checking 'before' because it is sorted in decreasing order.
				break
			}

			if pollsData[storeId][i].status == 1 {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					uptimeLastDay += int(timeToCheckFrom.Sub(timeToCheckTill))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastDayTime) {
					timeToCheckTill = lastDayTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				uptimeLastDay += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			} else {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					downtimeLastDay += int((timeToCheckFrom.Sub(timeToCheckTill)))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastDayTime) {
					timeToCheckTill = lastDayTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				downtimeLastDay += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			}
			ind++
		}

		// calculating uptime and downtime for the last week
		uptimeLastWeek, downtimeLastWeek := uptimeLastDay, downtimeLastDay
		lastWeekTime := referenceTime.timestampOfPoll.Add(-time.Hour * 24 * 6)
		for i := ind; i < n-1; i++ {
			// we have to calculate startime and endtime everyday since they can be different for each day for some stores.
			startTime, endTime := FindStartAndEndtime(storeId, int(pollsData[storeId][i].timestampOfPoll.Weekday()), *pollsData[storeId][i].timestampOfPoll, timezone)

			*pollsData[storeId][i].timestampOfPoll = pollsData[storeId][i].timestampOfPoll.In(timezone)

			// checking if the timestamp lies within the business hours
			if pollsData[storeId][i].timestampOfPoll.Before(startTime) && pollsData[storeId][i].timestampOfPoll.After(endTime) {
				timeToCheckFrom = endTime // since the polls are in descending order, we'll be starting again from the endtime
				continue
			}

			if pollsData[storeId][i].timestampOfPoll.Before(lastWeekTime) { // we are checking 'before' because it is sorted in decreasing order.
				break
			}

			if pollsData[storeId][i].status == 1 {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					uptimeLastWeek += int(timeToCheckFrom.Sub(timeToCheckTill))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastWeekTime) {
					timeToCheckTill = lastWeekTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				uptimeLastWeek += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			} else {
				// extrapolation
				if pollsData[storeId][i+1].status != pollsData[storeId][i].status {
					timeToCheckTill = pollsData[storeId][i].timestampOfPoll.Truncate(time.Hour)
					if timeToCheckTill.Before(startTime) {
						timeToCheckTill = startTime
					}
					downtimeLastWeek += int((timeToCheckFrom.Sub(timeToCheckTill)))
					timeToCheckFrom = timeToCheckTill
					continue
				}

				if pollsData[storeId][i+1].timestampOfPoll.In(timezone).Before(lastWeekTime) {
					timeToCheckTill = lastWeekTime
				} else {
					timeToCheckTill = pollsData[storeId][i+1].timestampOfPoll.In(timezone)
				}

				downtimeLastWeek += int(timeToCheckFrom.Sub(timeToCheckTill))
				timeToCheckFrom = timeToCheckTill

			}
			ind++
		}

		ref := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		uptimeLastHourFinal := time.Date(0, 1, 1, 0, 0, 0, uptimeLastHour, time.UTC).Sub(ref).String()
		uptimeLastDayFinal := time.Date(0, 1, 1, 0, 0, 0, uptimeLastDay, time.UTC).Sub(ref).String()
		uptimeLastWeekFinal := time.Date(0, 1, 1, 0, 0, 0, uptimeLastWeek, time.UTC).Sub(ref).String()
		downtimeLastHourFinal := time.Date(0, 1, 1, 0, 0, 0, downtimeLastHour, time.UTC).Sub(ref).String()
		downtimeLastDayFinal := time.Date(0, 1, 1, 0, 0, 0, downtimeLastDay, time.UTC).Sub(ref).String()
		downtimeLastWeekFinal := time.Date(0, 1, 1, 0, 0, 0, downtimeLastWeek, time.UTC).Sub(ref).String()

		values := []string{storeId, uptimeLastHourFinal, uptimeLastDayFinal, uptimeLastWeekFinal, downtimeLastHourFinal, downtimeLastDayFinal, downtimeLastWeekFinal}

		// storing them in a csv file

		csvWriter.Write(values)
		csvWriter.Flush()
	}
	defer reportFile.Close()
}
