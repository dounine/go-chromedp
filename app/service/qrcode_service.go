package service

import (
	"github.com/skip2/go-qrcode"
	qrcodeDecode "github.com/tuotoo/qrcode"
	"os"
)

//go:generate mockgen -destination=./qrcode_service_mock.go -package=service go-chromedp/app/service IQrcodeService

const (
	Low RecoveryLevel = iota
	Medium
	High
	Highest
)

type (
	RecoveryLevel  qrcode.RecoveryLevel
	IQrcodeService interface {
		Create(url string, recovery RecoveryLevel, size int, qrcodePath string) error
		Parser(qrcodePath string) (string, error)
	}
	QrcodeService struct {
		IQrcodeService
	}
)

func NewQrcodeService() *QrcodeService {
	return &QrcodeService{}
}

// Create 二维码生成
func (s *QrcodeService) Create(url string, recovery RecoveryLevel, size int, qrcodePath string) error {
	return qrcode.WriteFile(url, qrcode.RecoveryLevel(recovery), size, qrcodePath)
}

// Parser 二维码解析
func (s *QrcodeService) Parser(qrcodePath string) (*string, error) {
	file, err := os.Open(qrcodePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	qrmatrix, err := qrcodeDecode.Decode(file)
	if err != nil {
		return nil, err
	}
	context := qrmatrix.Content
	return &context, nil
}
