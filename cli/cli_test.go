package cli

import (
	"os"
	"testing"
)

func TestSetDefaultValueFromEnv(t *testing.T) {
	err := os.Setenv("KEY1", "value1")
	if err != nil {
		t.Errorf("unable to set env value: %v", err)
	}

	cases := []struct{
		key string
		defValue string
		expected string
	}{
		{"MISSING_KEY", "default", "default"},
		{"KEY1", "bad value", "value1"},
	}

	for _, c := range cases {
		var value = ""
		SetDefaultValueFromEnv(&value, c.key, c.defValue)
		if c.expected != value {
			t.Errorf("SetDefaultValueFromEnv(*value, %v, %v): %v, wants %v", c.key, c.defValue, value, c.expected)
		}
	}
}
