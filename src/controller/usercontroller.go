package controller

import (
	"fmt"
	"net/http"
)

func test(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "email",
		Value: "email",
	}
	fmt.Fprintln(w, "user")
	http.SetCookie(w, &cookie)
}
