package scp

import "fmt"

type status byte

var statuses = map[string]byte{
	"New":       0,
	"Started":   1,
	"Completed": 2,
	"Errored":   3,
	"Cancelled": 4,
}

func getStatus(s string) (*status, error) {
	if stat, ok := statuses[s]; ok {
		output := status(stat)
		return &output, nil
	}
	output := status(4)
	return &output, fmt.Errorf("unable to lookup status byte")
}

func getStatusString(s *status) (string, error) {

	for k, v := range statuses {
		val := status(v)
		if &val == s {
			return k, nil
		}
	}

	return "", fmt.Errorf("unable to lookup status string")
}

//ProcessingError type used to give feedback on unsucessful upload/download
type ProcessingError struct {
	JobID       string
	Action      string
	Source      string
	Destination string
	Err         error
}
