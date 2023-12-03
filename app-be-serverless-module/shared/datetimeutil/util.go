package datetimeutil

import "time"

// We remove time to avoid cluttering the table with many sort key.
//
// You can see the available value for timeZoneName in:
//
//	https://en.wikipedia.org/wiki/List_of_tz_database_time_zones (search TZ database name)
//
// for testing purpose (checking date difference):
//
//	above 10:00 UTC+0, use Pacific/Kiritimati or Etc/GMT-14 (UTC+14) (sign is intentionally inverted)
//	below 10:00 UTC+0, use Etc/GMT+12 (UTC-12) (sign is intentionally inverted)
func GetLocalDateString(t time.Time, timeZoneName string) (string, error) {
	if timeZoneName == "" {
		timeZoneName = "Asia/Jakarta" //default to WIB
	}
	timeLocation, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return "", err
	}
	return t.In(timeLocation).Format("2006-01-02"), nil
}
