/**
 * 工具函数库-AES加解密 @woylin, 2016-1-12
 */
package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

//AesKey解密，base64->[]byte
func AesKeyDecode(k string) []byte {
	data, _ := base64.StdEncoding.DecodeString(k + "=")
	return data
}

//AES加密,传入明文和密钥，[]byte
func AesDecrypt(a, key []byte) ([]byte, error) {
	k := len(key) //PKCS#7
	if len(a)%k != 0 {
		return nil, errors.New("aesKey error:length not mutiple")
	}
	pc, _ := aes.NewCipher(key)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(pc, iv)
	deMsg := make([]byte, len(a))
	blockMode.CryptBlocks(deMsg, a)
	return deMsg, nil
}

//Aes解密，传入密文和密钥，[]byte
func AesEncrypt(plainData []byte, key []byte) ([]byte, error) {
	k := len(key)
	if len(plainData)%k != 0 {
		plainData = PKCS7Pad(plainData, k)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)

	return cipherData, nil
}

//from github.com/vgorin/cryptogo
func PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := PadLength(len(message), blocksize)

		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

		// apply padding
		padded = append(message, padding...)
		return padded
	}
	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}

// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(slice_length, blocksize int) (padlen int) {
	padlen = blocksize - slice_length%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}
