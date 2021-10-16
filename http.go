package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func InvokeHttpReq(url string, method string, payload io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, targetUrl, payload)
	if err != nil {
		// fmt.Printf("create http req failed: %v\n", err.Error())
		return nil, err
	}
	// log.Printf("req value: %+v\n", req)

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
	return resp, nil
}
