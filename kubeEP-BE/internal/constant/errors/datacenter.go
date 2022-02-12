package errors

type DatacenterError string

const (
	DatacenterMismatch DatacenterError = "datacenter mismatch"
)
