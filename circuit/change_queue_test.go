package circuit_test

import (
	"slices"
	"testing"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/DylanSp/go-logical-circuits/circuit/wire"
	"github.com/stretchr/testify/assert"
)

func TestChangeQueue(t *testing.T) {
	t.Run("Next change from an empty queue is empty", func(t *testing.T) {
		queue := circuit.NewChangeQueue()

		nextChanges, ok := queue.GetNextChanges()

		assert.False(t, ok)
		assert.Empty(t, nextChanges)
	})

	t.Run("Next change from a single-change queue is that change", func(t *testing.T) {
		queue := circuit.NewChangeQueue()

		ch := circuit.Change{
			Time:   12,
			Wire:   nil,
			Signal: false,
		}
		queue.AddChange(ch)

		nextChanges, ok := queue.GetNextChanges()

		assert.True(t, ok)
		assert.Len(t, nextChanges, 1)
		assert.Equal(t, ch.Time, nextChanges[0].Time)
	})

	t.Run("Next change from a queue with two changes is the change with the lowest time", func(t *testing.T) {
		queue := circuit.NewChangeQueue()

		laterChange := circuit.Change{
			Time:   20,
			Wire:   nil,
			Signal: false,
		}
		queue.AddChange(laterChange)

		earlierChange := circuit.Change{
			Time:   10,
			Wire:   nil,
			Signal: false,
		}
		queue.AddChange(earlierChange)

		nextChanges, ok := queue.GetNextChanges()

		assert.True(t, ok)
		assert.Len(t, nextChanges, 1)
		assert.Equal(t, earlierChange.Time, nextChanges[0].Time)
	})

	t.Run("Next changes from a queue with two changes of equal time returns both changes", func(t *testing.T) {
		queue := circuit.NewChangeQueue()

		change1 := circuit.Change{
			Time:   5,
			Wire:   wire.New("wire1"),
			Signal: false,
		}
		queue.AddChange(change1)

		change2 := circuit.Change{
			Time:   5,
			Wire:   wire.New("wire2"),
			Signal: true,
		}
		queue.AddChange(change2)

		// this should *not* be popped on first call to GetNextChanges()
		laterChange := circuit.Change{
			Time:   100,
			Wire:   wire.New("later"),
			Signal: true,
		}
		queue.AddChange(laterChange)

		nextChanges, ok := queue.GetNextChanges()

		assert.True(t, ok)
		assert.Len(t, nextChanges, 2)
		assert.True(t, slices.ContainsFunc(nextChanges, func(ch circuit.Change) bool {
			return ch.Wire.Name() == change1.Wire.Name()
		}))
		assert.True(t, slices.ContainsFunc(nextChanges, func(ch circuit.Change) bool {
			return ch.Wire.Name() == change2.Wire.Name()
		}))
	})
}
