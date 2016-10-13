package shared

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DecryptKMS decrypts the provided string and returns it
func DecryptKMS(toDecrypt string) string {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := kms.New(sess)

	decodedString, err := base64.StdEncoding.DecodeString(toDecrypt)
	if err != nil {
		panic(err)
	}

	params := &kms.DecryptInput{
		CiphertextBlob: []byte(decodedString),
	}
	resp, err := svc.Decrypt(params)

	if err != nil {
		panic(err)
	}

	return string(resp.Plaintext[:])
}
