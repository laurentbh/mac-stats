package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	batCycle = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bat_cycle",
		Help: "Current battery cycles count.",
	})
	batCharge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bat_charge",
		Help: "Full battery charge in mAh.",
	})
	ssdSpare = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_spare",
		Help: "Percentage of available spare.",
	})
	ssdSpareThreshold = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_spare_threshold",
		Help: "Percentage of Available Spare Threshol.",
	})
	ssdPercentUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_percent_used",
		Help: "Percentage Used.",
	})
	ssdUnitWrite = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_unit_write",
		Help: "Data Units Written.",
	})
	ssdUnitRead = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_unit_read",
		Help: "Data Units Read.",
	})
	ssdUnitWriteTB = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_unit_write_tb",
		Help: "Data Units Written in TB.",
	})
	ssdUnitReadTB = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_unit_read_tb",
		Help: "Data Units Read in TB.",
	})
	ssdHostRead = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_host_read",
		Help: "Host Read Commands.",
	})
	ssdHostWrite = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_host_write",
		Help: "Host Read Commands.",
	})
	ssdPowerCycle = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_power_cycle",
		Help: "Power Cycles.",
	})
	ssdPowerHours = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_power_hours",
		Help: "Power On Hours.",
	})
	ssdMediaErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ssd_media_errors",
		Help: "Media and Data Integrity Error.",
	})
)

func init() {
	prometheus.MustRegister(ssdHostRead)
	prometheus.MustRegister(ssdHostWrite)
	prometheus.MustRegister(ssdMediaErrors)
	prometheus.MustRegister(ssdPercentUsed)
	prometheus.MustRegister(ssdPowerCycle)
	prometheus.MustRegister(ssdPowerHours)
	prometheus.MustRegister(ssdSpare)
	prometheus.MustRegister(ssdSpareThreshold)
	prometheus.MustRegister(ssdUnitRead)
	prometheus.MustRegister(ssdUnitReadTB)
	prometheus.MustRegister(ssdUnitWrite)
	prometheus.MustRegister(ssdUnitWriteTB)
	prometheus.MustRegister(batCharge)
	prometheus.MustRegister(batCycle)
}

func PushCounter(host string, battery BatteryInfo, ssd SsdInfo) error {
	p := push.New("http://localhost:9091", "mac_stats."+host)

	batCharge.SetToCurrentTime()
	batCharge.Set(battery.fullCharge)
	p.Collector(batCharge)

	batCycle.SetToCurrentTime()
	batCycle.Set(battery.cycle)
	p.Collector(batCycle)

	p.Grouping("battery", "info")
	err := p.Push()

	p = push.New("http://localhost:9091", "mac_stats."+host)
	ssdSpare.SetToCurrentTime()
	ssdSpare.Set(ssd.AvailableSpare)
	p.Collector(ssdSpare)

	ssdSpareThreshold.SetToCurrentTime()
	ssdSpareThreshold.Set(ssd.AvailableSpareThreshold)
	p.Collector(ssdSpareThreshold)

	ssdPercentUsed.SetToCurrentTime()
	ssdPercentUsed.Set(ssd.PercentageUsed)
	p.Collector(ssdPercentUsed)

	ssdUnitRead.SetToCurrentTime()
	ssdUnitRead.Set(ssd.UnitRead)
	p.Collector(ssdUnitRead)

	ssdUnitReadTB.SetToCurrentTime()
	ssdUnitReadTB.Set(ssd.UnitReadTB)
	p.Collector(ssdUnitReadTB)

	ssdUnitWrite.SetToCurrentTime()
	ssdUnitWrite.Set(ssd.UnitWrite)
	p.Collector(ssdUnitWrite)

	ssdUnitWriteTB.SetToCurrentTime()
	ssdUnitWriteTB.Set(ssd.UnitWriteTB)
	p.Collector(ssdUnitWriteTB)

	ssdHostRead.SetToCurrentTime()
	ssdHostRead.Set(ssd.HostRead)
	p.Collector(ssdHostRead)

	ssdHostWrite.SetToCurrentTime()
	ssdHostWrite.Set(ssd.HostWrite)
	p.Collector(ssdHostWrite)

	ssdMediaErrors.SetToCurrentTime()
	ssdMediaErrors.Set(ssd.MediaErrors)
	p.Collector(ssdMediaErrors)

	ssdPowerCycle.SetToCurrentTime()
	ssdPowerCycle.Set(ssd.PowerCycle)
	p.Collector(ssdPowerCycle)

	ssdPowerHours.SetToCurrentTime()
	ssdPowerHours.Set(ssd.PowerHours)
	p.Collector(ssdPowerHours)

	p.Grouping("ssd", "info")
	err = p.Push()
	if err != nil {
		return err
	}
	return nil

}
