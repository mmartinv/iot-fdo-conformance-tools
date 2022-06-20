package fdoshared

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
)

type HashType int

const (
	HASH_SHA256      HashType = -16
	HASH_SHA384      HashType = -43
	HASH_HMAC_SHA256 HashType = 5
	HASH_HMAC_SHA384 HashType = 6
)

type HashOrHmac struct {
	_    struct{} `cbor:",toarray"`
	Type HashType
	Hash []byte
}

var HmacToHashAlg map[HashType]HashType = map[HashType]HashType{
	HASH_HMAC_SHA256: HASH_SHA256,
	HASH_HMAC_SHA384: HASH_SHA384,
}

func GenerateFdoHash(data []byte, hashType HashType) (HashOrHmac, error) {
	switch hashType {
	case HASH_SHA256:
		hashDigest := sha256.Sum256(data)

		return HashOrHmac{
			Type: hashType,
			Hash: hashDigest[:],
		}, nil
	case HASH_SHA384:
		hashDigest := sha512.Sum384(data)

		return HashOrHmac{
			Type: hashType,
			Hash: hashDigest[:],
		}, nil
	default:
		return HashOrHmac{}, fmt.Errorf("Error generating hash. %d is unknown hashing algorithm", hashType)
	}
}

func GenerateFdoHmac(data []byte, hashType HashType, key []byte) (HashOrHmac, error) {
	switch hashType {
	case HASH_HMAC_SHA256:
		macInst := hmac.New(sha256.New, key)
		macInst.Write(data)

		return HashOrHmac{
			Type: hashType,
			Hash: macInst.Sum(nil),
		}, nil
	case HASH_HMAC_SHA384:
		macInst := hmac.New(sha512.New384, key)
		macInst.Write(data)

		return HashOrHmac{
			Type: hashType,
			Hash: macInst.Sum(nil),
		}, nil
	default:
		return HashOrHmac{}, fmt.Errorf("Error generating hmac. %d is unknown hmac algorithm", hashType)
	}
}

func VerifyHash(data []byte, fdoHashB HashOrHmac) (bool, error) {
	switch fdoHashB.Type {
	case HASH_SHA256:
		if len(fdoHashB.Hash) != sha256.New().Size() {
			return false, errors.New("Failed to verify hash. The input hash does not match expected hash size.")
		}

		fdoHashA, _ := GenerateFdoHash(data, fdoHashB.Type)
		if bytes.Equal(fdoHashB.Hash, fdoHashA.Hash) {
			return true, nil
		} else {
			return false, nil
		}
	case HASH_SHA384:
		if len(fdoHashB.Hash) != sha512.New384().Size() {
			return false, errors.New("Failed to verify hash. The input hash does not match expected hash size.")
		}

		fdoHashA, _ := GenerateFdoHash(data, fdoHashB.Type)
		if bytes.Equal(fdoHashB.Hash, fdoHashA.Hash) {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, fmt.Errorf("Error verifying hash. %d is an unknown hash algorithm", fdoHashB.Type)
	}
}

func VerifyHMac(data []byte, inputHmac HashOrHmac, key []byte) (bool, error) {
	switch inputHmac.Type {
	case HASH_HMAC_SHA256:
		macInst := hmac.New(sha256.New, key)
		macInst.Write(data)
		computedMac := macInst.Sum(nil)

		if bytes.Equal(inputHmac.Hash, computedMac) {
			return true, nil
		} else {
			return false, nil
		}
	case HASH_HMAC_SHA384:
		macInst := hmac.New(sha512.New384, key)
		macInst.Write(data)
		computedMac := macInst.Sum(nil)

		if bytes.Equal(inputHmac.Hash, computedMac) {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, fmt.Errorf("Error verifying hmac. %d is unknown hmac algorithm", inputHmac.Type)
	}
}
