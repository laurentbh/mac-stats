package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
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
	fmt.Printf("\u001b[32msaving battery stats: \u001b[34m%s\u001b[0m\n", fileName)
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
	fmt.Printf("\u001b[32msaving SSD stats: \u001b[34m%s\u001b[0m\n", fileName)
	return err
}

func (r *Recovery) LoadBattery() ([]RecoveryBattery, error) {
	ret := make([]RecoveryBattery, 0)
	files, err := ioutil.ReadDir(r.WorkDir)
	if err != nil {
		return nil, err
	}
	sort.Sort(filesByTime(files))
	for _, f := range files {
		if strings.HasPrefix(f.Name(), battPrefix) && strings.HasSuffix(f.Name(), ".json") {
			fmt.Printf("\u001b[32mrecovering: \u001b[34m%s\u001b[0m\n", f.Name())
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
	sort.Sort(filesByTime(files))
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ssdPrefix) && strings.HasSuffix(f.Name(), ".json") {
			fmt.Printf("\u001b[32mrecovering: \u001b[34m%s\u001b[0m\n", f.Name())
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

type filesByTime []fs.FileInfo
func (a filesByTime) Len() int {
	return len(a)
}

func (a filesByTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a filesByTime) Less (i, j int) bool {
	return a[i].ModTime().Before( a[j].ModTime())
}
