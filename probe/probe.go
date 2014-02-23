package probe

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"github.com/mrmorphic/hwio/devices/tmp102"
	"github.com/tracyde/onewire"
	"log"
)

type Probe struct {
	Target string
	Bus    string
	Id     string
	Name   string
	Type   string
	Value  float32
}

func NewThermal(target, bus, id, name string) *Probe {
	result := new(Probe)
	result.Target = target
	result.Bus = bus
	result.Id = id
	result.Name = name
	result.Type = "temp"

	return result
}

func (p *Probe) ReadData() (v float32, err error) {
	// ensure Id is not empty
	if p.Id == "" {
		log.Printf("Probe Error: No id:%v for probe named: %v\n", p.Id, p.Name)
	}

	switch p.Bus {
	case "i2c":
		v, err = p.readI2cData()
	case "1wire":
		v, err = p.read1WireData()
	default:
		err = fmt.Errorf("Unsupported bus type: %v", p.Bus)
	}
	return
}

func (p *Probe) readI2cData() (v float32, err error) {
	// Get the i2c module from the driver
	m, err := hwio.GetModule("i2c")
	if err != nil {
		log.Printf("Error getting i2c module: %v\n", err)
		return
	}

	// Assert that it is an i2c module
	i2c := m.(hwio.I2CModule)

	// Enable the i2c bus, needed for Raspberry Pi
	i2c.Enable()
	defer i2c.Disable()

	// Get a temp device on this bus
	temp := tmp102.NewTMP102(i2c)

	// Get the temperature sensor value
	v, err = temp.GetTemp()
	if err != nil {
		log.Printf("Error while reading i2c temp: %v\n", err)
		return
	}

	return
}

func (p *Probe) read1WireData() (v float32, err error) {
	names, err := onewire.ScanSlaves()
	if err != nil {
		log.Printf("Error scanning 1wire names: %v", err)
		return
	}

	devices := make([]*onewire.DS18B20, len(names))
	for i := range names {
		devices[i], err = onewire.NewDS18B20(names[i])
		if err != nil {
			log.Printf("Error opening device %v: %v\n", devices[i], err)
			return
		}
	}

	for i := range devices {
		log.Printf("attempting read on device %x\n", devices[i].Id)
		value, err := devices[i].Read()
		if err != nil {
			log.Printf("Error on read: %v\n", err)
			return 0, err
		}
		v = float32(value) / 1000
	}

	return
}
