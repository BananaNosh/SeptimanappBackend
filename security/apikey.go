package security

import (
	"SeptimanappBackend/database"
	"SeptimanappBackend/types"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func CreateApikey() string {
	uuidWithHyphen := uuid.New()
	unhashed := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	fmt.Printf("ApiKeyHash: %s\n", unhashed)
	key := HashKey(unhashed)
	return key
}

func HashKey(unhashed string) string {
	h := sha256.New()
	h.Write([]byte(unhashed))
	key := fmt.Sprintf("%x", h.Sum(nil))
	return key
}

func StoreNewApiKey() {
	key := CreateApikey()
	info := types.ApiKeyInfo{
		ApiKeyHash: key,
	}
	if err := database.StoreSecurityInfo(info); err != nil {
		fmt.Println(err)
	}
}

func ValidateApikey(key string) (bool, error) {
	hasKey, err := database.HasApiKeyInfo(
		types.ApiKeyInfo{
			ApiKeyHash: HashKey(key),
		})
	return hasKey, err
}
