---
title: OpenWeatherMap
---

Convenience extension for the
[OpenWeatherMap](https://barista.run/modules/weather/openweathermap) provider
of the built-in [weather module](https://barista.run/modules/weather).

## Configuration

* `NewFromConfig(string) (openweathermap.Provider, error)`: Creates a new openweathermap provider from a config file.

* `New(Config) (openweathermap.Provider, error)`: Creates a new openweathermap provider from config.

## Examples

<div class="module-example-out">14°C, few clouds</div>
Load openweathermap API key and config from a file and create a new  weather module:

```go
configFile := configdir.LocalConfig("i3/barista/openweathermap.json")

owm, err := openweathermap.NewFromConfig(configFile)
if err != nil {
  log.Fatal(err)
}

weather.New(owm).Output(func(info weather.Weather) bar.Output {
  return outputs.Textf("%.0f°C, %s", info.Temperature.Celsius(), info.Description)
})
```
