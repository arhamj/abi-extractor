package main

import (
	"fmt"
	"github.com/arhamj/abi-extractor/pkg/asm"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/arhamj/abi-extractor/pkg/service"
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
		Name:     "node",
		Usage:    "Provide a custom RPC endpoint",
		Required: false,
	}
)

var (
	defaultFlags = []cli.Flag{
		ContractAddressFlag,
		NodeRpcEndpointFlag,
	}
)

type app struct {
	logger       *zap.Logger
	chainGateway external.ChainGateway

	bytecodeService service.BytecodeService

	bytecode *external.EthCodeResp
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	a := app{}
	zap.ReplaceGlobals(logger)
	cliApp := &cli.App{
		Usage: "Bytecode disassembler CLI",
		Commands: []*cli.Command{
			{
				Name:        "bytecode",
				Aliases:     []string{"bc"},
				Flags:       defaultFlags,
				Description: "fetch the bytecode of a smart contract from node",
				Action:      a.PrintByteCodeHandler,
			},
			{
				Name:        "hex-events",
				Aliases:     []string{"he"},
				Description: "extract the event signature from contract bytecode",
				Flags:       defaultFlags,
				Action:      a.PrintEventSignatures,
			},
			{
				Name:        "hex-functions",
				Aliases:     []string{"hf"},
				Description: "extract the function signature from contract bytecode",
				Flags:       defaultFlags,
				Action:      a.PrintFunctionSignatures,
			},
			{
				Name:        "text-events",
				Aliases:     []string{"te"},
				Description: "extract the event signature (in text) from contract bytecode",
				Flags:       defaultFlags,
				Action:      a.PrintDecodedEventsSignatures,
			},
			{
				Name:        "text-functions",
				Aliases:     []string{"tf"},
				Description: "extract the function signature (in text) from contract bytecode",
				Flags:       defaultFlags,
				Action:      a.PrintDecodedFunctionsSignatures,
			},
		},
	}

	err = cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *app) setupApp(c *cli.Context) error {
	contract := c.String(ContractAddressFlag.Name)
	a.chainGateway = external.NewChainGatewayWithOpts()
	if NodeRpcEndpointFlag.IsSet() {
		endpoint := c.String(NodeRpcEndpointFlag.Name)
		a.chainGateway = external.NewChainGatewayWithOpts(external.WithEthEndpoint(endpoint))
	}
	resp, err := a.chainGateway.EthGetCode(contract)
	if err != nil {
		return err
	}
	a.logger = zap.L()
	signDecoder := service.NewSignDecoder(external.NewSamczsunGateway())
	parser, err := asm.NewSolidityParserStr(resp.Result[2:])
	if err != nil {
		return err
	}
	a.bytecodeService = service.NewBytecodeService(parser, signDecoder)
	a.bytecode = resp
	return nil
}

func (a *app) PrintByteCodeHandler(c *cli.Context) error {
	err := a.setupApp(c)
	if err != nil {
		return err
	}
	fmt.Printf("\nBytecode: %s\n", a.bytecode.Result)
	return nil
}

func (a *app) PrintEventSignatures(c *cli.Context) error {
	err := a.setupApp(c)
	if err != nil {
		return err
	}

	res := a.bytecodeService.GetEventSigns()
	fmt.Println("\nEvent signatures (in hex):")
	for _, k := range res.List() {
		fmt.Printf("- %s\n", k)
	}
	return nil
}

func (a *app) PrintFunctionSignatures(c *cli.Context) error {
	err := a.setupApp(c)
	if err != nil {
		return err
	}
	res := a.bytecodeService.GetFunctionSigns()
	fmt.Println("\nFunction signatures (in hex):")
	for _, k := range res.List() {
		fmt.Printf("- %s\n", k)
	}
	return nil
}

func (a *app) PrintDecodedEventsSignatures(c *cli.Context) error {
	err := a.setupApp(c)
	if err != nil {
		return err
	}
	res := a.bytecodeService.GetDecodedEventSigns()
	fmt.Println("\nEvent signatures (<in hex>: <in text>):")
	for hexSign, textSign := range res {
		fmt.Printf("- %s: %s\n", hexSign, textSign)
	}
	return nil
}

func (a *app) PrintDecodedFunctionsSignatures(c *cli.Context) error {
	err := a.setupApp(c)
	if err != nil {
		return err
	}
	res := a.bytecodeService.GetDecodedFunctionSigns()
	fmt.Println("\nFunction signatures (<in hex>: <in text>):")
	for hexSign, textSign := range res {
		fmt.Printf("- %s: %s\n", hexSign, textSign)
	}
	return nil
}
