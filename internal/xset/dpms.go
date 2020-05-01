package xset

import (
	"errors"
	"os/exec"
	"regexp"
)

var dpmsRegexp = regexp.MustCompile(`(?m)^\s*DPMS is\s+(.*)$`)

// SetDPMS enables or disables DPMS.
func SetDPMS(enabled bool) error {
	arg := "-dpms"
	if enabled {
		arg = "+dpms"
	}

	return exec.Command("xset", arg).Run()
}

// GetDPMS retrieves the current DPMS status.
func GetDPMS() (bool, error) {
	out, err := exec.Command("xset", "-q").Output()
	if err != nil {
		return false, err
	}

	return parseDPMSStatus(out)
}

func parseDPMSStatus(raw []byte) (bool, error) {
	match := dpmsRegexp.FindStringSubmatch(string(raw))
	if match == nil {
		return false, errors.New("failed to match DPMS status")
	}

	return match[1] == "Enabled", nil
}
