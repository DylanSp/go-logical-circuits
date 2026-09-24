package circuit

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
