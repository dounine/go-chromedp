package util

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
)

type P12Util struct{}

func NewP12Util() *P12Util {
	return &P12Util{}
}

func (util *P12Util) P12ToPem(
	ctx *context.Context,
	info struct {
		P12Path  string
		Password string
		PemPath  string
	}) (err error) {
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"pkcs12", "-in",
		info.P12Path, "-out",
		info.PemPath, "-nokeys",
		"-passin", "pass:"+info.Password,
	).Output()
	return err
}
func (util *P12Util) PemToTextFile(
	ctx *context.Context,
	info struct {
		PemPath  string
		TextPath string
	}) (err error) {
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"x509",
		"-noout",
		"-text",
		"-in", info.PemPath,
		">", info.TextPath,
	).Output()
	return err
}
func (util *P12Util) DerToPem(
	ctx *context.Context,
	info struct {
		DerPath string
		PemPath string
	}) (err error) {
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"x509",
		"-inform", "der",
		"-in", info.DerPath,
		"-out", info.PemPath,
	).Output()
	return err
}
func (util *P12Util) ChangeP12Password(
	ctx *context.Context,
	info struct {
		P12Path     string
		Password    string
		NewP12Path  string
		NewPassword string
	}) (err error) {
	tempDir := os.TempDir()
	tmpPemPath := tempDir + "/" + uuid.New().String() + ".pem"
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"pkcs12",
		"-in", info.P12Path,
		"-nodes",
		"-out", tmpPemPath,
		"-password",
		"pass:"+info.Password,
	).Output()
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(tmpPemPath)
	_, err = exec.CommandContext(
		*ctx,
		"openssl",
		"pkcs12",
		"-export",
		"-in", tmpPemPath,
		"-out", info.NewP12Path,
		"-password", "pass:"+info.NewPassword,
	).Output()
	return err
}
func (util *P12Util) OscpCheck(
	ctx *context.Context,
	info struct {
		AppRootCaPemPath string
		IssuerPath       string
		P12PemPath       string
		OcspUrl          string
		Header           string
		ResultPath       string
	}) (bytes []byte, err error) {
	bytes, err = exec.CommandContext(
		*ctx,
		"openssl",
		"ocsp",
		"-CAfile", info.AppRootCaPemPath,
		"-issuer", info.IssuerPath,
		"-cert", info.P12PemPath,
		"-text", "-no_nonce",
		"-url", info.OcspUrl,
		info.Header, ">", info.ResultPath,
	).Output()
	return
}
