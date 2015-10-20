
package ga

type ISimulator interface {
	OnBeginSimulation()
	Simulate( *IGenome )
	OnEndSimulation()
	ExitFunc( *IGenome ) bool
}

type NullSimulator struct {
}
func ( ns *NullSimulator ) Simulate( *IGenome ) {
}
func ( ns *NullSimulator ) OnBeginSimulation() {
}
func ( ns *NullSimulator ) OnEndSimulation() {
}
func ( ns *NullSimulator ) ExitFunc( *IGenome ) bool {
	return false;
}