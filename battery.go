package main

import (
	"context"
	"os/exec"
)

const (
	BATTERY_CMD       = "system_profiler"
	BATTERY_ARG       = "SPPowerDataType"
	BATTERY_BUF_SIZE  = 10000
	battery_key_cycle = "cycle"
	battery_key_full  = "full"
)

type BatteryInfo struct {
	fullCharge string
	cycle      string
}

var batteryExp = []Expression{
	{Pat: "^Cycle Count: (.*)$", Key: battery_key_cycle},
	{Pat: "^Full Charge Capacity \\(mAh\\): (.*)$", Key: battery_key_full},
}

func Battery() (*BatteryInfo, error) {
	return systemCall(context.Background())
}

func systemCall(ctx context.Context) (*BatteryInfo, error) {
	exp, err := compilePattern(batteryExp)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(context.Background(), BATTERY_CMD, BATTERY_ARG)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	data := make(chan string, 1000)
	done := make(chan bool)
	go parseToLines(stdout, data)
	go processLines(ctx, exp, data, done)

	<-done
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	ret := BatteryInfo{}
	for e := 0; e < len(exp); e++ {
		if exp[e].expression.Key == battery_key_cycle {
			ret.cycle = exp[e].expression.Value
		}
		if exp[e].expression.Key == battery_key_full {
			ret.fullCharge = exp[e].expression.Value
		}

	}
	return &ret, nil
}
