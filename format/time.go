package format

import "fmt"
import "time"

// import "strings"
// import "errors"

var days = []string{
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}

var monthes = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"November",
	"October",
	"December",
}

func tryToParseStrTime(timeString string) (time.Time, error) {
	// Format used "2006 Jan 2 15:04:05 +0000"
	// var dayName string
	// var find = func(t string, strs []string) (string, int, int, error) {
	// 	// first check if the standard 3 characters works
	// 	for _, str := range strs {
	// 		if index := strings.Index(t, str[:3]); index >= 0 {
	// 			return str, index, 3, nil
	// 		}
	// 	}

	// 	// else do all combinations until 4 characters
	// 	// TODO do lower/upper ?
	// 	for _, str := range strs {
	// 		for i := len(str); i > 3; i-- {
	// 			if index := strings.Index(t, str[:3]); index >= 0 {
	// 				return str[:3], index, i, nil
	// 			}
	// 		}
	// 	}
	// 	return "", 0, 0, errors.New("")
	// }

	// var findDayEnglishName = func(t string) (string, int, int, error) {
	// 	val, index, length, err := find(t, days)
	// 	if err != nil {
	// 		return "", 0, 0, errors.New("Could not find the day")
	// 	}
	// 	return val, index, length, nil
	// }

	// var findDayNumber = func(t string) (string, error) {
	// 	return "", nil
	// }

	// var findMonthEnglishName = func(t string) (string, int, int, error) {
	// 	val, index, length, err := find(t, monthes)
	// 	if err != nil {
	// 		return "", 0, 0, errors.New("Could not find the month")
	// 	}
	// 	return val, index, length, nil
	// }

	// valDay, valIndex, valLength, errDay := findDayEnglishName(timeString)
	// if errDay == nil {
	// 	dayName = valDay
	// }

	// valMonth, _, _, errMonth := findMonthEnglishName(timeString)
	// if errMonth == nil {
	// }

	// if errDay != nil {
	// 	return time.Time{}, errDay
	// }
	// if errMonth != nil {
	// 	return time.Time{}, errMonth
	// }
	return time.Time{}, nil
}

func parseRssTime(str string) time.Time {
	var feedTime, err = time.Parse(time.RFC1123, str)
	if err != nil {
		// fmt.Println(err)
		feedTime = time.Time{}
	}
	return feedTime
}

func parseAtomTime(str string) time.Time {
	var feedTime, err = time.Parse(time.RFC3339, str)
	if err != nil {
		// fmt.Println(err)
		feedTime = time.Time{}
	}
	return feedTime
}

func timeToString(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
