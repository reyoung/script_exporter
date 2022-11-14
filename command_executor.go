package main

import (
	"github.com/hoisie/mustache"
	"os/exec"
	"time"
)

type commandExecutor struct {
	command     string
	interval    time.Duration
	matrix      map[string][]string
	valueSetter func(spec map[string]string, stdout string)
}

type commandExecutorMatrixItem struct {
	key   string
	value []string
}

func matrixCall(acc map[string]string, items []commandExecutorMatrixItem, callback func(spec map[string]string)) {
	if acc == nil {
		acc = map[string]string{}
	}
	if len(items) == 0 {
		callback(acc)
		return
	}

	item := items[0]
	for _, v := range item.value {
		acc[item.key] = v
		matrixCall(acc, items[1:], callback)
	}
}

func (e *commandExecutor) exec() {
	items := make([]commandExecutorMatrixItem, 0, len(e.matrix))
	for k, v := range e.matrix {
		items = append(items, commandExecutorMatrixItem{
			key:   k,
			value: v,
		})
	}

	getCommand := func(spec map[string]string) string {
		return e.command
	}

	if len(e.matrix) != 0 {
		tmpl := panicT(mustache.ParseString(e.command))
		getCommand = func(spec map[string]string) string {
			return tmpl.Render(spec)
		}
	}

	for {
		matrixCall(nil, items, func(spec map[string]string) {
			command := getCommand(spec)
			cmd := exec.Command("bash", "-c", command)
			buf := panicT(cmd.Output())
			e.valueSetter(spec, string(buf))
		})
		if e.interval == 0 {
			e.interval = 30 * time.Second
		}
		time.Sleep(e.interval)
	}
}
