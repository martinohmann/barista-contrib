---
title: Microphone Amplitude (micamp)
---

Show the amplitude passing through the microphone: `micamp.New(ctx, "Blue Microphones")`.

It is used to give a visual indication that audio is passing through the microphone. The only provider currently supported is pulseaudio.

## Configuration

Default output will be:

- When the mic is muted or doesn't receive any audio:

  ```
  NaN  .......... 🎙
  ```

- When the mic is receiving audio:

  ```
  0%   .......... 🎙
  50%  :::::..... 🎙
  100% :::::::::: 🎙
  ```

- When the amplitude isn't between 0-100:
  `ERR 200% (amp=2.000) 🎙`

Parameters:

- `ctx`: context which when done will stop the stream and gracefully shut down the pulse audio client.

- `micSourceNamePrefix`: The prefix of the microphone name as seen by the pulse audio description (which can be found using `pactl list sources`). If it's empty (`""`) then the pulse audio default source will be used.

## Examples

<div class="module-example-out">1%   .......... 🎙</div>

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
defer cancel()

micamp.New(ctx, "Blue Microphones")
```
