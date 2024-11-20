package pinboard

import (
	"log/slog"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var entries = sync.Map{}

func _l(msg string) func() {
	pc, _, _, _ := runtime.Caller(1)
	var c *atomic.Uint32
	v, ok := entries.Load(pc)
	if !ok {
		c = &atomic.Uint32{}
	} else {
		c = v.(*atomic.Uint32)
	}
	var cnt string
	id := c.Add(1)
	cnt = strconv.Itoa(int(id))
	// debug.PrintStack()
	stack := string(debug.Stack())
	st := strings.Split(stack, "\n")
	slog.Debug("++++" + msg + "(" + cnt + ") " + strings.Join(st[7:8], " "))
	return func() {
		c.Add(^(id - 1))
		slog.Debug("----" + msg + "(" + cnt + ") " + strings.Join(st[7:10], " "))
	}
}
