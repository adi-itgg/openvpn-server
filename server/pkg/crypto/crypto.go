package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMD5(dataBytes []byte) string {
	md5Hash := md5.Sum(dataBytes)
	return hex.EncodeToString(md5Hash[:])
}
