package common

import (
	"fmt"
	"github.com/pkg/errors"
	. "github.com/sjhitchner/trackly/domain"
	"github.com/sjhitchner/trackly/infrastructure/security"
	"strings"
)

var DeviceIdEmptyError = errors.New("Device Id empty")
var SignatureEmptyError = errors.New("Signature empty")
var DeviceTokenEmptyError = errors.New("Device token empty")
var SignatureInvalidError = errors.New("Signatures invalid")

type defaultDeviceTokenEncrypter struct {
	encryptionKey string
	secretKey     string
}

func NewDeviceTokenEncrypter(encryptionKey, secretKey string) DeviceTokenEncrypter {
	return &defaultDeviceTokenEncrypter{
		encryptionKey: encryptionKey,
		secretKey:     secretKey,
	}
}

func (t defaultDeviceTokenEncrypter) Validate() error {
	if len(t.encryptionKey) != security.AESKeyLength {
		return errors.Errorf(
			"Invalid encryption key size %d should be %d",
			len(t.encryptionKey),
			security.AESKeyLength,
		)
	}
	return nil
}

func (t defaultDeviceTokenEncrypter) Decrypt(deviceToken string) (DeviceId, error) {
	if deviceToken == "" {
		return "", DeviceTokenEmptyError
	}

	decrypted, err := security.DecryptAES(t.encryptionKey, deviceToken)
	if err != nil {
		return "", errors.Wrap(err, "unable to decrypt device token")
	}

	deviceId, signature := splitDecrypted(decrypted)
	if deviceId == "" {
		return "", DeviceIdEmptyError
	}

	if signature == "" {
		return "", SignatureEmptyError
	}

	calculatedSignature, err := security.HmacSha1(t.secretKey, string(deviceId))
	if err != nil {
		return "", errors.Wrap(err, "unable to hmac")
	}

	if signature != calculatedSignature {
		return "", SignatureInvalidError
	}

	return deviceId, nil
}

func splitDecrypted(decrypted string) (DeviceId, string) {
	s := strings.SplitN(decrypted, ":", 2)
	if len(s) != 2 {
		return "", ""

	}
	return DeviceId(s[0]), s[1]

}

func (t defaultDeviceTokenEncrypter) Encrypt(deviceId DeviceId) (string, error) {
	if deviceId == "" {
		return "", DeviceIdEmptyError
	}

	signature, err := security.HmacSha1(t.secretKey, string(deviceId))
	if err != nil {
		return "", errors.Wrap(err, "unable to hmac")
	}

	toEncrypt := fmt.Sprintf("%s:%s", deviceId, signature)

	encrypted, err := security.EncryptAES(t.encryptionKey, toEncrypt)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}
