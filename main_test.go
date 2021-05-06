package sdk

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"testing"
)

const (
	value1           = "Hello"
	key1             = "val1"
	keyFileName      = "pub.pem"
	keySize          = 1024
	privateKeyPath   = "./privateKey_test.pem"
	publicKeyPath    = "./publicKey_test.pem"
	privateKeyString = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA8C58A0KbO6Jl8nRsCxFXecDmg1X57lJcc/nSn3/QoWMA7XUk
7fVIVBhrNbSMD2G4FoYsgPuAGTKqUrCRttVYdPhxdiIhHGPHDH1LvNnCnqc1QDwD
N1W62I0/k8e9AkbImd8T58ZYDsgfQf+JhNjDYsAGH66mM0Yt2R+vsa/Pt0lkunvD
iKCSxwzlymbGOlIh0d+OmYuUDRnrUUL+jle2AleJ+0LcTMTre16rIyT/tpSVlkxZ
qzrFVhWLwJLkbX+ha9kcNpDDfGxDVmIFmcy614/CHX3JH91jtrBwIhXIuPOzNYSZ
8J85WZW4Ow6MAyEc5t5E9uEKVTM2WDrHpBb/BwIDAQABAoIBAQDnLGayoJ5XJLUp
S0Ne178xciisysj3yRAxlIhUeqptW6Rd6b20x7xpLOOr2m5gs7aC/3vAXdHq7ugf
FNH7f5dXZnWWtbzW3XaNn9+REquPFvNbMygJT5u6qSFDdSGIGmckKyG2mSLSf24O
kQ1k71oIJzj9r5VKjsa8UBJEXSr5hm5uYG7F+msgH58DEVPlybV3Jdq4q5w0SJ9Z
/UWxKQ2izZ7co869vhIEifN0VliVSJ9xa1z/Jc77HSOj7NDIbJUARfQ/gtt+dQ16
ORsKKWOaYsOdoa+VXVCIwmnDi5N5YBVs6nEr+afR6/SyoorHYkvfg+78rr+pBSAR
8/KEY6HZAoGBAPrGYkaKgnj+JT1QMUhMQnPAogYW2F4iXqWiL5wXYJvJULLDykt7
k/zPVehYjQNdEqQHZCHPtqijFfEPV2PXeWLRa9bQ6bzIA8uboeXjWD4FeLvgJ2C+
wlJafKCcmKXddrmVTUvTqJ7Pb4Ezw1t60juhy/05G0gkzQsdMsuPuvmbAoGBAPUv
mKdVEW4SUYJCuq054LxBj/4OVg8AVtQFkvEH1LqEoTH+dbkny06LlVcHBiz+wsgg
NfzYoqjAPp20EUFobDXPK4biuVA2mn3oyWFJjNOnuQghE+ChLRh8hhKEiaseUuHU
2I+1GOf3sts8uLt44LWidtM79iQ1Y0M/doBPHc0FAoGBAOVtVj/fPJrhOMStd0kD
q9AmrpUPlYgZvamfhhsyMAqW1aOXCJ6iQrQKJDhbuzcWkZVLxcpBNIV4HvzZ4kPP
wJgtrJFttEooW4CNtEKUCglEDD8mRiB2pWWer2JpoiYtRQ9ojr0OubgBY6w65UHu
TiSMVAopktIgCQ9f+TbPGmp9AoGAYvghcHoAHSQ7zo7M95uDQbpdOzniNw/1/IN7
etukXN2oi5uhPWn4wO3LDGQDdCopycpmwHdZwTBIljPXO0XBWD8V3M6r6tr/pY9P
qnub4tuy7rsbYPLuVxH8tIDXaUFGR245NFjvgsMTaTergdEbM3Yu7LkpdBgwxzZY
yRYme1kCgYAe8QG1p6SNLKrr1vZJ9INCocDfFXgxOCCWvbbil+sW7JEGuD0S4Uzb
rwqw25i402PrlUSAH6pagbsfs36y8yMK5TkrvGXfFS79q6bNoc/F1Kp0WTVtjEAW
d536+JxVtD1hhF8iaZn3a1iS82cjpQ1nC1DB4plaMK++VMGC0a+c0Q==
-----END RSA PRIVATE KEY-----
`
)

func TestEncryptDecryptMessage(t *testing.T) {
	// write private on file
	privateKeyBytes := []byte(privateKeyString)

	privateKey, err := bytesToPrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(privateKeyPath, privateKeyToBytes(privateKey), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// write public on file
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	err = ioutil.WriteFile(publicKeyPath, pubBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := &StorageClient{
		Bbolt: &BboltStorage{
			SotoragePath: "./mytest.db",
			BucketName:   "test",
		},
	}

	avocado := Avocado{}
	err = avocado.New(storage)
	if err != nil {
		t.Fatal(err)
	}

	_, err = avocado.EecryptAndStoreValue([]byte(key1), []byte(value1), publicKeyPath)
	if err != nil {
		t.Fatal(err)
	}

	keys, err := avocado.GetAllKeys()
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) == 0 {
		t.Fatal("At least one key should be present")
	}

	valueDecrypted, err := avocado.FindAndDecryptValueBy([]byte(key1), privateKeyPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(valueDecrypted) != value1 {
		t.Fatal("Should be expected message equals to ", value1)
	}

	err = avocado.Delete([]byte(key1))
	if err != nil {
		t.Fatal(err)
	}

	_, err = avocado.FindAndDecryptValueBy([]byte(key1), privateKeyPath)
	if err == nil {
		t.Fatal("The error should be ", ErrorEncryptedValueNotFound)
	}
}

// privateKeyToBytes converts a private key to bytes form
func privateKeyToBytes(privateKey *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return privBytes
}
