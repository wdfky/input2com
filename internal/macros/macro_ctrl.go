package macros

import (
	"encoding/json"
	"input2com/internal/input"
	"input2com/internal/logger"
	"input2com/internal/serial"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type MacroMouseKeyboard struct {
	MouseBtnArgs map[byte]chan bool
	KeyArgs      map[byte]chan bool
	Ctrl         *serial.ComMouseKeyboard
	Macros       map[string]Macro
}

func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

type Macro struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Fn          func(*MacroMouseKeyboard, chan bool)
}

var Macros = make(map[string]Macro)

func downDragMacro(data [][4]int32) func(mk *MacroMouseKeyboard, ch chan bool) {
	return func(mk *MacroMouseKeyboard, ch chan bool) {
		counter := int32(0)
		mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
		for {
			select {
			case <-ch:
				mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
				return
			default:
				for _, row := range data {
					if row[0] > counter {
						mk.Ctrl.MouseMove(row[1], row[2], 0)
						time.Sleep(time.Duration(row[3]) * time.Millisecond)
						break
					}
				}
				counter++
			}
		}
	}
}

func loadPlugins() {
	pluginsDir := "plugins"
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		logger.Logger.Warnf("Failed to read plugins directory: %v", err)
		return
	}
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
		}
		files = append(files, info)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			loadPlugin(pluginsDir + "/" + file.Name())
		}
	}
}

func NewMacroMouseKeyboard(controler *serial.ComMouseKeyboard) *MacroMouseKeyboard {
	mouseBtnArgs := make(map[byte]chan bool)
	keyArgs := make(map[byte]chan bool)
	for i := 0; i < 8; i++ {
		mouseBtnArgs[byte(1<<i)] = make(chan bool, 1)
	}
	for i := 0; i < 256; i++ {
		keyArgs[byte(i)] = make(chan bool, 1)
	}

	// Load macros from json
	file, err := os.ReadFile("config/macros.json")
	if err != nil {
		logger.Logger.Fatalf("Failed to read macros file: %v", err)
	}
	var macroData map[string][][4]int32
	err = json.Unmarshal(file, &macroData)
	if err != nil {
		logger.Logger.Fatalf("Failed to unmarshal macros file: %v", err)
	}

	for name, data := range macroData {
		Macros[name] = Macro{
			Name:        name,
			Description: "Data driven macro",
			Fn:          downDragMacro(data),
		}
	}

	Macros["btn_left_hold_autofire"] = Macro{
		Name:        "左键按住连发",
		Description: "按住左键 = 连点左键",
		Fn: func(mk *MacroMouseKeyboard, ch chan bool) {
			for {
				select {
				case <-ch:
					return
				default:
					mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
					time.Sleep(8 * time.Millisecond)
					mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
					time.Sleep(8 * time.Millisecond)
				}
			}
		},
	}

	Macros["btn_left"] = Macro{
		Name:        "左键",
		Description: "普通的左键功能，用于其他按键映射",
		Fn: func(mk *MacroMouseKeyboard, ch chan bool) {
			mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
			<-ch // 等待信号停止
			mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
		},
	}

	loadPlugins()

	return &MacroMouseKeyboard{
		MouseBtnArgs: mouseBtnArgs,
		KeyArgs:      keyArgs,
		Ctrl:         controler,
		Macros:       Macros,
	}
}

func (mk *MacroMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	for dx != 0 || dy != 0 || Wheel != 0 {
		stepDx := clamp(dx, -128, 127)
		stepDy := clamp(dy, -128, 127)
		stepWheel := clamp(Wheel, -128, 127)
		if err := mk.Ctrl.MouseMove(stepDx, stepDy, stepWheel); err != nil {
			return err
		}
		dx -= stepDx
		dy -= stepDy
		Wheel -= stepWheel
	}
	return nil
}

func (mk *MacroMouseKeyboard) MouseBtnDown(keyCode byte) error {
	// This logic will be moved to server package
	return mk.Ctrl.MouseBtnDown(keyCode)
}

func (mk *MacroMouseKeyboard) MouseBtnUp(keyCode byte) error {
	// This logic will be moved to server package
	return mk.Ctrl.MouseBtnUp(keyCode)
}

func (mk *MacroMouseKeyboard) KeyDown(keyCode uint16) error {
	return mk.Ctrl.KeyDown(input.Linux2hid[keyCode])
}

func (mk *MacroMouseKeyboard) KeyUp(keyCode uint16) error {
	return mk.Ctrl.KeyUp(input.Linux2hid[keyCode])
}
