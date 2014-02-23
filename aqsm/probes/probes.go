package probes

import (
	"github.com/tracyde/aquadyno/probe"
)

type Probes int

func (*Probes) Gather(_ *struct{}, reply *[]probe.Probe) (err error) {
	var probes []probe.Probe

	// Create ambient temp probe
	ambientProbe := *probe.NewThermal("tank1", "i2c", "0x48", "ambient temp")

	// Read ambient probe data
	ambientProbe.Value, err = ambientProbe.ReadData()
	if err != nil {
		return err
	}

	// Add ambient probe to probes
	probes = append(probes, ambientProbe)

	// Create aquarium temp probe
	aquaProbe := *probe.NewThermal("tank1", "1wire", "28-000005b3c94a", "aquarium temp")

	// Read aquarium temp data
	aquaProbe.Value, err = aquaProbe.ReadData()
	if err != nil {
		return err
	}

	// Add aquarium temp probe to probes
	probes = append(probes, aquaProbe)

	// Format reply
	*reply = probes
	return nil
}
