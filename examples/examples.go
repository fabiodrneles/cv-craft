// Package examples embute os YAMLs de exemplo no binário; o modelo mínimo é
// usado pelo comando init.
package examples

import _ "embed"

// Minimal é o modelo gerado por "cv-craft init".
//
//go:embed minimal.yaml
var Minimal []byte
