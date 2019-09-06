package main

import (
  "encoding/base64"
  "encoding/json"
  "flag"
  "fmt"
	"io"
  "io/ioutil"
  "log"
  "os"
  "strings"
  "text/template"
)

type SecretData struct{
  Name string
  Secrets map[string]interface{}
}

type Formatter = func(map[string]interface{}, *os.File, map[string]string)

const FILENAME = ".envrc"
const TERRAFORM = "terraform"
const KUBERNETES = "kubernetes"

var exportTypes = []string {TERRAFORM, KUBERNETES}

func check(err error) {
  if err != nil {
    panic(err)
  }
}

func checkStdin() {
	fi, err := os.Stdin.Stat()
  if err != nil {
    panic(err)
  }
  if fi.Mode() & os.ModeNamedPipe == 0 {
  	os.Exit(0)
  }
}

func getTemplate() string {
   data, err := ioutil.ReadFile("./k8Secrets.tmpl")
   check(err)
   return string(data)
}

func getSecrets(input io.Reader) map[string]interface{} {
	var secrets = make(map[string]interface{})

  jsonErr := json.NewDecoder(input).Decode(&secrets)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  return secrets;
}

func encodeValues(secrets map[string]interface{}) map[string]interface{} {
  var encodedSecrets = make(map[string]interface{})

  for k, v := range secrets {
    lowercaseKey := strings.ToLower(k)
    str := base64.StdEncoding.EncodeToString([]byte(v.(string)))
    encodedSecrets[lowercaseKey] = str
  }

  return encodedSecrets;
}

func kubernetesSecrets(secrets map[string]interface{}, file *os.File, flags map[string]string) {
  secretTmpl := getTemplate()
  name := flags["name"]

  tmpl, err := template.New("secret").Parse(secretTmpl)

  encodedSecrets := encodeValues(secrets);

  data := SecretData{name, encodedSecrets};

  err = tmpl.Execute(file, data)
  check(err)
}

func terraformSecrets(secrets map[string]interface{}, file *os.File, flags map[string]string) {
  const EXPORT_TF_PREFIX = "export TF_VAR_"

  for k, v := range secrets {
    lowercaseKey := strings.ToLower(k)
    prefix := strings.Join([]string{EXPORT_TF_PREFIX, lowercaseKey}, "")
    fmt.Fprintf(file, "%s=%s\n", prefix, v)
  }
}

func getFormatter(exportType string) Formatter {
  return map[string]Formatter{
    TERRAFORM: terraformSecrets,
    KUBERNETES: kubernetesSecrets,
  }[exportType]
}

func main() {
	checkStdin()

	filename := flag.String("filename", FILENAME, "The filename")
  exportType := flag.String("export-type", TERRAFORM, "The export type")
  name := flag.String("name", "", "The name for the kubernetes secrets config")

  flag.Parse()

  secrets := getSecrets(os.Stdin)

  flags := map[string]string{
    "filename": *filename,
    "name": *name,
  }

  if len(secrets) == 0 {
		os.Exit(0)
	}

  file, err := os.Create(*filename)
  if err != nil {
    log.Fatal("Cannot create file", err)
    return;
  }
  defer file.Close()

  formatter := getFormatter(*exportType)

  if (formatter == nil) {
    fmt.Printf("Invalid export-type. Please use one of: %s\n", strings.Join(exportTypes, ", "))
    os.Exit(0)
  }

  formatter(secrets, file, flags)
}
