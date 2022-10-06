package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arhamj/abi-extractor/pkg/asm"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/arhamj/abi-extractor/pkg/scraper"
	"github.com/arhamj/abi-extractor/pkg/service"
	"github.com/arhamj/abi-extractor/pkg/util"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	// HexStringFlag provides a custom RPC endpoint
	HexStringFlag = &cli.StringFlag{
		Name:     "hex",
		Usage:    "Provide the hex string flag to be decoded",
		Required: true,
	}
)

var (
	defaultFlags = []cli.Flag{
		ContractAddressFlag,
		NodeRpcEndpointFlag,
	}
	hexFlags = []cli.Flag{
		HexStringFlag,
	}
)

type app struct {
	logger       *zap.Logger
	chainGateway external.ChainGateway

	scraperDb *sql.DB

	bytecodeParser asm.BytecodeParser

	bytecodeService service.BytecodeService
	signDecoder     service.SignDecoderService

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
			{
				Name:        "decode-hex-event",
				Aliases:     []string{"dhe"},
				Description: "extract the text signature of a given hex event",
				Flags:       hexFlags,
				Action:      a.PrintDecodedEventSignature,
			},
			{
				Name:        "decode-hex-function",
				Aliases:     []string{"dhf"},
				Description: "extract the text signature of a given hex function",
				Flags:       hexFlags,
				Action:      a.PrintDecodedFunctionSignature,
			},
			{
				Name:        "sync-4byte-events",
				Aliases:     []string{"s4e"},
				Description: "scrape 4byte database to local SQLite",
				Action:      func(c *cli.Context) error { return a.Sync4Byte(c, scraper.Event) },
			},
			{
				Name:        "sync-4byte-function",
				Aliases:     []string{"s4f"},
				Description: "scrape 4byte database to local SQLite",
				Action:      func(c *cli.Context) error { return a.Sync4Byte(c, scraper.Function) },
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
	parser, err := asm.NewSolidityParserStr(resp.Result[2:])
	if err != nil {
		return err
	}
	a.bytecodeParser = parser
	a.logger = zap.L()
	scraperDb, err := util.NewSQLiteDB("db/scraper.db", scraper.FourByteMigrations)
	if err != nil {
		return err
	}
	a.scraperDb = scraperDb
	a.signDecoder = service.NewSignDecoder(external.NewSamczsunGateway(), service.WithScraperDbOpt(scraperDb))
	a.bytecodeService = service.NewBytecodeService(a.signDecoder)
	return nil
}

func (a *app) setupAppWithoutContract(c *cli.Context) error {
	a.logger = zap.L()
	scraperDb, err := util.NewSQLiteDB("db/scraper.db", scraper.FourByteMigrations)
	if err != nil {
		return err
	}
	a.scraperDb = scraperDb
	a.signDecoder = service.NewSignDecoder(external.NewSamczsunGateway(), service.WithScraperDbOpt(scraperDb))
	a.bytecodeService = service.NewBytecodeService(a.signDecoder)
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

	res := a.bytecodeService.GetEventSigns(a.bytecodeParser)
	fmt.Println("\nEvent signatures (<in hex>):")
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
	res := a.bytecodeService.GetFunctionSigns(a.bytecodeParser)
	fmt.Println("\nFunction signatures (<in hex>):")
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
	res := a.bytecodeService.GetDecodedEventSigns(a.bytecodeParser)
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
	res := a.bytecodeService.GetDecodedFunctionSigns(a.bytecodeParser)
	fmt.Println("\nFunction signatures (<in hex>: <in text>):")
	for hexSign, textSign := range res {
		fmt.Printf("- %s: %s\n", hexSign, textSign)
	}
	return nil
}

func (a *app) PrintDecodedEventSignature(c *cli.Context) error {
	err := a.setupAppWithoutContract(c)
	if err != nil {
		return err
	}
	hexString := c.String(HexStringFlag.Name)
	res, err := a.signDecoder.GetEventTextSignature(hexString)
	if err != nil {
		return err
	}
	fmt.Println("\nEvent signatures (<in text>):", res.Sign)
	return nil
}

func (a *app) PrintDecodedFunctionSignature(c *cli.Context) error {
	err := a.setupAppWithoutContract(c)
	if err != nil {
		return err
	}
	hexString := c.String(HexStringFlag.Name)
	res, err := a.signDecoder.GetFunctionTextSignature(hexString)
	if err != nil {
		return err
	}
	fmt.Println("\nFunction signatures (<in text>):", res.Sign)
	return nil
}

func (a *app) Sync4Byte(c *cli.Context, kind scraper.MappingKind) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	util.SetupDevLogger()
	scraper4Byte, err := scraper.NewFourByteScraper(ctx, external.NewFourByteGateway())
	if err != nil {
		return err
	}
	err = scraper4Byte.Start(kind)
	if err != nil {
		return err
	}
	return nil
}
