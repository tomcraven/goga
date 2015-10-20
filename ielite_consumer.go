
package ga

type IEliteConsumer interface {
	OnElite( *IGenome )
}

type NullEliteConsumer struct {
}

func ( nec *NullEliteConsumer ) OnElite( *IGenome ) {
}