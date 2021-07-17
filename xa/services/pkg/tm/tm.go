package tm

import "net/http"

func Register(tid string, xid string, callback string) error {
	_, err := http.Get("http://tm:9999/register?tid=" + tid + "&xid=" + xid + "&callback=" + callback)
	return err
}
