package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func UnmarshalSubscanEventsBodyByBatch(src []byte) error {
	subscanBody := SubscanEventsRespBody{}
	if err := json.Unmarshal(src, &subscanBody); err != nil {
		log.Printf("unmarshal subscan response body failed: %v\n", err.Error())
		return err
	}
	log.Printf("==============\n")
	log.Printf("%+v\n", subscanBody)

	/*
		var mutex sync.RWMutex
		eventsCollect := &EventsCollect{
			make([]*EventCollect, 0),
			&mutex,
		}
	*/

	// TODO: goroutines init
	maxProcess := runtime.NumCPU()
	log.Printf("maxProcess number: %d\n", maxProcess)
	runtime.GOMAXPROCS(maxProcess)

	for i, v := range subscanBody.Data.Events {
		wg.Add(1)

		// wait seconds
		time.Sleep(time.Second)
		go func(index int, value *SubscanEvent) {
			defer wg.Done()

			var extrinsicSinger = ""
			var err error

			if value.ExtrinsicHash != "" {
				extrinsicSinger, err = getAssociateAddrByHash(value.ExtrinsicHash)
				if err != nil {
					log.Printf("index: %d, extrinsic_idx: %d, getAssociateAddrByHash error: %v\n", index, value.ExtrinsicIdx, err.Error())
				}
			}
			log.Printf("index: %d, extrinsicSinger: %s\n", index, extrinsicSinger)
			eventCollect := &EventCollect{
				EventIndex:        value.EventIndex,
				BlockNum:          value.BlockNum,
				ExtrinsicHash:     value.ExtrinsicHash,
				ExtrinsicIdx:      value.ExtrinsicIdx,
				Section:           value.ModuleID,
				Metion:            value.EventID,
				AssociatedAddress: extrinsicSinger,
			}
			// TODO: add mutex
			/* 	eventsCollect.Mutex.Lock()
			eventsCollect.Events = append(eventsCollect.Events, eventCollect)
			eventsCollect.Mutex.Unlock()
			*/
			c <- eventCollect
		}(i, v)
	}

	// TODO: blocking wait
	wg.Wait()
	return nil
}

/* func unmarshalSubscanEvent(srcCollect []*EventCollect, extrHash string) (*EventCollect, error) {
	extrinsicSinger, err := getAssociateAddrByHash(extrHash)
	if err != nil {
		log.Printf("index: %d, extrinsic_idx: %d, getAssociateAddrByHash error: %v\n", i, v.ExtrinsicIdx, err.Error())
		continue
	}

	eventCollect := &EventCollect{
		ExtrinsicHash:     v.ExtrinsicHash,
		Section:           v.ModuleID,
		Metion:            v.EventID,
		AssociatedAddress: extrinsicSinger,
	}
	eventsCollect = append(eventsCollect, eventCollect)

	return eventsCollect, nil
}
*/

// TODO:
// obtain the associated account address by passing the extrinsic hash.
func getAssociateAddrByHash(extrHash string) (string, error) {
	// create get-extrinsic request
	extrinsicURL := targetUrl + "/api/scan/extrinsic"
	hash := fmt.Sprintf(`{
		"hash": "%s"
	}`, extrHash)
	method := "POST"
	payload := strings.NewReader(hash)

	respBody, err := HttpRetry(extrinsicURL, method, payload, InvokeHttpReq)
	// respBody, err := InvokeHttpReq(extrinsicURL, method, payload)
	if err != nil {
		return "", err
	}

	log.Printf("associate_addr body: %v", string(respBody))
	temp := SubscanExtrinscRespBody{}
	if err := json.Unmarshal(respBody, &temp); err != nil {
		return "", err
	}

	result := temp.Data.AccountId
	return result, nil
}
