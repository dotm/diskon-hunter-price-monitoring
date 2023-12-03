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

const JwtToken = "jwt=eyJhbGciOiJIUzUxMiIsImtpZCI6IjhiYmUwN2RlLTE4NjctNDJmMS04YjJiLWMwODMzNTIyYzM3NS0yYmE2ZTlhNy1kMTBjLTRiY2YtYjhhMi01Y2RhMjhhZmRjMDAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiZXRlbi1mcm9udGVuZCJdLCJldGNjIjp7InVzZXJfaWQiOiI5M2YwMzJkNi1lMGMxLTQxZWItYjY4Ny0wYzQwYmUzMDU5ZGMifSwiZXhwIjoxNzA2OTc2NzE0LCJpYXQiOjE2NzU0NDA3MTQsImlzcyI6ImV0ZW50ZWNoL2ludmVudG9yeS1tYW5hZ2VtZW50LWJhY2tlbmQifQ.ez5TNmA7dqSrCwZ7MjVDvbqrl12P2qmqbUpGy9g-zK_iT4f_zq61jpEEJG-YaV8WhRGiHFQ1xBFAncFtnaYPCA;"
