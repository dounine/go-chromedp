package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/pkg/errors"
	"os/exec"
)

func RsaGenerateP12(
	ctx *context.Context, info struct {
		CertPemPath       string
		PrivateKeyPemPath string
		P12Path           string
		Password          string
	}) (err error) {
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"pkcs12",
		"-export",
		"-in", info.CertPemPath,
		"-inkey", info.PrivateKeyPemPath,
		"-out", info.P12Path,
		"-passout", "pass:"+info.Password,
	).Output()
	return
}
func RsaGeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func RsaGenerateRSAKey() (*struct {
	PrivateKey []byte
	PublicKey  []byte
}, error) {
	privateKey, err := RsaGeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := pem.Block{
		Type:  "Rsa Private Key",
		Bytes: x509PrivateKey,
	}
	privateKeyPem := pem.EncodeToMemory(&privateBlock)
	publicKey := privateKey.PublicKey
	x509PublicKey, err2 := x509.MarshalPKIXPublicKey(&publicKey)
	if err2 != nil {
		return nil, err2
	}
	publicBlock := pem.Block{
		Type:  "Rsa Public Key",
		Bytes: x509PublicKey,
	}
	publicKeyPem := pem.EncodeToMemory(&publicBlock)
	return &struct {
		PrivateKey []byte
		PublicKey  []byte
	}{
		PrivateKey: privateKeyPem,
		PublicKey:  publicKeyPem,
	}, nil
}

func RsaGenerateCSR(privateKey *rsa.PrivateKey) ([]byte, error) {
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:       []string{"CN"},
			Organization:  []string{"lake"},
			CommonName:    "ipadump.com",
			Locality:      []string{"guang zhou"},
			Province:      []string{"guang dong"},
			StreetAddress: []string{"dian he dang dong"},
		},
		Version: 1,
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "CSR请求生成失败")
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	}), nil
}
