package circuit

type ChangeQueue struct {
	// TODO - internal fields
}

func (cq *ChangeQueue) AddChange(ch Change) {
	panic("unimplemented")
}

// pops the change(s) with the smallest time
// if there are multiple changes with the same time - returns all of them
// if there are no queued changes - returns ([], false)
func (cq *ChangeQueue) GetNextChanges() ([]Change, bool) {
	panic("unimplemented")
}
