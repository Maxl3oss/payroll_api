package utils

import "time"

func ToThaiTime(dateInfo string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateInfo)
	if err != nil {
		return time.Time{}, err
	}
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, err
	}
	return date.In(location), nil
}
