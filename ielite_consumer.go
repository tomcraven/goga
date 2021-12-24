package goga

// IEliteConsumer - an interface to the elite consumer
type IEliteConsumer interface {
	OnElite(Genome)
}

// NullEliteConsumer - a null implementation of the elite consumer
type NullEliteConsumer struct {
}

// OnElite - null implementation on OnElite from the EliteConsumer interface
func (nec *NullEliteConsumer) OnElite(Genome) {
}
