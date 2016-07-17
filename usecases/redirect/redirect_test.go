package redirect

import (
	"testing"
)

func TestRedirect(t *testing.T) {

}

// Base62 valid output string containing 0-9a-zA-Z
//
// Examples
// 100 b 10 = 1 * 10^2 + 0 * 10^1 + 0 * 10^0
// 100 b 62 = 1 * 62^1 + 38 * 62^0 => 1C

// 111 b 10 = 1 * 10^2 + 1 * 10^1 + 1 * 10^0
// 111 b 62 = 1 * 62^1 + 49 * 62^0 => 1N
//
// 0        1         2         3         4         5         6
// 012345679abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ

const Base62Alphabet = "012345679abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func base62(num int, alphabet string) string {

}
