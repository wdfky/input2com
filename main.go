package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
	"github.com/kenshaw/evdev"
)

type event_pack struct {
	//表示一个动作 由一系列event组成
	dev_name string
	events   []*evdev.Event
}

func dev_reader(event_reader chan *event_pack, index int) {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", index), os.O_RDONLY, 0)
	if err != nil {
		logger.Errorf("读取设备失败 : %v", err)
		return
	}
	d := evdev.Open(fd)
	defer d.Close()
	event_ch := d.Poll(context.Background())
	events := make([]*evdev.Event, 0)
	dev_name := d.Name()
	logger.Infof("开始读取设备 : %s", dev_name)
	d.Lock()
	defer d.Unlock()
	for {
		select {
		case <-global_close_signal:
			logger.Infof("释放设备 : %s", dev_name)
			return
		case event := <-event_ch:
			if event == nil {
				logger.Warnf("移除设备 : %s", dev_name)
				return
			} else if event.Type == evdev.SyncReport {
				pack := &event_pack{
					dev_name: dev_name,
					events:   events,
				}
				event_reader <- pack
				events = make([]*evdev.Event, 0)
			} else {
				events = append(events, &event.Event)
			}
		}
	}
}

var global_close_signal = make(chan bool) //仅会在程序退出时关闭  不用于其他用途

type dev_type uint8

const (
	type_mouse    = dev_type(0)
	type_keyboard = dev_type(1)
	type_joystick = dev_type(2)
	type_touch    = dev_type(3)
	type_unknown  = dev_type(4)
)

func check_dev_type(dev *evdev.Evdev) dev_type {
	abs := dev.AbsoluteTypes()
	key := dev.KeyTypes()
	rel := dev.RelativeTypes()
	_, MTPositionX := abs[evdev.AbsoluteMTPositionX]
	_, MTPositionY := abs[evdev.AbsoluteMTPositionY]
	_, MTSlot := abs[evdev.AbsoluteMTSlot]
	_, MTTrackingID := abs[evdev.AbsoluteMTTrackingID]
	if MTPositionX && MTPositionY && MTSlot && MTTrackingID {
		return type_touch //触屏检测这几个abs类型即可
	}
	_, RelX := rel[evdev.RelativeX]
	_, RelY := rel[evdev.RelativeY]
	_, HWheel := rel[evdev.RelativeHWheel]
	_, MouseLeft := key[evdev.BtnLeft]
	_, MouseRight := key[evdev.BtnRight]
	_, MouseMiddle := key[evdev.BtnMiddle]
	if RelX && RelY && HWheel && MouseLeft && MouseRight && MouseMiddle {
		return type_mouse //鼠标 检测XY 滚轮 左右中键
	}
	keyboard_keys := true
	for i := evdev.KeyEscape; i <= evdev.KeyScrollLock; i++ {
		_, ok := key[i]
		keyboard_keys = keyboard_keys && ok
	}
	if keyboard_keys {
		return type_keyboard //键盘 检测keycode(1-70)
	}

	axis_count := 0
	for i := evdev.AbsoluteX; i <= evdev.AbsoluteRZ; i++ {
		_, ok := abs[i]
		if ok {
			axis_count++
		}
	}
	LS_RS := axis_count >= 4

	key_count := 0
	for i := evdev.BtnA; i <= evdev.BtnZ; i++ {
		_, ok := key[i]
		if ok {
			key_count++
		}
	}
	A_B_X_Y := key_count >= 4

	if LS_RS && A_B_X_Y {
		return type_joystick //手柄 检测LS,RS A,B,X,Y
	}
	return type_unknown
}

func get_possible_device_indexes(skipList map[int]bool) map[int]dev_type {
	files, _ := os.ReadDir("/dev/input")
	result := make(map[int]dev_type)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) <= 5 {
			continue
		}
		if file.Name()[:5] != "event" {
			continue
		}
		index, _ := strconv.Atoi(file.Name()[5:])
		reading, exist := skipList[index]
		if exist && reading {
			continue
		} else {
			fd, err := os.OpenFile(fmt.Sprintf("/dev/input/%s", file.Name()), os.O_RDONLY, 0)
			if err != nil {
				logger.Errorf("读取设备/dev/input/%s失败 : %v ", file.Name(), err)
			}
			d := evdev.Open(fd)
			defer d.Close()
			devType := check_dev_type(d)
			if devType != type_unknown {
				result[index] = devType
			}
		}
	}
	return result
}

func get_dev_name_by_index(index int) string {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", index), os.O_RDONLY, 0)
	if err != nil {
		return "读取设备名称失败"
	}
	d := evdev.Open(fd)
	defer d.Close()
	return d.Name()
}

func auto_detect_and_read(event_chan chan *event_pack) {
	//自动检测设备并读取 循环检测 自动管理设备插入移除
	devices := make(map[int]bool)
	for {
		select {
		case <-global_close_signal:
			return
		default:
			auto_detect_result := get_possible_device_indexes(devices)
			devTypeFriendlyName := map[dev_type]string{
				type_mouse:    "鼠标",
				type_keyboard: "键盘",
				type_joystick: "手柄",
				type_touch:    "触屏",
				type_unknown:  "未知",
			}
			for index, devType := range auto_detect_result {
				devName := get_dev_name_by_index(index)
				if devName == "input2com-virtual-device" {
					continue //跳过生成的虚拟设备
				}
				if devType == type_mouse || devType == type_keyboard || devType == type_joystick {
					logger.Infof("检测到设备 %s(/dev/input/event%d) : %s", devName, index, devTypeFriendlyName[devType])
					localIndex := index
					go func() {
						devices[localIndex] = true
						dev_reader(event_chan, localIndex)
						devices[localIndex] = false
					}()
				}
			}
			time.Sleep(time.Duration(400) * time.Millisecond)
		}
	}
}

func main() {
	//如果有参数-n 则测试模式
	parser := argparse.NewParser("input2com", " ")

	var debug *bool = parser.Flag("d", "debug", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "调试模式",
	})

	var badurate *int = parser.Int("b", "badurate", &argparse.Options{
		Required: false,
		Help:     "波特率",
		Default:  2000000,
	})

	var ttyPath *string = parser.String("t", "tty", &argparse.Options{
		Required: false,
		Default:  "/dev/ttyUSB*",
		Help:     "串口设备路径，可以使用通配符来匹配第一个设备",
	})

	go serve() //启动配置服务器

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *debug {
		logger.WithDebug()
	}

	matches, err := filepath.Glob(*ttyPath)
	if err != nil {
		logger.Fatalf("无法匹配设备路径: %v", err)
	}
	if len(matches) == 0 {
		logger.Fatalf("没有找到匹配的设备路径: %s", *ttyPath)
	}

	devpath := matches[0] // 取第一个匹配的设备路径
	logger.Infof("使用设备路径: %s", devpath)
	logger.Infof("波特率: %d", *badurate)

	events_ch := make(chan *event_pack) //主要设备事件管道
	go auto_detect_and_read(events_ch)
	comKB := NewComMouseKeyboard(devpath, *badurate)
	macroKB := NewMacroMouseKeyboard(comKB)

	handel_rel_event := func(x, y, HWhell, Wheel int32) {
		if x != 0 || y != 0 || HWhell != 0 || Wheel != 0 {
			macroKB.MouseMove(x, y, Wheel)
		}
	}
	handel_key_events := func(events []*evdev.Event, dev_name string) {
		for _, event := range events {
			if event.Value == 0 {
				logger.Debugf("%v 按键释放: %v", dev_name, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					macroKB.MouseBtnUp(MOUSE_BTN_LEFT)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					macroKB.MouseBtnUp(MOUSE_BTN_RIGHT)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					macroKB.MouseBtnUp(MOUSE_BTN_MIDDLE)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					macroKB.MouseBtnUp(MOUSE_BTN_BACK)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					macroKB.MouseBtnUp(MOUSE_BTN_FORWARD)
				} else {
					macroKB.KeyUp(event.Code) // 其他按键释放
				}
			} else if event.Value == 1 {
				logger.Debugf("%v 按键按下: %v", dev_name, event.Code)
				if event.Code == uint16(evdev.BtnLeft) { // 鼠标左键释放
					macroKB.MouseBtnDown(MOUSE_BTN_LEFT)
				} else if event.Code == uint16(evdev.BtnRight) { // 鼠标右键释放
					macroKB.MouseBtnDown(MOUSE_BTN_RIGHT)
				} else if event.Code == uint16(evdev.BtnMiddle) { // 鼠标中键释放
					macroKB.MouseBtnDown(MOUSE_BTN_MIDDLE)
				} else if event.Code == uint16(evdev.BtnSide) { // 鼠标后退键释放
					macroKB.MouseBtnDown(MOUSE_BTN_BACK)
				} else if event.Code == uint16(evdev.BtnExtra) { // 鼠标前进键释放
					macroKB.MouseBtnDown(MOUSE_BTN_FORWARD)
				} else {
					macroKB.KeyDown(event.Code) // 其他按键释放
				}
			} else if event.Value == 2 {
				logger.Debugf("%v 按键重复: %v", dev_name, event.Code)
			}
		}
	}

	handel_abs_events := func(events []*evdev.Event, dev_name string) {
		if len(events) == 0 {
			return
		}
		for _, event := range events {
			if event.Type != evdev.EventAbsolute {
				continue
			}
		}
	}

	go func() {
		for {
			key_events := make([]*evdev.Event, 0)
			abs_events := make([]*evdev.Event, 0)
			var x int32 = 0
			var y int32 = 0
			var HWhell int32 = 0
			var Wheel int32 = 0
			select {
			case <-global_close_signal:
				return
			case event_pack := <-events_ch:
				if event_pack == nil {
					continue
				}
				for _, event := range event_pack.events {
					switch event.Type {
					case evdev.EventKey:
						key_events = append(key_events, event)
					case evdev.EventAbsolute:
						abs_events = append(abs_events, event)
					case evdev.EventRelative:
						switch event.Code {
						case uint16(evdev.RelativeX):
							x = event.Value
						case uint16(evdev.RelativeY):
							y = event.Value
						case uint16(evdev.RelativeHWheel):
							HWhell = event.Value
						case uint16(evdev.RelativeWheel):
							Wheel = event.Value
						}
					}
				}
				var perfPoint time.Time

				perfPoint = time.Now()
				handel_rel_event(x, y, HWhell, Wheel)
				rel_sin := time.Since(perfPoint)
				perfPoint = time.Now()
				handel_key_events(key_events, event_pack.dev_name)
				key_sin := time.Since(perfPoint)
				perfPoint = time.Now()
				handel_abs_events(abs_events, event_pack.dev_name)
				abs_sin := time.Since(perfPoint)
				logger.Debugf("")
				logger.Debugf("handel rel_event\t%v \n", rel_sin)
				logger.Debugf("handel key_events\t%v \n", key_sin)
				logger.Debugf("handel abs_events\t%v \n", abs_sin)
			}
		}
	}()

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-exitChan
	close(global_close_signal)
	logger.Info("已停止")
	time.Sleep(time.Millisecond * 40)
}
