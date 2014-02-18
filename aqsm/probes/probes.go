package probes

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"github.com/mrmorphic/hwio/devices/tmp102"
	"github.com/tracyde/aquadyno/probe"
	"github.com/tracyde/onewire"
)

type Probes int

const (
	ID1 = "28-000005b3c94a"
)

func (*Probes) Gather(_ *struct{}, reply *[]probe.Probe) error {
	var probes []probe.Probe

	// Create ambient temp probe
	aT, err := getAmbient()
	if err != nil {
		return err
	}
	probes = append(probes, *probe.NewThermal("tank1", "ambient temp", aT))

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

func getAmbient() (float32, error) {
	// Get the i2c module from the driver
	m, err := hwio.GetModule("i2c")
	if err != nil {
		fmt.Printf("Error getting i2c module: %v\n", err)
		return 0, err
	}

	// Assert that it is an i2c module
	i2c := m.(hwio.I2CModule)

	// Enable the i2c bus, needed for Raspberry Pi
	i2c.Enable()
	defer i2c.Disable()

	// Get a temp device on this bus
	temp := tmp102.NewTMP102(i2c)

	// Get the temperature sensor value
	t, err := temp.GetTemp()
	if err != nil {
		fmt.Printf("Error while reading i2c temp: %v\n", err)
		return 0, err
	}

	return t, nil
}
