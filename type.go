package main

import "sync"

type EventsCollect struct {
	Events []*EventCollect
	Mutex  *sync.RWMutex
}

type EventCollect struct {
	EventIndex        string `json:"event_index" csv:"event_index" binding:"event_index"`
	BlockNum          int    `json:"block_num" csv:"block_num" binding:"block_num"`
	ExtrinsicHash     string `json:"extrinsic_hash" csv:"extrinsic_hash" binding:"require"`
	ExtrinsicIdx      int    `json:"extrinsic_idx" csv:"extrinsic_idx"`
	Section           string `json:"section" csv:"section"`
	Metion            string `json:"metion" csv:"metion"`
	AssociatedAddress string `json:"associated_address" csv:"Associated_address"`
}

type SubscanEventsRespBody struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	GeneratedAt int    `json:"generated_at"`
	Data        struct {
		Count  int             `json:"count"`
		Events []*SubscanEvent `json:"events"`
	} `json:"data"`
}

type SubscanExtrinscRespBody struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	GeneratedAt int    `json:"generated_at"`
	Data        struct {
		AccountId string `json:"account_id"`
	} `json:"data"`
}

type SubscanEvent struct {
	EventIndex     string `json:"event_index"`
	BlockNum       int    `json:"block_num"`
	ExtrinsicIdx   int    `json:"extrinsic_idx"`
	ModuleID       string `json:"module_id"`
	EventID        string `json:"event_id"`
	Params         string `json:"params"`
	EventIdx       int    `json:"event_idx"`
	ExtrinsicHash  string `json:"extrinsic_hash"`
	Finalized      bool   `json:"finalized"`
	BlockTimestamp int    `json:"block_timestamp"`
}
