package external

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	defaultEthEndpoint = "https://rpc.ankr.com/eth"
)

type ChainGateway struct {
	ethEndpoint string
	logger      *zap.Logger
	httpclient  *resty.Client
}

type EthReq struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type EthCodeResp struct {
	Result string `json:"result"`
}

type ChainGatewayOpt func(gateway *ChainGateway)

func WithEthEndpoint(endpoint string) func(gateway *ChainGateway) {
	return func(gateway *ChainGateway) {
		gateway.ethEndpoint = endpoint
	}
}

func NewChainGatewayWithOpts(opts ...ChainGatewayOpt) ChainGateway {
	chainGateway := ChainGateway{
		ethEndpoint: defaultEthEndpoint,
		logger:      zap.L().With(zap.String("loc", "ChainGateway")),
		httpclient:  resty.New(),
	}
	for _, opt := range opts {
		opt(&chainGateway)
	}
	return chainGateway
}

func (g ChainGateway) EthGetCode(contract string) (*EthCodeResp, error) {
	req := EthReq{
		Jsonrpc: "2.0",
		Method:  "eth_getCode",
		Params:  []string{contract, "latest"},
		Id:      1,
	}
	resp, err := g.httpclient.R().
		SetBody(req).
		SetHeader("Accept", "application/json").
		SetResult(&EthCodeResp{}).
		Post(g.ethEndpoint)
	if err != nil {
		g.logger.Error("EthGetCode: error making RPC call", zap.String("contract", contract), zap.Error(err))
		return nil, errors.New("error when fetching bytecode for contract")
	}
	return resp.Result().(*EthCodeResp), nil
}
