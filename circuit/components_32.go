package circuit

import "github.com/DylanSp/go-logical-circuits/circuit/wire"

func (cir *Circuit) Not32(in *wire.Wire32) *wire.Wire32 {
	// don't think we necessary need this? the individual Components will be the single-bit Not gates
	// may be helpful for later optimization/tooling, though
	// gate := Component{
	// inputWires: map[*wire.Wire]struct{}{},
	// }

	outWires := [32]*wire.Wire{}

	for i, singleIn := range in.AllWires() {
		outWires[i] = cir.Not(singleIn)
		// gate.inputWires[singleIn] = struct{}{}
	}

	out := cir.addInternalWire32FromSingleWires(outWires)

	return out
}

func binaryGateFactory32(
	gate func(*wire.Wire, *wire.Wire) *wire.Wire,
) func(*Circuit, *wire.Wire32, *wire.Wire32) *wire.Wire32 {
	return func(cir *Circuit, in1, in2 *wire.Wire32) *wire.Wire32 {
		outWires := [32]*wire.Wire{}

		for i := range 32 {
			outWires[i] = gate(in1.Wire(i), in2.Wire(i))
		}

		out := cir.addInternalWire32FromSingleWires(outWires)
		return out
	}
}

func (cir *Circuit) And32(in1, in2 *wire.Wire32) *wire.Wire32 {
	return binaryGateFactory32(cir.And)(cir, in1, in2)
}
