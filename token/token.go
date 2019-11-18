package token

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/duoflow/yc-snapshot/config"
)

// GetSignedToken - to form JWT signed token
func GetSignedToken(conf *config.Configuration) string {
	issuedAt := time.Now()
	token := jwt.NewWithClaims(ps256WithSaltLengthEqualsHash, jwt.StandardClaims{
		Issuer:    conf.ServiceAccountID,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: issuedAt.Add(time.Hour).Unix(),
		Audience:  "https://iam.api.cloud.yandex.net/iam/v1/tokens",
	})
	token.Header["kid"] = conf.KeyID

	privateKey := loadPrivateKey(conf)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signed
}

// По умолчанию Go RSA PSS использует PSSSaltLengthAuto,
// но на странице https://tools.ietf.org/html/rfc7518#section-3.5 сказано, что
// размер значения соли должен совпадать с размером вывода хеш-функции.
// После исправления https://github.com/dgrijalva/jwt-go/issues/285
// можно будет заменить на jwt.SigningMethodPS256.
var ps256WithSaltLengthEqualsHash = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

func loadPrivateKey(conf *config.Configuration) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(conf.PrivateRSAKeyFile)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		panic(err)
	}
	return rsaPrivateKey
}

// GetIAMToken - to get IAM token
func GetIAMToken(conf *config.Configuration) string {
	jot := GetSignedToken(conf)
	fmt.Println(jot)
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("%s: %s", resp.Status, body))
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data.IAMToken
}
