package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

/*
重置navicat-premium试用时间
参考项目：https://github.com/ismailatak/navicat-premium-trial-reset-go
*/

func main() {
	// Detect Navicat Premium version
	cmd := exec.Command("defaults", "read", "/Applications/Navicat Premium.app/Contents/Info.plist")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Please make sure Navicat Premium is installed")
		fmt.Println("Error reading Info.plist:", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`CFBundleShortVersionString = "([^\.]+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("Error detecting Navicat Premium version")
		os.Exit(1)
	}

	version := matches[1]
	fmt.Printf("Detected Navicat Premium version %s\n", version)

	var file string

	// Select plist file by version
	switch version {
	case "17", "16":
		file = os.Getenv("HOME") + "/Library/Preferences/com.navicat.NavicatPremium.plist"
	case "15":
		file = os.Getenv("HOME") + "/Library/Preferences/com.prect.NavicatPremium15.plist"
	default:
		fmt.Printf("Unsupported Navicat Premium version: %s\n", version)
		os.Exit(1)
	}

	// File exists
	exists := exec.Command("ls", "-l", "-a", file)
	output, err = exists.Output()
	if err != nil || len(output) == 0 {
		fmt.Println("File does not exist or is empty:", err)
		os.Exit(1)
	}

	fmt.Println("Resetting trial time...")

	// Delete hash from plist file
	cmd = exec.Command("defaults", "read", file)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error reading plist file (%s): %+v\n", file, err)
		os.Exit(1)
	}

	re = regexp.MustCompile(`([0-9A-Z]{32}) = `)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		hash := matches[1]
		fmt.Printf("deleting %s array...\n", hash)
		err = exec.Command("defaults", "delete", file, hash).Run()
		if err != nil {
			fmt.Println("Error deleting hash:", err)
			os.Exit(1)
		}
	}

	// Delete hidden file in Application Support
	appSupport := os.Getenv("HOME") + "/Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium/"
	cmd = exec.Command("ls", "-a", appSupport)
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error reading hidden file:", err)
		os.Exit(1)
	}

	re = regexp.MustCompile(`\.([0-9A-Z]{32})`)
	matches = re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		hash2 := matches[1]
		fmt.Printf("deleting %s folder...\n", hash2)
		err = exec.Command("rm", appSupport+"."+hash2).Run()
		if err != nil {
			fmt.Println("Error deleting hidden file:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Done")
}
