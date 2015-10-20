
package ga_test

import (
	. "gopkg.in/check.v1"
	"github.com/tomcraven/goga"
)

type SimulatorTestSuite struct {
}
func ( s *SimulatorTestSuite ) SetUpTest( t *C ) {
}
func ( s *SimulatorTestSuite ) TearDownTest( t *C ) {
}
var _ = Suite( &SimulatorTestSuite{} )

func ( s *SimulatorTestSuite ) TestShouldReturnFalseFromExitFunctionFromNullSimulator( t *C ) {
	nullSimulator := ga.NullSimulator{}
	genome := ga.NewGenome( ga.Bitset{} )
	t.Assert( nullSimulator.ExitFunc( &genome ), IsFalse )
}