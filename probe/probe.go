package probe

type Probe struct {
	Target string
	Name   string
	Value  float32
	Type   string
}

func NewThermal(target, name string, value float32) *Probe {
	result := new(Probe)
	result.Target = target
	result.Name = name
	result.Type = "temp"
	result.Value = value

	return result
}
