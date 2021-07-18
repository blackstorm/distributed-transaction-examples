package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.SetPrefix("api service: ")

	http.HandleFunc("/order", func(rw http.ResponseWriter, r *http.Request) {
		// open a new transaction
		tid, _ := newTransaction()
		log.Printf("new transaction id %s", tid)

		// call customer service to reduce moneny
		resp1, err := http.Get("http://customer:4000/reduce?tid=" + tid)
		if err != nil || resp1.StatusCode != 200 {
			// ensure rpc service is support empty (not exist) rollback
			rollback(tid)
			rw.Write([]byte("failed"))
			return
		}

		// call merchant service to add moneny
		resp2, err := http.Get("http://merchant:5000/add?tid=" + tid)
		if err != nil || resp2.StatusCode != 200 {
			// ditto
			rollback(tid)
			rw.Write([]byte("failed"))
			return
		}

		// finished the transaction
		done(tid)
		rw.Write([]byte("ok"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newTransaction() (string, error) {
	resp, err := http.Get("http://tm:9999/new")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// transaction id
	tid := readHttpResponse2String(resp)
	return tid, nil
}

func rollback(tid string) {
	// TODO catch error
	log.Printf("do rollback tid=%s", tid)
	http.Get("http://tm:9999/rollback?tid=" + tid)
}

func done(tid string) {
	// TODO catch error
	http.Get("http://tm:9999/done?tid=" + tid)
}

func readHttpResponse2String(resp *http.Response) string {
	// TODO catch error
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
