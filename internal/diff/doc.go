// Package diff compares two environment variable maps and produces a
// structured list of additions, removals, and modifications.
//
// It is used by the rotate command to show what changed after re-sealing
// a vault, and can be used anywhere two env snapshots need to be compared.
package diff
