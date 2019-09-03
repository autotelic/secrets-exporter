package main

import (
	"io"
  "strings"
  "fmt"
  "log"
  "os"
  "encoding/json"
  "flag"
)

const EXPORT_TF_PREFIX = "export TF_VAR_"
const FILENAME = ".envrc"

func checkStdin() {
	fi, err := os.Stdin.Stat()
  if err != nil {
    panic(err)
  }
  if fi.Mode() & os.ModeNamedPipe == 0 {
  	os.Exit(0)
  }
}

func getSecrets(input io.Reader) map[string]interface{} {
	var secrets = make(map[string]interface{})

  jsonErr := json.NewDecoder(input).Decode(&secrets)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  return secrets;
}

func main() {
	checkStdin()

	filename := flag.String("filename", FILENAME, "The file the secrets should be written to")
	flag.Parse()

  secrets := getSecrets(os.Stdin)

  if len(secrets) == 0 {
		os.Exit(0)
	}

  file, err := os.Create(*filename)
  if err != nil {
    log.Fatal("Cannot create file", err)
    return;
  }

  defer file.Close()

  for k, v := range secrets {
  	lowercaseKey := strings.ToLower(k)
  	prefix := strings.Join([]string{EXPORT_TF_PREFIX, lowercaseKey}, "")
	  fmt.Fprintf(file, "%s=%s\n", prefix, v)
	}
}
