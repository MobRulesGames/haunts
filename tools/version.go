// +build ignore

package main

import (
  "fmt"
  "text/template"
  "os"
  "io/ioutil"
  "path/filepath"
  "strings"
)

var outputTemplate = template.Must(template.New("output").Parse(outputTemplateStr))

const outputTemplateStr = `
package main
func Version() string {
  return "{{.}}"
}
`

func read(path string) (string, error) {
  f, err := os.Open(path)
  if err != nil {
    return "", err
  }
  data, err := ioutil.ReadAll(f)
  if err != nil {
    return "", err
  }
  return string(data), nil
}

func main() {
  head, err := read(filepath.Join("..", ".git", "HEAD"))
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }
  head = strings.TrimSpace(head)

  target := filepath.Join("..", "GEN_version.go")
  os.Remove(target) // Don't care about errors on this one
  f, err := os.Create(target)
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }
  outputTemplate.Execute(f, head)
  f.Close()
}
