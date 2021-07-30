package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/blackstorm/dt/tcc/services/pkg/database"
	"github.com/blackstorm/dt/tcc/services/pkg/tcc"
)

type AccountRequest struct {
	Amount uint
}

func main() {
	log.SetPrefix("account :")

	db, err := database.Connect2DB("account")
	if err != nil {
		panic(err)
	}

	// do try
	http.HandleFunc("/try", func(rw http.ResponseWriter, r *http.Request) {
		tid := r.URL.Query().Get("tid")

		var accountRequest AccountRequest
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&accountRequest)

		// 幂等检查
		tccChecker := tcc.NewTCCChecker(tid, db)
		if isTried, err := tccChecker.IsTried(); err != nil {
			log.Printf("tcc check is tried error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if isTried {
				return
			}
		}

		// 悬挂问题
		if isCancel, err := tccChecker.IsCancel(); err != nil {
			log.Printf("tcc check is cancel error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if isCancel {
				// return or catch exception??
				return
			}
		}

		// ----------------------------- 事务支持 ----------------------------
		// 冻结和减少账户
		res, err := db.Exec(fmt.Sprintf("update account set balance = balance - %d where balance - %d > 0 and id = 1", accountRequest.Amount, accountRequest.Amount))
		if err != nil {
			log.Printf("tcc check is cancel error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if rows, _ := res.RowsAffected(); rows != 1 {
			log.Printf("balance insufficient")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// 记录到冻结库
		_, err = db.Exec(fmt.Sprintf("insert into account_trading (account_id, trading_balance) values(1, %d)", accountRequest.Amount))
		if err != nil {
			log.Printf("tcc check is cancel error %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 插入 try log
		tccLoger := tcc.NewTCCLoger(tid, db)
		err = tccLoger.LogTry()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/confirm", func(rw http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/cancel", func(rw http.ResponseWriter, r *http.Request) {

	})

	log.Fatal(http.ListenAndServe(":5000", nil))

}
