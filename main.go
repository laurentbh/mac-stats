package main

import (
	"fmt"
	"os"
)

func main() {

	// now := time.Now()
	// fmt.Printf("it's %v\n%s", now, now.Format(time.RFC3339))

	// err := battery(context.TODO())
	// if err != nil {
	// 	panic(err)
	// }
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", name)
	batteryInfo, err := Battery()
	if err != nil {
		panic(err)
	}
	fmt.Printf("battery: %v", *batteryInfo)

}
