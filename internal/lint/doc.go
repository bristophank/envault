// Package lint provides static analysis for .env files.
//
// It checks for common issues such as empty values, duplicate keys,
// placeholder values (e.g. "changeme", "TODO"), and keys that do not
// conform to the conventional UPPER_SNAKE_CASE naming style.
//
// Usage:
//
//	results, err := lint.Lint(content)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, r := range results {
//		fmt.Printf("line %d: %s\n", r.Line, r.Message)
//	}
package lint
