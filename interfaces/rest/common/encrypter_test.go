package common

import (
	. "github.com/sjhitchner/trackly/domain"
	. "gopkg.in/check.v1"
	"testing"
)

type PixelTrackingSuite struct{}

var _ = Suite(&PixelTrackingSuite{})

func Test(t *testing.T) {
	TestingT(t)
}

func (s PixelTrackingSuite) TestDefaultDeviceTokenEncrypter(c *C) {
	encrypter := DefaultDeviceTokenEncrypter{
		//encryptionKey: "5764328pn3a41194q5p4nqv89ppg99pn0qspq44r",
		encryptionKey: "5764328pn3a41194q5p4nqv89ppg99pn",
		secretKey:     "592baef1a02bfaf13d63b280c38775b9",
	}

	deviceId := DeviceId("add52c6e-fb76-4690-a4dd-d403e4df90d6")

	deviceToken, err := encrypter.Encrypt(deviceId)
	c.Assert(err, IsNil)

	newDeviceId, err := encrypter.Decrypt(deviceToken)
	c.Assert(err, IsNil)

	c.Assert(deviceId, Equals, newDeviceId)
}
