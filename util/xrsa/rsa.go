package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// 公钥加密
func PublicEncrypt(encryptStr string, publicStr string) (string, error) {
	// pem 解码
	block, _ := pem.Decode([]byte(publicStr))

	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}

	//返回密文
	return base64.StdEncoding.EncodeToString(encryptedStr), nil
}

// 私钥解密
func PrivateDecrypt(decryptStr string, privateStr string) (string, error) {
	// pem 解码
	block, _ := pem.Decode([]byte(privateStr))

	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.StdEncoding.DecodeString(decryptStr)

	//对密文进行解密
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)

	//返回明文
	return string(decrypted), nil
}
