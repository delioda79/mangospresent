package main

import (
	"log"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
	"os"
	"sync"
)

func publisher(uu []string, w *sync.WaitGroup) {
	pbs, err := pub.NewSocket()
	if err != nil {
		os.Exit(1)
	}

	for _, u:= range uu {
		err := pbs.Dial(u)
		if err != nil {
			log.Fatalf("%s not reachable: %v", u, err)
		}
	}

	for i:=0; i<=0; i++ {
		w.Add(len(uu))
		pbs.Send([]byte("hello"))
	}

	w.Done()

}

func subscriber(u string, w1,w2 *sync.WaitGroup) {
	sbs, err := sub.NewSocket()
	if err != nil {
		os.Exit(1)
	}

	err = sbs.Listen(u)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", u, err)
	}

	sbs.SetOption(mangos.OptionSubscribe, []byte(""))
	log.Printf("listening to %s", u)
	w1.Done()
	for {
		bts, err := sbs.Recv()
		if err != nil {
			log.Printf("Subscriber %s received an error: %v\n", err)
		} else {
			log.Printf("%s received %s \n", u, string(bts))
		}
		w2.Done()
	}
}

func main() {

	w1 := &sync.WaitGroup{}
	w2 := &sync.WaitGroup{}

	urls := []string{"inproc://url1", "inproc://url2", "inproc://url3", "inproc://url4"}

	w1.Add(len(urls))
	for _, u := range urls {
		go subscriber(u, w1,w2)
	}

	w1.Wait()

	w2.Add(1)
	go publisher(urls, w2)

	w2.Wait()

}
