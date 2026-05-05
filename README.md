# envault

Minimal secrets manager that encrypts `.env` files using [age](https://github.com/FiloSottile/age) encryption for safe team sharing.

---

## Installation

```bash
go install github.com/yourusername/envault@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envault/releases).

---

## Usage

**Encrypt a `.env` file for your team:**

```bash
# Encrypt using a recipient's public key
envault encrypt .env --recipient age1ql3z7hjy... --out .env.age

# Decrypt when you need it
envault decrypt .env.age --identity ~/.age/key.txt --out .env

# Run a command with secrets loaded directly (no plaintext file written)
envault run --identity ~/.age/key.txt --file .env.age -- go run main.go
```

**Typical workflow:**

1. Each team member generates an age key pair: `age-keygen -o ~/.age/key.txt`
2. Share public keys with the team.
3. Encrypt the `.env` file for all recipients and commit `.env.age` to version control.
4. Never commit the plaintext `.env` file.

---

## Requirements

- Go 1.21+
- [age](https://github.com/FiloSottile/age) (bundled, no separate install needed)

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)