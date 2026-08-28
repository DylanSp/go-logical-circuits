package wire

type Signal bool

const (
	Low  Signal = false
	High Signal = true
)

type Wire struct {
	name   string
	signal Signal
}

func New(name string) Wire {
	return Wire{
		name:   name,
		signal: Low,
	}
}

func (w Wire) Name() string {
	return w.name
}

// implement stringer interface
func (w Wire) String() string {
	return w.Name()
}

func (w Wire) Signal() Signal {
	return w.signal
}

func (w *Wire) SetLow() {
	w.signal = Low
}

func (w *Wire) SetHigh() {
	w.signal = High
}
