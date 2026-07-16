package security

import (
	"SeptimanappBackend/database"
	"SeptimanappBackend/types"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

// CreateApikey returns a new random api key. Only the hash of a key is stored,
// so the value returned here is the only copy of it there will ever be.
func CreateApikey() string {
	uuidWithHyphen := uuid.New()
	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

func HashKey(unhashed string) string {
	h := sha256.New()
	h.Write([]byte(unhashed))
	key := fmt.Sprintf("%x", h.Sum(nil))
	return key
}

// StoreNewApiKey creates an api key, stores its hash and returns the key itself.
// It is returned instead of printed so that a caller cannot hand out a key whose
// hash never reached the database.
func StoreNewApiKey(repository database.Repository) (string, error) {
	key := CreateApikey()
	err := repository.StoreSecurityInfo(types.ApiKeyInfo{ApiKeyHash: HashKey(key)})
	if err != nil {
		return "", err
	}
	return key, nil
}

func ValidateApikey(repository database.Repository, key string) (bool, error) {
	hasKey, err := repository.HasApiKeyInfo(
		types.ApiKeyInfo{
			ApiKeyHash: HashKey(key),
		})
	return hasKey, err
}
