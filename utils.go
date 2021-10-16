package main

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func UnmarshalSubscanBodyByBatch(src []byte) ([]*EventCollect, error) {
	subscanBody := SubscanRespBody{}
	if err := json.Unmarshal(src, &subscanBody); err != nil {
		log.Printf("unmarshal subscan response body failed: %v\n", err.Error())
		return nil, err
	}
	log.Printf("==============\n")
	log.Printf("%+v\n", subscanBody)

	var mutex sync.RWMutex
	eventsCollect := &EventsCollect{
		make([]*EventCollect, 0),
		&mutex,
	}

	// TODO: goroutines init
	maxProcess := runtime.NumCPU()
	log.Printf("maxProcess number: %d\n", maxProcess)
	runtime.GOMAXPROCS(maxProcess)

	for i, v := range subscanBody.Data.Events {
		wg.Add(1)
		go func(index int, value *SubscanEvent) {
			defer wg.Done()
			extrinsicSinger, err := getAssociateAddrByHash(value.ExtrinsicHash)
			if err != nil {
				log.Printf("index: %d, extrinsic_idx: %d, getAssociateAddrByHash error: %v\n", index, value.ExtrinsicIdx, err.Error())
			}

			eventCollect := &EventCollect{
				ExtrinsicHash:     value.ExtrinsicHash,
				Section:           value.ModuleID,
				Metion:            value.EventID,
				AssociatedAddress: extrinsicSinger,
			}
			// TODO: add mutex
			eventsCollect.Mutex.Lock()
			eventsCollect.Events = append(eventsCollect.Events, eventCollect)
			eventsCollect.Mutex.Unlock()
		}(i, v)
	}

	// TODO: blocking wait
	wg.Wait()
	return eventsCollect.Events, nil
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
	return "", nil
}
