package circuit

import "slices"

// TODO - optimize if necessary

type ChangeQueue struct {
	agenda []Change // invariant - agenda is sorted by time, lowest time first
}

func NewChangeQueue() *ChangeQueue {
	return &ChangeQueue{
		agenda: []Change{},
	}
}

func (cq *ChangeQueue) Length() int {
	return len(cq.agenda)
}

func (cq *ChangeQueue) sortAgenda() {
	slices.SortFunc(cq.agenda, func(ch1 Change, ch2 Change) int {
		return ch1.Time - ch2.Time
	})
}

func (cq *ChangeQueue) AddChange(ch Change) {
	cq.agenda = append(cq.agenda, ch)
	cq.sortAgenda()
}

// pops the change(s) with the smallest time
// if there are multiple changes with the same time - returns all of them
// if there are no queued changes - returns ([], false)
func (cq *ChangeQueue) GetNextChanges() ([]Change, bool) {
	if len(cq.agenda) == 0 {
		return []Change{}, false
	}

	nextChanges := []Change{cq.agenda[0]}
	cq.agenda = cq.agenda[1:]

	// pop any other changes with the same time
	for len(cq.agenda) > 0 && cq.agenda[0].Time == nextChanges[0].Time {
		nextChanges = append(nextChanges, cq.agenda[0])
		cq.agenda = cq.agenda[1:]
	}

	return nextChanges, true
}
