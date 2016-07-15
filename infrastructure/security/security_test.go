package security

import (
	. "gopkg.in/check.v1"
	"testing"
)

type InfrastructureSuite struct{}

var _ = Suite(&InfrastructureSuite{})

func Test(t *testing.T) {
	TestingT(t)
}

const (
	TestKey     = "asfhethnahaye243nsdiussdfb3242hl"
	TestText    = "Hello my name is Stephen"
	TestHmacKey = "234sfdhhk34589sscfsdh9sfkhsf79"
)

func (s *InfrastructureSuite) TestEncryption(c *C) {

	encryptedText, err := EncryptAES(TestKey, TestText)
	if err != nil {
		c.Fatal(err)
	}

	decryptedText, err := DecryptAES(TestKey, encryptedText)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(TestText, Equals, decryptedText)
}

func (s *InfrastructureSuite) BenchmarkEncryption(c *C) {
	for i := 0; i < c.N; i++ {
		encryptedText, err := EncryptAES(TestKey, TestText)
		if err != nil {
			c.Fatal(err)
		}

		if _, err := DecryptAES(TestKey, encryptedText); err != nil {
			c.Fatal(err)

		}
	}
}

func (s *InfrastructureSuite) TestHMAC(c *C) {

	signature1 := HmacSha1(TestHmacKey, TestText)
	signature2 := HmacSha1(TestHmacKey, TestText)

	c.Assert(signature1, Equals, signature2)
}
