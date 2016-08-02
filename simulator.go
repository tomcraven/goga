package ga

// Simulator - a Simulator interface
type Simulator interface {
	OnBeginSimulation()
	Simulate(*Genome)
	OnEndSimulation()
	ExitFunc(*Genome) bool
}

// NullSimulator - a null implementation of the Simulator interface
type NullSimulator struct {
}

// Simulate - a null implementation of Simulator's 'Simulate'
func (ns *NullSimulator) Simulate(*Genome) {
}

// OnBeginSimulation - a null implementation of Simulator's 'OnBeginSimulation'
func (ns *NullSimulator) OnBeginSimulation() {
}

// OnEndSimulation - a null implementation of Simulator's 'OnEndSimulation'
func (ns *NullSimulator) OnEndSimulation() {
}

// ExitFunc - a null implementation of Simulator's 'ExitFunc'
func (ns *NullSimulator) ExitFunc(*Genome) bool {
	return false
}

// SimulatorSwitch - takes an array of simulators and calls into each on the simulator callbacks
type SimulatorSwitch struct {
	Simulators []Simulator
}

// Simulate - calls Simulate on its array of simulators
func (ss *SimulatorSwitch) Simulate(g *Genome) {
	for _, s := range ss.Simulators {
		s.Simulate(g)
	}
}

// OnBeginSimulation - calls OnBeginSumulation on each of its array of simulators
func (ss *SimulatorSwitch) OnBeginSimulation() {
	for _, s := range ss.Simulators {
		s.OnBeginSimulation()
	}
}

// OnEndSimulation - calls OnEndSimulation on each of its array of simulators
func (ss *SimulatorSwitch) OnEndSimulation() {
	for _, s := range ss.Simulators {
		s.OnEndSimulation()
	}
}

// ExitFunc - calls ExitFunc on each of its array of simulators
// returns true if any one of the simulators returns true
// returns false if none of the simulators return false
func (ss *SimulatorSwitch) ExitFunc(g *Genome) bool {
	for _, s := range ss.Simulators {
		if s.ExitFunc(g) {
			return true
		}
	}

	return false
}
