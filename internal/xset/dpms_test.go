package xset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseXsetQuery(t *testing.T) {
	given := []byte(`Keyboard Control:
  auto repeat:  on    key click percent:  0    LED mask:  00000000
  XKB indicators:
    00: Caps Lock:   off    01: Num Lock:    off    02: Scroll Lock: off
    03: Compose:     off    04: Kana:        off    05: Sleep:       off
    06: Suspend:     off    07: Mute:        off    08: Misc:        off
    09: Mail:        off    10: Charging:    off    11: Shift Lock:  off
    12: Group 2:     off    13: Mouse Keys:  off
  auto repeat delay:  660    repeat rate:  25
  auto repeating keys:  00ffffffdffffbbf
                        fadfffefffedffff
                        9fffffffffffffff
                        fff7ffffffffffff
  bell percent:  50    bell pitch:  400    bell duration:  100
Pointer Control:
  acceleration:  2/1    threshold:  4
Screen Saver:
  prefer blanking:  yes    allow exposures:  yes
  timeout:  1200    cycle:  1200
Colors:
  default colormap:  0x22    BlackPixel:  0x0    WhitePixel:  0xffffff
Font Path:
  /usr/share/fonts/TTF,built-ins
DPMS (Energy Star):
  Standby: 1200    Suspend: 1200    Off: 1200
  DPMS is Enabled
  Monitor is On
`)

	enabled, err := parseDPMSStatus(given)
	require.NoError(t, err)
	assert.True(t, enabled)

	enabled, err = parseDPMSStatus([]byte(`  DPMS is Disabled`))
	require.NoError(t, err)
	assert.False(t, enabled)

	_, err = parseDPMSStatus([]byte(`invalid`))
	require.Error(t, err)
}
