package main

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit"
)

func clockCircuit() {
	cir := circuit.NewCircuit()

	clockOut := cir.Clock(2)
	cir.Initialize()

	fmt.Println("Done initializing")

	initialChange := circuit.Change{
		Time:   0,
		Wire:   clockOut,
		Signal: circuit.Low,
	}
	cir.Propagate(initialChange)
}

func adderCircuit() {
	cir := circuit.NewCircuit()

	in1 := cir.AddInputWire32("input1")
	in2 := cir.AddInputWire32("input2")
	out, overflow := cir.FullAdder32(in1, in2)
	cir.Initialize()

	fmt.Printf("Initial output value: %v\n", out.AsUint32())
	fmt.Printf("Initial overflow value: %v\n", overflow.Signal())

	initialChanges := []circuit.Change{}

	inputValue1 := uint32(4)
	inputValue2 := uint32(7)
	initialChanges = append(initialChanges, circuit.Change32(1, in1, inputValue1)...)
	initialChanges = append(initialChanges, circuit.Change32(1, in2, inputValue2)...)

	cir.Propagate(initialChanges...)

	fmt.Printf("Output value: %v\n", out.AsUint32())
	fmt.Printf("Overflow: %v\n", overflow.Signal())
}

func main() {
	// adderCircuit()
	clockCircuit()
}
