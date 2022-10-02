package service

import (
	"errors"
	"github.com/arhamj/abi-extractor/pkg/external"
	"go.uber.org/zap"
)

type SignDecoderService struct {
	logger         *zap.Logger
	decoderGateway external.SamczsunGateway
}

type TextSignature struct {
	Sign     string
	Verified bool
}

func NewSignDecoder(decoderGateway external.SamczsunGateway) SignDecoderService {
	return SignDecoderService{
		logger:         zap.L().With(zap.String("loc", "SignDecoderService")),
		decoderGateway: decoderGateway,
	}
}

func (s SignDecoderService) GetEventTextSignature(eventSign string) (*TextSignature, error) {
	resp, err := s.decoderGateway.GetEventTextSignature(eventSign)
	if err != nil {
		return nil, err
	}
	result, ok := resp.Result.Event[eventSign]
	if !ok || len(result) == 0 {
		s.logger.Error("GetEventTextSignature: event text signature not found", zap.String("sign", eventSign))
		return nil, errors.New("text signature not found for event")
	}
	return &TextSignature{
		Sign: result[0].Name,
		// As per documentation, Filtered field in response is true when the obtained result is likely a spam
		Verified: !result[0].Filtered,
	}, nil
}

func (s SignDecoderService) GetFunctionTextSignature(functionSign string) (*TextSignature, error) {
	resp, err := s.decoderGateway.GetFunctionTextSignature(functionSign)
	if err != nil {
		return nil, err
	}
	result, ok := resp.Result.Function[functionSign]
	if !ok || len(result) == 0 {
		s.logger.Error("GetFunctionTextSignature: function text signature not found", zap.String("sign", functionSign))
		return nil, errors.New("text signature not found for function")
	}
	return &TextSignature{
		Sign: result[0].Name,
		// As per documentation, Filtered field in response is true when the obtained result is likely a spam
		Verified: !result[0].Filtered,
	}, nil
}
