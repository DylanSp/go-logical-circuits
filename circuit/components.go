package circuit

import (
	"fmt"
	"maps"
	"slices"
)

func (cir *Circuit) Not(in *Wire) *Wire {
	gate := cir.addInternalWire()
	gate.SetInputs(in)

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   gate,
				Signal: !ch.Signal,
			},
		}
	}

	return gate
}

func binaryGateFactory(truthTable func(Signal, Signal) Signal) func(*Circuit, *Wire, *Wire) *Wire {
	return func(cir *Circuit, in1, in2 *Wire) *Wire {
		gate := cir.addInternalWire()
		gate.SetInputs(in1, in2)

		gate.processChangeLogic = func(ch Change) []Change {
			currentInputs := slices.Collect(maps.Values(gate.inputWires))

			return []Change{
				{
					Time:   ch.Time + gateDelay,
					Wire:   gate,
					Signal: truthTable(currentInputs[0].Signal(), currentInputs[1].Signal()),
				},
			}
		}

		return gate
	}
}

func (cir *Circuit) And(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 && s2
	})(cir, in1, in2)
}

func (cir *Circuit) Or(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 || s2
	})(cir, in1, in2)
}

func (cir *Circuit) Xor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 != s2
	})(cir, in1, in2)
}

func (cir *Circuit) Nand(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return !(s1 && s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Nor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return !(s1 || s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Xnor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 == s2
	})(cir, in1, in2)
}

func (cir *Circuit) HalfAdder(in1, in2 *Wire) (*Wire, *Wire) {
	sum := cir.Xor(in1, in2)
	carry := cir.And(in1, in2)

	return sum, carry
}

func (cir *Circuit) FullAdder(in1, in2, inCarry *Wire) (*Wire, *Wire) {
	sum := cir.Xor(cir.Xor(in1, in2), inCarry)
	carry := cir.Or(cir.And(in1, in2), cir.And(inCarry, cir.Xor(in1, in2)))
	return sum, carry
}

func (cir *Circuit) Noop(in *Wire) *Wire {
	gate := cir.addInternalWire()
	gate.SetInputs(in)

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   gate,
				Signal: ch.Signal,
			},
		}
	}

	return gate
}

// TODO - doesn't currently work if period=1; try and fix this?
// currently needs a manual change to be applied to its output wire (which is also the loopback wire) to start it
func (cir *Circuit) Clock(period int) *Wire {
	if period <= 0 {
		panic(fmt.Sprintf("Unable to create clock with period %v; period must be at least 1", period))
	}

	in := cir.addInternalWire() // dummy wire for initial input when constructing no-ops; will be irrelevant once we set up loopback
	var out *Wire

	var initialNoop *Wire

	// chain of (period - 1) no-ops to delay clock signal, followed by NOT gate to invert signal after `period` ticks
	for i := range period - 1 {
		out = cir.Noop(in)
		if i == 0 {
			initialNoop = out
		}

		in = out
	}

	// set up loopback, connect to first Noop component
	loopbackWire := cir.Not(out)
	initialNoop.SetInputs(loopbackWire)

	return out
}

// see https://en.wikipedia.org/wiki/Flip-flop_(electronics)#Gated_D_latch
// diagram (input stage is on the left, output stage on the right; gates numbered top to bottom)
// https://commons.wikimedia.org/wiki/File:D-Type_Transparent_Latch.svg
func (cir *Circuit) gatedDLatch(input, clock *Wire) *Wire {
	// input stage
	inputStageGate1 := cir.Nand(input, clock)
	inputStageGate2 := cir.Nand(inputStageGate1, clock)

	// dummy wires for constructing output stage gates;
	// will be ignored after output gates' inputs are reset
	dummy1 := cir.addInternalWire()
	dummy2 := cir.addInternalWire()

	// output stage (SR NAND latch)
	outputStageGate1 := cir.Nand(dummy1, dummy2)
	outputStageGate2 := cir.Nand(dummy1, dummy2)

	// set output stage inputs to their proper values, looping back
	outputStageGate1.SetInputs(inputStageGate1, outputStageGate2)
	outputStageGate2.SetInputs(outputStageGate1, inputStageGate2)

	return outputStageGate1
}

// Master–slave edge-triggered D flip-flop, using Wikipedia's terminology
// stores a value on the *falling* edge of the clock
// see https://en.wikipedia.org/wiki/Flip-flop_(electronics)#Master%E2%80%93slave_edge-triggered_D_flip-flop
// diagram of rising-edge version: https://commons.wikimedia.org/wiki/File:D-Type_Flip-flop_Diagram.svg
// (fallingEdgeFlipFlop doesn't have the first NOT on C/clock input)
func (cir *Circuit) fallingEdgeFlipFlop(input, clock *Wire) *Wire {
	latch1 := cir.gatedDLatch(input, clock)
	notClock := cir.Not(clock)
	latch2 := cir.gatedDLatch(latch1, notClock)
	return latch2
}

// stores a value on the *rising* edge of the clock
func (cir *Circuit) FlipFlop(input, clock *Wire) *Wire {
	return cir.fallingEdgeFlipFlop(input, cir.Not(clock))
}

func (cir *Circuit) Mux(ifTrue, ifFalse, selector *Wire) *Wire {
	return cir.Or(cir.And(ifTrue, selector), cir.And(ifFalse, cir.Not(selector)))
}
