package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tracyde/aquadyno/probe"
)

const (
	USER      = "tracyde"
	APIKEY    = "1234-1234-1234-1234"
	UPDATEURL = "http://www.mytankstats.com/api.php?data="
)

type Update struct {
	User   string
	Apikey string
	Date   time.Time
	Probes *[]probe.Probe
}

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
	resp, err := http.Get(UPDATEURL + b64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Got status code: %v\n", resp.StatusCode)
	return nil
}

func gatherProbes() (p *[]probe.Probe, _ error) {
	p1 := probe.NewThermal("tank1", "aquarium temp", 25.556)
	p2 := probe.NewThermal("tank1", "ambient temp", 22.197)

	p = &[]probe.Probe{*p1, *p2}
	return
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
