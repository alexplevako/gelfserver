package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var port = flag.Int("address", 12201, "port to listen on")

func main() {
	flag.Parse()

	reader, err := gelf.NewReader(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listen on %s", reader.Addr())

	go func() {
		for {
			message, err := reader.ReadMessage()
			if err != nil {
				log.Println(err)
			}

			var buffer bytes.Buffer
			if err := message.MarshalJSONBuf(&buffer); err == nil {
				log.Print(buffer.String())
			} else {
				log.Print(err)
			}
		}
	}()

	terminate := make(chan os.Signal)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	s := <-terminate
	log.Println(s)
}
