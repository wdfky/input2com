package main

import (
	"sync"

	"go.bug.st/serial"
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
		logger.Error("Value must be in the range of -128 to 127")
		return 0x00 // Return a default value if out of range
	}
	if value >= 0 {
		return byte(value)
	}
	return byte(0x100 + value)
}

type comMouseKeyboard struct {
	port            serial.Port
	mouseButtonByte byte
	keyBytes        []byte
	mu              sync.Mutex
}

func NewComMouseKeyboard(portName string, baudRate int) *comMouseKeyboard {
	port, err := OpenSerialWritePipe(portName, baudRate)
	if err != nil {
		logger.Error("Failed to open serial port")
		return nil
	}
	port.Write([]byte{0x57, 0xAB, 0x02, 0x00, 0x00, 0x00, 0x00})
	port.Write([]byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return &comMouseKeyboard{port: port, mouseButtonByte: 0x00, keyBytes: []byte{0x57, 0xAB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
}

func (mk *comMouseKeyboard) MouseMove(dx, dy, Wheel int32) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, intToByte(dx), intToByte(dy), intToByte(Wheel)})
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) MouseBtnDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte |= keyCode
	_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) MouseBtnUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	mk.mouseButtonByte &^= keyCode
	_, err := mk.port.Write([]byte{0x57, 0xAB, 0x02, mk.mouseButtonByte, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) KeyDown(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= KeyLeftCtrl && keyCode <= KeyRightGui {
		mk.keyBytes[3] |= specialKeysMap[keyCode]
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
	_, err := mk.port.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (mk *comMouseKeyboard) KeyUp(keyCode byte) error {
	mk.mu.Lock()
	defer mk.mu.Unlock()
	if keyCode >= KeyLeftCtrl && keyCode <= KeyRightGui {
		mk.keyBytes[3] &^= specialKeysMap[keyCode]
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
	_, err := mk.port.Write(mk.keyBytes)
	if err != nil {
		return err
	}
	return nil
}
