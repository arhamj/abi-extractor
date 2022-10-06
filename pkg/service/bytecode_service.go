package service

import (
	"github.com/arhamj/abi-extractor/pkg/asm"
	"go.uber.org/zap"
	"sync"
)

type BytecodeService struct {
	logger      *zap.Logger
	signDecoder SignDecoderService
}

func NewBytecodeService(signDecoder SignDecoderService) BytecodeService {
	return BytecodeService{
		logger:      zap.L().With(zap.String("loc", "BytecodeService")),
		signDecoder: signDecoder,
	}
}

func (b BytecodeService) GetFunctionSigns(bytecodeParser asm.BytecodeParser) asm.FunctionSigns {
	return bytecodeParser.GetFunctionSigns()
}

func (b BytecodeService) GetEventSigns(bytecodeParser asm.BytecodeParser) asm.EventSigns {
	return bytecodeParser.GetEventSigns()
}

func (b BytecodeService) GetDecodedFunctionSigns(bytecodeParser asm.BytecodeParser) map[string]string {
	res := make(map[string]string, 0)
	writeLock := new(sync.Mutex)
	signs := bytecodeParser.GetFunctionSigns()
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

func (b BytecodeService) GetDecodedEventSigns(bytecodeParser asm.BytecodeParser) map[string]string {
	res := make(map[string]string, 0)
	writeLock := new(sync.Mutex)
	signs := bytecodeParser.GetEventSigns()
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

func (b BytecodeService) GetABI(bytecodeParser asm.BytecodeParser) string {
	//events := b.GetEventSigns(bytecodeParser)
	//functions := b.GetDecodedFunctionSigns(bytecodeParser)
	return ""
}
