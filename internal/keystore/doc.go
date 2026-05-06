// Package keystore manages age encryption keys for envault.
//
// It stores the user's private identity key and a list of public recipient
// keys in a local directory (default: ~/.envault). The private key is written
// with mode 0600 so that only the owning user can read it.
//
// Typical usage:
//
//	s := keystore.New(os.UserHomeDir())
//	s.Init()                        // create ~/.envault/
//	s.SavePrivateKey(privKey)       // write identity.txt
//	s.AddRecipient(pubKey)          // append to recipients.txt
//	recipients, _ := s.LoadRecipients()
package keystore
