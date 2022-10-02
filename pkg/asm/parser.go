package asm

type BytecodeParser interface {
	// GetFunctionSigns returns a list of signatures prefixed with 0x
	GetFunctionSigns() FunctionSigns
	// GetEventSigns returns a list of event signatures prefixed with 0x
	GetEventSigns() EventSigns
}

type FunctionSigns struct {
	Signatures map[string]bool
}

func NewFunctionSigns(signs []string) FunctionSigns {
	signMap := make(map[string]bool)
	for _, s := range signs {
		signMap[s] = true
	}
	return FunctionSigns{Signatures: signMap}
}

func (f FunctionSigns) List() []string {
	res := make([]string, 0)
	if len(f.Signatures) == 0 {
		return res
	}
	for k := range f.Signatures {
		res = append(res, k)
	}
	return res
}

type EventSigns struct {
	Signatures map[string]bool
}

func NewEventSigns(signs []string) EventSigns {
	signMap := make(map[string]bool)
	for _, s := range signs {
		signMap[s] = true
	}
	return EventSigns{Signatures: signMap}
}

func (e EventSigns) List() []string {
	res := make([]string, 0)
	if len(e.Signatures) == 0 {
		return res
	}
	for k := range e.Signatures {
		res = append(res, k)
	}
	return res
}
