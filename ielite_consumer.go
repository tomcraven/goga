package ga

// EliteConsumer - an interface to the elite consumer
type EliteConsumer interface {
	OnElite(*Genome)
}

// NullEliteConsumer - a null implementation of the elite consumer
type NullEliteConsumer struct {
}

// OnElite - null implementation of OnElite from the EliteConsumer interface
func (nec *NullEliteConsumer) OnElite(*Genome) {
}
