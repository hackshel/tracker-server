package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func CheckClientConnectable(ip string, port uint16) bool {
	address := fmt.Sprintf("%s:%d", ip, port)

	// 设置连接超时时间，比如 3 秒
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func main() {
	s := "12"

	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Print("errors %v", err)
	}

	fmt.Printf("value id %v, type %T", d, d)

	dt := time.Now()
	fmt.Printf("dt val %v, type %T\n", dt, dt)

	fmt.Printf("dt val %v, type %T\n", dt.Unix(), dt)

	r := CheckClientConnectable("192.168.11.239", 25114)
	fmt.Printf("error %v", r)

}
