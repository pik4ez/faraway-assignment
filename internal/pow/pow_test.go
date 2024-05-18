package pow

import "testing"

func TestVerifyPow(t *testing.T) {
	tests := []struct {
		name       string
		challenge  string
		nonce      string
		difficulty int
		expected   bool
	}{
		{
			// hash: 0000e6311749b7443090fc19195ee485c70b74ea2bf22fb9c3776118997e21aa
			name:       "hash-with-4-zeros-on-difficulty-4-should-pass",
			challenge:  "WnImAsM6A16QsF9k",
			nonce:      "WVia24BIAhyajNOu",
			difficulty: 4,
			expected:   true,
		},
		{
			// hash: 0000e6311749b7443090fc19195ee485c70b74ea2bf22fb9c3776118997e21aa
			name:       "hash-with-4-zeros-on-difficulty-5-should-fail",
			challenge:  "WnImAsM6A16QsF9k",
			nonce:      "WVia24BIAhyajNOu",
			difficulty: 5,
			expected:   false,
		},
	}

	for _, test := range tests {
		actual := VerifyPoW(test.challenge, test.nonce, test.difficulty)
		if actual != test.expected {
			t.Fatalf("Test \"%s\" failed, actual %t, expected %t\n", test.name, actual, test.expected)
		}
	}
}
