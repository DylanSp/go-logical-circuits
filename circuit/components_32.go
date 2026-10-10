package circuit

import "math"

func (cir *Circuit) Not32(in *Wire32) *Wire32 {
	outWires := [32]*Wire{}

	for i, singleIn := range in.AllWires() {
		outWires[i] = cir.Not(singleIn)
	}

	out := cir.addInternalWire32FromSingleWires(outWires)

	return out
}

func binaryGateFactory32(
	gate func(*Wire, *Wire) *Wire,
) func(*Circuit, *Wire32, *Wire32) *Wire32 {
	return func(cir *Circuit, in1, in2 *Wire32) *Wire32 {
		outWires := [32]*Wire{}

		for i := range 32 {
			outWires[i] = gate(in1.Wire(i), in2.Wire(i))
		}

		out := cir.addInternalWire32FromSingleWires(outWires)
		return out
	}
}

func (cir *Circuit) And32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.And)(cir, in1, in2)
}

func (cir *Circuit) Or32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.Or)(cir, in1, in2)
}

func (cir *Circuit) Xor32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.Xor)(cir, in1, in2)
}

func (cir *Circuit) Nand32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.Nand)(cir, in1, in2)
}

func (cir *Circuit) Nor32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.Nor)(cir, in1, in2)
}

func (cir *Circuit) Xnor32(in1, in2 *Wire32) *Wire32 {
	return binaryGateFactory32(cir.Xnor)(cir, in1, in2)
}

// ripple adder
// TODO - should overflow output be a Wire32 for consistency?
func (cir *Circuit) FullAdder32(in1, in2 *Wire32) (*Wire32, *Wire) {
	sumOutWires := [32]*Wire{}

	// half adder for bit 0 (least significant bit)
	sum0, carry := cir.HalfAdder(in1.Wire(0), in2.Wire(0))
	sumOutWires[0] = sum0

	// for bits 1 through 31: full adders, passing carries through
	for i := range 31 {
		bitNumber := i + 1
		sum, carryOut := cir.FullAdder(in1.Wire(bitNumber), in2.Wire(bitNumber), carry)
		sumOutWires[bitNumber] = sum
		carry = carryOut
	}

	sumOut := cir.addInternalWire32FromSingleWires(sumOutWires)

	// overflow detection - instead of connecting carry from last full adder to another adder, use it as overflow output
	return sumOut, carry
}

// multiplexer between 2 32-bit wires with a single-bit selector
func (cir *Circuit) Mux32(ifTrue, ifFalse *Wire32, selector *Wire) *Wire32 {
	outWires := [32]*Wire{}

	for i := range 32 {
		outWires[i] = cir.Mux(ifTrue.Wire(i), ifFalse.Wire(i), selector)
	}

	return cir.addInternalWire32FromSingleWires(outWires)
}

// multiplexer with an n-bit selector between 2^n 32-bit wires
// selectors is a big-endian index into inputs; selectors[0] is most-significant bit
// 2^len(selectors) must equal len(inputs)
// returns (output wire, true) if construction is valid; returns (nil, false) if invalid
func (cir *Circuit) Mux32N(inputs []*Wire32, selectors []*Wire) (*Wire32, bool) {
	if powInt(2, len(selectors)) != len(inputs) {
		return nil, false
	}

	if len(selectors) == 0 {
		return nil, false
	}

	// base case
	if len(selectors) == 1 {
		return cir.Mux32(inputs[1], inputs[0], selectors[0]), true
	}

	// recursive case
	// use the most significant selector bit to select between a Mux32N of the top half of `inputs` and a Mux32N of the bottom half
	// both those use the remaining (n - 1) bits of the selector
	halfLength := len(inputs) / 2
	bottomHalf := inputs[:halfLength]
	topHalf := inputs[halfLength:]
	remainingSelectors := selectors[1:]

	topHalfMux, _ := cir.Mux32N(topHalf, remainingSelectors)
	bottomHalfMux, _ := cir.Mux32N(bottomHalf, remainingSelectors)

	return cir.Mux32(topHalfMux, bottomHalfMux, selectors[0]), true
}

func powInt(base, exponent int) int {
	return int(math.Pow(float64(base), float64(exponent)))
}
