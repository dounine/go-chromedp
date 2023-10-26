package util

import (
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"go-chromedp/app/middleware"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

func CopyFields(src any, des any) {
	sourcePtr := reflect.ValueOf(src)
	destPtr := reflect.ValueOf(des)
	if sourcePtr.Kind() != reflect.Ptr || destPtr.Kind() != reflect.Ptr {
		panic("Both source and destination must be pointers")
		return
	}
	sourceStruct := sourcePtr.Elem()
	destStruct := destPtr.Elem()
	if sourceStruct.Kind() != reflect.Struct || destStruct.Kind() != reflect.Struct {
		panic("Both source and destination must be structs")
		return
	}

	for i := 0; i < sourceStruct.NumField(); i++ {
		field := sourceStruct.Type().Field(i)
		fieldValue := sourceStruct.FieldByName(field.Name)
		destField := destStruct.FieldByName(field.Name)
		if destField.IsValid() && destField.Type() == sourceStruct.Field(i).Type() {
			destField.Set(fieldValue)
		}
	}
}

func FileToBase64(path string) (str string, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(bytes)
	return
}

func Long_uid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func P8JwtToken(
	info struct {
		Iss        string
		Kid        string
		PrivateKey []byte
	}) (token string, err error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": info.Iss,
		"aud": "appstoreconnect-v1",
		"iat": now.Unix(),
		"exp": now.Add(20 * time.Minute).Unix(),
	}
	var pk *ecdsa.PrivateKey
	pk, err = jwt.ParseECPrivateKeyFromPEM(info.PrivateKey)
	if err != nil {
		return
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t.Header["alg"] = "ES256"
	t.Header["kid"] = info.Kid
	t.Header["typ"] = "JWT"
	token, err = t.SignedString(pk)
	return
}

func Secure_uid() string {
	uid := Long_uid()
	numToUpper := 8
	rand.New(rand.NewSource(time.Now().UnixNano()))
	upperIndexes := rand.Perm(len(uid))[:numToUpper]
	result := []int32(uid)
	for _, idx := range upperIndexes {
		result[idx] = int32(strings.ToUpper(string(result[idx]))[0])
	}
	return string(result)
}
func Short_uid() string {
	uid, err := shortid.Generate()
	if err != nil {
		middleware.Logger.Error("short_uid error", err)
		return Long_uid()
	}
	return uid
}
func Domain() string {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:8000"
	}
	return domain
}
