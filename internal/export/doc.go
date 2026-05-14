// Package export serialises decrypted environment variables into
// multiple output formats: dotenv, shell (export KEY=VALUE), and JSON.
//
// Usage:
//
//	e := export.New()
//	output, err := e.Export(vars, export.FormatShell)
package export
