package main

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit"
)

func clockCircuit() {
	cir := circuit.NewCircuit()
	sim := circuit.NewSimulation(cir)

	clockOut := cir.Clock(2)
	sim.Initialize()

	fmt.Println("Done initializing")

	initialChange := circuit.Change{
		Time:   0,
		Wire:   clockOut,
		Signal: circuit.Low,
	}
	sim.Schedule(initialChange)
	for range 10 {
		tick, ok := sim.Tick()
		if !ok {
			break
		}

		for _, ch := range tick.Changes {
			// only print changes from output wire
			if ch.Wire == clockOut {
				fmt.Println(ch)
			}
		}
	}
}

func adderCircuit() {
	cir := circuit.NewCircuit()
	sim := circuit.NewSimulation(cir)

	in1 := cir.AddInputWire32("input1")
	in2 := cir.AddInputWire32("input2")
	out, overflow := cir.FullAdder32(in1, in2)
	sim.Initialize()

	fmt.Printf("Initial output value: %v\n", out.AsUint32())
	fmt.Printf("Initial overflow value: %v\n", overflow.Signal())

	initialChanges := []circuit.Change{}

	inputValue1 := uint32(4)
	inputValue2 := uint32(7)
	initialChanges = append(initialChanges, circuit.Change32(1, in1, inputValue1)...)
	initialChanges = append(initialChanges, circuit.Change32(1, in2, inputValue2)...)

	sim.Schedule(initialChanges...)
	sim.RunUntilStable()

	fmt.Printf("Output value: %v\n", out.AsUint32())
	fmt.Printf("Overflow: %v\n", overflow.Signal())
}

func andCircuit() {
	cir := circuit.NewCircuit()
	in1 := cir.AddInputWire("input1")
	in2 := cir.AddInputWire("input2")
	_ = cir.And(in1, in2)
	sim := circuit.NewSimulation(cir)
	sim.Initialize()

	sim.Schedule(
		circuit.Change{Time: 10, Wire: in1, Signal: circuit.High},
		circuit.Change{Time: 10, Wire: in2, Signal: circuit.High},
	)

	sim.Tick() // run input changes

	outputTick, _ := sim.Tick()
	fmt.Printf("Number of changes in outputTick: %v\n", len(outputTick.Changes))
	fmt.Printf("outputTick time: %v\n", outputTick.Time)
	for i, ch := range outputTick.Changes {
		fmt.Printf("Change %v\n", i)
		fmt.Println(ch)
		fmt.Println()
	}
}

func binaryInputChangeCircuit() {
	cir := circuit.NewCircuit()
	in1 := cir.AddInputWire("input1")
	in2 := cir.AddInputWire("input2")
	in3 := cir.AddInputWire("input3")
	gate := cir.And(in1, in2)

	sim := circuit.NewSimulation(cir)
	sim.Initialize()
	sim.RunUntilStable()

	gate.SetInputs(in2, in3) // remove in1, add in3

	// if changing inputs with .SetInputs works correctly, this should drive gate high;
	// if it doesn't, it'll be held low by in1
	sim.Schedule(
		circuit.Change{Time: 10, Wire: in2, Signal: circuit.High},
		circuit.Change{Time: 11, Wire: in3, Signal: circuit.High},
	)

	// run input changes
	sim.Tick() // in2 change
	sim.Tick() // in3 change (change from in2 will propagate to gate)

	outputTick, _ := sim.Tick() // change from in3 propagates to gate
	fmt.Printf("Number of changes in outputTick: %v\n", len(outputTick.Changes))
	fmt.Printf("outputTick time: %v\n", outputTick.Time)
	for i, ch := range outputTick.Changes {
		fmt.Printf("Change %v\n", i)
		fmt.Println(ch)
		fmt.Println()
	}
}

func main() {
	// adderCircuit()
	// clockCircuit()
	// andCircuit()
	binaryInputChangeCircuit()
}
