package tcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TCC struct {
	Id string
}

func NewTransaction() (*TCC, error) {
	resp, err := http.Get("http://tm:9999/new")
	if err != nil {
		return nil, err
	}
	tid := readHttpResponse2String(resp)

	return &TCC{
		Id: tid,
	}, nil
}

type BusinessFunc func() error

func (t *TCC) DoBusiness(fn BusinessFunc) error {
	cancel := func() {
		err := t.Cancel()
		if err != nil {
			log.Printf("tx cancel error %v", err)
		}
	}

	// catch runtime exception
	defer func() {
		if e := recover(); e != nil {
			cancel()
			panic(e)
		}
	}()

	// catch exception
	err := fn()
	if err != nil {
		cancel()
	}
	return err
}

func (t *TCC) CallBranch(service string, v interface{}) error {
	resp, err := httpPost("http://tm:9999/branch?tid="+t.Id+"&service="+service, v)
	if err == nil {
		branchId := readHttpResponse2String(resp)
		url := fmt.Sprintf("%s/tcc/try?tid=%s&branchId=%s", service, t.Id, branchId)
		// TODO with retry
		resp, err = httpPost(url, v)
		if resp.StatusCode != 200 {
			err = fmt.Errorf("try service %s error, response status code=%d", service, resp.StatusCode)
		}
		if err == nil {
			return nil
		}
	}

	return err
}

func (t *TCC) Cancel() error {
	// TODO retry
	http.Get("http://tm:9999/cancel?tid=" + t.Id)
	return nil
}

func httpPost(url string, v interface{}) (*http.Response, error) {
	value, _ := json.Marshal(v)
	return http.Post(url, "application/json", bytes.NewReader(value))
}

func readHttpResponse2String(resp *http.Response) string {
	// TODO catch error
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
