package i3bargo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

const (
	DEFAULT_BAR_Interval         = 2
	DEFAULT_BAR_Battery          = 0
	DEFAULT_BAR_Date_format      = "2006-01-02 15:04:05" // Following format: Mon Jan 2 15:04:05 -0700 MST 2006
	DEFAULT_BAR_NetworkInterface = "eth0"
	DEFAULT_BAR_BarSize          = 10
	DEFAULT_BAR_BarStart         = ""
	DEFAULT_BAR_BarEnd           = ""
	DEFAULT_BAR_BarEmpty         = " "
	DEFAULT_BAR_BarFull          = "#"
	DEFAULT_BAR_ColorBad         = "#d00000"
	DEFAULT_BAR_ColorGood        = "#00d000"
	DEFAULT_BAR_Color            = "#cccccc"

	// Block defaults
	DEFAULT_Color       = "#ffffff"
	DEFAULT_BorderColor = "#ffffff"
	DEFAULT_Label       = ""
)

var DEFAULT_TEMPLATES = map[string]string{
	"battery":  "{{if .battery}}{{.prefix}} {{.bar}} ({{.remaining}} {{.wattage}}W){{else}}No battery{{end}}",
	"clock":    "{{.time}}",
	"cpu":      "{{.bar}}",
	"disk":     "{{.bar}}",
	"hostname": "{{.hostname}}",
	"ip":       "{{.ip}}",
	"loadavg":  "{{.fifteen}} {{.five}} {{.one}}",
	"memory":   "{{.bar}}",
	"uptime":   "{{.uptime}}",
}

var DEFAULT_Borders = []int{0, 0, 0, 0}

type Config struct {
	Order   *[]string
	Modules map[string]*Block

	Interval         int
	Battery          int // battery ID
	DateFormat       string
	NetworkInterface string
	BarSize          int
	BarStart         string
	BarEnd           string
	BarEmpty         string
	BarFull          string
	ColorBad         string
	ColorGood        string
	Color            string
}

type YAMLconfig struct {
	Order            []string
	DateFormat       string "yaml:dateformat,omitempty"
	NetworkInterface string "yaml:network-interface,omitempty"
	Interval         int    "yaml:,omitempty"
	BarSize          int    "yaml:barsize,omitempty"
	BarStart         string "yaml:barstart,omitempty"
	BarEnd           string "yaml:barend,omitempty"
	BarEmpty         string "yaml:barempty,omitempty"
	BarFull          string "yaml:barfull,omitempty"
	ColorBad         string "yaml:color-bad,omitempty"
	ColorGood        string "yaml:color-good,omitempty"
	Color            string "yaml:color-normal,omitempty"

	Modules map[string]map[string]string
}

func ReadConfig() *Config {

	config := Config{}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}

	configFile := filepath.Join(configHome, "i3bargo.conf")

	if _, err := os.Stat(configFile); err == nil {

		// file, err := os.Open(configFile)
		// if err != nil {
		// log.Fatal(err)
		// }
		// defer file.Close()

		var c YAMLconfig

		bytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(bytes, &c)
		if err != nil {
			log.Fatal(err)
		}

		config.Order = &c.Order
		config.Modules = make(map[string]*Block)

		// Set Defaults
		config.Interval = DEFAULT_BAR_Interval
		config.Battery = DEFAULT_BAR_Battery
		config.DateFormat = DEFAULT_BAR_Date_format
		config.NetworkInterface = DEFAULT_BAR_NetworkInterface
		config.ColorBad = DEFAULT_BAR_ColorBad
		config.ColorGood = DEFAULT_BAR_ColorGood
		config.Color = DEFAULT_BAR_Color
		config.BarSize = DEFAULT_BAR_BarSize
		config.BarStart = DEFAULT_BAR_BarStart
		config.BarEnd = DEFAULT_BAR_BarEnd
		config.BarEmpty = DEFAULT_BAR_BarEmpty
		config.BarFull = DEFAULT_BAR_BarFull

		if c.DateFormat != "" {
			config.DateFormat = c.DateFormat
		}

		if c.NetworkInterface != "" {
			config.NetworkInterface = c.NetworkInterface
		}

		if c.Interval != 0 {
			config.Interval = c.Interval
		}

		if c.BarSize != 0 {
			config.BarSize = c.BarSize
		}

		if c.BarStart != "" {
			config.BarStart = c.BarStart
		}

		if c.BarEnd != "" {
			config.BarEnd = c.BarEnd
		}

		if c.BarEmpty != "" {
			config.BarEmpty = c.BarEmpty
		}

		if c.BarFull != "" {
			config.BarFull = c.BarFull
		}

		if c.Color != "" {
			config.Color = c.Color
		}

		for k, v := range c.Modules {
			block := Block{}
			block.Name = k

			// Set Defaults
			block.Label = DEFAULT_Label
			block.Interval = config.Interval
			block.Color = config.Color
			block.Borders = DEFAULT_Borders
			block.BorderColor = DEFAULT_BorderColor

			if t, ok := DEFAULT_TEMPLATES[k]; ok {
				temp := template.New(k)
				block.Template, err = temp.Parse(t)
			}

			// LABEL
			if d, ok := v["label"]; ok {
				block.Label = d
			}

			// COMMAND
			if d, ok := v["command"]; ok {
				block.Command = NewCommand(d)
			} else if d, ok := Presets[block.Name]; ok {
				block.Command = d
			}

			// INTERVAL
			if d, ok := v["interval"]; ok {
				i, err := strconv.Atoi(d)
				if err != nil {
					log.Fatal("interval for " + k + " has to be a number!")
				} else {
					block.Interval = i
				}
			}

			// COLOR
			if d, ok := v["color"]; ok {
				block.Color = d
			}

			// BORDER COLOR
			if d, ok := v["border-color"]; ok {
				block.BorderColor = d
			}

			// BORDER
			if d, ok := v["borders"]; ok {
				parts := strings.Split(d, " ")
				iparts := make([]int, 4)
				for i, e := range DEFAULT_Borders {
					iparts[i] = e

				}

				for i, e := range parts[:len(parts)] {
					b, err := strconv.Atoi(e)
					if err != nil {
						log.Fatal("interval for " + k + " has to be a number!")
						break
					}

					iparts[i] = b
				}

				block.Borders = iparts
			}

			if d, ok := v["template"]; ok {
				t := template.New(block.Name)
				block.Template, err = t.Parse(d)
				if err != nil {
					fmt.Println("Bad template")
					panic(err)
				}
			}

			if d, ok := v["onclick"]; ok {
				split := strings.Split(d, " ")
				block.Onclick = exec.Command(split[0], split[1:]...)
			}

			if block.Command != nil {
				block.Result = block.Command(&config, &block)
			} else {
				block.Result = &StatusInfo{block.Label, STATUS_GOOD}
			}

			config.Modules[k] = &block
		}

		return &config

	} else {
		log.Fatal(err)
	}

	return nil
}
