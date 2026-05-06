// Package recipients provides a simple ordered, deduplicated list of age
// public keys used as encryption recipients for envault vaults.
//
// A recipients file (.envault-recipients) is stored alongside the sealed
// vault file so that any team member can inspect who is authorised to
// decrypt the secrets, and so that the vault can be re-sealed for the
// correct set of keys after a rotation.
//
// File format: one age public key per line; lines starting with '#' and
// blank lines are treated as comments and ignored on load.
package recipients
