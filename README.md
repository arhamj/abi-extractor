# ABI Extractor

A Golang based SDK + CLI tool to extract, decode, and generate ABI from the EVM bytecode.

## Note

- The logic for the extraction of function hashes and the event hashes is Solidity specific. But the application can be
  extended for other cases.
- The event extraction logic is temporary and needs a revisit

## CLI usage

### Requirements

- Go 1.19+

### Installation

```
make install
```

### Overview

```
NAME:abi-extractor                                                                                                                                                       ─╯
   abi-extractor - Bytecode disassembler CLI

USAGE:
   abi-extractor [global options] command [command options] [arguments...]

COMMANDS:
   bytecode, bc              
   hex-events, he            
   hex-functions, hf         
   text-events, te           
   text-functions, tf        
   decode-hex-event, dhe     
   decode-hex-function, dhf  
   sync-4byte-events, s4e    
   sync-4byte-function, s4f  
   help, h                   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

To get the detailed usage information for a specific command enter

```
abi-extractor <command_name> 
```

Example

```
NAME:abi-extractor hex-events                                                                                                                                            ─╯
   abi-extractor hex-events

USAGE:
   abi-extractor hex-events [command options] [arguments...]

DESCRIPTION:
   extract the event signature from contract bytecode

OPTIONS:
   --contract value  Provide the contract address
   --node value      Provide a custom RPC endpoint
```

## SDK Usage

### Installation

```
go get github.com/arhamj/abi-extractor
```

### Usage example

```go
package main

import (
	"github.com/arhamj/abi-extractor/pkg/asm"
	"github.com/arhamj/abi-extractor/pkg/service"
)

func main() {
	bytecodeService := service.NewDefaultBytecodeService()

	// Needs to be initialised for every new contract
	parser, err := asm.NewSolidityParserStr(byteCodeString)
	if err != nil {
		panic(err)
	}

	// Fetch decode information
	//	Get signatures in hex
	bytecodeService.GetFunctionSigns(parser)
	bytecodeService.GetEventSigns(parser)
	//	Get signatures in text (best-effort)
	//	(text signature is returned based on scraped data from 4byte(https://www.4byte.directory/) or Eth Sign Database(https://sig.eth.samczsun.com/)
	bytecodeService.GetDecodedEventSigns(parser)
	bytecodeService.GetDecodedFunctionSigns(parser)
}
```

### Advanced users

- If your usage is frequent we recommend you to scrape signature data in a local SQLite DB for better speeds

```
>> abi-extractor sync-4byte-function
>> abi-extractor sync-4byte-events
```

- If the sync is time-consuming you can use the following sync backup and save it as `/db/scraper.db` in the project
  directory
    - [backup]()