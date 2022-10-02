package asm

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

type SolidityParser struct {
	Disassembler
}

type SolidityParserOpt func(parser *SolidityParser)

func WithDisassembler(d Disassembler) SolidityParserOpt {
	return func(parser *SolidityParser) {
		parser.Disassembler = d
	}
}

func NewSolidityParserWithOpts(opts ...SolidityParserOpt) SolidityParser {
	p := SolidityParser{}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func NewSolidityParser(evmCode []byte) (SolidityParser, error) {
	d, err := NewDisassembler(evmCode)
	if err != nil {
		return SolidityParser{}, err
	}
	return SolidityParser{Disassembler: d}, nil
}

// GetFunctionSignatures For Solidity we are looking for the following pattern in the Bytecode
//
//	DUP1 PUSH4 <BYTE4> EQ PUSH1 <BYTE1> JUMPI
//	Ref: https://github.com/ethereum/solidity/blob/242096695fd3e08cc3ca3f0a7d2e06d09b5277bf/libsolidity/codegen/ContractCompiler.cpp#L333
func (p SolidityParser) GetFunctionSignatures() FunctionSignatures {
	functionSelectors := make([]string, 0)
	if len(p.Instructions) == 0 {
		return NewFunctionSignatures(functionSelectors)
	}
	for i := 0; i < len(p.Instructions); i++ {
		inst := p.Instructions[i]
		if inst.Op == vm.JUMPI {
			var (
				// Jump destination
				dest *big.Int
				sign string
			)

			// Validate i-1 to be a valid PUSH
			{
				prevInst := p.Instructions[i-1]
				if prevInst.Op.IsPush() && len(prevInst.Arg) > 0 {
					dest = big.NewInt(0)
					dest.SetBytes(prevInst.Arg)
				} else {
					continue
				}
			}

			// Validate i-2 to be a valid EQ
			{
				prevInst := p.Instructions[i-2]
				if prevInst.Op != vm.EQ {
					continue
				}
			}

			// Validate i-3 to be a PUSH4
			{
				prevInst := p.Instructions[i-3]
				if prevInst.Op == vm.PUSH4 && len(prevInst.Arg) > 0 {
					sign = hexutil.Encode(prevInst.Arg)
				} else {
					continue
				}
			}

			// Validate i-4 to be a DUP1
			{
				prePrev := p.Instructions[i-4]
				if prePrev.Op != vm.DUP1 {
					continue
				}
			}

			instAtJumpDest := p.GetOperationAtOffset(dest)
			if instAtJumpDest.Op != vm.JUMPDEST {
				continue
			}

			functionSelectors = append(functionSelectors, sign)
		}
	}
	return NewFunctionSignatures(functionSelectors)
}

// GetEventSignatures For Solidity we are looking for the following pattern in the Bytecode
//
//	PUSH32 <BYTE32>
//
//	Note: This is not the correct way to extract event hashes. Correct by building the stack step by step and
//	analyse it when a LOG instruction is encountered
func (p SolidityParser) GetEventSignatures() EventSignatures {
	eventSignatures := make([]string, 0)
	if len(p.Instructions) == 0 {
		return NewEventSignatures(eventSignatures)
	}
	var (
		push32Found bool
		push32Index int
		interSign   string
	)
	for i := 0; i < len(p.Instructions); i++ {
		inst := p.Instructions[i]
		if inst.Op == vm.PUSH32 && len(inst.Arg) > 0 {
			interSign = hexutil.Encode(inst.Arg)
			push32Found = true
			push32Index = i
		}
		if push32Found && (inst.Op == vm.LOG0 || inst.Op == vm.LOG1 ||
			inst.Op == vm.LOG2 || inst.Op == vm.LOG3 || inst.Op == vm.LOG4) {
			if i-push32Index > 25 {
				push32Found = false
				continue
			}
			eventSignatures = append(eventSignatures, interSign)
			push32Found = false
		}
	}
	return NewEventSignatures(eventSignatures)
}
