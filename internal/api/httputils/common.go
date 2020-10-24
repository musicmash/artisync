package httputils

import "net/http"

func GetUserName(r *http.Request) string {
	return r.Header.Get("X-User-Name")
}
