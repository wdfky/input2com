package serial

import (
	"context"
	"input2com/internal/input"
	"input2com/internal/logger"
	"sync"
	"time"

	"go.bug.st/serial"
	"golang.org/x/time/rate"
)

func OpenSerialWritePipe(portName string, baudRate int) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	return port, nil
}

func intToByte(value int32) byte {
	if value < -128 || value > 127 {
		logger.Logger.Error("Value must be in the range of -128 to 127")
		return 0x00 // Return a default value if out of range
	}
	if value >= 0 {
		return byte(value)
	}
	return byte(0x100 + value)
}

type ComMouseKeyboard struct {
	serial.Port
	mouseButtonByte byte
	keyBytes        []byte
	mu              sync.Mutex
	limiter         *rate.Limiter
}

func (mk *ComMouseKeyboard) Write(p []byte) (n int, err error) {
	// 等待限流器许可
	if err = mk.limiter.Wait(context.TODO()); err != nil {
		return 0, err
	}
	return mk.Port.Write(p)
}
func NewComMouseKeyboard(portName string, baudRate int) *ComMouseKeyboard {
	port, err := OpenSerialWritePipe(portName, baudRate)
	if err != nil {
		logger.Logger.Error("Failed to open serial port")
		return nil
	}
	port.Write([]byte{0x57, 0xAB, 0x02, 0x00, 0x00, 0x00, 0x00})
	port.Write([]byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return &ComMouseKeyboard{
		Port:            port,
		mouseButtonByte: 0x00,
		keyBytes:        []byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		limiter:         rate.NewLimiter(rate.Every(time.Millisecond), 1),
	}
}

//func (mk *ComMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
//	mk.mu.Lock()
//	defer mk.mu.Unlock()
//	_, err := mk.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, intToByte(dx), intToByte(dy), intToByte(Wheel)})
//	if err != nil {
//		return err
//	}
//	return nil
//}

// 大范围鼠标移动方法
func (mk *ComMouseKeyboard) MouseMoveLarge(dx, dy, wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()

	// 处理X方向移动
	var xBytes [2]byte
	if dx >= 0 {
		// 向右移动
		xBytes = int16ToBytes(int16(dx))
	} else {
		// 向左移动，使用补码表示
		xBytes = int16ToBytes(int16(65536 + dx))
	}

	// 处理Y方向移动
	var yBytes [2]byte
	if dy >= 0 {
		// 向下移动
		yBytes = int16ToBytes(int16(dy))
	} else {
		// 向上移动，使用补码表示
		yBytes = int16ToBytes(int16(65536 + dy))
	}

	// 构建数据包
	data := []byte{
		0x57, 0xAB, 0x22, // 头部和命令字节
		mk.mouseButtonByte,   // 鼠标按键状态
		xBytes[0], xBytes[1], // X方向移动距离（低字节在前）
		yBytes[0], yBytes[1], // Y方向移动距离（低字节在前）
		intToByte(wheel), // 滚轮滚动齿数
	}

	_, err := mk.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// 将int16转换为2字节（小端序）
func int16ToBytes(value int16) [2]byte {
	return [2]byte{
		byte(value & 0xFF),        // 低字节
		byte((value >> 8) & 0xFF), // 高字节
	}
}

// 原有的小范围移动方法（保持兼容性）
func (mk *ComMouseKeyboard) MouseMove(dx, dy, wheel int32) error {
	// 如果移动范围在单字节范围内，使用小范围移动
	if dx >= -127 && dx <= 127 && dy >= -127 && dy <= 127 {
		return mk.mouseMoveSmall(dx, dy, wheel)
	}
	// 否则使用大范围移动
	return mk.MouseMoveLarge(dx, dy, wheel)
}

// 原有的小范围移动实现（重命名）
func (mk *ComMouseKeyboard) mouseMoveSmall(dx, dy, wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	_, err := mk.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, intToByte(dx), intToByte(dy), intToByte(wheel)})
	if err != nil {
		return err
	}
	return nil
}

func (mk *ComMouseKeyboard) IsMouseBtnPressed(keyCode byte) bool {
	return mk.mouseButtonByte&keyCode != 0
}
func (mk *ComMouseKeyboard) MouseBtnDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte |= keyCode
	_, err := mk.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}
	return nil
}

func (mk *ComMouseKeyboard) MouseBtnUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte &^= keyCode
	_, err := mk.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}
	return nil
}

func (mk *ComMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= input.KeyLeftCtrl && keyCode <= input.KeyRightGui {
		mk.keyBytes[3] |= input.SpecialKeysMap[keyCode]
	} else {
		for i := 0; i < 7; i++ {
			if i == 6 {
				return nil // No space to add new key, ignore
			}
			if mk.keyBytes[i+5] == 0x00 {
				mk.keyBytes[i+5] = keyCode
				break
			}
		}
	}
	_, err := mk.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (mk *ComMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= input.KeyLeftCtrl && keyCode <= input.KeyRightGui {
		mk.keyBytes[3] &^= input.SpecialKeysMap[keyCode]
	} else {
		for i := 0; i < 7; i++ {
			if i == 6 {
				return nil // No space to add new key, ignore
			}
			if mk.keyBytes[i+5] == keyCode {
				mk.keyBytes[i+5] = 0x00
				break
			}
		}
	}
	_, err := mk.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}
