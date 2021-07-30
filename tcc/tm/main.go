package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rs/xid"
)

type Branch struct {
	Id      string
	Service string
	ReqBody []byte
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

		// todo should read form body
		reqBody, _ := ioutil.ReadAll(r.Body)

		setLogPrefix(tid)
		log.Printf("register branch= %s service = %s", branchId, service)
		t := transactions[tid]
		if t != nil {
			// TODO is need save request body
			t.Branchs[branchId] = &Branch{
				Id:      branchId,
				Service: service,
				ReqBody: reqBody,
			}
		}
	})

	http.HandleFunc("/confirm", func(rw http.ResponseWriter, r *http.Request) {
		getAndRangeTransactions(r, handleStep2(r, "confirm"))
	})

	http.HandleFunc("/cancel", func(rw http.ResponseWriter, r *http.Request) {
		getAndRangeTransactions(r, handleStep2(r, "cancel"))
	})

	http.HandleFunc("/transactions", func(rw http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(transactions)
		rw.Write(bytes)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}

func handleStep2(r *http.Request, action string) func(*Transaction, *Branch) error {
	return func(t *Transaction, b *Branch) error {
		url := fmt.Sprintf("%s/tcc/%s?tid=%s", b.Service, action, t.Id)
		// TODO whit retry
		http.Post(url, "application/json", bytes.NewReader(b.ReqBody))
		return nil
	}
}

func getAndRangeTransactions(r *http.Request, fn func(*Transaction, *Branch) error) error {
	tid := r.URL.Query().Get("tid")
	setLogPrefix(tid)
	if transaction, ok := transactions[tid]; ok {
		for _, branch := range transaction.Branchs {
			err := fn(transaction, branch)
			if err != nil {
				return err
			}
		}
	}
	return errors.New("transaction not found")
}
