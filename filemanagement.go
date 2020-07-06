package scp

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// REMOTE

func (c *client) runCommand(cmd string) ([]string, error) {

	output := new(strings.Builder)

	sess, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return nil, err
	}

	//go io.Copy(os.Stdout, sessStdOut)
	go io.Copy(output, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return nil, err
	}

	//go io.Copy(os.Stderr, sessStderr)
	go io.Copy(output, sessStderr)

	//log.Printf("Executing: %s", cmd)
	err = sess.Run(cmd)
	time.Sleep(1 * time.Second)
	if err != nil {
		return nil, err
	}

	outRaw := strings.Split(output.String(), "\n")
	out := []string{}

	for _, line := range outRaw {
		if len(line) > 0 {
			out = append(out, line)
		}
	}

	return out, nil
}

func (c *client) newRemoteDirectory(path string, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("mkdir -p \"%s\"", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) getRemoteSubDirectories(path string, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("find \"%s\" -type d", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) getRemoteFiles(path string, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("find \"%s\" -type f", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) getRemoteSubDirectoriesLimited(path string, depth int, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("find \"%s\" -maxdepth %d -type d", path, depth)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) getRemoteFilesLimited(path string, depth int, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("find \"%s\" -maxdepth %d -type f", path, depth)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) newRemoteFile(path string, sudo bool) ([]string, error) {
	splitStrings := strings.Split(path, "/")
	fileName := splitStrings[len(splitStrings)-1]
	dir := strings.TrimRight(path, "/")
	dir = strings.TrimRight(dir, fmt.Sprintf("/%s", fileName))

	s, err := c.newRemoteDirectory(dir, sudo)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("touch \"%s\"", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}

	s2, err2 := c.runCommand(cmd)
	if err2 != nil {
		return s, err2
	}

	out := concatResponses(s, s2)

	return out, nil
}

func (c *client) overwriteRemoteFile(path string, content []byte, sudo bool) ([]string, error) {
	s, err := c.removeRemoteFile(path, sudo)
	if err != nil {
		return nil, err
	}

	s2, err := c.newRemoteFile(path, sudo)
	if err != nil {
		return s, err
	}

	concat1 := concatResponses(s, s2)

	s3, err := c.setRemoteFileContent(path, content, sudo)
	if err != nil {
		return concat1, err
	}

	out := concatResponses(concat1, s3)

	return out, nil
}

func (c *client) removeRemoteFile(path string, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("rm \"%s\" -f", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) removeRemoteDirectory(path string, sudo bool) ([]string, error) {
	cmd := fmt.Sprintf("rm -r \"%s\" -f", path)
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	return c.runCommand(cmd)
}

func (c *client) setRemoteFileContent(path string, content []byte, sudo bool) ([]string, error) {
	return nil, nil
}

// LOCAL

func newLocalDirectory(path string) ([]string, error) {
	return nil, os.MkdirAll(path, 0777)
}

func getLocalSubdirectory(path string) ([]string, error) { return nil, nil }

func getLocalFiles(path string) ([]string, error) { return nil, nil }

func getLocalSubdirectoryLimited(path string, depth int) ([]string, error) { return nil, nil }

func getLocalFilesLimited(path string, depth int) ([]string, error) { return nil, nil }

func newLocalFile(path string) ([]string, error) { return nil, nil }

func overwriteLocalFile(path string) ([]string, error) { return nil, nil }

func removeLocalFile(path string) ([]string, error) { return nil, nil }

func removeLocalDirectory(path string) ([]string, error) { return nil, nil }

func setLocalFileContent(path string, content []byte) ([]string, error) { return nil, nil }

// UTIL

func concatResponses(r1, r2 []string) []string {
	for _, v := range r2 {
		r1 = append(r1, v)
	}

	return r1
}
