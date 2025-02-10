# Renamatic

Renamatic is a CLI tool that automates the renaming of function calls in `.gno` files. Using Go's AST parser, it updates function names based on a user-supplied YAML mapping. It is designed primarily for transforming calls to functions from the `std` package.

## Features

- **Recursive Processing:** Scans directories recursively for `.gno` files.
- **YAML-Based Mapping:** Uses a YAML file to map old function names to new ones.
- **Chained Calls Support:** Handles renaming in chained function calls (e.g., `std.PrevRealm().Addr()` becomes `std.PrevRealm().Address()`).
- **AST-Based Transformation:** Leverages Go's AST parsing for reliable code transformations.

## Installation

Build the CLI tool using Go:

```bash
go build -o renamatic ./cmd/renamatic
```

Or run it directly:

```bash
go run ./cmd/renamatic -mapping=path/to/mapping.yaml -dir=path/to/your/files
```

## YAML Mapping Configuration

Create a YAML file (e.g., `mapping.yaml`) that specifies the function renaming rules. For example:

```yaml
GetCallerAt: CallerAt
GetOrigSend: OriginSend
origSend: originSend
GetOrigCaller: OriginCaller
origCaller: originCaller
Orig: Origin
Addr: Address
PrevRealm: PreviousRealm
GetChainID: ChainID
GetBanker: NewBanker
GetChainDomain: ChainDomain
GetHeight: Height
```

Each key is the original function name, and each value is the new function name.

## Usage

Run Renamatic by specifying the YAML mapping file and the target directory:

```bash
renamatic -mapping=mapping.yaml -dir=./your_gno_directory
```

Renamatic will recursively search for `.gno` files in the specified directory and apply the renaming rules defined in the YAML file.

## License

Renamatic is licensed under the [Apache License](LICENSE).
