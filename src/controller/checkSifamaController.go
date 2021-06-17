package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sifamaGO/src/selenium"
)

func checkSifama(w http.ResponseWriter, r *http.Request) {
	var request Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	errorHandle(err)
	if request.StartDigitacao {
		user := request.User
		pass := request.Passd
		message, err := selenium.CheckSifama(user, pass)
		fmt.Println(message)
		fmt.Println(err)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(message))
		}

	}
}
