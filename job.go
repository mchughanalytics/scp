package scp

import (
	"crypto/rand"
	"fmt"
)

type job struct {
	uuid       string
	remotePath string
	localPath  string
	direction  *direction
	overwrite  bool
	status     *status
	message    string
}

func newJob(rp, lp string, d *direction, overwrite bool) (*job, error) {

	s, _ := getStatus("New")
	myUUID, _ := newUUID()

	j := &job{
		uuid:       myUUID,
		remotePath: rp,
		localPath:  lp,
		direction:  d,
		overwrite:  overwrite,
		status:     s,
	}

	return j, nil
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}

func (j *job) GetID() string {
	return j.uuid
}

func (j *job) GetStatus() string {
	s, _ := getStatusString(j.status)
	return s
}

func (j *job) GetMessage() string {
	return j.message
}

type direction string

var directions = map[string]string{
	"upload":   "UP",
	"send":     "UP",
	"set":      "UP",
	"put":      "UP",
	"push":     "UP",
	"up":       "UP",
	"u":        "UP",
	"download": "DOWN",
	"get":      "DOWN",
	"fetch":    "DOWN",
	"receive":  "DOWN",
	"pull":     "DOWN",
	"down":     "DOWN",
	"d":        "DOWN",
}

func newDirection(s string) (*direction, error) {
	if _, ok := directions[s]; !ok {
		d := direction("")
		return &d, fmt.Errorf("invalid action")
	}

	d := direction(directions[s])

	return &d, nil
}
