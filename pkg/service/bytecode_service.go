package service

import "github.com/arhamj/abi-extractor/pkg/asm"

type BytecodeService struct {
	bytecodeParser asm.BytecodeParser
	signDecoder    SignDecoderService
}

func NewBytecodeService(bytecodeParser asm.BytecodeParser, signDecoder SignDecoderService) BytecodeService {
	return BytecodeService{
		bytecodeParser: bytecodeParser,
		signDecoder:    signDecoder,
	}
}

func (b BytecodeService) GetFunctionSigns() {

}

func (b BytecodeService) GetEventSign() {

}

func (b BytecodeService) GetDecodedFunctionSigns() {

}

func (b BytecodeService) GetDecodedEventSigns() {

}

func (b BytecodeService) GetABI() {

}
