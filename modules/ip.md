---
title: IP
---

Show public IP address: `ip.New(someProvider)`.

IP supports displaying currently available updates using pluggable providers, with the ability to add custom providers fairly easily. Provider is just

```go
type Provider interface {
	// GetIP retrieves the current public client IP. Must return nil for both
	// return values if there is no internet connection.
	GetIP() (net.IP, error)
}
```

The following IP providers are available in barista-contrib:

* [ipify](https://godoc.org/github.com/martinohmann/barista-contrib/modules/ip/ipify): Uses [ipify](https://ipify.org) to fetch the public IP address.

## Configuration

* `Output(func(ip.Info) bar.Output)`: Sets the output format.

* `Every(time.Duration)`: Sets the interval to wait before checking the IP address again. Defaults to 10 minutes.

* `Refresh()`: Fetches the IP address and refreshes the module.

## Examples

<div class="module-example-out">online: 1.2.3.4</div>
Show online status:

```go
ip.New(ipify.Provider).Output(func(info ip.Info) bar.Output {
    if info.Connected() {
        return outputs.Textf("online: %s", info.IP)
    }

    return outputs.Text("offline")
})
```

## Data: `type Info struct`

### Fields

* `IP net.IP`: The current public IP address, `nil` if not connected.

### Methods

* `Connected() bool`: Returns `true` if connected to the internet.
