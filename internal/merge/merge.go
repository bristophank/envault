// Package merge provides utilities for merging two .env variable sets,
// applying overrides and producing a unified result with conflict tracking.
package merge

// Strategy controls how conflicting keys are resolved.
type Strategy int

const (
	// KeepBase retains the value from the base map on conflict.
	KeepBase Strategy = iota
	// TakeIncoming overwrites the base value with the incoming value on conflict.
	TakeIncoming
)

// Conflict describes a key whose value differed between base and incoming.
type Conflict struct {
	Key      string
	BaseVal  string
	IncomingVal string
}

// Result holds the merged environment and any conflicts that were encountered.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
	Added     []string // keys present only in incoming
	Removed   []string // keys present only in base
}

// Merge combines base and incoming according to the given strategy.
// Keys present in both maps with differing values are recorded as conflicts.
func Merge(base, incoming map[string]string, strategy Strategy) Result {
	result := Result{
		Env: make(map[string]string, len(base)),
	}

	// Copy base into result.
	for k, v := range base {
		result.Env[k] = v
	}

	for k, inVal := range incoming {
		baseVal, exists := base[k]
		if !exists {
			result.Env[k] = inVal
			result.Added = append(result.Added, k)
			continue
		}
		if baseVal != inVal {
			result.Conflicts = append(result.Conflicts, Conflict{
				Key:         k,
				BaseVal:     baseVal,
				IncomingVal: inVal,
			})
			if strategy == TakeIncoming {
				result.Env[k] = inVal
			}
		}
	}

	for k := range base {
		if _, exists := incoming[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	return result
}
