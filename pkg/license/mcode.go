package license

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/meilihao/goutil/crypto"
	"github.com/meilihao/goutil/hardware"
	"golang.org/x/crypto/sha3"
)

var (
	DefaultMcoder = &defaultMcoder{}
)

type Mcoder interface {
	Generate() (string, error)
}

// defaultMcoder use uuid v4 for mcode
type defaultMcoder struct {
}

func (m *defaultMcoder) Generate() (string, error) {
	return strings.ReplaceAll(uuid.NewString(), "-", ""), nil
}

// AdvanceMocoder use hardware info for mocde
type AdvanceMocoder struct {
	VirtType  string
	MachineID string
	RealMACs  string
}

func (m *AdvanceMocoder) Generate() (string, error) {
	var err error
	m.VirtType, err = hardware.VirtualInfo()
	if err != nil {
		return "", err
	}

	m.MachineID, err = hardware.MachineID()
	if err != nil {
		return "", err
	}

	var rMACs []string
	rMACs, err = hardware.RealMACs()
	if err != nil {
		return "", err
	}
	m.RealMACs = strings.Join(rMACs, "_")

	tmp := []string{
		fmt.Sprintf("virt_type(%s)", m.VirtType),
		fmt.Sprintf("machine_id(%s)", m.MachineID),
		fmt.Sprintf("rmacs(%s)", m.RealMACs),
	}
	input := strings.Join(tmp, "@")
	fmt.Printf("AdvanceMocoder input: %s\n", tmp)
	return crypto.HashString(sha3.New256(), input), nil
}
