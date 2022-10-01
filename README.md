# ABI Extractor

## Goals

- To be able to fetch the Bytecode of a smart contract
- Parse the EVM bytecode
- Extract the method signatures
- Extract the event signatures
- Use 4bytes to check if they have been found
- Come up with a format to indicate the extracted signatures
- Expose a CLI which uses defaults

## High level commands for the CLI

- Fetch Bytecode
- Fetch function signatures
- Fetch event signatures
- Fetch decoded function signatures
- Fetch decoded event signatures
- Fetch best effort ABI
- Fetch completed decoded json

## SDK high level requirements

- All of the above CLI features must be available as methods
- The constructor of the SDK must be able to accept the endpoint

## Quality requirements

- Avoid using geth as a library and rely on a standard HTTP call. To support a variety of chains
- Performing contract address validation should be optional at an SDK level and should happen by default for the CLI

## Project planning

- Build the core functionality and the SDK first
- Have good quality tests
- Document all the useful resources in the `Resources` section

## Existing projects

- https://github.com/shazow/whatsabi
  * Extracts methods and decodes them but not events

## Resources

