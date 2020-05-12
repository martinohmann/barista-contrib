package pacman

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/martinohmann/barista-contrib/modules/updates"
)

// New creates a new *updates.Module with the pacman provider.
func New() *updates.Module {
	return updates.New(Provider)
}

// Provider is an updates.Provider which checks for pacman updates.
var Provider = updates.ProviderFunc(func() (updates.Info, error) {
	out, err := exec.Command("checkupdates").Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ProcessState.ExitCode() == 2 {
				// exit code 2 is not an error but signals that
				// there are no updates available right now.
				err = nil
			}
		}

		return updates.Info{}, err
	}

	details, err := ParsePackageDetails(out)
	if err != nil {
		return updates.Info{}, err
	}

	info := updates.Info{
		Updates:        len(details),
		PackageDetails: details,
	}

	return info, nil
})

// ParsePackageDetails parses package details from pacman compatible output of
// the form "packageName currentVersion -> targetVersion" and returns the
// package details. Returns an error if raw contains malformed lines.
func ParsePackageDetails(raw []byte) (updates.PackageDetails, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))

	details := updates.PackageDetails{}

	for scanner.Scan() {
		var detail updates.PackageDetail

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		_, err := fmt.Sscanf(line, "%s %s -> %s", &detail.PackageName, &detail.CurrentVersion, &detail.TargetVersion)
		if err != nil {
			return nil, err
		}

		details = append(details, detail)
	}

	return details, nil
}
