package main

import (
	"log"
	"net/http"

	"github.com/blackstorm/dt/tcc/services/pkg/database"
	"github.com/blackstorm/dt/tcc/services/pkg/tcc"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("api-order: ")

	db, err := database.Connect2DB("order")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/order", func(rw http.ResponseWriter, r *http.Request) {
		tx, err := tcc.NewTransaction()
		if err != nil {
			panic(err)
		}

		/*
			defer func() {
				if err := recover(); err != nil {
					log.Printf("%v", err)
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}()
		*/

		if err = tx.CallBranch("http://order:4000", nil); err != nil {
			panic(err)
		}

		if err = tx.CallBranch("http://account:5000", nil); err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/tcc/try", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		log.Printf("tcc %s try ", tid)

		// 幂等问题
		tccChecker := tcc.NewTCCChecker(tid, db)
		if isTried, err := tccChecker.IsTried(); err != nil {
			log.Printf("check is tried error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			if isTried {
				return
			}
		}

		// TODO 悬挂问题

		// 流程处理，注意本地事务加持
		tcc.NewTCCLoger(tid, db).LogTry()
		// TODO update order status;
	})

	http.HandleFunc("/tcc/confirm", func(rw http.ResponseWriter, r *http.Request) {
		// TODO
	})

	http.HandleFunc("/tcc/cancel", func(rw http.ResponseWriter, r *http.Request) {
		// 处理空回滚
	})

	log.Fatal(http.ListenAndServe(":4000", nil))

}
