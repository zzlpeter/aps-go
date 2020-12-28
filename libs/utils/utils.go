package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
)

type HostManager struct {
}

// 获取本机hostname
func (h HostManager) LocalHostName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return name, nil
}

// 获取本机IP
func (h HostManager) LocalIp() (string, error) {
	addrS, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	for _, address := range addrS {
		// 检查IP地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("未查询到IP地址")
}

// 获取出口公网IP
func (h HostManager) ExternalIp() (string, error) {
	rsp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		fmt.Printf("获取出口公网IP异常: %v \n", err.Error())
		return "", err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err1 := ioutil.ReadAll(rsp.Body)
	if err1 != nil {
		return "", err
	}

	return string(body), nil
}

// 计算字符串MD5值
func MD5(s string) string {
	data := []byte(s)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 生成UUID
func GenUUID() string {
	u1 := uuid.NewV4()
	s := u1.String()
	s = strings.Replace(s, "-", "", 4)
	return s
}

// 随时间自增ID
func GenAutoIncrementId() string {
	nanoSecond := StampNanoSecond()
	randInt := rand.Intn(10000)
	s := fmt.Sprintf("%d%05d", nanoSecond, randInt)
	return s
}

// 环境变量
type Environ struct {
}

// 环境变量 - 设置
func (e Environ) Set(env, value string) error {
	return os.Setenv(env, value)
}

// 环境变量 - 读取
func (e Environ) Get(env string) string {
	return os.Getenv(env)
}