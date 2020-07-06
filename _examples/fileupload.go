package main

import (
	"log"
	"time"

	"github.com/mchughanalytics/scp"
)

func main() {

	c, err := scp.NewClientFromCredentials("10.0.0.74:22", "[USER]", "[PASSWORD]")
	if err != nil {
		log.Print(err)
	}

	path := "/home/pi/minecraft/forge/1.12.2/test1/test2/test3.txt"

	response, err := c.NewRemoteFile(path, true)
	if err != nil {
		log.Print(err)
	}
	time.Sleep(1 * time.Second)
	for i, s := range response {
		log.Printf("line %d: %s", i, s)
	}
	time.Sleep(1 * time.Second)

	c.Close()

}
