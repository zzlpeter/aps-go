package utils

import "testing"

func TestHostIpManager_LocalHostName(t *testing.T) {
	m := HostManager{}
	hostName, _ := m.LocalHostName()
	if hostName == "" {
		t.Errorf("HostIpManager.LocalHostName failed")
	}
}

func TestHostIpManager_LocalIp(t *testing.T) {
	m := HostManager{}
	ip, _ := m.LocalIp()
	if ip == "" {
		t.Errorf("HostIpManager.LocalIp failed")
	}
}

func TestHostManager_ExternalIp(t *testing.T) {
	m := HostManager{}
	extIp, _ := m.ExternalIp()
	if extIp == "" {
		t.Errorf("HostIpManager.ExternalIp failed")
	}
}

func TestGenUUID(t *testing.T) {
	uuid := GenUUID()
	if len(uuid) != 32 {
		t.Errorf("TestGenUUID failed")
	}
}

func TestMD5(t *testing.T) {
	str := "hello world"
	md5 := MD5(str)
	if md5 != "5eb63bbbe01eeed093cb22bb8f5acdc3" {
		t.Errorf("TestMD5 failed")
	}
}

func TestGenAutoIncrementId(t *testing.T) {
	id := GenAutoIncrementId()
	if len(id) != 24 {
		t.Errorf("TestGenAutoIncrementId failed")
	}
}