// Package vault provides high-level operations for sealing (encrypting) and
// unsealing (decrypting) .env files using age encryption.
//
// A vault ties together the env parser and the age crypto layer. Sealing reads
// a plaintext .env file, validates its contents, and writes an age-encrypted
// ciphertext to the destination path. Unsealing reverses the process using the
// caller's private key.
//
// Example usage:
//
//	v := vault.New(cryptoService)
//
//	// Encrypt
//	if err := v.Seal("path/to/.env", "path/to/.env.age", recipients); err != nil {
//		log.Fatal(err)
//	}
//
//	// Decrypt
//	if err := v.Unseal("path/to/.env.age", "path/to/.env", privateKey); err != nil {
//		log.Fatal(err)
//	}
package vault
