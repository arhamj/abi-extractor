package external

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	samczsunBaseUrl = "https://sig.eth.samczsun.com"
)

type SamczsunResp struct {
	Ok     bool `json:"ok"`
	Result struct {
		Event    map[string][]SamczsunTextResult `json:"event"`
		Function map[string][]SamczsunTextResult `json:"function"`
	} `json:"result"`
}

type SamczsunTextResult struct {
	Filtered bool   `json:"filtered"`
	Name     string `json:"name"`
}

type SamczsunGateway struct {
	logger     *zap.Logger
	httpclient *resty.Client
}

func NewSamczsunGateway(logger *zap.Logger) SamczsunGateway {
	return SamczsunGateway{
		logger:     logger.With(zap.String("loc", "SamczsunGateway")),
		httpclient: resty.New(),
	}
}

func (g *SamczsunGateway) GetEventTextSignature(eventSign string) (*SamczsunResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"event": eventSign,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&SamczsunResp{}).
		Get(samczsunBaseUrl + "/api/v1/signatures")
	if err != nil {
		g.logger.Error("Error in GetEventTextSignature", zap.String("sign", eventSign), zap.Error(err))
		return nil, errors.New("error when fetching event text signature")
	}
	return resp.Result().(*SamczsunResp), nil
}

func (g *SamczsunGateway) GetFunctionTextSignature(functionSign string) (*SamczsunResp, error) {
	resp, err := g.httpclient.R().
		SetQueryParams(map[string]string{
			"function": functionSign,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&SamczsunResp{}).
		Get(samczsunBaseUrl + "/api/v1/signatures")
	if err != nil {
		g.logger.Error("Error in GetFunctionTextSignature", zap.String("sign", functionSign), zap.Error(err))
		return nil, errors.New("error when fetching function text signature")
	}
	return resp.Result().(*SamczsunResp), nil
}
