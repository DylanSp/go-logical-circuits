package main

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit"
)

func main() {
	cir := circuit.NewCircuit()

	in1 := cir.AddInputWire32("input1")
	in2 := cir.AddInputWire32("input2")
	out, overflow := cir.FullAdder32(in1, in2)
	cir.Initialize()

	fmt.Printf("Initial output value: %v\n", out.AsUint32())
	fmt.Printf("Initial overflow value: %v\n", overflow.Signal())

	// initialChanges := []circuit.Change{
	// 	{
	// 		Time:   1,
	// 		Wire:   in1,
	// 		Signal: circuit.High,
	// 	},
	// }
	// cir.Propagate(initialChanges...)
	// fmt.Printf("Output value: %v\n", out.Signal())

	// clockOut := cir.Clock(2)
	// cir.Initialize()

	// fmt.Println("Done initializing")

	// initialChange := circuit.Change{
	// 	Time:   0,
	// 	Wire:   clockOut,
	// 	Signal: wire.Low,
	// }
	// cir.Propagate(initialChange)
}
