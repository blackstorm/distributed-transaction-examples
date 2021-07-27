package tcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (t *TCC) CallBranch(service string, v interface{}) error {

	resp, err := httpPost("http://tm:9999/branch?tid="+t.Id+"&service="+service, nil)
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

	t.Cancel()
	return err
}

/*
func (t *TCC) CallBranchs(v interface{}, services ...string) error {
	len := len(services)
	counter := make(chan string)
	errs := make(chan error)

	defer func() {
		close(counter)
		close(errs)
	}()

	ctx, cancel := context.WithCancel(context.Background())

	for _, srv := range services {
		go func(ctx context.Context, srv string, v interface{}) {
			err := t.CallBranch(ctx, srv, v)
			counter <- srv
			if err != nil {
				errs <- err
			}
		}(ctx, srv, v)
	}

	for {
		select {
		case e := <-errs:
			// log.Printf("%v", e)
			cancel()
			return e
		case <-counter:
			len--
			if len == 0 {
				ctx.Done()
				return nil
			}
		}
	}
}
*/

func (t *TCC) Cancel() {
	// TODO retry
	http.Get("http://tm:9999/cancel?tid=" + t.Id)
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
