package goga

// ISimulator - a Simulator interface
type ISimulator interface {
	OnBeginSimulation()
	Simulate(Genome)
	OnEndSimulation()
	ExitFunc(Genome) bool
}

// NullSimulator - a null implementation of the Simulator interface
type NullSimulator struct {
}

// Simulate - a null implementation of Simulator's 'Simulate'
func (ns *NullSimulator) Simulate(Genome) {
}

// OnBeginSimulation - a null implementation of Simulator's 'OnBeginSimulation'
func (ns *NullSimulator) OnBeginSimulation() {
}

// OnEndSimulation - a null implementation of Simulator's 'OnEndSimulation'
func (ns *NullSimulator) OnEndSimulation() {
}

// ExitFunc - a null implementation of Simulator's 'ExitFunc'
func (ns *NullSimulator) ExitFunc(Genome) bool {
	return false
}
