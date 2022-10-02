package service

import "github.com/arhamj/abi-extractor/pkg/asm"

type BytecodeService struct {
	solidityParser asm.SolidityParser
	signDecoder    SignDecoderService
}

func (b BytecodeService) GetBytecode() {

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
