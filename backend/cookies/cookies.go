package cookies

import (
	"net/http"
)

func GetHostNameFromUserCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

