package evaluator

type Callback interface {
	_cb()
}

type Return struct {
	Value OBJECT
}
type Break struct{}
type Continue struct{}

func (s Return) _cb()   {}
func (s Break) _cb()    {}
func (s Continue) _cb() {}
