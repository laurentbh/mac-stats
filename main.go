package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	recovery := NewRecovery()

	activateRecovery := false

	db, err := Connect()
	if err != nil {
		fmt.Printf("Recovery mode active (%v)\n", err)
		activateRecovery = true
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	fmt.Printf("hostname: %s\nnow:%v\n", hostname, now)
	batteryInfo, err := Battery()
	if err != nil {
		panic(err)
	}

	fmt.Printf("battery: %v\n", *batteryInfo)
	if activateRecovery {
		recovery.SaveBattery(now, hostname, *batteryInfo)
	} else {
		recoveredBat, err := recovery.LoadBattery()
		if err != nil {
			panic(err)
		}
		for _, r := range recoveredBat {
			err = db.InsertBattery(r.Host, r.Stamp, r.BatteryInfo)
			if err != nil {
				fmt.Printf("error insert battery %v\n", err)
			} else {
				if errDelete := os.Remove(r.FileName); errDelete != nil {
					panic(errDelete)
				}
			}
		}
		err = db.InsertBattery(hostname, now, *batteryInfo)
		if err != nil {
			fmt.Printf("error insert battery %v\n", err)
		}
	}

	now = time.Now()
	ssd, err := Ssd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ssd: %v\n", *ssd)

	if activateRecovery {
		recovery.SaveSsd(now, hostname, *ssd)
	} else {
		recoveredSsd, err := recovery.LoadSsd()
		if err != nil {
			panic(err)
		}
		for _, r := range recoveredSsd {
			err = db.InsertSSD(r.Host, r.Stamp, r.SsdInfo)
			if err != nil {
				fmt.Printf("error insert ssd %v\n", err)
			} else {
				if errDelete := os.Remove(r.FileName); errDelete != nil {
					panic(errDelete)
				}
			}
		}
		err = db.InsertSSD(hostname, now, *ssd)
		if err != nil {
			fmt.Printf("error insert ssd %v\n", err)
		}
	}

	// err = PushCounter(hostname, *batteryInfo, *ssd)
	// if err != nil {
	// 	panic(err)
	// }
}
