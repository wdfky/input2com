package main

import (
	"input2com/internal/input"
	"input2com/internal/macros"
	"time"
)

type PluginImpl struct{}

func (p *PluginImpl) RegisterMacros() map[string]macros.Macro {
	return map[string]macros.Macro{
		"Example Plugin Macro": {
			Name:        "Example Plugin Macro",
			Description: "A macro loaded from a plugin.",
			Fn: func(mk *macros.MacroMouseKeyboard, ch chan bool) {
				for {
					select {
					case <-ch:
						return
					default:
						mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
						time.Sleep(100 * time.Millisecond)
						mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
						time.Sleep(100 * time.Millisecond)
					}
				}
			},
		},
	}
}

var Plugin PluginImpl // 必须导出名为 "Plugin" 的变量
