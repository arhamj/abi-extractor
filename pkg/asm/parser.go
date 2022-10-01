package asm

type BytecodeParser interface {
	// GetFunctionSignatures returns a list of signatures prefixed with 0x
	GetFunctionSignatures() FunctionSignatures
	// GetEventSignatures returns a list of event signatures prefixed with 0x
	GetEventSignatures() EventSignatures
}

type FunctionSignatures struct {
	Signatures map[string]bool
}

func NewFunctionSignatures(signs []string) FunctionSignatures {
	signMap := make(map[string]bool)
	for _, s := range signs {
		signMap[s] = true
	}
	return FunctionSignatures{Signatures: signMap}
}

func (f FunctionSignatures) List() []string {
	res := make([]string, 0)
	if len(f.Signatures) == 0 {
		return res
	}
	for k := range f.Signatures {
		res = append(res, k)
	}
	return res
}

type EventSignatures struct {
	Signatures map[string]bool
}

func NewEventSignatures(signs []string) EventSignatures {
	signMap := make(map[string]bool)
	for _, s := range signs {
		signMap[s] = true
	}
	return EventSignatures{Signatures: signMap}
}

func (e EventSignatures) List() []string {
	res := make([]string, 0)
	if len(e.Signatures) == 0 {
		return res
	}
	for k := range e.Signatures {
		res = append(res, k)
	}
	return res
}
