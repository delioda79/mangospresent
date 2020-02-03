package main

import (
	"fmt"
	"log"
	"nanomsg.org/go/mangos/v2/protocol/rep"
	"nanomsg.org/go/mangos/v2/protocol/req"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
	"os"
	"sync"
)

func requester(name string, uu []string, w *sync.WaitGroup) {
	rqs, err := req.NewSocket()
	if err != nil {
		os.Exit(1)
	}

	for _, u:= range uu {
		err := rqs.Dial(u)
		if err != nil {
			log.Fatalf("%s not reachable: %v", u, err)
		}
	}

	for i:=0; i<=4; i++ {
		rqs.Send([]byte(name))
		bts, err := rqs.Recv()
		if err != nil {
			log.Printf("Requester received an error: %v\n", err)
		} else {
			log.Printf("%s received the message: %s", name, string(bts))
		}
	}


	w.Done()

}

func replier(u string, w *sync.WaitGroup) {
	rps, err := rep.NewSocket()
	if err != nil {
		os.Exit(1)
	}

	err = rps.Listen(u)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", u, err)
	}
	log.Printf("listening to %s", u)
	w.Done()
	for {
		bts, err := rps.Recv()
		if err != nil {
			log.Printf("Replier %s received an error: %v\n", err)
			rps.Send([]byte(""))
		} else {
			rps.Send([]byte(fmt.Sprintf("%s greets you %s", u, string(bts))))
		}
	}

}

func main() {

	w := &sync.WaitGroup{}

	urls := []string{"inproc://url1", "inproc://url2", "inproc://url3", "inproc://url4"}

	for _, u := range urls {
		go replier(u, w)
		w.Add(1)
	}

	w.Wait()

	w.Add(1)
	go requester("req1", urls, w)

	w.Add(1)
	go requester("req2", urls, w)

	w.Wait()

}
