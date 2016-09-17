package main

import (
	"errors"
	"flag"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func main() {
	lowBatteryLimit := flag.Int("limit", 10, "low battery percentage limmit,must over 1%")

	if *lowBatteryLimit < 1 {
		log.Fatalln("low battery percentage must over 1%!")
	}

	if err := supportACPI(); err != nil {
		log.Fatalln(err)
	}

	for {
		if baterryPercentage, err := fetchBattery(); err != nil {
			log.Println(err.Error())
		} else if baterryPercentage <= 10 {
			// TODO email or sms notify
			shutdownOutput, err := shutdown()
			if err != nil {
				log.Println(err)
			}
			log.Println(shutdownOutput)
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

// fetchBattery 获取电池电量信息
func fetchBattery() (int, error) {
	var (
		batteryPer int = 0
		err error = nil
		batteryStr string = ""
		regexpIns  *regexp.Regexp
	)

	regexpIns, err = regexp.Compile(`^Battery\s\d{1}:\s(\w+),\s([0-9]+){1,3}%([0-9a-z,:\s]+)?`)
	if err != nil {
		return batteryPer, err
	}

	batteryStr, err = fetchBatteryFromACPI()
	if err != nil {
		return batteryPer, err
	}

	batteryStr = regexpIns.ReplaceAllString(batteryStr, `${2}`)
	batteryPer, err = strconv.Atoi(batteryStr)

	return batteryPer, err
}

// supportACPI 是否支持ACPI命令
func supportACPI() error {
	if _, err := exec.LookPath("acpi"); err != nil {
		return errors.New("acpi not found in $PATH,you can use `apt-get install acpi` to intall acpi if u use Ubuntu!")
	}

	return nil
}

// fetchBatteryFromACPI  通过系统的ACPI命令获取系统电池的电量信息
func fetchBatteryFromACPI() (string, error) {
	acpiCmd := exec.Command("acpi", "b")
	outputByte, err := acpiCmd.Output()
	if err != nil {
		return "", err
	}

	return string(outputByte), err
}

// shutdown 关机命令
func shutdown() (string, error) {
	shutdownCmd := exec.Command("shutdown", "h", "now")
	outputByte, err := shutdownCmd.Output()
	if err != nil {
		return "", err
	}

	return string(outputByte), nil
}
