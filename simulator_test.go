package ga_test

import . "gopkg.in/check.v1"
import ga "github.com/tomcraven/goga"
import "github.com/tomcraven/bitset"

type FakeSimulator struct {
	OnBeginSimulationCalls int
	SimulateCalls          int
	OnEndSimulationCalls   int
	ExitFuncCalls          int
	ExitFuncRet            bool
}

func (fs *FakeSimulator) Reset() {
	fs.OnBeginSimulationCalls = 0
	fs.SimulateCalls = 0
	fs.OnEndSimulationCalls = 0
	fs.ExitFuncCalls = 0
}
func (fs *FakeSimulator) OnBeginSimulation() {
	fs.OnBeginSimulationCalls++
}
func (fs *FakeSimulator) Simulate(*ga.Genome) {
	fs.SimulateCalls++
}
func (fs *FakeSimulator) OnEndSimulation() {
	fs.OnEndSimulationCalls++
}
func (fs *FakeSimulator) ExitFunc(*ga.Genome) bool {
	fs.ExitFuncCalls++
	return fs.ExitFuncRet
}

type SimulatorTestSuite struct {
}

func (s *SimulatorTestSuite) SetUpTest(t *C) {
}
func (s *SimulatorTestSuite) TearDownTest(t *C) {
}

var _ = Suite(&SimulatorTestSuite{})

func (s *SimulatorTestSuite) TestShouldReturnFalseFromExitFunctionFromNullSimulator(t *C) {
	nullSimulator := ga.NullSimulator{}
	genome := ga.NewGenome(bitset.Create(0))
	t.Assert(nullSimulator.ExitFunc(&genome), IsFalse)
}

func (s *SimulatorTestSuite) TestSimulatorSwitch_OnBeginSimulation(t *C) {
	fs1 := FakeSimulator{}
	fs2 := FakeSimulator{}
	fs1.Reset()
	fs2.Reset()
	ss := ga.SimulatorSwitch{
		Simulators: []ga.Simulator{
			&fs1, &fs2,
		},
	}

	ss.OnBeginSimulation()
	t.Assert(fs2.OnBeginSimulationCalls, Equals, 1)
	t.Assert(fs1.OnBeginSimulationCalls, Equals, 1)
}

func (s *SimulatorTestSuite) TestSimulatorSwitch_Simulate(t *C) {
	fs1 := FakeSimulator{}
	fs2 := FakeSimulator{}
	fs1.Reset()
	fs2.Reset()
	ss := ga.SimulatorSwitch{
		Simulators: []ga.Simulator{
			&fs1, &fs2,
		},
	}

	ss.Simulate(nil)
	t.Assert(fs2.SimulateCalls, Equals, 1)
	t.Assert(fs1.SimulateCalls, Equals, 1)
}

func (s *SimulatorTestSuite) TestSimulatorSwitch_OnEndSimulation(t *C) {
	fs1 := FakeSimulator{}
	fs2 := FakeSimulator{}
	fs1.Reset()
	fs2.Reset()
	ss := ga.SimulatorSwitch{
		Simulators: []ga.Simulator{
			&fs1, &fs2,
		},
	}

	ss.OnEndSimulation()
	t.Assert(fs2.OnEndSimulationCalls, Equals, 1)
	t.Assert(fs1.OnEndSimulationCalls, Equals, 1)
}

func (s *SimulatorTestSuite) TestSimulatorSwitch_OnExitFunc(t *C) {
	fs1 := FakeSimulator{}
	fs2 := FakeSimulator{}
	fs1.Reset()
	fs2.Reset()
	ss := ga.SimulatorSwitch{
		Simulators: []ga.Simulator{
			&fs1, &fs2,
		},
	}

	ss.ExitFunc(nil)
	t.Assert(fs2.ExitFuncCalls, Equals, 1)
	t.Assert(fs1.ExitFuncCalls, Equals, 1)
}
