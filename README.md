# Secrets Exporter

# Setup

You will need to have [go installed](https://golang.org/doc/install)

## Running

To build a standalone binary: `go build main.go`

To run: `go run main.go`

## Usage

Secrets exporter will read JSON data from the Stdin, parse, format and write
the results to a file in `dotenv` format. Currently this package only formats variable for consumption as [terraform environment variables](https://www.terraform.io/docs/configuration/variables.html#environment-variables).

### Flags

- `filename`: The filename the variables will be written to. (defaults to `.envrc`)

### Example

With a JSON input of
```json
"{"HELLO_SECRET":"Bar"}"
```

The resulting output will be a `.envrc` file with the contents:

```bash
export TF_VAR_hello_secret=bar
```

```bash
$ echo '{"HELLO_SECRET":"Bar"}"' | go run main.go
```

Usage with chamber:

```bash
$ chamber export <service_name> | go run main.go --filename='.env'
```
