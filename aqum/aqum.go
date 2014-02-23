package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/tracyde/aquadyno/probe"
)

const (
	USER        = "tracyde"
	APIKEY      = "1234-1234-1234-1234"
	UPDATEURL   = "http://www.mytankstats.com/api.php?data="
	AQSMADDRESS = "127.0.0.1:2191"
)

type Update struct {
	User   string
	Apikey string
	Date   time.Time
	Probes *[]probe.Probe
}

type Empty struct{}

func sendUpdate(u *Update) error {
	j, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("json: %v\n", string(j))

	b64 := base64.StdEncoding.EncodeToString(j)
	fmt.Printf("base64: %v\n", b64)

	fmt.Printf("Submitting update to website: %v\n", UPDATEURL+b64)
	// resp, err := http.Get(UPDATEURL + b64)
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(20 * time.Second)
				c, err := net.DialTimeout(netw, addr, 10*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	resp, err := client.Get(UPDATEURL + b64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Got status code: %v\n", resp.StatusCode)
	return nil
}

func gatherProbes() (*[]probe.Probe, error) {
	client, err := rpc.DialHTTP("tcp", AQSMADDRESS)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply []probe.Probe
	err = client.Call("Probes.Gather", Empty{}, &reply)
	if err != nil {
		log.Fatal("probes error:", err)
	}
	return &reply, nil
}

func main() {
	probes, err := gatherProbes()
	if err != nil {
		fmt.Println(err)
		return
	}

	update := &Update{
		User:   USER,
		Apikey: APIKEY,
		Date:   time.Now(),
		Probes: probes}

	sendUpdate(update)
}
