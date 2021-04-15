package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

type RecoveryBattery struct {
	FileName string
	Stamp    time.Time
	Host     string
	BatteryInfo
}
type RecoverySsd struct {
	FileName string
	Stamp    time.Time
	Host     string
	SsdInfo
}

type Recovery struct {
	WorkDir string
}

const (
	battPrefix = "battery_"
	ssdPrefix  = "ssd_"
)

func (r *Recovery) SaveBattery(stamp time.Time, host string, bat BatteryInfo) error {
	rec := RecoveryBattery{
		Stamp:       stamp,
		Host:        host,
		BatteryInfo: bat,
	}
	byte, _ := json.Marshal(rec)

	tmpTime := stamp.UnixNano() / int64(time.Millisecond)
	fileName := r.WorkDir + "/" + battPrefix + strconv.FormatInt(tmpTime, 10) + ".json"
	err := ioutil.WriteFile(fileName, byte, 0644)
	return err
}
func (r *Recovery) SaveSsd(stamp time.Time, host string, ssd SsdInfo) error {
	rec := RecoverySsd{
		Stamp:   stamp,
		Host:    host,
		SsdInfo: ssd,
	}
	byte, _ := json.Marshal(rec)

	tmpTime := stamp.UnixNano() / int64(time.Millisecond)
	fileName := r.WorkDir + "/" + ssdPrefix + strconv.FormatInt(tmpTime, 10) + ".json"
	err := ioutil.WriteFile(fileName, byte, 0644)
	return err
}

func (r *Recovery) LoadBattery() ([]RecoveryBattery, error) {
	ret := make([]RecoveryBattery, 0)
	files, err := ioutil.ReadDir(r.WorkDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), battPrefix) && strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile(r.WorkDir + "/" + f.Name())
			ret = append(ret, RecoveryBattery{})
			json.Unmarshal(bytes, &(ret[len(ret)-1]))
			ret[len(ret)-1].FileName = r.WorkDir + "/" + f.Name()
		}
	}
	return ret, nil
}
func (r *Recovery) LoadSsd() ([]RecoverySsd, error) {
	ret := make([]RecoverySsd, 0)
	files, err := ioutil.ReadDir(r.WorkDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ssdPrefix) && strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile(r.WorkDir + "/" + f.Name())
			ret = append(ret, RecoverySsd{})
			json.Unmarshal(bytes, &(ret[len(ret)-1]))
			ret[len(ret)-1].FileName = r.WorkDir + "/" + f.Name()
		}
	}
	return ret, nil
}

func NewRecovery() *Recovery {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	workDir := homeDir + "/macStats"
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		err = os.MkdirAll(workDir, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		info, _ := os.Stat(workDir)
		if !info.IsDir() {
			panic(fmt.Errorf("%s is not a directory", workDir))
		}
	}
	return &Recovery{WorkDir: workDir}
}
