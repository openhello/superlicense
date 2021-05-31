package license

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"

	"github.com/farmerx/gorsa"
	"github.com/meilihao/goutil/crypto"
	"gopkg.in/yaml.v3"
)

var (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
)

func init() {
	if err := gorsa.RSA.SetPrivateKey(privateKey); err != nil {
		panic(err)
	}
}

func Generate(l *LicenseSection) (*LicenseDisplay, error) {
	raw, err := yaml.Marshal(l)
	if err != nil {
		return nil, err
	}

	aesKey := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	ciphertextData, err := crypto.AESEncrypt(aesKey, raw)
	if err != nil {
		return nil, err
	}

	ciphertextKey, err := gorsa.RSA.PriKeyENCTYPT(aesKey)
	if err != nil {
		return nil, err
	}

	version := make([]byte, 4)
	aesKeyLen := make([]byte, 4)

	binary.BigEndian.PutUint32(version, uint32(0x01))
	binary.BigEndian.PutUint32(aesKeyLen, uint32(len(ciphertextKey)))

	buf := bytes.NewBuffer(nil)
	buf.Write(version)
	buf.Write(aesKeyLen)
	buf.Write(ciphertextKey)
	buf.Write(ciphertextData)

	return &LicenseDisplay{
		License: l,
		Raw:     base64.URLEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
