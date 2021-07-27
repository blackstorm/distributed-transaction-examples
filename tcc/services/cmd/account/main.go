package main

import (
	"log"
	"net/http"

	"github.com/blackstorm/dt/tcc/services/pkg/database"
	"github.com/blackstorm/dt/tcc/services/pkg/tcc"
)

const dB_NAME = "account"

func main() {
	log.SetPrefix("account service:")

	db, err := database.Connect2DB(dB_NAME)
	if err != nil {
		panic(err)
	}

	// do try
	http.HandleFunc("/try", func(rw http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				// SET response
			}
		}()

		tid := r.URL.Query().Get("tid")
		tcc := tcc.NewTCCFromTid(tid, db)
		if is, err := tcc.IsRepeatTry(); err != nil {
			panic(err)
		} else if is {
			log.Printf("%s is repeat try skip!", tid)
			return
		}

		//ensure return susccess
	})

	http.HandleFunc("/confirm", func(rw http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/cancel", func(rw http.ResponseWriter, r *http.Request) {

	})

	log.Fatal(http.ListenAndServe(":5000", nil))

}
