package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xid"
)

type Branch struct {
	Id      string
	Service string
}

type Transaction struct {
	Id      string
	Branchs map[string]*Branch
}

var transactions map[string]*Transaction = make(map[string]*Transaction)

func setLogPrefix(prefix string) {
	log.SetPrefix(fmt.Sprintf("%s ", prefix))
}

func main() {
	// create a new global transaction
	http.HandleFunc("/new", func(rw http.ResponseWriter, r *http.Request) {
		tid := xid.New().String()
		setLogPrefix(tid)
		log.Print("new transaction")
		t := &Transaction{
			Id:      tid,
			Branchs: make(map[string]*Branch),
		}
		transactions[tid] = t
		rw.Write([]byte(tid))
	})

	// register a new transaction branch
	http.HandleFunc("/branch", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		service := r.URL.Query().Get("service")
		branchId := xid.New().String()
		setLogPrefix(tid)
		log.Printf("register branch= %s service = %s", branchId, service)
		t := transactions[tid]
		if t != nil {
			// TODO is need save request body
			t.Branchs[branchId] = &Branch{
				Id:      branchId,
				Service: service,
			}
		}
	})

	http.HandleFunc("/confirm", func(rw http.ResponseWriter, r *http.Request) {
		confirmOrCancel(r, "confirm")
	})

	http.HandleFunc("/cancel", func(rw http.ResponseWriter, r *http.Request) {
		confirmOrCancel(r, "cancel")
	})

	http.HandleFunc("/transactions", func(rw http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(transactions)
		rw.Write(bytes)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func confirmOrCancel(r *http.Request, status string) error {
	tid := r.URL.Query().Get("tid")
	setLogPrefix(tid)
	log.Printf("transaction %s", status)
	if transaction, ok := transactions[tid]; ok {

		for _, branch := range transaction.Branchs {
			log.Printf("call branch %s", branch.Service)
		}
		return nil
	} else {
		return errors.New("transaction not exist")
	}
}
