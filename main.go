package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

var targetUrl = "https://polkadot.api.subscan.io/api/scan/events"

func main() {
	// TODO: http request
	method := "POST"
	payload := strings.NewReader(`{
		"row": 2,
		"page": 0,
		"module": "council"
	}`)

	// get target total count.
	resp, err := InvokeHttpReq(targetUrl, method, payload)
	if err != nil {
		log.Printf("invoke http req error: %v\n", err.Error())
		return
	}

	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("decode result body failed: %v\n", err.Error())
		return
	}

	// get the event count of resp body
	tmp := SubscanRespBody{}
	if err = json.Unmarshal(resultBody, &tmp); err != nil {
		log.Printf("unmarshal resultBody failed: %v\n", err)
		return
	}

	// TODO: check target count, use 100 to pagging.
	targetCount := tmp.Data.Count
	row := 100
	var totalPageNum int
	
	// compute total page number
	if targetCount <= row {
		totalPageNum = 1
	} else {
		totalPageNum = targetCount / row + 1
	}

	// i is page number, default value is 1.
	for p := 0; p < totalPageNum; p++ {
		// TODO: create batch request.
		str := fmt.Sprintf(`{
			"row": %d,
			"page": %d,
			"module": "council"
		}`, row, p)
		payload := strings.NewReader(str)	
		resp, err := InvokeHttpReq(targetUrl, method, payload)
		if err != nil {
			log.Printf("batch invoke http request error: %v", err.Error())
			return
		}

		// TODO: unmarshal resp
		// TODO: unmarshal subscanBody to target struct
		res, err := UnmarshalSubscanBodyByBatch(resultBody)
		if err != nil {
			fmt.Printf("unmarshal subscan body by batch error: %v\n", err.Error())
			return
		}
		log.Printf("events collect: %v\n", res)

		clients := res
		clientsFile, err := os.OpenFile("test.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("open file failed: %v\n", err.Error())
			return
		}
		defer clientsFile.Close()

		// TODO: use append model
		if err := gocsv.MarshalFile(&clients, clientsFile); err != nil {
			panic(err)
		}
	}
	fmt.Printf("work success!\n")
}
