package i3bargo

import "fmt"

type ClickMessage struct {
	Name     string `json:"name,omitempty"`
	Instance string `json:"instance,omitempty"`
	Button   int    `json:"button"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func Header() string {
	return "{\"version\":1,\"click_events\":true}\n["
}

func Send(input string) {
	fmt.Println(input)
}
