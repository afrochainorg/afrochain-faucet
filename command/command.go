package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/afrochainorg/afrochain-faucet/config"
)

const (
	faucetExecCommand = "tx bank send %s %s %s --from %s --node %s --home %s --keyring-backend %s --chain-id %s --fees %s --broadcast-mode sync --yes --output json"
)

func ExecuteTransfer(c config.Config, recipient, amount string) ([]byte, error) {
	transferCommand := fmt.Sprintf(faucetExecCommand, c.FaucetWalletAddress, recipient, amount, c.FaucetWalletAlias, c.ChainNode, c.ChainHome, c.KeyringBackend, c.ChainID, c.Fees)

	cmd := exec.Command(c.BinaryName, strings.Split(transferCommand, " ")...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return out.Bytes(), fmt.Errorf("%s: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}
