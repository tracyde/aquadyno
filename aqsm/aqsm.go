package main

import (
	"github.com/tracyde/aquadyno/aqsm/probes"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	probes := new(probes.Probes)
	rpc.Register(probes)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2191")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
