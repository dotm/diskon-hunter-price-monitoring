package shared

import (
	"diskon-hunter/price-monitoring-e2e-test/shared/envhelper"
)

func init() {
	envhelper.SetLocalEnvVar()
}

func GetBackendUrl() string {
	backend_url := envhelper.GetEnvVar("backend_url")
	if backend_url == "" {
		panic("empty backend_url")
	}

	return backend_url
}

const JwtToken = "jwt=eyJhbGciOiJIUzUxMiIsImtpZCI6IjhiYmUwN2RlLTE4NjctNDJmMS04YjJiLWMwODMzNTIyYzM3NS0yYmE2ZTlhNy1kMTBjLTRiY2YtYjhhMi01Y2RhMjhhZmRjMDAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiZGlza29uLWh1bnRlci1mcm9udGVuZCJdLCJkaGNjIjp7InVzZXJfaWQiOiJhYzU0ZWQzNi1jNzcwLTQ2MDQtYmY1NC04ODVhODg0YzI1YTcifSwiZXhwIjoxNzM2NzU1MjI5LCJpYXQiOjE3MDUyMTkyMjksImlzcyI6ImRpc2tvbi1odW50ZXIvcHJpY2UtbW9uaXRvcmluZyJ9.BCfGCIhSwCAlO73sr_AKAX-SBwnlWdy9Sy8_wdBvvrUSAkyN6WHrz_4XXXT5S6jNomymSpJFuJaNsel8a2Lvpw;"
