package scp

import "fmt"

func (c *client) UploadFile(sourcePath, destinationPath string, overwrite bool) (string, error) {
	return c.newUploadJob(sourcePath, destinationPath, overwrite)
}

func (c *client) DownloadFile(sourcePath, destinationPath string, overwrite bool) (string, error) {
	return c.newDownloadJob(sourcePath, destinationPath, overwrite)
}

func (c *client) UploadDirectory(sourcePath, destinationPath string, overwrite bool) *[]ProcessingError {

	return nil
}

func (c *client) DownloadDirectory(sourcePath, destinationPath string, overwrite bool) *[]ProcessingError {

	return nil
}

//NewJob creates a new job of type upload or download
func (c *client) NewJob(local, remote, action string, overwrite bool) (string, error) {

	d, err := newDirection(action)
	if err != nil {
		return "", fmt.Errorf("unable to create job: %s", err)
	}

	j, err := newJob(remote, local, d, overwrite)
	if err != nil {
		return "", fmt.Errorf("unable to create job: %s", err)
	}

	c.jobs <- j
	c.jobCounter++
	c.jobLedger[j.uuid] = j

	return j.uuid, nil
}

//Jobs returns the number of jobs currently awaiting processing
func (c *client) Jobs() int {
	return c.jobCounter
}

//Process the next job from the channel
func (c *client) ProcessNextJob() *ProcessingError {

	var err error
	bad, _ := getStatus("Errored")
	j := <-c.jobs
	c.jobCounter--

	pe := &ProcessingError{
		JobID:       j.uuid,
		Source:      "",
		Destination: "",
		Action:      "",
		Err:         err,
	}

	upload, _ := newDirection("UP")
	download, _ := newDirection("DOWN")

	switch j.direction {
	case upload:
		fmt.Println("Upload")
		pe.Action = "Upload"
		pe.Source = j.localPath
		pe.Destination = j.remotePath
		done := c.upload(j)
		if done.status == bad {
			err = fmt.Errorf("upload failed")
		}
	case download:
		fmt.Println("Download")
		pe.Action = "Download"
		pe.Source = j.remotePath
		pe.Destination = j.localPath
		done := c.download(j)
		if done.status == bad {
			err = fmt.Errorf("download failed")
		}
	default:
		c.jobCounter--
		err = fmt.Errorf("unable to process job: invalid direction")
	}

	if err != nil {
		pe.Err = err
		return pe
	}

	return nil
}

//Get a job from the ledger
func (c *client) GetLedgerJob(uuid string) *job {

	if j, ok := c.jobLedger[uuid]; ok {
		return j
	}

	return nil
}

func (c *client) ClearLedger() {
	c.jobLedger = make(map[string]*job)
}

func (c *client) CancelAllJobs() {
	for _, j := range c.jobLedger {
		n, _ := getStatus("New")
		s, _ := getStatus("Cancelled")

		if j.status == n {
			j.status = s
		}

	}
}
