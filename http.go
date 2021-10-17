package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type hanldeFunc func(string, string, io.Reader) ([]byte, error)


func InvokeHttpReq(url string, method string, payload io.Reader) ([]byte, error) {
	// setting timeout
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second * 10)
				if err != nil {
					return nil, err
				}

				_ = conn.SetDeadline(time.Now().Add(time.Second * 10))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		// fmt.Printf("create http req failed: %v\n", err.Error())
		return nil, err
	}
	log.Printf("req value: %+v\n", req)

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		// err = errors.New(fmt.Sprintf("invoke req failed: %v\n", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf(" request api interface error: %v\n", resp.StatusCode))
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// retry http request
func HttpRetry(url, method string, payload io.Reader, callBack hanldeFunc) ([]byte, error) {
	// setting retry param
	errStr := ""
	attempts := 3
	sleepTime := 3

	for i := 0; i < attempts; i++ {
		resp, err := callBack(url, method, payload)
		if err == nil && resp != nil {
			return resp, nil
		}
		if i != 0 {
			errStr += "|" + err.Error()
		} else {
			errStr  += err.Error()
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	return nil, errors.New("retry http request error: " + errStr)
}
