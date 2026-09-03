package main

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit"
)

func main() {
	cir := circuit.NewCircuit()

	in1 := cir.AddInputWire("input1")
	out1 := cir.AddOutputWire("output1")
	cir.MkNotGate(in1, out1)

	in2 := cir.AddInputWire("input2")
	in3 := cir.AddInputWire("input3")
	out2 := cir.AddOutputWire("output2")
	cir.MkAndGate(in2, in3, out2)

	in4 := cir.AddInputWire("input4")
	out3 := cir.AddOutputWire("output3")
	cir.Mk3AndGate(in2, in3, in4, out3)

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
		{
			Time:   0,
			Wire:   in4,
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

	// 3-way (compound) AND gate testing

	fmt.Printf("Initial state of out3 wire: %v\n\n", out3.Signal())

	change5 := circuit.Change{
		Time:   50,
		Wire:   in4,
		Signal: true,
	}
	cir.Propagate(change5)

	fmt.Printf("Intermediate state of out3 wire: %v\n\n", out3.Signal())

	change6 := circuit.Change{
		Time:   60,
		Wire:   in2,
		Signal: false,
	}
	cir.Propagate(change6)
	fmt.Printf("Final state of out3 wire: %v\n\n", out3.Signal())
	fmt.Printf("Additionally, state of out2 wire: %v\n\n", out2.Signal())
}
