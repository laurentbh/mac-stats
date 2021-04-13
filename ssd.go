package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type SsdInfo struct {
	AvailableSpare          string
	AvailableSpareThreshold string
	PercentageUsed          string
	UnitRead                string
	UnitReadTB              string
	UnitWrite               string
	UnitWriteTB             string
	HostRead                string
	HostWrite               string
	PowerCycle              string
	PowerHours              string
	MediaErrors             string
}

const (
	SSD_CMD           = "smartctl"
	SSD_ARG1          = "--all"
	SSD_ARG2          = "/dev/disk0"
	spare             = "spare"
	spare_threshold   = "sparethreshold"
	percentage_used   = "percent"
	data_unit_read    = "data_unit_read"
	data_unit_readTB  = "data_unit_readTB"
	data_unit_write   = "data_unit_write"
	data_unit_writeTB = "data_unit_writeTB"
	host_write        = "host_write"
	host_read         = "host_read"
	power_cycle       = "power_cycle"
	power_hours       = "power_hours"
	media_errors      = "media_errors"
)

var ssdExp = []Expression{
	{Pat: "^Available Spare:\\s*([0-9]{1,3})%$", Key: spare},
	{Pat: "^Available Spare Threshold:\\s*([0-9]{1,3})%$", Key: spare_threshold},
	{Pat: "^Percentage Used:\\s*([0-9]{1,3})%$", Key: percentage_used},
	{Pat: "^Data Units Read:\\s*([0-9].*) \\[.*$", Key: data_unit_read},
	{Pat: "^Data Units Read:\\s*[0-9].* \\[(.*) TB]$", Key: data_unit_readTB},
	{Pat: "^Data Units Written:\\s*([0-9].*) \\[.*$", Key: data_unit_write},
	{Pat: "^Data Units Written:\\s*[0-9].* \\[(.*) TB]$", Key: data_unit_writeTB},
	{Pat: "^Host Read Commands:\\s*([0-9].*)$", Key: host_read},
	{Pat: "^Host Write Commands:\\s*([0-9].*)$", Key: host_write},
	{Pat: "^Power Cycles:\\s*([0-9].*)$", Key: power_cycle},
	{Pat: "^Power On Hours:\\s*([0-9].*)$", Key: power_hours},
	{Pat: "^Media and Data Integrity Errors:\\s*([0-9].*)$", Key: media_errors},
}

func Ssd() (*SsdInfo, error) {
	return ssdSystemCall(context.Background())
}

func ssdSystemCall(ctx context.Context) (*SsdInfo, error) {
	exp, err := compilePattern(ssdExp)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(context.Background(), SSD_CMD, SSD_ARG1, SSD_ARG2)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
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
		errBuf := make([]byte, 0, 51200)
		stderr.Read(errBuf)
		if exiterr, ok := err.(*exec.ExitError); ok {
			// TODO: for some reason command always return 4
			if exiterr.ExitCode() != 4 {
				fmt.Printf("\n%s return error: %s\n%v\n", SSD_CMD, string(errBuf), exiterr)
				return nil, err
			}
		}
	}

	ret := SsdInfo{}
	for e := 0; e < len(exp); e++ {

		if exp[e].expression.Key == media_errors {
			ret.MediaErrors = exp[e].expression.Value
		}
		if exp[e].expression.Key == power_hours {
			ret.PowerHours = exp[e].expression.Value
		}
		if exp[e].expression.Key == power_cycle {
			ret.PowerCycle = exp[e].expression.Value
		}
		if exp[e].expression.Key == host_write {
			ret.HostWrite = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == host_read {
			ret.HostRead = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == data_unit_writeTB {
			ret.UnitWriteTB = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == data_unit_write {
			ret.UnitWrite = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == data_unit_readTB {
			ret.UnitReadTB = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == data_unit_read {
			ret.UnitWrite = strings.ReplaceAll(exp[e].expression.Value, ",", "")
		}
		if exp[e].expression.Key == percentage_used {
			ret.PercentageUsed = exp[e].expression.Value
		}
		if exp[e].expression.Key == spare {
			ret.AvailableSpare = exp[e].expression.Value
		}
		if exp[e].expression.Key == spare_threshold {
			ret.AvailableSpareThreshold = exp[e].expression.Value
		}
	}
	return &ret, nil
}
