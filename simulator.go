package ga

// ISimulator - a Simulator interface
type ISimulator interface {
	OnBeginSimulation()
	Simulate(*IGenome)
	OnEndSimulation()
	ExitFunc(*IGenome) bool
}

// NullSimulator - a null implementation of the Simulator interface
type NullSimulator struct {
}

// Simulate - a null implementation of Simulator's 'Simulate'
func (ns *NullSimulator) Simulate(*IGenome) {
}

// OnBeginSimulation - a null implementation of Simulator's 'OnBeginSimulation'
func (ns *NullSimulator) OnBeginSimulation() {
}

// OnEndSimulation - a null implementation of Simulator's 'OnEndSimulation'
func (ns *NullSimulator) OnEndSimulation() {
}

// ExitFunc - a null implementation of Simulator's 'ExitFunc'
func (ns *NullSimulator) ExitFunc(*IGenome) bool {
	return false
}
