package core

import (
	"encoding/json"
	"os"

	"github.com/gliderlabs/ssh"
)

type AuthorizedUser struct {
	Username string `json:"username"`
	PubKey   string `json:"pub_key"`
}

// LoadAuthorizedKeys charge les utilisateurs depuis le JSON
func loadAuthorizedKeys(path string) ([]AuthorizedUser, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var users []AuthorizedUser
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// publicKeyHandler est la fonction de callback pour ton serveur SSH
func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	users, err := loadAuthorizedKeys("./data/users.json")
	if err != nil {
		return false
	}

	for _, u := range users {
		authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(u.PubKey))
		if err != nil {
			continue
		}

		// Comparaison simplifiée
		if ssh.KeysEqual(key, authorizedKey) {
			ctx.SetValue("username", u.Username)
			return true
		}
	}

	return false
}
