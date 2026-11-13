package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	hostname, _ := os.Hostname()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("Host Name   :", hostname)
	fmt.Println("OS / Arch   :", runtime.GOOS, "/", runtime.GOARCH)
	fmt.Println("CPU Cores   :", runtime.NumCPU())
	fmt.Println("App Memory  :", m.Alloc/1024, "KB")
	fmt.Println("System Time :", time.Now().Format(time.RFC1123))
}
