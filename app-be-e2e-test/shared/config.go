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

const JwtToken = "jwt=eyJhbGciOiJIUzUxMiIsImtpZCI6IjhiYmUwN2RlLTE4NjctNDJmMS04YjJiLWMwODMzNTIyYzM3NS0yYmE2ZTlhNy1kMTBjLTRiY2YtYjhhMi01Y2RhMjhhZmRjMDAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiZGlza29uLWh1bnRlci1mcm9udGVuZCJdLCJkaGNjIjp7InVzZXJfaWQiOiI4YmE0MDljNy05ZmMwLTQxMDgtYmRmMC1iZDYxNDVjNzVkNzQifSwiZXhwIjoxNzMzNzM2MTIwLCJpYXQiOjE3MDIyMDAxMjAsImlzcyI6ImRpc2tvbi1odW50ZXIvcHJpY2UtbW9uaXRvcmluZyJ9.z8g-gkVXo-UzPJSGVAFVqlseeYA5cabz5M0ZjhIelRswoFM05pri1ySGqM-zVuH1nCyAF07PC6-M_i3U_GD0iw;"
