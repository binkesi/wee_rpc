package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	day5httpdebug "wee_rpc/day5-http-debug"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	// pick a free port
	var foo Foo
	if err := day5httpdebug.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	day5httpdebug.HandleHTTP()
	addr <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := day5httpdebug.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()
	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {
	log.SetFlags(0)
	ch := make(chan string)
	go call(ch)
	startServer(ch)
}
