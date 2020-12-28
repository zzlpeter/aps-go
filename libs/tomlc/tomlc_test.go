package tomlc

import "testing"

func TestConfig_BasicConf(t *testing.T) {
	basicMap := Config{}.BasicConf()
	if basicMap == nil {
		t.Errorf("BasicConf fails")
	}
}
