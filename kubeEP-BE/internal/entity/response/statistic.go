package response

import "time"

type NodePoolStatus struct {
	CreatedAt time.Time `json:"created_at"`
	Count     int32     `json:"count"`
}

type HPAStatus struct {
	CreatedAt time.Time `json:"created_at"`
	Count     int32     `json:"count"`
}
