package ssmctl

import (
	"testing"
)

func TestGetParameter(t *testing.T) {
	m := NewMockSSMManager()
	cases := []struct {
		query           string
		withDescryption bool
		expected        string
	}{
		{
			"/service1/dev/key1",
			false,
			"dev_value1",
		},
		{
			"/service1/dev/key2",
			true,
			"dev_value2",
		},
		{
			"/service1/prod/key2",
			false,
			dummyEncryptedValue,
		},
	}

	for _, c := range cases {
		result, err := m.GetParameter(c.query, c.withDescryption)
		if err != nil {
			t.Fatalf("Failed: error = %s", err)
		}

		if result.Value != c.expected {
			t.Errorf("want: %s\nget : %s", c.expected, result)
		}
	}

}
