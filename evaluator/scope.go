package evaluator

type Scope struct {
	current map[string]OBJECT
	prev    *Scope
}

func CreateScope() *Scope {
	return &Scope{
		current: make(map[string]OBJECT),
		prev:    nil,
	}
}

func (scope *Scope) Child() *Scope {
	return &Scope{
		current: make(map[string]OBJECT),
		prev:    scope,
	}
}

func (scope *Scope) Get(identifier string) OBJECT {
	if obj, ok := scope.current[identifier]; ok {
		return obj
	}

	if scope.prev != nil {
		return scope.prev.Get(identifier)
	}

	return Null_
}

func (scope *Scope) Set(identifier string, value OBJECT) bool {
	if _, ok := scope.current[identifier]; ok {
		scope.current[identifier] = value
		return true
	}

	if scope.prev != nil {
		return scope.prev.Set(identifier, value)
	}

	return false
}

func (scope *Scope) Init(identifier string, value OBJECT) {
	scope.current[identifier] = value
}

func (scope *Scope) SetOrInit(identifier string, value OBJECT) {
	if !scope.Set(identifier, value) {
		scope.Init(identifier, value)
	}
}
