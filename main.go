package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	db, err := Connect()
	if err != nil {
		panic(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	fmt.Printf("hostname: %s\nnow:%v", hostname, now)
	batteryInfo, err := Battery()
	if err != nil {
		panic(err)
	}
	fmt.Printf("battery: %v\n", *batteryInfo)
	err = db.InsertBattery(hostname, now, *batteryInfo)
	if err != nil {
		fmt.Printf("error insert battery %v\n", err)
	}

	now = time.Now()
	ssd, err := Ssd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ssd: %v", *ssd)
	err = db.InsertSSD(hostname, now, *ssd)
	if err != nil {
		fmt.Printf("error insert ssd %v\n", err)
	}

	// err = PushCounter(hostname, *batteryInfo, *ssd)
	// if err != nil {
	// 	panic(err)
	// }
}
