<!-- untitled -->
# About

The barista-contrib project provides [barista](https://barista.run) modules
that are community maintained and complement or enrich the built-in modules. For
documentation of the built-in barista modules and customization examples [visit
the official documentation](https://barista.run).

# Installation

To make use of modules and module providers from pacman-contrib, add it to your `go.mod`:

```sh
go get -u github.com/martinohmann/barista-contrib
```

# Modules

Modules available in barista-contrib:

- [dpms](/modules/dpms): Show and toggle Display Power Management Signaling (DPMS) status.
- [ip](/modules/ip): Shows current public ip in the bar if connected to the internet.
- [keyboard](/modules/keyboard): Shows and control current keyboard layout.
- [updates](/modules/updates): Show available package updates.

# Extensions for built-in modules

The following extensions which provide functionality on top of existing built-in
modules are available:

- [openweathermap](/modules/weather/openweathermap): Provides convenience
  functionality to load OpenWeatherMap API key and configuration from a
  configuration file.
