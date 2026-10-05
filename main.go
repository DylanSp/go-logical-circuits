package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/DylanSp/go-logical-circuits/metrics"
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

func notCircuit() {
	cir := circuit.NewCircuit()
	in := cir.AddInputWire("input")
	out := cir.Not(in)

	sim := circuit.NewSimulation(cir)
	sim.Initialize()

	fmt.Println("Sim initialized")
	fmt.Printf("Input: %v\n", in.Signal())
	fmt.Printf("Output: %v\n", out.Signal())
	fmt.Println()

	fmt.Println("Setting input to high")
	sim.Schedule(
		circuit.Change{Time: 10, Wire: in, Signal: circuit.High},
	)
	sim.Tick() // run input changes
	sim.Tick() // run output changes
	fmt.Printf("Input: %v\n", in.Signal())
	fmt.Printf("Output: %v\n", out.Signal())
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

func flipFlopCircuit() {
	// IMPORTANT NOTE FOR DEBUG/TEST CODE:
	// if manually using .Tick(), make sure the number of Tick() calls matches up with the time on scheduled changes
	// need to run for 100 Ticks before changes scheduled at Time=100 will be executed!

	// OTHER IMPORTANT NOTE:
	// memory leak or something somewhere; simulation slows down as it is stepped forward
	// debug/profile (check change queue? maybe have logger use sync.Once (or a package-local mutable global var) so it only checks once?)

	// set up circuit
	cir := circuit.NewCircuit()
	input := cir.AddInputWire("input")
	// input := cir.AddWire("input")

	clock := cir.AddInputWire("clock") // will be manually controlled for testing instead of an actual clock
	// clock := cir.AddWire("clock")
	output := cir.FlipFlop(input, clock)

	sim := circuit.NewSimulation(cir)
	// manually initialize; sim.Initialize() will loop forever beacuse we don't currently discard no-op changes
	sim.Schedule(
		circuit.Change{
			Time:   0,
			Wire:   input,
			Signal: circuit.Low,
		},
		// run clock through a full cycle to stabilize flip-flop in a known & valid state
		circuit.Change{
			Time:   5,
			Wire:   clock,
			Signal: circuit.Low,
		},
		circuit.Change{
			Time:   10,
			Wire:   clock,
			Signal: circuit.High,
		},
		circuit.Change{
			Time:   15,
			Wire:   clock,
			Signal: circuit.Low,
		},
	)

	// let the initial values propagate until circuit stabilizes
	// (though Tick doesn't have any logic to ignore no-op changes after initialization, so changes will still be happening)
	for range 100 {
		sim.Tick()
	}

	fmt.Printf("Output value after 100 ticks: %v\n", output.Signal())
	sim.Tick()

	// fmt.Printf("Output value for the next several ticks:\n")
	// for i := range 20 {
	// 	sim.Tick()
	// 	fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	// }

	// return

	// set input to high; should be ignored by flip-flop, because clock remains low
	sim.Schedule(
		circuit.Change{
			Time:   100,
			Wire:   input,
			Signal: circuit.High,
		},
	)
	fmt.Println("Setting input to high, clock is still low")
	// fmt.Printf("Output value for the next several ticks:\n")
	// for i := range 10 {
	for range 100 {
		sim.Tick()
		// fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	// return

	// set clock to high (clock raising edge); flip-flop should now accept and store High input
	sim.Schedule(
		circuit.Change{
			Time:   200,
			Wire:   clock,
			Signal: circuit.High,
		},
	)

	// NOTE - takes 5 ticks for output to change
	fmt.Println("Clock rising edge")
	// fmt.Printf("Output value for the next several ticks:\n")
	// for i := range 10 {
	for range 100 {
		sim.Tick()
		// fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	// return

	// set input to low, then back to high;
	// should be ignored by flip-flop and output should remain high, because clock hasn't changed (still high)
	sim.Schedule(
		circuit.Change{
			Time:   200,
			Wire:   input,
			Signal: circuit.Low,
		},
		circuit.Change{
			Time:   202,
			Wire:   input,
			Signal: circuit.High,
		},
	)
	fmt.Println("Toggling input while clock remains high")
	// fmt.Printf("Output value for the next several ticks:\n")
	// for i := range 10 {
	for range 100 {
		sim.Tick()
		// fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	// return

	// set clock to low (clock falling edge); output should remain high
	sim.Schedule(
		circuit.Change{
			Time:   300,
			Wire:   clock,
			Signal: circuit.Low,
		})

	// fmt.Printf("Output value for the next several ticks after clock falling edge:\n")
	// for i := range 10 {
	for range 100 {
		sim.Tick()
		// fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	// return

	// set input to low; should be ignored by flip-flop (remaining high), because clock is still low
	sim.Schedule(
		circuit.Change{
			Time:   400,
			Wire:   input,
			Signal: circuit.Low,
		},
	)

	fmt.Println("Setting input to low, clock is still low")
	// fmt.Printf("Output value for the next several ticks:\n")
	// for i := range 10 {
	for range 100 {
		sim.Tick()
		// fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	// return

	// set clock to high (rising edge); flip-flop should now grab Low value after 5 ticks
	sim.Schedule(
		circuit.Change{
			Time:   500,
			Wire:   clock,
			Signal: circuit.High,
		},
	)
	fmt.Println("Clock rising edge")
	fmt.Printf("Output value for the next several ticks:\n")
	for i := range 10 {
		sim.Tick()
		fmt.Printf("After %v more ticks: %v\n", i+1, output.Signal())
	}

	dumpMetrics()
}

func dumpMetrics() {
	workingDirectory, _ := os.Getwd()
	metricsFileName := filepath.Join(workingDirectory, "metrics.log")
	fmt.Printf("Writing metrics to: %v\n", metricsFileName)
	metrics.WriteToFile(metricsFileName)
}

func main() {
	// adderCircuit()
	// clockCircuit()
	// andCircuit()
	// notCircuit()
	// binaryInputChangeCircuit()
	flipFlopCircuit()
}
