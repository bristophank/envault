// Package config manages envault's local user configuration.
//
// Configuration is stored as a JSON file, typically at
// ~/.envault/config.json, and holds defaults such as the path to the
// user's age identity file and the shared recipients list.
//
// Example usage:
//
//	mgr, err := config.New("") // uses $HOME/.envault
//	cfg, err := mgr.Load()
//	cfg.DefaultIdentityFile = "/path/to/key.age"
//	err = mgr.Save(cfg)
package config
