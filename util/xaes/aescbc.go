package xaes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"esfgit.leju.com/golang/frame/encode"
)

// AesCbcEncrypt aes-cbc 加密
func AesCbcEncrypt(plainTextOrg, keyOrg, ivAesOrg string) (string, error) {
	plainText := []byte(plainTextOrg)
	key := []byte(keyOrg)
	ivAes := []byte(ivAesOrg)

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", encode.ErrAesKeyLengthSixteen
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	paddingText := PKCS5Padding(plainText, blockSize)

	var iv []byte
	if len(ivAes) != blockSize {
		return "", encode.ErrAesIv
	} else {
		iv = ivAes
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AesCbcDecrypt aes-cbc 解密
func AesCbcDecrypt(cipherTextOrg, keyOrg, ivAesOrg string) (string, error) {
	cipherText, _ := base64.StdEncoding.DecodeString(cipherTextOrg)
	key := []byte(keyOrg)
	ivAes := []byte(ivAesOrg)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", encode.ErrAesKeyLengthSixteen
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	var iv []byte
	if len(ivAes) != block.BlockSize() {
		return "", encode.ErrAesIv
	} else {
		iv = ivAes
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddingText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(paddingText, cipherText)

	plainText, err := PKCS5UnPadding(paddingText)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
