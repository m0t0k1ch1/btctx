package btc

import (
	"bytes"
	"regexp"

	"github.com/m0t0k1ch1/base58"
)

type Pkh []byte

func (pkh Pkh) Bytes() []byte {
	return []byte(pkh)
}

func (pkh Pkh) Address() (Address, error) {
	if len(pkh) != PkhLength {
		return "", ErrInvalidPkhLength
	}

	b := pkh

	if isTestnet() {
		b = append([]byte{AddressVersionTest}, b...)
	} else {
		b = append([]byte{AddressVersionMain}, b...)
	}

	doubleHashedBytes, err := Sha256Double(b)
	if err != nil {
		return "", err
	}
	checksumBytes := doubleHashedBytes[0:4]
	b = append(b, checksumBytes...)

	address, err := base58.NewBitcoinBase58().EncodeToString(b)
	if err != nil {
		return "", nil
	}

	return Address(address), nil
}

type Address string

func (address Address) String() string {
	return string(address)
}

func (address Address) Pkh() (Pkh, error) {
	b, err := base58.NewBitcoinBase58().DecodeString(address.String())
	if err != nil {
		return nil, err
	}

	return b[1 : len(b)-4], nil
}

// ref. https://bitcointalk.org/index.php?topic=1026.0
func (address Address) IsValid() (bool, error) {
	s := address.String()

	if ok := regexp.MustCompile(`^[a-zA-Z1-9]{27,35}$`).MatchString(s); !ok {
		return false, nil
	}

	b, err := base58.NewBitcoinBase58().DecodeString(s)
	if err != nil {
		return false, err
	}

	doubleHashedBytes, err := Sha256Double(b[:len(b)-4])
	if err != nil {
		return false, err
	}
	if !bytes.Equal(b[len(b)-4:], doubleHashedBytes[:4]) {
		return false, nil
	}

	return true, nil
}
