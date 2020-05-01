package xkbmap

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
)

var xkbInfoRegexp = regexp.MustCompile(`([^:]*?)\s*:\s*(.*)$`)

// Info contains information about the current keyboard layout.
type Info struct {
	Rules  string
	Model  string
	Layout string
}

// Query retrieves keyboard information using setxkbmap -query.
func Query() (Info, error) {
	output, err := exec.Command("setxkbmap", "-query").Output()
	if err != nil {
		return Info{}, err
	}

	return parseQueryOutput(output), nil
}

// SetLayout sets the keyboard layout.
func SetLayout(layout string) error {
	return exec.Command("setxkbmap", layout).Run()
}

func parseQueryOutput(raw []byte) Info {
	scanner := bufio.NewScanner(bytes.NewReader(raw))

	info := Info{}

	for scanner.Scan() {
		submatches := xkbInfoRegexp.FindStringSubmatch(scanner.Text())
		if submatches == nil {
			continue
		}

		key := submatches[1]
		value := submatches[2]

		switch key {
		case "rules":
			info.Rules = value
		case "model":
			info.Model = value
		case "layout":
			info.Layout = value
		}
	}

	return info
}
