package main

import "strconv"

var specialKeysMap = map[byte]byte{
	KEY_LEFT_CTRL:   byte(1 << 0),
	KEY_LEFT_SHIFT:  byte(1 << 1),
	KEY_LEFT_ALT:    byte(1 << 2),
	KEY_LEFT_GUI:    byte(1 << 3),
	KEY_RIGHT_CTRL:  byte(1 << 4),
	KEY_RIGHT_SHIFT: byte(1 << 5),
	KEY_RIGHT_ALT:   byte(1 << 6),
	KEY_RIGHT_GUI:   byte(1 << 7),
}

const (
	MOUSE_BTN_LEFT    = byte(1 << 0) // 左键
	MOUSE_BTN_RIGHT   = byte(1 << 1) // 右键
	MOUSE_BTN_MIDDLE  = byte(1 << 2) // 中键
	MOUSE_BTN_BACK    = byte(1 << 3) // 后退键
	MOUSE_BTN_FORWARD = byte(1 << 4) // 前进键

	KEY_LEFT_CTRL   = byte(0xe0)
	KEY_LEFT_SHIFT  = byte(0xe1)
	KEY_LEFT_ALT    = byte(0xe2)
	KEY_LEFT_GUI    = byte(0xe3)
	KEY_RIGHT_CTRL  = byte(0xe4)
	KEY_RIGHT_SHIFT = byte(0xe5)
	KEY_RIGHT_ALT   = byte(0xe6)
	KEY_RIGHT_GUI   = byte(0xe8)

	KEY_A           = byte(0x04)
	KEY_B           = byte(0x05)
	KEY_C           = byte(0x06)
	KEY_D           = byte(0x07)
	KEY_E           = byte(0x08)
	KEY_F           = byte(0x09)
	KEY_G           = byte(0x0A)
	KEY_H           = byte(0x0B)
	KEY_I           = byte(0x0C)
	KEY_J           = byte(0x0D)
	KEY_K           = byte(0x0E)
	KEY_L           = byte(0x0F)
	KEY_M           = byte(0x10)
	KEY_N           = byte(0x11)
	KEY_O           = byte(0x12)
	KEY_P           = byte(0x13)
	KEY_Q           = byte(0x14)
	KEY_R           = byte(0x15)
	KEY_S           = byte(0x16)
	KEY_T           = byte(0x17)
	KEY_U           = byte(0x18)
	KEY_V           = byte(0x19)
	KEY_W           = byte(0x1A)
	KEY_X           = byte(0x1B)
	KEY_Y           = byte(0x1C)
	KEY_Z           = byte(0x1D)
	KEY_1           = byte(0x1E)
	KEY_2           = byte(0x1F)
	KEY_3           = byte(0x20)
	KEY_4           = byte(0x21)
	KEY_5           = byte(0x22)
	KEY_6           = byte(0x23)
	KEY_7           = byte(0x24)
	KEY_8           = byte(0x25)
	KEY_9           = byte(0x26)
	KEY_0           = byte(0x27)
	KEY_RETURN      = byte(0x28)
	KEY_ENTER       = byte(0x28)
	KEY_ESC         = byte(0x29)
	KEY_ESCAPE      = byte(0x29)
	KEY_BCKSPC      = byte(0x2A)
	KEY_BACKSPACE   = byte(0x2A)
	KEY_TAB         = byte(0x2B)
	KEY_SPACE       = byte(0x2C)
	KEY_MINUS       = byte(0x2D)
	KEY_DASH        = byte(0x2D)
	KEY_EQUALS      = byte(0x2E)
	KEY_EQUAL       = byte(0x2E)
	KEY_LBRACKET    = byte(0x2F)
	KEY_RBRACKET    = byte(0x30)
	KEY_BACKSLASH   = byte(0x31)
	KEY_HASH        = byte(0x32)
	KEY_NUMBER      = byte(0x32)
	KEY_SEMICOLON   = byte(0x33)
	KEY_QUOTE       = byte(0x34)
	KEY_BACKQUOTE   = byte(0x35)
	KEY_TILDE       = byte(0x35)
	KEY_COMMA       = byte(0x36)
	KEY_PERIOD      = byte(0x37)
	KEY_STOP        = byte(0x37)
	KEY_SLASH       = byte(0x38)
	KEY_CAPS_LOCK   = byte(0x39)
	KEY_CAPSLOCK    = byte(0x39)
	KEY_F1          = byte(0x3A)
	KEY_F2          = byte(0x3B)
	KEY_F3          = byte(0x3C)
	KEY_F4          = byte(0x3D)
	KEY_F5          = byte(0x3E)
	KEY_F6          = byte(0x3F)
	KEY_F7          = byte(0x40)
	KEY_F8          = byte(0x41)
	KEY_F9          = byte(0x42)
	KEY_F10         = byte(0x43)
	KEY_F11         = byte(0x44)
	KEY_F12         = byte(0x45)
	KEY_PRINT       = byte(0x46)
	KEY_SCROLL_LOCK = byte(0x47)
	KEY_SCROLLLOCK  = byte(0x47)
	KEY_PAUSE       = byte(0x48)
	KEY_INSERT      = byte(0x49)
	KEY_HOME        = byte(0x4A)
	KEY_PAGEUP      = byte(0x4B)
	KEY_PGUP        = byte(0x4B)
	KEY_DEL         = byte(0x4C)
	KEY_DELETE      = byte(0x4C)
	KEY_END         = byte(0x4D)
	KEY_PAGEDOWN    = byte(0x4E)
	KEY_PGDOWN      = byte(0x4E)
	KEY_RIGHT       = byte(0x4F)
	KEY_LEFT        = byte(0x50)
	KEY_DOWN        = byte(0x51)
	KEY_UP          = byte(0x52)
	KEY_NUM_LOCK    = byte(0x53)
	KEY_NUMLOCK     = byte(0x53)
	KEY_KP_DIVIDE   = byte(0x54)
	KEY_KP_MULTIPLY = byte(0x55)
	KEY_KP_MINUS    = byte(0x56)
	KEY_KP_PLUS     = byte(0x57)
	KEY_KP_ENTER    = byte(0x58)
	KEY_KP_RETURN   = byte(0x58)
	KEY_KP_1        = byte(0x59)
	KEY_KP_2        = byte(0x5A)
	KEY_KP_3        = byte(0x5B)
	KEY_KP_4        = byte(0x5C)
	KEY_KP_5        = byte(0x5D)
	KEY_KP_6        = byte(0x5E)
	KEY_KP_7        = byte(0x5F)
	KEY_KP_8        = byte(0x60)
	KEY_KP_9        = byte(0x61)
	KEY_KP_0        = byte(0x62)
	KEY_KP_PERIOD   = byte(0x63)
	KEY_KP_STOP     = byte(0x63)
	KEY_APPLICATION = byte(0x65)
	KEY_POWER       = byte(0x66)
	KEY_KP_EQUALS   = byte(0x67)
	KEY_KP_EQUAL    = byte(0x67)
	KEY_F13         = byte(0x68)
	KEY_F14         = byte(0x69)
	KEY_F15         = byte(0x6A)
	KEY_F16         = byte(0x6B)
	KEY_F17         = byte(0x6C)
	KEY_F18         = byte(0x6D)
	KEY_F19         = byte(0x6E)
	KEY_F20         = byte(0x6F)
	KEY_F21         = byte(0x70)
	KEY_F22         = byte(0x71)
	KEY_F23         = byte(0x72)
	KEY_F24         = byte(0x73)
	KEY_EXECUTE     = byte(0x74)
	KEY_HELP        = byte(0x75)
	KEY_MENU        = byte(0x76)
	KEY_SELECT      = byte(0x77)
	KEY_CANCEL      = byte(0x78)
	KEY_REDO        = byte(0x79)
	KEY_UNDO        = byte(0x7A)
	KEY_CUT         = byte(0x7B)
	KEY_COPY        = byte(0x7C)
	KEY_PASTE       = byte(0x7D)
	KEY_FIND        = byte(0x7E)
	KEY_MUTE        = byte(0x7F)
	KEY_VOLUME_UP   = byte(0x80)
	KEY_VOLUME_DOWN = byte(0x81)
)

var linux2hid = map[uint16]uint8{
	30:  4,
	48:  5,
	46:  6,
	32:  7,
	18:  8,
	33:  9,
	34:  10,
	35:  11,
	23:  12,
	36:  13,
	37:  14,
	38:  15,
	50:  16,
	49:  17,
	24:  18,
	25:  19,
	16:  20,
	19:  21,
	31:  22,
	20:  23,
	22:  24,
	47:  25,
	17:  26,
	45:  27,
	21:  28,
	44:  29,
	2:   30,
	3:   31,
	4:   32,
	5:   33,
	6:   34,
	7:   35,
	8:   36,
	9:   37,
	10:  38,
	11:  39,
	28:  40,
	1:   41,
	14:  42,
	15:  43,
	57:  44,
	12:  45,
	13:  46,
	26:  47,
	27:  48,
	43:  49,
	39:  51,
	40:  52,
	41:  53,
	51:  54,
	52:  55,
	53:  56,
	58:  57,
	59:  58,
	60:  59,
	61:  60,
	62:  61,
	63:  62,
	64:  63,
	65:  64,
	66:  65,
	67:  66,
	68:  67,
	87:  68,
	88:  69,
	99:  70,
	70:  71,
	119: 72,
	110: 73,
	102: 74,
	104: 75,
	111: 76,
	107: 77,
	109: 78,
	106: 79,
	105: 80,
	108: 81,
	103: 82,
	69:  83,
	98:  84,
	55:  85,
	74:  86,
	78:  87,
	96:  88,
	79:  89,
	80:  90,
	81:  91,
	75:  92,
	76:  93,
	77:  94,
	71:  95,
	72:  96,
	73:  97,
	82:  98,
	83:  99,
	86:  100,
	127: 101,
	29:  224,
	42:  225,
	56:  226,
	125: 227,
	97:  228,
	54:  229,
	100: 230,
	126: 232,
}

var mouse_valid_keys = map[string]bool{
	strconv.FormatUint(uint64(MOUSE_BTN_LEFT), 10):    true,
	strconv.FormatUint(uint64(MOUSE_BTN_RIGHT), 10):   true,
	strconv.FormatUint(uint64(MOUSE_BTN_MIDDLE), 10):  true,
	strconv.FormatUint(uint64(MOUSE_BTN_BACK), 10):    true,
	strconv.FormatUint(uint64(MOUSE_BTN_FORWARD), 10): true,
}

var keyboard_valid_keys = map[string]bool{
	strconv.FormatUint(uint64(KEY_LEFT_CTRL), 10):   true,
	strconv.FormatUint(uint64(KEY_LEFT_SHIFT), 10):  true,
	strconv.FormatUint(uint64(KEY_LEFT_ALT), 10):    true,
	strconv.FormatUint(uint64(KEY_LEFT_GUI), 10):    true,
	strconv.FormatUint(uint64(KEY_RIGHT_CTRL), 10):  true,
	strconv.FormatUint(uint64(KEY_RIGHT_SHIFT), 10): true,
	strconv.FormatUint(uint64(KEY_RIGHT_ALT), 10):   true,
	strconv.FormatUint(uint64(KEY_RIGHT_GUI), 10):   true,
	strconv.FormatUint(uint64(KEY_A), 10):           true,
	strconv.FormatUint(uint64(KEY_B), 10):           true,
	strconv.FormatUint(uint64(KEY_C), 10):           true,
	strconv.FormatUint(uint64(KEY_D), 10):           true,
	strconv.FormatUint(uint64(KEY_E), 10):           true,
	strconv.FormatUint(uint64(KEY_F), 10):           true,
	strconv.FormatUint(uint64(KEY_G), 10):           true,
	strconv.FormatUint(uint64(KEY_H), 10):           true,
	strconv.FormatUint(uint64(KEY_I), 10):           true,
	strconv.FormatUint(uint64(KEY_J), 10):           true,
	strconv.FormatUint(uint64(KEY_K), 10):           true,
	strconv.FormatUint(uint64(KEY_L), 10):           true,
	strconv.FormatUint(uint64(KEY_M), 10):           true,
	strconv.FormatUint(uint64(KEY_N), 10):           true,
	strconv.FormatUint(uint64(KEY_O), 10):           true,
	strconv.FormatUint(uint64(KEY_P), 10):           true,
	strconv.FormatUint(uint64(KEY_Q), 10):           true,
	strconv.FormatUint(uint64(KEY_R), 10):           true,
	strconv.FormatUint(uint64(KEY_S), 10):           true,
	strconv.FormatUint(uint64(KEY_T), 10):           true,
	strconv.FormatUint(uint64(KEY_U), 10):           true,
	strconv.FormatUint(uint64(KEY_V), 10):           true,
	strconv.FormatUint(uint64(KEY_W), 10):           true,
	strconv.FormatUint(uint64(KEY_X), 10):           true,
	strconv.FormatUint(uint64(KEY_Y), 10):           true,
	strconv.FormatUint(uint64(KEY_Z), 10):           true,
	strconv.FormatUint(uint64(KEY_1), 10):           true,
	strconv.FormatUint(uint64(KEY_2), 10):           true,
	strconv.FormatUint(uint64(KEY_3), 10):           true,
	strconv.FormatUint(uint64(KEY_4), 10):           true,
	strconv.FormatUint(uint64(KEY_5), 10):           true,
	strconv.FormatUint(uint64(KEY_6), 10):           true,
	strconv.FormatUint(uint64(KEY_7), 10):           true,
	strconv.FormatUint(uint64(KEY_8), 10):           true,
	strconv.FormatUint(uint64(KEY_9), 10):           true,
	strconv.FormatUint(uint64(KEY_0), 10):           true,
	strconv.FormatUint(uint64(KEY_RETURN), 10):      true,
	strconv.FormatUint(uint64(KEY_ENTER), 10):       true,
	strconv.FormatUint(uint64(KEY_ESC), 10):         true,
	strconv.FormatUint(uint64(KEY_ESCAPE), 10):      true,
	strconv.FormatUint(uint64(KEY_BCKSPC), 10):      true,
	strconv.FormatUint(uint64(KEY_BACKSPACE), 10):   true,
	strconv.FormatUint(uint64(KEY_TAB), 10):         true,
	strconv.FormatUint(uint64(KEY_SPACE), 10):       true,
	strconv.FormatUint(uint64(KEY_MINUS), 10):       true,
	strconv.FormatUint(uint64(KEY_DASH), 10):        true,
	strconv.FormatUint(uint64(KEY_EQUALS), 10):      true,
	strconv.FormatUint(uint64(KEY_EQUAL), 10):       true,
	strconv.FormatUint(uint64(KEY_LBRACKET), 10):    true,
	strconv.FormatUint(uint64(KEY_RBRACKET), 10):    true,
	strconv.FormatUint(uint64(KEY_BACKSLASH), 10):   true,
	strconv.FormatUint(uint64(KEY_HASH), 10):        true,
	strconv.FormatUint(uint64(KEY_NUMBER), 10):      true,
	strconv.FormatUint(uint64(KEY_SEMICOLON), 10):   true,
	strconv.FormatUint(uint64(KEY_QUOTE), 10):       true,
	strconv.FormatUint(uint64(KEY_BACKQUOTE), 10):   true,
	strconv.FormatUint(uint64(KEY_TILDE), 10):       true,
	strconv.FormatUint(uint64(KEY_COMMA), 10):       true,
	strconv.FormatUint(uint64(KEY_PERIOD), 10):      true,
	strconv.FormatUint(uint64(KEY_STOP), 10):        true,
	strconv.FormatUint(uint64(KEY_SLASH), 10):       true,
	strconv.FormatUint(uint64(KEY_CAPS_LOCK), 10):   true,
	strconv.FormatUint(uint64(KEY_CAPSLOCK), 10):    true,
	strconv.FormatUint(uint64(KEY_F1), 10):          true,
	strconv.FormatUint(uint64(KEY_F2), 10):          true,
	strconv.FormatUint(uint64(KEY_F3), 10):          true,
	strconv.FormatUint(uint64(KEY_F4), 10):          true,
	strconv.FormatUint(uint64(KEY_F5), 10):          true,
	strconv.FormatUint(uint64(KEY_F6), 10):          true,
	strconv.FormatUint(uint64(KEY_F7), 10):          true,
	strconv.FormatUint(uint64(KEY_F8), 10):          true,
	strconv.FormatUint(uint64(KEY_F9), 10):          true,
	strconv.FormatUint(uint64(KEY_F10), 10):         true,
	strconv.FormatUint(uint64(KEY_F11), 10):         true,
	strconv.FormatUint(uint64(KEY_F12), 10):         true,
	strconv.FormatUint(uint64(KEY_PRINT), 10):       true,
	strconv.FormatUint(uint64(KEY_SCROLL_LOCK), 10): true,
	strconv.FormatUint(uint64(KEY_SCROLLLOCK), 10):  true,
	strconv.FormatUint(uint64(KEY_PAUSE), 10):       true,
	strconv.FormatUint(uint64(KEY_INSERT), 10):      true,
	strconv.FormatUint(uint64(KEY_HOME), 10):        true,
	strconv.FormatUint(uint64(KEY_PAGEUP), 10):      true,
	strconv.FormatUint(uint64(KEY_PGUP), 10):        true,
	strconv.FormatUint(uint64(KEY_DEL), 10):         true,
	strconv.FormatUint(uint64(KEY_DELETE), 10):      true,
	strconv.FormatUint(uint64(KEY_END), 10):         true,
	strconv.FormatUint(uint64(KEY_PAGEDOWN), 10):    true,
	strconv.FormatUint(uint64(KEY_PGDOWN), 10):      true,
	strconv.FormatUint(uint64(KEY_RIGHT), 10):       true,
	strconv.FormatUint(uint64(KEY_LEFT), 10):        true,
	strconv.FormatUint(uint64(KEY_DOWN), 10):        true,
	strconv.FormatUint(uint64(KEY_UP), 10):          true,
	strconv.FormatUint(uint64(KEY_NUM_LOCK), 10):    true,
	strconv.FormatUint(uint64(KEY_NUMLOCK), 10):     true,
	strconv.FormatUint(uint64(KEY_KP_DIVIDE), 10):   true,
	strconv.FormatUint(uint64(KEY_KP_MULTIPLY), 10): true,
	strconv.FormatUint(uint64(KEY_KP_MINUS), 10):    true,
	strconv.FormatUint(uint64(KEY_KP_PLUS), 10):     true,
	strconv.FormatUint(uint64(KEY_KP_ENTER), 10):    true,
	strconv.FormatUint(uint64(KEY_KP_RETURN), 10):   true,
	strconv.FormatUint(uint64(KEY_KP_1), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_2), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_3), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_4), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_5), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_6), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_7), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_8), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_9), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_0), 10):        true,
	strconv.FormatUint(uint64(KEY_KP_PERIOD), 10):   true,
	strconv.FormatUint(uint64(KEY_KP_STOP), 10):     true,
	strconv.FormatUint(uint64(KEY_APPLICATION), 10): true,
	strconv.FormatUint(uint64(KEY_POWER), 10):       true,
	strconv.FormatUint(uint64(KEY_KP_EQUALS), 10):   true,
	strconv.FormatUint(uint64(KEY_KP_EQUAL), 10):    true,
	strconv.FormatUint(uint64(KEY_F13), 10):         true,
	strconv.FormatUint(uint64(KEY_F14), 10):         true,
	strconv.FormatUint(uint64(KEY_F15), 10):         true,
	strconv.FormatUint(uint64(KEY_F16), 10):         true,
	strconv.FormatUint(uint64(KEY_F17), 10):         true,
	strconv.FormatUint(uint64(KEY_F18), 10):         true,
	strconv.FormatUint(uint64(KEY_F19), 10):         true,
	strconv.FormatUint(uint64(KEY_F20), 10):         true,
	strconv.FormatUint(uint64(KEY_F21), 10):         true,
	strconv.FormatUint(uint64(KEY_F22), 10):         true,
	strconv.FormatUint(uint64(KEY_F23), 10):         true,
	strconv.FormatUint(uint64(KEY_F24), 10):         true,
	strconv.FormatUint(uint64(KEY_EXECUTE), 10):     true,
	strconv.FormatUint(uint64(KEY_HELP), 10):        true,
	strconv.FormatUint(uint64(KEY_MENU), 10):        true,
	strconv.FormatUint(uint64(KEY_SELECT), 10):      true,
	strconv.FormatUint(uint64(KEY_CANCEL), 10):      true,
	strconv.FormatUint(uint64(KEY_REDO), 10):        true,
	strconv.FormatUint(uint64(KEY_UNDO), 10):        true,
	strconv.FormatUint(uint64(KEY_CUT), 10):         true,
	strconv.FormatUint(uint64(KEY_COPY), 10):        true,
	strconv.FormatUint(uint64(KEY_PASTE), 10):       true,
	strconv.FormatUint(uint64(KEY_FIND), 10):        true,
	strconv.FormatUint(uint64(KEY_MUTE), 10):        true,
	strconv.FormatUint(uint64(KEY_VOLUME_UP), 10):   true,
	strconv.FormatUint(uint64(KEY_VOLUME_DOWN), 10): true,
}
