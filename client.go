package scp

import (
	"golang.org/x/crypto/ssh"
)

type client struct {
	remoteHost   string
	client       *ssh.Client
	clientConfig *ssh.ClientConfig

	jobs       chan *job
	jobCounter int
	jobLedger  map[string]*job
}

//NewClient returns new SCP client
func NewClient(host string, config *ssh.ClientConfig) (*client, error) {

	l := make(map[string]*job)

	c := &client{
		remoteHost:   host,
		clientConfig: config,
		jobs:         make(chan *job, 50),
		jobCounter:   0,
		jobLedger:    l,
	}

	err := c.Dial()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) Dial() error {
	//log.Print("Dialing...")
	conn, err := ssh.Dial("tcp", c.remoteHost, c.clientConfig)
	if err != nil {
		return err
	}
	c.client = conn
	return nil
}

//NewClientFromCredentials returns new SCP client
func NewClientFromCredentials(host, username, password string) (*client, error) {

	a := []ssh.AuthMethod{}
	a = append(a, ssh.Password(password))

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            a,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}

	return NewClient(host, config)
}

//newUpload creates a new upload job
func (c *client) newUploadJob(src, dest string, overwrite bool) (string, error) {

	d, _ := newDirection("UP")

	j, err := newJob(dest, src, d, overwrite)
	if err != nil {
		return "", err
	}

	c.jobs <- j
	c.jobCounter++
	c.jobLedger[j.uuid] = j

	return j.uuid, nil
}

//newDownload creates a new download job
func (c *client) newDownloadJob(src, dest string, overwrite bool) (string, error) {

	d, _ := newDirection("DOWN")

	j, err := newJob(src, dest, d, overwrite)
	if err != nil {
		return "", err
	}

	c.jobs <- j
	c.jobCounter++
	c.jobLedger[j.uuid] = j

	return j.uuid, nil
}

func (c *client) upload(j *job) *job {

	return j
}

func (c *client) download(j *job) *job {

	return j
}

func (c *client) Close() {
	c.client.Close()
}
