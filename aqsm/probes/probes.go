package probes

import (
	"fmt"
	"github.com/tracyde/aquadyno/probe"
	"github.com/tracyde/onewire"
)

type Probes int

const (
	ID = "28-000005b3c94a"
)

func (*Probes) Gather(_ *struct{}, reply *[]probe.Probe) error {
	var probes []probe.Probe
	probes = append(probes, *probe.NewThermal("tank1", "ambient temp", 22.197))

	names, err := onewire.ScanSlaves()
	if err != nil {
		fmt.Printf("Error scanning 1wire names: %v", err)
		return err
	}

	devices := make([]*onewire.DS18B20, len(names))
	for i := range names {
		devices[i], err = onewire.NewDS18B20(names[i])
		if err != nil {
			fmt.Printf("Error opening device %v: %v", devices[i], err)
			return err
		}
	}

	for i := range devices {
		fmt.Printf("attempting read on device %x\n", devices[i].Id)
		value, err := devices[i].Read()
		if err != nil {
			fmt.Printf("Error on read: %v", err)
			return err
		}
		probes = append(probes, *probe.NewThermal("tank1", "aquarium temp", float32(value)/1000))
	}

	*reply = probes
	return nil
}
