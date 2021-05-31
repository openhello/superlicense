package license

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/farmerx/gorsa"
	"github.com/meilihao/goutil/crypto"
	"github.com/meilihao/goutil/file"
	"gopkg.in/yaml.v3"
)

var (
	ErrLicenseLen           = errors.New("less than license min len")
	ErrLicenseMcodeNotMatch = errors.New("mcode not match")

	publicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA5ptoMDnSnLzs0PfhmOeJ
VjRGh40lFI4OmAgMfPxaBDsGQ3yL9bJx5g/8ysV6pxHe/KWDJRXJJqNrQ4AYEoWT
YwEh3yHxPXO8X6vqIZmdDFr0cCbQ0kdrCDAU6HZmLhKCed0gkkx2uo4PsiJ80J6q
aUs6gvnfWK22RrjnFmVJmHk2ZZgUPTLo/Y2xcIUx8ilMhGIYGTulc67m4Y5u4Br3
kTHiEufBsjdY75tYvr/0qjD8utObQaK/Ox7N2VPIeHl3z4DJaC9IVerCFYhb16nn
B3/kbcB36eATG1SUXEsCbn6t+K5XYhfMI1966rajqguIfHbE7Ns1fBKu+wYdrPQL
pxVEnOnJNidhec/C38yuvlJ9QHNnEzcyh/TyZLmWL4tO1JJGjd+IZLd/5u1ZUNFP
aXzphIKAORbMTxYysmQ4Yl+aTqJGinNLS+V5kbXb8iCIT6wrT/xT+XSIOBrv1o8G
M+q8yQ5xcKnvXU9j3ymn9yjSIEQGiq8To685qedq0IuqSgdL5D76yZzIy9QbXfgD
asZR8hgOH97ur3jEYfw0bs4yaO7izyiM17336EN4tYfoPFDmvGM+t9NDrcHBuOkY
3uHHGzYd6Fl5DX6OqzSR93UnpPf2uzNpdl0wNV4hDG1cGmdlKqgmBIxfhE8KU5iI
rcUEmXqEIozlqWDtGyrjWikCAwEAAQ==
-----END PUBLIC KEY-----`
)

func init() {
	if err := gorsa.RSA.SetPublicKey(publicKey); err != nil {
		panic(err)
	}
}

type LicenseDisplay struct {
	License *LicenseSection `yaml:"license"`
	Raw     string          `yaml:"raw"`
}

type TimeAt int64

func (t TimeAt) MarshalYAML() (interface{}, error) {
	tmp := time.Unix(int64(t), 0)

	return tmp.Format("2006-01-02 15:04:05"), nil
}

func (t *TimeAt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp string
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	tmpT, err := time.ParseInLocation("2006-01-02 15:04:05", tmp, time.UTC)
	if err != nil {
		return err
	}

	*t = TimeAt(tmpT.Unix())

	return nil
}

type LicenseFeature struct {
	Code    string `yaml:"code"`
	Content string `yaml:"content"`
	Type    string `yaml:"type"`
}

// core license ParentId is "", others ParentId = core license ParentId
// core license Mcode/User/Product is not "", others Mcode/User/Product is ""
type LicenseSection struct {
	Id        string            `yaml:"id"`
	ParentId  string            `yaml:"parent_id"`
	Mcode     string            `yaml:"mcode"`
	User      string            `yaml:"user"`
	Product   string            `yaml:"product"`
	Category  string            // core is must, others is expand
	Features  []*LicenseFeature `yaml:"features"`
	SignedAt  TimeAt            `yaml:"signed_at"`
	ExpiredAt TimeAt            `yaml:"expired_at"`
}

func Load(fs []string) ([]*LicenseSection, error) {
	var ls []*LicenseSection
	var l, parent *LicenseSection
	var err error

	for _, v := range fs {
		l, err = ParseLicenseWithYamlContent(file.GetFileValue(v))
		if err != nil {
			return nil, err
		}

		if l.ParentId == "" {
			if parent != nil {
				panic("duplicate parent license")
			}

			tmp := make([]*LicenseSection, len(ls)+1)
			tmp = append(tmp, l)
			if len(ls) > 0 {
				tmp = append(tmp, ls...)
			}

			ls = tmp
			parent = l
		} else {
			ls = append(ls, l)
		}
	}

	return ls, nil
}

func ParseLicenseWithYamlContent(raw string) (*LicenseSection, error) {
	ld := &LicenseDisplay{}

	err := yaml.Unmarshal([]byte(raw), ld)
	if err != nil {
		return nil, err
	}

	return ParseLicenseWithRaw(ld.Raw)
}

func ParseLicenseWithRaw(raw string) (*LicenseSection, error) {
	// rsa加密限制: 明文长度需要小于密钥长度
	// data = version(4) + encryptedAesKeyLen(4) + rsa.encrypt(aesKey) + aes.encrypt(data)
	data, _ := base64.URLEncoding.DecodeString(raw)
	if len(data) < 8 {
		return nil, ErrLicenseLen
	}

	var l *LicenseSection
	var err error

	version := binary.BigEndian.Uint32(data[:4])
	switch version {
	case 0x01:
		fmt.Printf("license version: %x\n", version)

		fallthrough
	default:
		l, err = ParseWithVersion01(data[4:])
	}

	return l, err
}

func ParseWithVersion01(data []byte) (*LicenseSection, error) {
	encryptedAesKeyLen := binary.BigEndian.Uint32(data[:4])

	aesKey, err := gorsa.RSA.PubKeyDECRYPT(data[4 : 4+encryptedAesKeyLen])
	if err != nil {
		return nil, err
	}

	yamlData, err := crypto.AESDecrypt(aesKey, data[4+encryptedAesKeyLen:])
	if err != nil {
		return nil, err
	}

	var l *LicenseSection
	if err = yaml.Unmarshal(yamlData, &l); err != nil {
		return nil, err
	}

	return l, nil
}

// LoadAndCheck
// if mcode is empty, skip mcode check
func LoadAndCheck(fs []string, mcode string) ([]*LicenseSection, error) {
	ls, err := Load(fs)
	if err != nil {
		return nil, err
	}

	for _, l := range ls {
		if mcode != "" && l.Mcode != mcode {
			return nil, ErrLicenseMcodeNotMatch
		}
	}

	return ls, err
}
