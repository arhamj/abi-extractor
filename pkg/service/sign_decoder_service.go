package service

import (
	"database/sql"
	"errors"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/arhamj/abi-extractor/pkg/scraper"
	"go.uber.org/zap"
)

type SignDecoderService struct {
	logger         *zap.Logger
	scraperDb      *sql.DB
	decoderGateway external.SamczsunGateway
}

type TextSignature struct {
	Sign     string
	Verified bool
}

func NewSignDecoder(scraperDb *sql.DB, decoderGateway external.SamczsunGateway) SignDecoderService {
	return SignDecoderService{
		logger:         zap.L().With(zap.String("loc", "SignDecoderService")),
		scraperDb:      scraperDb,
		decoderGateway: decoderGateway,
	}
}

func (s SignDecoderService) GetEventTextSignature(eventSign string) (*TextSignature, error) {
	textSignFromDb, err := s.fetchTextSignatureFromDb(scraper.Event, eventSign)
	if err == nil {
		s.logger.Debug("GetEventTextSignature: text sign fetched from db", zap.String("sign", eventSign))
		return textSignFromDb, nil
	}
	resp, err := s.decoderGateway.GetEventTextSignature(eventSign)
	if err != nil {
		return nil, err
	}
	result, ok := resp.Result.Event[eventSign]
	if !ok || len(result) == 0 {
		s.logger.Debug("GetEventTextSignature: event text signature not found", zap.String("sign", eventSign))
		return nil, errors.New("text signature not found for event")
	}
	return &TextSignature{
		Sign: result[0].Name,
		// As per documentation, Filtered field in response is true when the obtained result is likely a spam
		Verified: !result[0].Filtered,
	}, nil
}

func (s SignDecoderService) GetFunctionTextSignature(functionSign string) (*TextSignature, error) {
	textSignFromDb, err := s.fetchTextSignatureFromDb(scraper.Function, functionSign)
	if err == nil {
		s.logger.Debug("GetFunctionTextSignature: text sign fetched from db", zap.String("sign", functionSign))
		return textSignFromDb, nil
	}
	resp, err := s.decoderGateway.GetFunctionTextSignature(functionSign)
	if err != nil {
		return nil, err
	}
	result, ok := resp.Result.Function[functionSign]
	if !ok || len(result) == 0 {
		s.logger.Debug("GetFunctionTextSignature: function text signature not found", zap.String("sign", functionSign))
		return nil, errors.New("text signature not found for function")
	}
	return &TextSignature{
		Sign: result[0].Name,
		// As per documentation, Filtered field in response is true when the obtained result is likely a spam
		Verified: !result[0].Filtered,
	}, nil
}

func (s SignDecoderService) fetchTextSignatureFromDb(kind scraper.MappingKind, hexSign string) (*TextSignature, error) {
	var textSign string
	row := s.scraperDb.QueryRow("SELECT string_sign FROM sign_mapping_fourbyte WHERE kind = ? AND hex_sign = ? ORDER BY created_at limit 1", kind, hexSign)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(&textSign)
	if err != nil {
		return nil, err
	}
	return &TextSignature{
		Sign:     textSign,
		Verified: true,
	}, nil
}
