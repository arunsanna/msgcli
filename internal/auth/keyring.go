package auth

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/99designs/keyring"
)

const (
	serviceName = "msgcli"
)

// TokenData holds OAuth tokens for an account
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"` // Unix timestamp
	Email        string `json:"email"`
}

// AccountInfo holds metadata about a configured account
type AccountInfo struct {
	Alias string `json:"alias"`
	Email string `json:"email"`
}

func openKeyring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: serviceName,
		// Use default backends for each platform:
		// macOS: Keychain
		// Linux: Secret Service or file-based fallback
		// Windows: Credential Manager
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		},
		FileDir:                      "~/.msgcli/keyring",
		FilePasswordFunc:             fileKeyringPassword,
		KeychainTrustApplication:     true,
		KeychainAccessibleWhenUnlocked: true,
	})
}

func fileKeyringPassword(prompt string) (string, error) {
	// Check for environment variable first (for CI/agent use)
	if pw := getEnvPassword(); pw != "" {
		return pw, nil
	}
	return "", errors.New("file keyring requires MSGCLI_KEYRING_PASSWORD environment variable")
}

func getEnvPassword() string {
	return getEnv("MSGCLI_KEYRING_PASSWORD")
}

func getEnv(key string) string {
	// Using os.Getenv but wrapped for testability
	return envGetter(key)
}

var envGetter = os.Getenv

// SaveToken stores tokens for an account in the keyring
func SaveToken(alias string, token *TokenData) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return kr.Set(keyring.Item{
		Key:  "token:" + alias,
		Data: data,
	})
}

// LoadToken retrieves tokens for an account from the keyring
func LoadToken(alias string) (*TokenData, error) {
	kr, err := openKeyring()
	if err != nil {
		return nil, err
	}

	item, err := kr.Get("token:" + alias)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, errors.New("account not found - run 'msgcli auth add " + alias + "' first")
		}
		return nil, err
	}

	var token TokenData
	if err := json.Unmarshal(item.Data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// DeleteToken removes tokens for an account from the keyring
func DeleteToken(alias string) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}

	return kr.Remove("token:" + alias)
}

// ListAccounts returns all configured account aliases
func ListAccounts() ([]AccountInfo, error) {
	kr, err := openKeyring()
	if err != nil {
		return nil, err
	}

	keys, err := kr.Keys()
	if err != nil {
		return nil, err
	}

	var accounts []AccountInfo
	for _, key := range keys {
		if len(key) > 6 && key[:6] == "token:" {
			alias := key[6:]
			item, err := kr.Get(key)
			if err != nil {
				continue
			}
			var token TokenData
			if err := json.Unmarshal(item.Data, &token); err != nil {
				continue
			}
			accounts = append(accounts, AccountInfo{
				Alias: alias,
				Email: token.Email,
			})
		}
	}

	return accounts, nil
}

// GetDefaultAccount returns the first configured account alias
func GetDefaultAccount() (string, error) {
	accounts, err := ListAccounts()
	if err != nil {
		return "", err
	}
	if len(accounts) == 0 {
		return "", errors.New("no accounts configured - run 'msgcli auth add <alias>' first")
	}
	return accounts[0].Alias, nil
}
