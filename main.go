package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	externalScreen = "DP-1-1"
	internalScreen = "eDP-1"
	hdmiScreen     = "HDMI-1"
)

type ScreenConfig struct {
	Scale      float64
	Resolution *R
}

var screens = map[string]*ScreenConfig{
	externalScreen: &ScreenConfig{
		Scale:      1,
		Resolution: &R{3440, 1440},
	},
	internalScreen: &ScreenConfig{
		Scale:      1,
		Resolution: &R{1920, 1080},
	},
	hdmiScreen: &ScreenConfig{
		Scale:      1,
		Resolution: &R{2560, 1440},
	},
}

type R struct {
	W float64
	H float64
}

func main() {
	if len(os.Args) == 0 {
		log.Fatal("Please provide a command")
	}

	cmd := strings.ToUpper(os.Args[1])

	switch cmd {
	case "ON":
		must(enableExternalScreens())
	default:
		must(disableExternalScreens())
	}
}

func isConnected(output string) bool {
	shellOut, err := sh("xrandr", "-q")
	must(err)
	re := regexp.MustCompile(fmt.Sprintf("%s connected", output))
	return re.Match(shellOut)
}

func XrandrRun(command string) error {
	args := strings.Split(command, " ")
	out, err := sh(args...)

	if err != nil {
		log.Print(string(out))
	}

	return err

}

func XrandrOff(target string) error {
	template := "xrandr --output %s --off"
	command := fmt.Sprintf(template, target)

	return XrandrRun(command)
}

func XrandrOn(target string, conf *ScreenConfig) error {
	scale := conf.Scale
	resolution := conf.Resolution
	template := "xrandr --output %s --fb %s --panning %s --auto --scale %s --mode %s --pos 0x0"
	fb := fmt.Sprintf("%.fx%.f", resolution.W*scale, resolution.H*scale)
	mode := fmt.Sprintf("%.fx%.f", resolution.W, resolution.H)
	sc := fmt.Sprintf("%.1fx%.1f", scale, scale)
	command := fmt.Sprintf(template, target, fb, fb, sc, mode)
	args := strings.Split(command, " ")
	out, err := sh(args...)

	if err != nil {
		log.Print(string(out))
	}

	return err
}

func enableExternalScreens() error {
	externalConnected := isConnected(externalScreen)
	hdmiConnected := isConnected(hdmiScreen)

	if (externalConnected || hdmiConnected) && isConnected(internalScreen) {
		err := XrandrOff(internalScreen)
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 2)
	}

	if externalConnected {
		err := XrandrOn(externalScreen, screens[externalScreen])
		if err != nil {
			return err
		}
	}

	if hdmiConnected {
		err := XrandrOn(hdmiScreen, screens[hdmiScreen])
		if err != nil {
			return err
		}
	}

	return nil
}

func disableExternalScreens() error {
	if isConnected(externalScreen) {
		err := XrandrOff(externalScreen)
		if err != nil {
			return err
		}
	}

	if isConnected(hdmiScreen) {
		err := XrandrOff(hdmiScreen)
		if err != nil {
			return err
		}
	}

	if isConnected(internalScreen) {
		time.Sleep(time.Second * 2)

		err := XrandrOn(internalScreen, screens[internalScreen])
		if err != nil {
			return err
		}
	}

	return nil
}

func sh(args ...string) ([]byte, error) {
	log.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
