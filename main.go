package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var address = flag.String("address", "127.0.0.1:12201", "address to listen on")

func main() {
	flag.Parse()

	reader, err := gelf.NewReader(*address)
	if err != nil {
		log.Fatal(err)
	}

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
