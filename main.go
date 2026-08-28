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

	in2 := cir.AddWire("input2")
	in3 := cir.AddWire("input3")
	out2 := cir.AddWire("output2")
	andGate1 := circuit.MkAndGate(in2, in3, out2)
	cir.AddComponents(andGate1)

	// TODO - need a way to initialize circuit with initial values of input wires
	// TODO - have Circuit track its input wires, add an Initialize method that sets all input wires to low, then calls Propagate()?
	// TODO - may need to ignore/eliminate no-op check in Propagate upon initialization, to make sure initial values are correct

	fmt.Println("Initializing circuit")

	// initialize circuit
	initializingChanges := []circuit.Change{
		{
			Time:   0,
			Wire:   in1,
			Signal: false,
		},
		{
			Time:   0,
			Wire:   in2,
			Signal: false,
		},
		{
			Time:   0,
			Wire:   in3,
			Signal: false,
		},
	}
	for _, initializer := range initializingChanges {
		cir.Propagate(initializer)
	}
	fmt.Printf("Circuit initialized\n\n")

	// NOT gate testing

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

	// AND gate testing

	fmt.Printf("Initial state of out2 wire: %v\n\n", out2.Signal())

	change3 := circuit.Change{
		Time:   30,
		Wire:   in2,
		Signal: true,
	}
	cir.Propagate(change3)

	fmt.Printf("Intermediate state of out2 wire: %v\n\n", out2.Signal())

	change4 := circuit.Change{
		Time:   40,
		Wire:   in3,
		Signal: true,
	}
	cir.Propagate(change4)

	fmt.Printf("Final state of out2 wire: %v\n\n", out2.Signal())
}
