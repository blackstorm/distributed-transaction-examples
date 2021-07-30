package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/blackstorm/dt/tcc/services/pkg/database"
	"github.com/blackstorm/dt/tcc/services/pkg/tcc"
	"github.com/rs/xid"
)

type OrderRequest struct {
	Id string
}

type AccountRequest struct {
	Amount uint
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("api-order: ")

	db, err := database.Connect2DB("order")
	if err != nil {
		panic(err)
	}

	// The order api. create a new order and enable a gloabl transaction!
	http.HandleFunc("/order", func(rw http.ResponseWriter, r *http.Request) {
		// create a new transaction
		tx, err := tcc.NewTransaction()
		if err != nil {
			panic(err)
		}

		// create a new order
		orderId := xid.New().String()
		if _, err := db.Exec(fmt.Sprintf("insert into `order` (id, status) values('%s', 'CREATED')", orderId)); err != nil {
			panic(err)
		}

		// do business
		err = tx.DoBusiness(func() error {
			if err = tx.CallBranch("http://order:4000", &OrderRequest{Id: orderId}); err != nil {
				return err
			}
			if err = tx.CallBranch("http://account:5000", &AccountRequest{Amount: 10}); err != nil {
				return err
			}
			return nil
		})

		// catch error for response
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/tcc/try", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")

		// get request body
		decoder := json.NewDecoder(r.Body)
		var orderREquest OrderRequest
		decoder.Decode(&orderREquest)

		log.Printf("tcc %s try order id = %s", tid, orderREquest.Id)

		logAndResponseStatusError := func(message string, err error) {
			log.Printf(message, err)
			rw.WriteHeader(http.StatusInternalServerError)
		}

		// 幂等问题
		tccChecker := tcc.NewTCCChecker(tid, db)
		if isTried, err := tccChecker.IsTried(); err != nil {
			logAndResponseStatusError("check is tried error %v", err)
			return
		} else {
			if isTried {
				return
			}
		}

		// 悬挂问题
		if isCancel, err := tccChecker.IsCancel(); err != nil {
			logAndResponseStatusError("check is cancel error %v", err)
			return
		} else {
			if isCancel {
				return
			}
		}

		if isConfirm, err := tccChecker.IsConfirm(); err != nil {
			logAndResponseStatusError("check is confirm error %v", err)
			return
		} else {
			if isConfirm {
				return
			}
		}

		// TODO 事务支持
		if _, err := db.Exec(fmt.Sprintf("update `order` set status = 'UPDATING' where id = '%s'", orderREquest.Id)); err != nil {
			logAndResponseStatusError("order update status error %v", err)
			return
		}

		// 事务落库
		if err := tcc.NewTCCLoger(tid, db).LogTry(); err != nil {
			log.Printf("log try error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/tcc/confirm", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		log.Printf("transaction %s confirm", tid)

		var orderRequest OrderRequest
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&orderRequest)

		// 幂等
		checker := tcc.NewTCCChecker(tid, db)
		if isConfirm, err := checker.IsConfirm(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if isConfirm {
				return
			}
		}

		// 支持空回滚
		if _, err := db.Exec(fmt.Sprintf("update `order` set status = 'FINISHED' where id = '%s'", orderRequest.Id)); err != nil {
			log.Printf("order update status error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		tccLoger := tcc.NewTCCLoger(tid, db)
		if err := tccLoger.LogConfirm(); err != nil {
			log.Printf("tcc log confirm error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/tcc/cancel", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")
		log.Printf("transaction %s cancel", tid)

		var orderRequest OrderRequest
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&orderRequest)

		// 幂等
		checker := tcc.NewTCCChecker(tid, db)
		if isCanceled, err := checker.IsCancel(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if isCanceled {
				return
			}
		}

		// 支持空回滚
		if _, err := db.Exec(fmt.Sprintf("update `order` set status = 'CANCELED' where id = '%s'", orderRequest.Id)); err != nil {
			log.Printf("order update status error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		tccLoger := tcc.NewTCCLoger(tid, db)
		if err := tccLoger.LogCancel(); err != nil {
			log.Printf("tcc log cancel error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":4000", nil))

}
