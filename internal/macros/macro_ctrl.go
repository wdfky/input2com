package macros

import (
	"encoding/json"
	"input2com/internal/input"
	"input2com/internal/logger"
	"input2com/internal/serial"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Fn          func(*MacroMouseKeyboard, chan bool) `json:"-"`
}

var (
	MouseConfigDict       = make(map[string]map[byte]string)
	MouseConfigDictSwitch = make(map[string]map[byte]string)
	KeyboardConfigDict    = make(map[byte]string)
	MousedictMutex        sync.RWMutex
	KeyboarddictMutex     sync.RWMutex
)
var Macros = make(map[string]Macro)

func downDragMacro(recoils []*Recoil, xMultiplier, yMultiplier float64) func(mk *MacroMouseKeyboard, ch chan bool) {
	return func(mk *MacroMouseKeyboard, ch chan bool) {
		mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
		defer mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)

		// 执行所有recoil移动
		for _, recoil := range recoils {
			select {
			case <-ch:
				return // 收到释放信号，立即返回
			default:
				if recoil.Count > 0 {
					moveCnt := recoil.Count
					actualStepTime := recoil.RelativeTime / float64(moveCnt)
					sleepDuration := time.Duration(actualStepTime * float64(time.Second))
					for i := 0; i < int(moveCnt); i++ {
						select {
						case <-ch:
							return // 收到释放信号，立即返回
						default:
							mk.Ctrl.MouseMove(recoil.Dx, recoil.Dx, 0)
							if sleepDuration > 0 {
								select {
								case <-ch:
									return
								case <-time.After(sleepDuration):
								}
							}
						}
					}
				} else {
					mk.Ctrl.MouseMove(recoil.Dx, recoil.Dy, 0)
					time.Sleep(time.Duration(recoil.RelativeTime * float64(time.Second)))
				}
			}
		}
		// recoils序列执行完毕，等待释放信号
		<-ch
	}
}
func downDragMacroWithRight(recoils []*Recoil) func(mk *MacroMouseKeyboard, ch chan bool) {
	return func(mk *MacroMouseKeyboard, ch chan bool) {
		mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
		defer mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
		for _, recoil := range recoils {
			select {
			case <-ch:
				return
			default:
				if mk.Ctrl.IsMouseBtnPressed(input.MouseBtnRight) {
					mk.Ctrl.MouseMove(recoil.Dx, recoil.Dy, 0)
				}
				time.Sleep(time.Duration(recoil.RelativeTime * float64(time.Second)))
			}
		}
		// recoils序列执行完毕，等待释放信号
		<-ch
	}
}

// ---------- 主函数 ----------
func downDragMacroWithForward(recoils []*Recoil, xMultiplier, yMultiplier float64) func(mk *MacroMouseKeyboard, ch chan bool) {
	return func(mk *MacroMouseKeyboard, ch chan bool) {
		mk.Ctrl.MouseBtnDown(input.MouseBtnLeft)
		defer mk.Ctrl.MouseBtnUp(input.MouseBtnLeft)
		for _, recoil := range recoils {
			select {
			case <-ch:
				return
			default:
				if mk.Ctrl.IsMouseBtnPressed(input.MouseBtnForward) {
					if recoil.Count > 0 {
						moveCnt := float64(recoil.Count) * xMultiplier
						actualStepTime := recoil.RelativeTime / moveCnt
						sleepDuration := time.Duration(actualStepTime * float64(time.Second))
						for i := 0; i < int(moveCnt); i++ {
							select {
							case <-ch:
								return // 收到释放信号，立即返回
							default:
								mk.Ctrl.MouseMove(recoil.Dx, recoil.Dx, 0)
								if sleepDuration > 0 {
									select {
									case <-ch:
										return
									case <-time.After(sleepDuration):
									}
								}
							}
						}
					} else {
						mk.Ctrl.MouseMove(recoil.Dx, recoil.Dy, 0)
						time.Sleep(time.Duration(recoil.RelativeTime * float64(time.Second)))
					}
				}
			}
		}
		// recoils序列执行完毕，等待释放信号
		<-ch
	}
}
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
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

type RecoilConfig struct {
	Name        string    `json:"name"`
	Recoils     []*Recoil `json:"recoil"`
	XMultiplier float64   `json:"x_multiplier"` // X 轴乘数因子
	YMultiplier float64   `json:"y_multiplier"` // Y 轴乘数因子
}
type Recoil struct {
	Dx           int32   `json:"dx"`
	Dy           int32   `json:"dy"`
	RelativeTime float64 `json:"relative_time"`
	Count        int32
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
	var recoils []*RecoilConfig
	err = json.Unmarshal(file, &recoils)
	if err != nil {
		logger.Logger.Fatalf("Failed to unmarshal macros file: %v", err)
	}
	logger.Logger.Infof("Loaded recoils: %v", recoils)

	for _, recoil := range recoils {
		Macros[recoil.Name] = Macro{
			Name:        recoil.Name,
			Description: "压枪宏仅按键按下",
			Fn:          downDragMacro(recoil.Recoils, recoil.XMultiplier, recoil.YMultiplier),
		}
		Macros[recoil.Name+"_withright"] = Macro{
			Name:        recoil.Name + "_withright",
			Description: "压枪宏右键按下",
			Fn:          downDragMacroWithRight(recoil.Recoils),
		}
		Macros[recoil.Name+"_forward"] = Macro{
			Name:        recoil.Name + "_forward",
			Description: "压枪宏前侧键按下",
			Fn:          downDragMacroWithForward(recoil.Recoils, recoil.XMultiplier, recoil.YMultiplier),
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
	Macros["switch"] = Macro{
		Name:        "切换",
		Description: "切换原生功能与宏功能",
		Fn: func(mk *MacroMouseKeyboard, ch chan bool) {
			MouseConfigDict, MouseConfigDictSwitch = MouseConfigDictSwitch, MouseConfigDict
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
	// 分别处理 dx, dy, Wheel 的拆分移动
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
func (mk *MacroMouseKeyboard) MouseBtnDown(keyCode byte, devName string) error {

	// 1. 优先根据设备名获取该设备的宏配置（外层 map 键为设备名）
	deviceMacroConfig, deviceExists := MouseConfigDict[strings.ToLower(devName)]
	if !deviceExists {
		// 设备无宏配置，直接调用底层控制器
		return mk.Ctrl.MouseBtnDown(keyCode)
	}
	// 2. 在设备配置中查找当前按键码对应的宏标识符
	macroID, keyExists := deviceMacroConfig[keyCode]
	if !keyExists {
		// 当前按键在该设备下无宏配置，直接调用底层控制器
		return mk.Ctrl.MouseBtnDown(keyCode)
	}
	if macroFunc, exists := mk.Macros[macroID]; exists { // 如果有宏函数，执行宏
		go macroFunc.Fn(mk, mk.MouseBtnArgs[keyCode])
		return nil
	}
	return mk.Ctrl.MouseBtnDown(keyCode) // 如果没有宏函数，直接调用控制器的MouseBtnDown
}

func (mk *MacroMouseKeyboard) MouseBtnUp(keyCode byte, devName string) error {

	// 1. 优先根据设备名获取该设备的宏配置（外层 map 键为设备名）
	deviceMacroConfig, deviceExists := MouseConfigDict[strings.ToLower(devName)]
	if !deviceExists {
		// 设备无宏配置，直接调用底层控制器
		return mk.Ctrl.MouseBtnUp(keyCode)
	}
	// 2. 在设备配置中查找当前按键码对应的宏标识符
	macroID, keyExists := deviceMacroConfig[keyCode]
	if !keyExists {
		// 当前按键在该设备下无宏配置，直接调用底层控制器
		return mk.Ctrl.MouseBtnUp(keyCode)
	}
	if _, exists := mk.Macros[macroID]; exists { // 如果有宏函数，执行宏
		mk.MouseBtnArgs[keyCode] <- true // 发送信号停止宏
		return nil
	}
	return mk.Ctrl.MouseBtnUp(keyCode)
}

func (mk *MacroMouseKeyboard) KeyDown(keyCode uint16) error {
	return mk.Ctrl.KeyDown(input.Linux2hid[keyCode])
}

func (mk *MacroMouseKeyboard) KeyUp(keyCode uint16) error {
	return mk.Ctrl.KeyUp(input.Linux2hid[keyCode])
}
