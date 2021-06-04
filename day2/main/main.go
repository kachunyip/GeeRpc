package main

import (
	"fmt"
	"geerpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	client, _ := geerpc.Dial("tcp", <-addr)

	//conn, _ := net.Dial("tcp", <-addr)
	defer func() {_ = client.Close()}()

	time.Sleep(time.Second)
	var wg sync.WaitGroup
	//_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	//cc := codec.NewGobCodec(conn)
	for i := 0; i<5;i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
		//h := &codec.Header{
		//	ServiceMethod: "Foo.Sum",
		//	Seq: uint64(i),
		//}
		//_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		//_ = cc.ReadHeader(h)
		//var reply string
		//_ = cc.ReadBody(&reply)
		//log.Println("reply:", reply)
	}
	wg.Wait()
}
