package types

import "testing"

func TestToString(t *testing.T) {
	cases := []struct {
		value    DriveMode
		expected string
	}{
		{DriveModeUser, "user"},
		{DriveModePilot, "pilot"},
		{DriveModeInvalid, ""},
	}

	for _, c := range cases {
		val := ToString(c.value)
		if val != c.expected {
			t.Errorf("ToString(%v): %v, wants %v", c.value, val, c.expected)
		}
	}
}

func TestParseString(t *testing.T) {
	cases := []struct {
		value    string
		expected DriveMode
	}{
		{"user", DriveModeUser},
		{"pilot", DriveModePilot},
		{"", DriveModeInvalid},
		{"invalid", DriveModeInvalid},
	}

	for _, c := range cases {
		val := ParseString(c.value)
		if val != c.expected {
			t.Errorf("ParseString(%v): %v, wants %v", c.value, val, c.expected)
		}
	}
}
