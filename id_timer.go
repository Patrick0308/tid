package tid

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type TIDTimer interface {
	Second() (uint32, error)
}

func NewDefaultTimer(address string) TIDTimer {
	client := &http.Client{}
	return &defaultTimer{client, address}
}

type defaultTimer struct {
	client *http.Client
	address string
}

func (dt defaultTimer) Second() (uint32, error) {
	var resp *http.Response
	var err error
	var body []byte
	var i64s uint64
	if resp, err = dt.client.Get(dt.address + "/getSecond"); err != nil {
		return 0, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("close client error：%v", err)
		}
	}()
	if body, err = ioutil.ReadAll(resp.Body);err != nil {
		return 0, err
	}
	if i64s, err = strconv.ParseUint(string(body), 10, 32) ;err != nil {
		return 0, err
	}
	return uint32(i64s), nil
}
