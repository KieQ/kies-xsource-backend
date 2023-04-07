package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func GenerateKey(cur int64)[aes.BlockSize]byte{
	key := [aes.BlockSize]byte{}
	for i:=0; i<len(key); i++{
		key[i] = uint8(cur % (256 - int64(i)))
	}
	return key
}

func Decrypt(cipherText, iv string, key []byte)(string, error){
	cipherTextBytes, err := hex.DecodeString(cipherText)
	if err != nil{
		return "", err
	}
	ivBytes, err := hex.DecodeString(iv)
	if err != nil{
		return "", err
	}
	decryptedBytes, err := decryptCBC(cipherTextBytes, key, ivBytes)
	if err != nil{
		return "", err
	}
	if len(decryptedBytes) == 0{
		return "", errors.New("empty decrypted text")
	}

	paddingLength := int(decryptedBytes[len(decryptedBytes)-1])
	if len(decryptedBytes) < paddingLength{
		return "", errors.New("padding is longer than decrypted text")
	}

	result := decryptedBytes[:len(decryptedBytes) - paddingLength]
	return string(result), nil
}

func decryptCBC(cipherText, key, iv []byte) ([]byte, error) {
    var block cipher.Block
	var err error
	var plainText []byte

    if block, err = aes.NewCipher(key); err != nil {
        return plainText, err
    }

    if len(cipherText) < aes.BlockSize {
        fmt.Printf("ciphertext too short")
        return plainText, errors.New("ciphertext too short")
    }

    cbc := cipher.NewCBCDecrypter(block, iv)
    cbc.CryptBlocks(cipherText, cipherText)

    plainText = cipherText

    return plainText, nil
}