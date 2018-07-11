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

func runXrandrOn(target, format string) error {
	command := fmt.Sprintf(format, target)
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
		err := runXrandrOn(internalScreen, "xrandr --output %s --off")
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 2)
	}

	if externalConnected {
		err := runXrandrOn(externalScreen, "xrandr --output %s --fb 6880x2880 --panning 6880x2880 --auto --scale 2x2 --mode 3440x1440 --pos 0x0")
		if err != nil {
			return err
		}
	}

	if hdmiConnected {
		err := runXrandrOn(hdmiScreen, "xrandr --output %s --fb 3840x2160 --panning 3840x2160 --auto --scale 2x2 --mode 1920x1080 --pos 0x0")
		if err != nil {
			return err
		}
	}

	return nil
}

func disableExternalScreens() error {
	if isConnected(externalScreen) {
		err := runXrandrOn(externalScreen, "xrandr --output %s --off")
		if err != nil {
			return err
		}
	}

	if isConnected(hdmiScreen) {
		err := runXrandrOn(hdmiScreen, "xrandr --output %s --off")
		if err != nil {
			return err
		}
	}

	if isConnected(internalScreen) {
		time.Sleep(time.Second * 2)

		err := runXrandrOn(internalScreen, "xrandr --output %s --auto --scale 1x1 --mode 3840x2160 --pos 0x0 --fb 3840x2160 --panning 3840x2160")
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
