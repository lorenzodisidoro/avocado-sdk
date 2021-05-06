package sdk

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	oaep "github.com/lorenzodisidoro/rsa-oaep"
)

// Error messages
var (
	ErrorFilePathNotProvided    = errors.New("file path not provided")
	ErrorValueNotProvided       = errors.New("value not provided")
	ErrorKeyNotProvided         = errors.New("value label not provided")
	ErrorEncryptedValueNotFound = errors.New("encrypted value not found")
)

// Avocado defines RSA keys and storage configurations
type Avocado struct {
	storage *StorageClient
}

// New constructs a key value to store
func (a *Avocado) New(config *StorageClient) error {
	a.storage = config
	err := a.storage.New()

	return err
}

// EecryptAndStoreValue encrypt and store value using RSA public key
// @param key and value in byte format
// @return encrypted value bytes
func (a *Avocado) EecryptAndStoreValue(key, value []byte, publicKeyPath string) ([]byte, error) {
	var err error

	if len(value) <= 0 {
		return nil, ErrorValueNotProvided
	}

	if len(key) <= 0 {
		return nil, ErrorKeyNotProvided
	}

	publicKey, err := getPublicKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	oaep := oaep.NewRSAOaep(sha256.New())
	cipherText, err := oaep.Encrypt(publicKey, value, key)
	if err != nil {
		return nil, err
	}

	err = a.storage.set(key, cipherText)
	if err != nil {
		return cipherText, err
	}

	return cipherText, nil
}

// FindAndDecryptValueBy find and decrypt value
// method read private key from file
// @param key in byte format
// @return decrypted value bytes
func (a *Avocado) FindAndDecryptValueBy(key []byte, privateKeyPath string) ([]byte, error) {
	encryptedValue, err := a.storage.get(key)
	if err != nil {
		return nil, err
	}

	if encryptedValue == nil {
		return nil, ErrorEncryptedValueNotFound
	}

	privateKey, err := getPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	oaep := oaep.NewRSAOaep(sha256.New())
	decriptedValue, err := oaep.Dencrypt(privateKey, encryptedValue, key)
	if err != nil {
		return nil, err
	}

	return decriptedValue, nil
}

// GetAllKeys return all keys
// @return all keys
func (a *Avocado) GetAllKeys() ([][]byte, error) {
	return a.storage.getAll()
}

// Delete remove a key
// @return error
func (a *Avocado) Delete(key []byte) error {
	return a.storage.delete(key)
}

// readFile data from file
func readFile(path string) ([]byte, error) {
	contentByte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return contentByte, nil
}

// getPrivateKeyFromFile return private key in rsa.PublicKey format
func getPrivateKeyFromFile(privateKeyPath string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := bytesToPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKeyFromFile(publicKeyPath string) (*rsa.PublicKey, error) {
	publicKeyByte, err := readFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := bytesToPublicKey(publicKeyByte)
	if err != nil {
		return nil, err
	}

	return publicKey, err
}

// Thanks to https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a#file-rsa_util-go-L68
func bytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return key, nil
}

// bytesToPrivateKey serialises bytes to rsa.PrivateKey
func bytesToPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error

	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}

	return key, nil
}
