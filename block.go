package i3bargo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"text/template"
	"time"
)

type Block struct {
	Name        string
	Command     StatusSource
	Interval    int
	Label       string
	Color       string
	Borders     []int
	BorderColor string
	Template    *template.Template
	Onclick     *exec.Cmd

	Result *StatusInfo

	Lock sync.Mutex // command result
	// BackgroundColor string
}

type StatusSource func(*Config, *Block) *StatusInfo

var Presets = map[string]StatusSource{
	"battery":  Battery,
	"clock":    Clock,
	"cpu":      CPU,
	"disk":     Disk,
	"hostname": Hostname,
	"ip":       IPAddress,
	"loadavg":  LoadAvg,
	"memory":   Memory,
	"uptime":   Uptime,
}

func (block *Block) ReadResult() *StatusInfo {
	block.Lock.Lock()
	si := block.Result
	block.Lock.Unlock()

	return si
}

func (block *Block) WriteResult(si *StatusInfo) {
	block.Lock.Lock()
	block.Result = si
	block.Lock.Unlock()
}

func LoadAvg(c *Config, b *Block) *StatusInfo {
	one, five, fifteen := ReadLoadAvg()
	data := make(map[string]string)
	data["one"] = fmt.Sprintf("%.2f", one)
	data["five"] = fmt.Sprintf("%.2f", five)
	data["fifteen"] = fmt.Sprintf("%.2f", fifteen)
	si := NewStatus(b.Template, data)
	cpu := float64(runtime.NumCPU())
	if one > cpu {
		si.Status = STATUS_BAD
	}
	return si
}

func Clock(c *Config, b *Block) *StatusInfo {
	data := make(map[string]string)
	data["time"] = time.Now().Format(c.DateFormat)
	si := NewStatus(b.Template, data)
	return si
}

func IPAddress(c *Config, b *Block) *StatusInfo {
	data := make(map[string]string)
	data["ip"] = IfaceAddr(c.NetworkInterface)
	si := NewStatus(b.Template, data)
	return si
}

func Hostname(c *Config, b *Block) *StatusInfo {
	data := make(map[string]string)
	data["hostname"], _ = os.Hostname()
	si := NewStatus(b.Template, data)
	return si
}

func Uptime(c *Config, b *Block) *StatusInfo {
	data := make(map[string]string)
	buf := new(syscall.Sysinfo_t)
	syscall.Sysinfo(buf)
	data["uptime"] = HumanDuration(buf.Uptime)
	si := NewStatus(b.Template, data)
	return si
}
