package main

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit"
)

func main() {
	cir := circuit.NewCircuit()
	in1 := cir.AddWire("input1")
	out1 := cir.AddWire("output1")
	notGate1 := circuit.MkNotGate(in1, out1)
	cir.AddComponents(notGate1)

	// TODO - need a way to initialize circuit with initial values of input wires

	fmt.Printf("Initial state of out1 wire: %v\n\n", out1.Signal())

	change1 := circuit.Change{
		Time:   10,
		Wire:   in1,
		Signal: true,
	}
	cir.Propagate(change1)

	fmt.Printf("Intermediate state of out1 wire: %v\n\n", out1.Signal())

	change2 := circuit.Change{
		Time:   20,
		Wire:   in1,
		Signal: false,
	}
	cir.Propagate(change2)

	fmt.Printf("Final state of out1 wire: %v\n\n", out1.Signal())
}
