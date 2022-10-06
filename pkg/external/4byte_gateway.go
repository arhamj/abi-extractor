package external

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

const (
	fourBytesBaseUrl = "https://www.4byte.directory"
)

type FourBytesResp struct {
	Count   int              `json:"count"`
	Results []TextSignResult `json:"results"`
}

type TextSignResult struct {
	CreatedAt     time.Time `json:"created_at"`
	TextSignature string    `json:"text_signature"`
	HexSignature  string    `json:"hex_signature"`
}

type FourByteGateway struct {
	logger     *zap.Logger
	httpclient *resty.Client
}

func NewFourByteGateway() FourByteGateway {
	return FourByteGateway{
		logger:     zap.L().With(zap.String("loc", "FourByteGateway")),
		httpclient: resty.New(),
	}
}

func (g *FourByteGateway) GetEventTextSignature(eventSign string) (*FourBytesResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"hex_signature": eventSign,
			"sort":          "id",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&FourBytesResp{}).
		Get(fourBytesBaseUrl + "/api/v1/event-signatures/")
	if err != nil {
		g.logger.Error("GetEventTextSignature: error making call to 4byte", zap.String("sign", eventSign), zap.Error(err))
		return nil, errors.New("error when fetching event text signature")
	}
	return resp.Result().(*FourBytesResp), nil
}

func (g *FourByteGateway) GetFunctionTextSignature(functionSign string) (*FourBytesResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"hex_signature": functionSign,
			"sort":          "id",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&FourBytesResp{}).
		Get(fourBytesBaseUrl + "/api/v1/signatures/")
	if err != nil {
		g.logger.Error("GetFunctionTextSignature: error making call to 4byte", zap.String("sign", functionSign), zap.Error(err))
		return nil, errors.New("error when fetching function text signature")
	}
	return resp.Result().(*FourBytesResp), nil
}

func (g *FourByteGateway) GetFunctionSignatures(pageNo int) (*FourBytesResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"page":     fmt.Sprintf("%d", pageNo),
			"ordering": "created_at",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&FourBytesResp{}).
		Get(fourBytesBaseUrl + "/api/v1/signatures/")
	if err != nil {
		g.logger.Error("GetFunctionSignatures: error making call to 4byte", zap.Int("page", pageNo), zap.Error(err))
		return nil, errors.New("error when fetching function text signature")
	}
	return resp.Result().(*FourBytesResp), nil
}

func (g *FourByteGateway) GetEventSignatures(pageNo int) (*FourBytesResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"page":     fmt.Sprintf("%d", pageNo),
			"ordering": "created_at",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&FourBytesResp{}).
		Get(fourBytesBaseUrl + "/api/v1/event-signatures/")
	if err != nil {
		g.logger.Error("GetEventSignatures: error making call to 4byte", zap.Int("page", pageNo), zap.Error(err))
		return nil, errors.New("error when fetching event text signature")
	}
	return resp.Result().(*FourBytesResp), nil
}
