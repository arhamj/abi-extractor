package service

import (
	"github.com/arhamj/abi-extractor/pkg/asm"
	"go.uber.org/zap"
	"sync"
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
	writeLock := new(sync.Mutex)
	signs := b.bytecodeParser.GetFunctionSigns()
	wg := new(sync.WaitGroup)
	for sign := range signs.Signatures {
		wg.Add(1)
		go func(sign string) {
			defer wg.Done()
			textSign, err := b.signDecoder.GetFunctionTextSignature(sign)
			if err != nil || !textSign.Verified {
				b.logger.Debug("text sign not found for function", zap.String("sign", sign))
				return
			}
			writeLock.Lock()
			res[sign] = textSign.Sign
			writeLock.Unlock()
		}(sign)
	}
	wg.Wait()
	return res
}

func (b BytecodeService) GetDecodedEventSigns() map[string]string {
	res := make(map[string]string, 0)
	writeLock := new(sync.Mutex)
	signs := b.bytecodeParser.GetEventSigns()
	wg := new(sync.WaitGroup)
	for sign := range signs.Signatures {
		wg.Add(1)
		func(sign string) {
			defer wg.Done()
			textSign, err := b.signDecoder.GetEventTextSignature(sign)
			if err != nil || !textSign.Verified {
				b.logger.Debug("text sign not found for event", zap.String("sign", sign))
				return
			}
			writeLock.Lock()
			res[sign] = textSign.Sign
			writeLock.Unlock()
		}(sign)
	}
	wg.Wait()
	return res
}

func (b BytecodeService) GetABI() {

}
