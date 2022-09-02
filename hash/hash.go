package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func MD5(plain string) string {
	hash := md5.Sum([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA1(plain string) string {
	hash := sha1.Sum([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA224(plain string) string {
	hash := sha256.Sum224([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA256(plain string) string {
	hash := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA384(plain string) string {
	hash := sha512.Sum384([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA512(plain string) string {
	hash := sha512.Sum512([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA512224(plain string) string {
	hash := sha512.Sum512_224([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func SHA512256(plain string) string {
	hash := sha512.Sum512_256([]byte(plain))
	return hex.EncodeToString(hash[:])
}
