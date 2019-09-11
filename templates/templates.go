package templates

const KubernetesSecretTmpl = `
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Name }}
type: Opaque
data:{{ range $key, $value := .Secrets }}
  {{ $key }}: {{ $value }}{{ end }}
`
