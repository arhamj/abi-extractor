package main

import (
	"fmt"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	// ContractAddressFlag provides the contract address
	ContractAddressFlag = &cli.StringFlag{
		Name:     "contract",
		Usage:    "Provide the contract address",
		Required: true,
	}
	// NodeRpcEndpointFlag provides a custom RPC endpoint
	NodeRpcEndpointFlag = &cli.StringFlag{
		Name:     "contract",
		Usage:    "Provide a custom RPC endpoint",
		Required: false,
	}
)

var (
	getCodeFlags = []cli.Flag{
		ContractAddressFlag,
	}
)

type app struct {
	logger       *zap.Logger
	chainGateway external.ChainGateway
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	a := app{
		logger:       zap.L(),
		chainGateway: external.NewChainGatewayWithOpts(),
	}
	zap.ReplaceGlobals(logger)
	cliApp := &cli.App{
		Usage: "fetch EVM bytecode given contract address",
		Commands: []*cli.Command{
			{
				Name:   "bytecode",
				Flags:  getCodeFlags,
				Action: a.GetByteCodeHandler,
			},
		},
	}

	err = cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (a app) GetByteCodeHandler(c *cli.Context) error {
	contract := c.String(ContractAddressFlag.Name)
	if NodeRpcEndpointFlag.IsSet() {
		endpoint := c.String(NodeRpcEndpointFlag.Name)
		a.chainGateway = external.NewChainGatewayWithOpts(external.WithEthEndpoint(endpoint))
	}
	resp, err := a.chainGateway.EthGetCode(contract)
	if err != nil {
		return err
	}
	fmt.Printf("Bytecode: %s\n", resp.Result)
	return nil
}
