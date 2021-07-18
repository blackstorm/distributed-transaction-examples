package callback

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/blackstorm/dt/xa/services/pkg/xa"
)

func CallbackHandleFunc(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		xid := r.URL.Query().Get("xid")
		tid := r.URL.Query().Get("tid")

		log.Printf("tm callback status=%s tid=%s xid=%s", status, tid, xid)

		var err error
		if status == "commit" {
			err = xa.Commit(db, xid)
		} else {
			err = xa.Rollback(db, xid)
		}

		if err != nil {
			log.Printf("callback error xid=%s", xid)
		}

	}
}
