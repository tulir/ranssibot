package main

import (
	"errors"
	"strconv"
	"strings"
)

// Date ...
type Date struct {
	Year  int
	Month int
	Day   int
}

// DateToString ...
func DateToString(date Date) string {
	return strconv.Itoa(date.Day) + "." + strconv.Itoa(date.Month) + "." + strconv.Itoa(date.Year)
}

// StringToDate ...
func StringToDate(str string) (Date, error) {
	dateparts := strings.Split(str, ".")
	day, err1 := strconv.Atoi(dateparts[0])
	month, err2 := strconv.Atoi(dateparts[1])
	year, err3 := strconv.Atoi(dateparts[2])
	if err1 == nil && err2 == nil && err3 == nil {
		return Date{year, month, day}, nil
	}
	return Date{1970, 1, 1}, errors.New("Failed to parse date")
}

// Time ...
type Time struct {
	Hours   int
	Minutes int
}

// TimeToString ...
func TimeToString(time Time) string {
	hourStr := strconv.Itoa(time.Hours)
	if len(hourStr) == 1 {
		hourStr = "0" + hourStr
	}
	minStr := strconv.Itoa(time.Minutes)
	if len(minStr) == 1 {
		minStr = "0" + minStr
	}
	return hourStr + ":" + minStr
}

// StringToTime ...
func StringToTime(str string) (Time, error) {
	timeparts := strings.Split(str, ":")
	hour, err1 := strconv.Atoi(timeparts[0])
	minute, err2 := strconv.Atoi(timeparts[1])
	if err1 == nil && err2 == nil {
		return Time{hour, minute}, nil
	}
	return Time{0, 0}, errors.New("Failed to parse time")
}
