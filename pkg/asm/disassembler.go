package asm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

// Instruction is to encapsulate an EVM instruction
type Instruction struct {
	PC  uint64
	Op  vm.OpCode
	Arg []byte
}

type Disassembler struct {
	OriginalByteCode []byte
	Instructions     []Instruction
	// JumpDestinations map of Program Counter to Instruction id in Instructions
	JumpDestinations map[uint64]int
}

// NewDisassembler accepts the EVM bytecode as a bytea
func NewDisassembler(evmCode []byte) (Disassembler, error) {
	d := Disassembler{
		OriginalByteCode: evmCode,
		Instructions:     make([]Instruction, 0),
		JumpDestinations: make(map[uint64]int, 0),
	}
	it := asm.NewInstructionIterator(evmCode)
	for it.Next() {
		d.Instructions = append(d.Instructions, Instruction{
			PC:  it.PC(),
			Op:  it.Op(),
			Arg: it.Arg(),
		})
		if it.Op() == vm.JUMPDEST {
			d.JumpDestinations[it.PC()] = len(d.Instructions) - 1
		}
	}
	return d, it.Error()
}

func (d Disassembler) GetOperationAtOffset(offset *big.Int) Instruction {
	defaultInst := Instruction{
		Op:  vm.STOP,
		Arg: nil,
	}

	// offset exceeds the length of the byte code
	if offset.Cmp(big.NewInt(int64(len(d.OriginalByteCode)))) > 1 {
		return defaultInst
	}
	instNumber := d.JumpDestinations[offset.Uint64()]
	return d.Instructions[instNumber]
}

func (d Disassembler) PrintDisassembled() {
	for _, inst := range d.Instructions {
		if inst.Arg != nil && 0 < len(inst.Arg) {
			fmt.Printf("%05x: %v %#x\n", inst.PC, inst.Op, inst.Arg)
		} else {
			fmt.Printf("%05x: %v\n", inst.PC, inst.Op)
		}
	}
}
