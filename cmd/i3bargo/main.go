package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	. "github.com/nananas/i3bargo"
)

type JsonAny interface{}

func main() {
	log.SetFlags(log.Lshortfile)
	// config := NewConfig()
	config := ReadConfig()

	// check order
	for _, e := range *config.Order {
		if _, ok := config.Modules[e]; !ok {
			log.Fatal("Order contains module <" + e + "> that is not specified")
		}
	}

	// start new goroutine to handle Stdin
	go func() {
		for {
			bio := bufio.NewReader(os.Stdin)
			line, _, _ := bio.ReadLine()

			var m ClickMessage

			if line[0] == ',' {
				line = line[1:]
			}

			err := json.Unmarshal(line, &m)
			if err != nil {
			} else {
				if config.Modules[m.Name].Onclick != nil {

					err := config.Modules[m.Name].Onclick.Start()
					if err != nil {
						log.Fatal(err)
					}
				}

			}

			time.Sleep(1)
		}
	}()

	// Start goroutines for each block, on its own interval speed
	for _, e := range *config.Order {
		block := config.Modules[e]

		if block.Command != nil {
			go func() {
				for {
					si := block.Command(config, block)
					block.WriteResult(si)
					Sleep(1)
				}

			}()
		}

	}

	// MAIN LOOP
	jsonArray := make([]map[string]JsonAny, len(*config.Order))
	Send(Header())
	for {

		for idx, e := range *config.Order {
			block := config.Modules[e]

			// r := block.ReadResult()
			si := block.ReadResult()
			log.Println(si)
			color := block.Color
			if si.IsBad() {
				color = config.ColorBad
			}

			jsonArray[idx] = map[string]JsonAny{
				"full_text":     block.Label + "  " + si.FullText,
				"color":         color,
				"border":        block.BorderColor,
				"border_top":    block.Borders[0],
				"border_right":  block.Borders[1],
				"border_bottom": block.Borders[2],
				"border_left":   block.Borders[3],
				"name":          block.Name,
			}

		}

		jsonData, _ := json.Marshal(jsonArray)
		Send(string(jsonData) + ",")

		Sleep(config.Interval)

	}

}
