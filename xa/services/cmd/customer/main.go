package main

import (
	"log"
	"net/http"

	"github.com/blackstorm/dt/xa/services/pkg/callback"
	"github.com/blackstorm/dt/xa/services/pkg/database"
	"github.com/blackstorm/dt/xa/services/pkg/tm"
	"github.com/blackstorm/dt/xa/services/pkg/xa"
	"github.com/pkg/errors"
)

func main() {
	log.SetPrefix("customer service: ")
	db := database.ConnectDB()

	db.Ping()

	http.HandleFunc("/reduce", func(w http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		XA := xa.NewXA(db)

		log.Printf("reduce tid=%s xid=%s", tid, XA.XID)

		defer func() {
			if err := recover(); err != nil {
				log.Printf("%+v\n", err)
			}
		}()

		// ensure the customer account balance > 0
		var balance int
		if err := db.QueryRow("select balance from customer_account where id = 1 limit 1").Scan(&balance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(errors.WithStack(err))
		}
		if balance <= 0 {
			log.Printf("balance is insufficient.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// xa start
		if err := XA.Start(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(errors.WithStack(err))
		}

		if _, err := db.Exec("update customer_account set balance = balance - 10 where id = 1"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(errors.WithStack(err))
		}
		if err := XA.End(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(errors.WithStack(err))
		}

		if err := tm.Register(tid, XA.XID, "http://customer:4000/callback"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(errors.WithStack(err))
		}

		XA.Prepare()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/callback", callback.CallbackHandleFunc(db))

	log.Fatal(http.ListenAndServe(":4000", nil))
}
