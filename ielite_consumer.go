package ga

// IEliteConsumer - an interface to the elite consumer
type IEliteConsumer interface {
	OnElite(*IGenome)
	OnBeginSimulation()
}

// NullEliteConsumer - a null implementation of the elite consumer
type NullEliteConsumer struct {
}

// OnElite - null implementation of OnElite from the EliteConsumer interface
func (nec *NullEliteConsumer) OnElite(*IGenome) {
}

// OnBeginSimulation - null implementation of OnBeginSimulation from the EliteConsumer interface
func (nec *NullEliteConsumer) OnBeginSimulation() {
}
