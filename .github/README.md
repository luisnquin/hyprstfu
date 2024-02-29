# Hyprstfu

> [!NOTE]  
> As the repository description indicates, you'll need to use a **Pulseaudio** or **Pipewire**(with *pipewire-pulse* module) backend.


## Purpose

Ever found yourself wishing to silence certain things without causing a big fuss?

Imagine being on Discord, wanting to mute yourself but still craving peace from all
the surrounding noise without anyone noticing. I've been in that situation myself a
few years back.

Creating a tool like this isn't rocket science. You could whip up a bash script,
tweak some grepping to handle all the Hyprland and pactl outputs. But I opted for
something with better readability and performance(?) for my program.

In my current circunstances I'd even try to do the same thing in Zig but your safe for now. :)

## Features

 - Integrated with Hyprland IPC server.
 - Easy to integrate as a bind in your `hypr/hyprland.conf`.
 - Really lightweight.
 - No need to have **pactl** in your $PATH.

## Install

```sh
# Requires go >=v1.22
$ go install github.com/luisnquin/hyprstfu@latest
```

## LICENSE

[MIT](./LICENSE)
