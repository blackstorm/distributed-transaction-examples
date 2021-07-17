package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blackstorm/dt/xa/services/pkg/callback"
	"github.com/blackstorm/dt/xa/services/pkg/database"
	"github.com/blackstorm/dt/xa/services/pkg/tm"
	"github.com/blackstorm/dt/xa/services/pkg/xa"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

func main() {
	log.SetPrefix("merchant service: ")

	db := database.ConnectDB()

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		xid := xid.New().String()

		log.Printf("add tid=%s xid=%s", tid, xid)

		// xa start
		if _, err := db.Exec(fmt.Sprintf("XA START '%s'", xid)); err != nil {
			log.Println(errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := db.Exec("update merchant_account set balance = balance + 10 where id = 1"); err != nil {
			log.Println(errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := db.Exec(fmt.Sprintf("XA END '%s'", xid)); err != nil {
			log.Println(errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := tm.Register(tid, xid, "http://merchant:5000/callback"); err != nil {
			log.Println(errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		xa.Prepare(db, xid)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/callback", callback.CallbackHandleFunc(db))

	log.Fatal(http.ListenAndServe(":5000", nil))
}
