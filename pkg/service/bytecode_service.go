package service

import (
	"github.com/arhamj/abi-extractor/pkg/asm"
	"go.uber.org/zap"
)

type BytecodeService struct {
	logger         *zap.Logger
	bytecodeParser asm.BytecodeParser
	signDecoder    SignDecoderService
}

func NewBytecodeService(bytecodeParser asm.BytecodeParser, signDecoder SignDecoderService) BytecodeService {
	return BytecodeService{
		logger:         zap.L().With(zap.String("loc", "BytecodeService")),
		bytecodeParser: bytecodeParser,
		signDecoder:    signDecoder,
	}
}

func (b BytecodeService) GetFunctionSigns() asm.FunctionSigns {
	return b.bytecodeParser.GetFunctionSigns()
}

func (b BytecodeService) GetEventSigns() asm.EventSigns {
	return b.bytecodeParser.GetEventSigns()
}

func (b BytecodeService) GetDecodedFunctionSigns() map[string]string {
	res := make(map[string]string, 0)
	signs := b.bytecodeParser.GetFunctionSigns()
	for sign := range signs.Signatures {
		textSign, err := b.signDecoder.GetFunctionTextSignature(sign)
		if err != nil || !textSign.Verified {
			b.logger.Debug("text sign not found for function", zap.String("sign", sign))
			continue
		}
		res[sign] = textSign.Sign
	}
	return res
}

func (b BytecodeService) GetDecodedEventSigns() map[string]string {
	res := make(map[string]string, 0)
	signs := b.bytecodeParser.GetEventSigns()
	for sign := range signs.Signatures {
		textSign, err := b.signDecoder.GetEventTextSignature(sign)
		if err != nil || !textSign.Verified {
			b.logger.Debug("text sign not found for event", zap.String("sign", sign))
			continue
		}
		res[sign] = textSign.Sign
	}
	return res
}

func (b BytecodeService) GetABI() {

}
