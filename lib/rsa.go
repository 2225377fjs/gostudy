package lib

import (
	"encoding/pem"
	"errors"
	"crypto/x509"
	"fmt"
	"crypto/rsa"
	"crypto/sha1"
	"crypto"
	"encoding/base64"
	"crypto/rand"
)




// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {return parsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCODpInJjTexIWIeaH51iKpLyfZ7u59WxJkja6aMzyRqrHQYmbz
IgJ4ChHcg7yoj4vLmMvPYImQR/Dx0BRG0fxuKi8Javm8Vh4gJ2dRcwuz4rCpGRwM
jDHPl7NGY+YMcNFS/lghPfJDosThyBBMbq9K3QdJJ/xz7pfJsCpwJbbHzwIDAQAB
AoGAPpvFZmO8YYITqDaTYN0zoYUa9z4K9kyxKogXL/bs9vZYMBNudDHFDMrDXIDj
IRdC8ZSCHv/ZITHTy399sEjUC8CYLn5qqPRoALOcRXoFEMIxipSw6q3ir6VErmjn
RQLB50mNrGJ/OyDIoTt6BTTs0YOPykD5I/Q+jY/ZfBcLFKECQQDUEPxblTswMHdQ
DWPnIyfBsytP+Fqwiz/JYJYOBGhTJdnedFr7XtpEfpE4LwGX9OiKK4+sHtt1UP4L
B+ufM1CDAkEAq3ycKE4jxoSj88Z5B8qHvUdCqsv94MdyoRVJ3V1aJQ2xZFbe069V
X4SD7mbw3t02QNZIxlIMNnZkR95hEZBxxQJBAMSu6E9MjjO4j8BIWwLxwRrOwPoP
npUk4Uk1cpaXkeakMXg3tHZ1V7y1IpzYRAMam14i3sLFb8dUEfpLI0ZpQl8CQC13
fG+zSAj6Yf3gQXavXA3zNtnR/B38w4ex/UOT3LK2TrIr1iiJ9Di/CbvLz1FHlXrb
VbA/UL2f5jan31So14ECQF5Q0UonkL/oYXScIzhTQ4NBGSU5dn3mxtrrLk9/1TFq
WcwVr40XHVtE/Yc5U46YNBF6kzQ0y0IEX3zDlqmXcBE=
-----END RSA PRIVATE KEY-----`))
}


func loadPrivateKey2(path string) (Signer, error) {return parsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDQeH5gEL7NJawkikgiGK91aSPvpAKBz+gMu0glvAFSUc9otWSH25UxQ1jZCtRl41PArJchsvVQQhn0ndGJdTSPhCYatuEZxkPhsmD27T0Apcs/RbDcuc44cMod/8XQ1obAa42BN56rBzf9AJ6ndyl1AJw0/Ojfeoa85O9oU49O/QIDAQABAoGBAKmGss4QD/jnf7sfMFV52YUTAQQpTVie50cjLSJyZmi42n99ssjACezLpX6qTdqlKEBwmV3wF4kyl8TSacjsJNZvQXUDNaGV4lrTKmnjuIEsoimE9Pre4VgyUXAL/aWYgws87sxz0mzBGRr4yifUB09y0SaHfVe3vz3agYvOIWmJAkEA+15RoJl3F4mfG8KwDgnsBZG0++SYI2lY/dXOGH7f89rjpSW2ttAwsCid6n1sO80f9ZdTX4WGzsWKuQbKSIQ4VwJBANRP1IQpnn8E4nFk/Bb1CKF6FHFAd2jdcjpLI6MAQ3yEIPJvPU4J20SsOdmSKOGnYqB5PtbAqe6YnhlGF8HBLssCQQCLCKNLmja17Sf1Od0ZFsHWXr5lKQ5BX+6aD907zUlf3u1VFiQWv9Z+SSj3X0IzXYTU2UuDJR7oVXkiWDAgpglnAkAyLyh5kOjg90ObMBaSSpseqB+a4XUYOXfdpZMn3VEWZpjvFTI1dwj4Q4ltDypQpGMgsWgUFPhV6Ic+TB4jc0lfAkBIfHKm8Ht2eySN7s56fIjbaNjqVg63mUhNF2eC5QTWRcga8PjDC+kwkG+9eu2n1Pts+LRgAXGjr5pgELCgWi6z
-----END RSA PRIVATE KEY-----`))
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {block, _ := pem.Decode(pemBytes)
	if block == nil {return nil, errors.New("ssh: no key found")}

	var rawkey interface{}
	switch block.Type {case"RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {return nil, err}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
	SignData(data string) string
}


func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}



type rsaPrivateKey struct {*rsa.PrivateKey}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {h := sha1.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA1, d)
}

func (r *rsaPrivateKey) SignData(data string) string{
	signed, _ := r.Sign([]byte(data))
	sig := base64.StdEncoding.EncodeToString(signed)
	return sig
}



var signer1 Signer
var signer2 Signer



func init() {
	signer1, _ = loadPrivateKey("")
	signer2, _ = loadPrivateKey2("")
}


type BidRequestData struct {
	PageIndex int
	StartDateTime string
}

func SingData1(data string) string{
	signed, _ := signer1.Sign([]byte(data))
	sig := base64.StdEncoding.EncodeToString(signed)
	return sig
}

func SingData2(data string) string{
	signed, _ := signer2.Sign([]byte(data))
	sig := base64.StdEncoding.EncodeToString(signed)
	return sig
}


func GetBidListRequestData(data BidRequestData) string {
	return "StartDateTime" + data.StartDateTime
}


func GetSigner(key string) (Signer, error) {
	return parsePrivateKey([]byte(key))
}