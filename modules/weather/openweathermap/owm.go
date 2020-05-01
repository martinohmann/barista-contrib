package openweathermap

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"barista.run/modules/weather"
	"barista.run/modules/weather/openweathermap"
)

// ErrAPIKeyMissing is returned by New if the API key is missing in the config.
var ErrAPIKeyMissing = errors.New("apiKey missing")

// Config for openweathermap.
type Config struct {
	APIKey      string  `json:"apiKey"`
	CityID      string  `json:"cityID"`
	CityName    string  `json:"cityName"`
	CountryCode string  `json:"countryCode"`
	ZipCode     string  `json:"zipCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// NewFromConfig creates a new weather.Provider from the config at given path.
func NewFromConfig(path string) (weather.Provider, error) {
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return New(config)
}

// New creates a new weather.Provider from given config.
func New(config Config) (weather.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyMissing
	}

	c := openweathermap.New(config.APIKey)

	switch {
	case config.CityID != "":
		return c.CityID(config.CityID), nil
	case config.ZipCode != "":
		return c.Zipcode(config.ZipCode, config.CountryCode), nil
	case config.CityName != "":
		return c.CityName(config.CityName, config.CountryCode), nil
	default:
		return c.Coords(config.Latitude, config.Longitude), nil
	}
}

// LoadConfig loads the config at path.
func LoadConfig(path string) (Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(buf, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
