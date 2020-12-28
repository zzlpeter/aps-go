package utils

import "testing"

func TestTimestampFormat(t *testing.T) {
	s := TimestampFormat(1604633050, "")
	if s != "2020-11-06 11:24:10" {
		t.Errorf("TimestampFormat fails")
	}
}