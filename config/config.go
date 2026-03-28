package config

import (
	"fmt"
)

type Config struct {
	FaucetPort          string
	BinaryName          string
	FaucetWalletAddress string
	FaucetWalletAlias   string
	ChainNode           string
	ChainHome           string
	KeyringBackend      string
	ChainID             string
}

func (c Config) IsValid() error {
	if c.FaucetPort == "" {
		return fmt.Errorf("port flag must not be empty")
	}

	if c.BinaryName == "" {
		return fmt.Errorf("cli flag must not be empty")
	}

	if c.FaucetWalletAddress == "" {
		return fmt.Errorf("address flag must not be empty")
	}

	if c.FaucetWalletAlias == "" {
		return fmt.Errorf("alias flag must not be empty")
	}
	if c.ChainNode == "" {
		return fmt.Errorf("node flag must not be empty")
	}
	if c.ChainHome == "" {
		return fmt.Errorf("home flag must not be empty")
	}
	if c.KeyringBackend == "" {
		return fmt.Errorf("keyring-backend flag must not be empty")
	}
	if c.ChainID == "" {
		return fmt.Errorf("chain-id flag must not be empty")
	}

	return nil
}
