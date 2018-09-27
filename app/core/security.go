package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	rsa "crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/revel/revel"
)

// GetPubKey returns the public rrsa key
func GetPubKey() (pubkey *rsa.PublicKey, err error) {
	fullPath := revel.BasePath + "/conf/certs/pubkey.pem"

	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return pubkey, err
	}

	block, _ := pem.Decode(keyBytes)
	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return pubkey, err
		}
		rawkey = rsa
	default:
		return pubkey, errors.New("ssh: unsupported key type " + block.Type)
	}

	pubkey, ok := rawkey.(*rsa.PublicKey)
	if !ok {
		return pubkey, errors.New("Error casting rawKey to pubKey")
	}
	return pubkey, nil
}

// GetPrivKey returns the private rrsa key
func GetPrivKey() (privkey *rsa.PrivateKey, err error) {
	fullPath := revel.BasePath + "/conf/certs/mykey.pem"

	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return privkey, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return privkey, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return privkey, err
		}
		rawkey = rsa
	default:
		return nil, errors.New("ssh: unsupported key type %q" + block.Type)
	}

	privkey, ok := rawkey.(*rsa.PrivateKey)
	if !ok {
		return privkey, errors.New("Error casting rawKey to pubKey")
	}

	return privkey, nil
}

// EncryptCipher uses Cipher crypt to return a cryp text
func EncryptCipher(text string) string {

	qrKey := "*.$pyCh4773r_MX*.S3cr3t_QR_k3y!."

	key := []byte(qrKey)

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		revel.ERROR.Print(err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, errIO := io.ReadFull(rand.Reader, iv); err != nil {
		revel.ERROR.Print(errIO)
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// DecryptCipher decripts string using Cipher crypt
func DecryptCipher(cryptoText string) string {
	qrKey := "*.$pyCh4773r_MX*.S3cr3t_QR_k3y!."
	key := []byte(qrKey)

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		revel.ERROR.Print("Cipher text too short")
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

// EncryptRSA encrypts data(string, int, float, bool) into enconded encrypted string
func EncryptRSA(data interface{}, cryptType string) (string, error) {

	pubkey, err := GetPubKey()
	if err != nil {
		revel.ERROR.Print(err)
		return "", err
	}

	dataStr := ""

	switch data.(type) {
	case string:
		dataStr = data.(string)
	case int16:
		dataStr = strconv.FormatInt(int64(data.(int16)), 10)
	case int32:
		dataStr = strconv.FormatInt(int64(data.(int32)), 10)
	case int64:
		dataStr = strconv.FormatInt(int64(data.(int64)), 10)
	case float32:
		dataStr = strconv.FormatFloat(float64(data.(float32)), 'f', 6, 32)
	case float64:
		dataStr = strconv.FormatFloat(data.(float64), 'f', 6, 64)
	case bool:
		dataStr = strconv.FormatBool(data.(bool))
	default:
		return "", errors.New("Unsupported type")
	}

	var bytesMsg = []byte(dataStr)

	if cryptType == "oaep" {
		/*
			var label = []byte("test")
			var sha1hash = sha1.New()

			log.Print(pubkey)
			encryptedmsg, err := rsa.EncryptOAEP(sha1hash, rand.Reader, &publicKey, bytesMsg, label)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(bytesMsg), encryptedmsg)
			//responseMsg = string(encryptedmsg)
		*/
	}
	if cryptType == "pkcs" {
	}

	// EncryptPKCS1v15
	encryptedPKCS1v15, errPKCS1v15 := rsa.EncryptPKCS1v15(rand.Reader, pubkey, bytesMsg)
	if errPKCS1v15 != nil {
		revel.ERROR.Print(errPKCS1v15)
		return "", errPKCS1v15
	}

	return base64.URLEncoding.EncodeToString(encryptedPKCS1v15), nil
}

// DecryptRSA receives data and returns the decrypted data
func DecryptRSA(data string, field reflect.Value, cryptType string) (interface{}, error) {

	privkey, err := GetPrivKey()
	if err != nil {
		revel.ERROR.Print(err)
		return "", err
	}

	if cryptType == "oaep" {
		/*
			decryptedmsg, err := rsa.DecryptOAEP(sha1hash, rand.Reader, privatekey, encryptedmsg, label)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		*/
	}

	ciphertext, _ := base64.URLEncoding.DecodeString(data)

	byteArr, err := rsa.DecryptPKCS1v15(rand.Reader, privkey, ciphertext)
	if err != nil {
		return "", err
	}

	var dataStr = string(byteArr)
	var response interface{}

	switch field.Kind() {
	//
	case reflect.String:
		response = dataStr
	//
	case reflect.Int16:
		r, _ := strconv.Atoi(dataStr)
		response = int16(r)
	//
	case reflect.Int32:
		r, _ := strconv.Atoi(dataStr)
		response = int32(r)
	//
	case reflect.Int64:
		r, _ := strconv.Atoi(dataStr)
		response = int64(r)
	//
	case reflect.Float32:
		response, _ = strconv.ParseFloat(dataStr, 32)
	//
	case reflect.Float64:
		response, _ = strconv.ParseFloat(dataStr, 64)
	//
	case reflect.Bool:
		response, _ = strconv.ParseBool(dataStr)
	//
	default:
		return "", errors.New("Unsupported type")
	}
	return response, nil

}
