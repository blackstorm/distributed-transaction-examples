package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type Resource struct {
	Xid      string
	Callback string
}

type Transaction struct {
	Id        string
	Resources []Resource
}

var transactions map[string]*Transaction = make(map[string]*Transaction)

func main() {
	log.SetPrefix("tm: ")

	// create a gloabl transaction id;
	http.HandleFunc("/new", func(rw http.ResponseWriter, r *http.Request) {
		guid := xid.New().String()
		log.Printf("new transaction id %s", guid)

		transactions[guid] = &Transaction{
			Id:        guid,
			Resources: make([]Resource, 0),
		}
		rw.Write([]byte(guid))
	})

	// rm register sub transaction
	http.HandleFunc("/register", func(rw http.ResponseWriter, r *http.Request) {
		xid := r.URL.Query().Get("xid")
		tid := r.URL.Query().Get("tid")
		callback := r.URL.Query().Get("callback")

		log.Printf("register resource tid=%s xid=%s callback=%s", tid, xid, callback)

		transactions[tid].Resources = append(transactions[tid].Resources, Resource{
			Xid:      xid,
			Callback: callback,
		})
	})

	http.HandleFunc("/done", func(rw http.ResponseWriter, r *http.Request) {
		commitOrRollback(r, "commit")
	})

	http.HandleFunc("/rollback", func(rw http.ResponseWriter, r *http.Request) {
		commitOrRollback(r, "rollback")
	})

	http.HandleFunc("/transactions", func(rw http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(transactions)
		rw.Write(b)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func commitOrRollback(r *http.Request, status string) {
	tid := r.URL.Query().Get("tid")
	t := transactions[tid]
	if t != nil {
		for _, r := range t.Resources {
			// r.Callback + "?status=ok&tid=" + tid + "&xid=" + r.Xid
			url := fmt.Sprintf("%s?status=%s&tid=%s&xid=%s", r.Callback, status, tid, r.Xid)
			log.Printf("callback url=%s", url)
			_, err := http.Get(url)
			if err != nil {
				log.Println(errors.WithStack(err))
			}
		}
		delete(transactions, tid)
	} else {
		log.Printf("do commit or rollback transaction tid=%s no exist", tid)
	}
}
