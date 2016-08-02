## i3bargo

Status bar replacement for [i3status](http://i3wm.org/i3status/). Works with `i3bar`.

Heavily inspired by [damog's mastodon bar](https://github.com/damog/mastodon). Most of the module presets and templating are straight copies.

It includes click events and borders.

## Config
Configuration can be done with a file called `i3bargo.conf` at `~/.config/` or in your `XDG_CONFIG_HOME` directory.

An example:

```
order: [cpu, memory, disk, ip, loadavg, uptime, clock]

modules:
  cpu:
    label: 
    borders: 0 0 2 0 
    border-color: '#ffffff'
    interval: 2
    onclick: gnome-system-monitor -p

  memory:
    label: 
    borders: 0 0 2 0 
    border-color: '#00ff00'
    interval: 10
    onclick: gnome-system-monitor -r

  disk:
    label: 
    borders: 0 0 2 0 
    interval: 100
    border-color: '#ffff00'
    interval: 20
    onclick: gnome-system-monitor -f

  ip:
    label: 
    borders: 0 0 2 0 
    border-color: '#00ffff'
    interval: 5

  loadavg:
    label: 
    borders: 0 0 2 0 
    border-color: '#ff0000'
    interval: 5

  uptime:
    label: 
    borders: 0 0 2 0 
    border-color: '#909090'
    interval: 60

  clock:
    label: 
    interval: 1
    borders: 0 0 2 0 
    border-color: '#ffffff'


dateformat: 15:04:05 Monday 02 Jan 

networkinterface: enp2s0

color: '#cccccc'

barempty: □
barfull: ■ 

interval: 1
```

This produces the following bar:

![](http://i.imgur.com/3oyWLiG.png)

Default settings:
```
// Bar defaults
Interval         = 2
Battery          = 0
Date_format      = 2006-01-02 15:04:05" // Following format: Mon Jan 2 15:04:05 -0700 MST 2006
NetworkInterface = eth0"
BarSize          = 10
BarStart         = ''
BarEnd           = ''
BarEmpty         = ' '
BarFull          = '#'
ColorBad         = '#d00000'
ColorGood        = '#00d000'
Color            = '#cccccc'

// Block defaults
Color       = '#ffffff'
BorderColor = '#ffffff'
Label       = ''
Borders = 0 0 0 0

```

Currently, ony the following modules are available, using these templates by default:
```
battery:  {{if .battery}}{{.prefix}} {{.bar}} ({{.remaining}} {{.wattage}}W){{else}}No battery{{end}}
clock:    {{.time}}
cpu:      {{.bar}}
disk:     {{.bar}}
hostname: {{.hostname}}
ip:       {{.ip}}
loadavg:  {{.fifteen}} {{.five}} {{.one}}
memory:   {{.bar}}
uptime:   {{.uptime}}

```

Todo: add custom command

Each block runs in its own goroutine, with a different speed (specified by the interval option). 
But only the global interval option will make the bar update. An interval of 1 is desired when using for example seconds in the clock module.
