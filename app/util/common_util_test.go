package util

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCopyFields(t *testing.T) {
	type A struct {
		A int
	}
	type B struct {
		A int
	}
	a := A{A: 1}
	b := B{}
	CopyFields(&a, &b)
	assert := assert2.New(t)
	assert.Equal(b.A, a.A)

	type C struct {
		A int
	}
	type D struct {
		A int16
	}
	c := C{A: 1}
	d := D{}
	CopyFields(&c, &d)
	assert.NotEqual(c.A, d.A)
}

func TestDomain(t *testing.T) {
	assert := assert2.New(t)
	os.Setenv("DOMAIN", "http://localhost:8080")
	domain := Domain()
	assert.Equal(domain, "http://localhost:8080")

	os.Setenv("DOMAIN", "")
	domain = Domain()
	assert.Equal(domain, "http://localhost:8000")
}

func TestFileToBase64(t *testing.T) {
	assert := assert2.New(t)
	content := "hello"
	path := "/tmp/base64.txt"
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		assert.Error(err, "os.WriteFile()")
		return
	}
	defer os.Remove(path)
	str, err := FileToBase64(path)
	if err != nil {
		assert.Error(err, "FileToBase64()")
		return
	}
	assert.Equal(str, "aGVsbG8=")
}

func TestLong_uid(t *testing.T) {
	assert := assert2.New(t)
	uid := Long_uid()
	if len(uid) != 32 {
		assert.Failf("Long_uid() = %v, want %v", uid, 32)
	}
}

func TestSecure_uid(t *testing.T) {
	assert := assert2.New(t)
	uid := Secure_uid()
	if len(uid) != 32 {
		assert.Failf("Secure_uid() = %v, want %v", uid, 32)
	}
}

func TestShort_uid(t *testing.T) {
	assert := assert2.New(t)
	uid := Short_uid()
	if len(uid) != 9 {
		assert.Failf("Short_uid() = %v, want %v", uid, 9)
	}
}
