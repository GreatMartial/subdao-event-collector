package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

const (
	// targetUrl = "https://polkadot.api.subscan.io"
	targetUrl = "https://kusama.api.subscan.io"

	targetModule = "treasury"
	targetCall = "proposed"

)

var (
	// append event collect of marshal from subScan.
	c chan interface{}
	// signal chan bool
)

func main() {
	// TODO: http request
	eventsURL := fmt.Sprintf(targetUrl + "/api/scan/events")
	log.Printf("%s", eventsURL)
	method := "POST"
	str := fmt.Sprintf(`{
		"row": 1,
		"page": 0,
		"module": "%s",
		"call": "%s"
	}`, targetModule, targetCall)

/* 	str := fmt.Sprintf(`{
		"row": 1,
		"page": 0,
		"module": "%s"
	}`, targetModule)	 */
	payload := strings.NewReader(str)
	

	// get target total count of events.
	respBody, err := HttpRetry(eventsURL, method, payload, InvokeHttpReq)
	// respBody, err := InvokeHttpReq(eventsURL, method, payload)
	if err != nil {
		log.Printf("invoke http req error: %v\n", err.Error())
		return
	}

	// get the event count of resp body
	tmp := SubscanEventsRespBody{}
	if err = json.Unmarshal(respBody, &tmp); err != nil {
		log.Printf("unmarshal resultBody failed: %v\n", err)
		return
	}

	log.Printf("total: %+v\n", tmp)
	// TODO: check target count, use 100 to pagging.
	targetCount := tmp.Data.Count
	row := 100
	var totalPageNum int

	// compute total page number
	if targetCount <= row {
		totalPageNum = 1
	} else {
		totalPageNum = targetCount/row + 1
	}


	c = make(chan interface{}, 100)
	csvFile, err := os.OpenFile(targetModule+"-"+"test.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("open file failed: %v\n", err.Error())
		return
	}
	defer csvFile.Close()

	/* 	go func() {
		if err := gocsv.MarshalChan(c, gocsv.DefaultCSVWriter(csvFile)); err != nil {
			panic(fmt.Sprintf("gocsv marshal chan error: %v", err.Error()))
		}
		log.Printf("=======================================")
		log.Printf("Groutine exit: gocsv marshal chan")
		}()
		*/
		
	go func() {
		// i is page number, default value is 1.
		for p := 0; p < totalPageNum; p++ {	
			// TODO: create batch request.
			str := fmt.Sprintf(`{
				"row": %d,
				"page": %d,
				"module": "%s",
				"call": "%s"
			}`, row, p, targetModule, targetCall)
			// str := fmt.Sprintf(`{
			// 	"row": %d,
			// 	"page": %d,
			// 	"module": "%s"
			// }`, row, p, targetModule)	
			// log.Printf(str)
			payload := strings.NewReader(str)
			log.Printf("str: %s", payload)
			log.Printf("batch request----row: %d, page: %d, module: %s\n", row, p, targetModule)
		
			// wait seconds
			time.Sleep(time.Second)
			respBody, err := HttpRetry(eventsURL, method, payload, InvokeHttpReq)
			// respBody, err := InvokeHttpReq(eventsURL, method, payload)
			if err != nil {
				log.Printf("batch invoke http request error: %v", err.Error())
				return
			}
		
			// TODO: unmarshal resp
			// TODO: unmarshal subscanBody to target struct
			if err := UnmarshalSubscanEventsBodyByBatch(respBody); err != nil {
				panic(fmt.Sprintf("unmarshal subscan body by batch error: %v\n", err.Error()))
			}
			// log.Printf("events collect: %v\n", res)
		}
		close(c)
	}()

	// async write file
/* 
	go func() {
		if err := gocsv.MarshalChan(c, gocsv.DefaultCSVWriter(csvFile)); err != nil {
			panic(fmt.Sprintf("gocsv marshal chan error: %v", err.Error()))
		}
		signal <- true
	}()	

	<- signal 
*/
	if err := gocsv.MarshalChan(c, gocsv.DefaultCSVWriter(csvFile)); err != nil {
		panic(fmt.Sprintf("gocsv marshal chan error: %v", err.Error()))
	}
	fmt.Printf("work success!\n")
}
