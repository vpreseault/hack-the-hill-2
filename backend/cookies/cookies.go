package cookies

import (
	"net/http"
)

func GetUserNameFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

