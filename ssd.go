package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type SsdInfo struct {
	AvailableSpare          float64
	AvailableSpareThreshold float64
	PercentageUsed          float64
	UnitRead                float64
	UnitReadTB              float64
	UnitWrite               float64
	UnitWriteTB             float64
	HostRead                float64
	HostWrite               float64
	PowerCycle              float64
	PowerHours              float64
	MediaErrors             float64
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

func Ssd(cfg Config) (*SsdInfo, error) {
	return ssdSystemCall(context.Background(), cfg)
}

func ssdSystemCall(ctx context.Context, cfg Config) (*SsdInfo, error) {
	exp, err := compilePattern(ssdExp)
	if err != nil {
		return nil, err
	}
	exe := cfg.Smart.Path + "/" + SSD_CMD
	cmd := exec.CommandContext(context.Background(), exe, SSD_ARG1, SSD_ARG2)

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
	done := make(chan int)
	go parseToLines(stdout, data)
	go processLines(ctx, exp, data, done)

	ok := <-done
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
	if ok != 0 {
		return nil, fmt.Errorf("some SSD regexps where not found")
	}

	ret := SsdInfo{}
	for e := 0; e < len(exp); e++ {

		if exp[e].expression.Key == media_errors {
			ret.MediaErrors, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
		if exp[e].expression.Key == power_hours {
			ret.PowerHours, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
		if exp[e].expression.Key == power_cycle {
			ret.PowerCycle, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
		if exp[e].expression.Key == host_write {
			ret.HostWrite, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == host_read {
			ret.HostRead, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == data_unit_writeTB {
			ret.UnitWriteTB, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == data_unit_write {
			ret.UnitWrite, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == data_unit_readTB {
			ret.UnitReadTB, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == data_unit_read {
			ret.UnitRead, _ = strconv.ParseFloat(strings.ReplaceAll(exp[e].expression.Value, ",", ""), 64)
		}
		if exp[e].expression.Key == percentage_used {
			ret.PercentageUsed, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
		if exp[e].expression.Key == spare {
			ret.AvailableSpare, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
		if exp[e].expression.Key == spare_threshold {
			ret.AvailableSpareThreshold, _ = strconv.ParseFloat(exp[e].expression.Value, 64)
		}
	}
	return &ret, nil
}
