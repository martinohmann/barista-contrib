package openweathermap

import (
	"testing"

	"barista.run/modules/weather/openweathermap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		given       Config
		expected    string
		expectedErr error
	}{
		{
			name:        "missing apiKey",
			given:       Config{},
			expectedErr: ErrAPIKeyMissing,
		},
		{
			name: "cityID takes precedence",
			given: Config{
				APIKey:      "secret",
				CityID:      "123",
				CityName:    "Berlin",
				CountryCode: "DE",
				ZipCode:     "12345",
				Latitude:    52.5,
				Longitude:   13.4,
			},
			expected: "appid=secret&id=123",
		},
		{
			name: "zipCode takes precedence",
			given: Config{
				APIKey:      "secret",
				CityName:    "Berlin",
				CountryCode: "DE",
				ZipCode:     "12345",
				Latitude:    52.5,
				Longitude:   13.4,
			},
			expected: "appid=secret&zip=12345%2CDE",
		},
		{
			name: "cityName takes precedence",
			given: Config{
				APIKey:      "secret",
				CityName:    "Berlin",
				CountryCode: "DE",
				Latitude:    52.5,
				Longitude:   13.4,
			},
			expected: "appid=secret&q=Berlin%2CDE",
		},
		{
			name: "use lat/lon if nothing else is configured",
			given: Config{
				APIKey:    "secret",
				Latitude:  52.5,
				Longitude: 13.4,
			},
			expected: "appid=secret&lat=52.500000&lon=13.400000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider, err := New(test.given)

			if test.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, test.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.IsType(t, provider, openweathermap.Provider(""))
				url := string(provider.(openweathermap.Provider))
				assert.Equal(t, "https://api.openweathermap.org/data/2.5/weather?"+test.expected, url)
			}
		})
	}
}

func TestNewFromConfig(t *testing.T) {
	provider, err := NewFromConfig("testdata/config.json")
	require.NoError(t, err)
	require.IsType(t, provider, openweathermap.Provider(""))
	url := string(provider.(openweathermap.Provider))
	assert.Equal(t, "https://api.openweathermap.org/data/2.5/weather?appid=secret&id=123", url)

	_, err = NewFromConfig("testdata/nonexistent.json")
	require.Error(t, err)
}
