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
- `export-type`: The format to export the secrets. (Defaults to `terraform`)
- `name`: The name of the secrets object (kubernetes only)

### Terraform

#### Example

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

### Kubernetes

#### Example

The exporter also supports generation of Kubernetes secrets objects.

```bash
$ echo '{"HELLO_SECRET":"world"}"' | go run main.go \
  --filename="Secret.yaml" \
  --export-type="kubernetes" \
  --name="demosecret"
```

Will create the following Kubernetes secret which may be applied using `kubectl`.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: DemoSecret
type: Opaque
data:
  hello_secret: d29ybGQ=
```
