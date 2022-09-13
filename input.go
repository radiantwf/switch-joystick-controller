package joystick

import (
	"strconv"
	"strings"
)

type dPadType byte

const (
	dPadTypeUp dPadType = iota
	dPadTypeTopRight
	dPadTypeRight
	dPadTypeBottomRight
	dPadTypeBottom
	dPadTypeBottomLeft
	dPadTypeLeft
	dPadTypeTopLeft
	dPadTypeCenter
)

type JoyStickInput struct {
	modify bool
	button struct {
		y, b, a, x                  bool
		l, r, zl, zr                bool
		minus, plus, lPress, rPress bool
		home, capture               bool
	}
	dPad  dPadType
	stick struct {
		left, right struct {
			x, y byte
		}
	}
}

// NewJoyStickInput  inputLine: "A","A|B|L","A|LStick@-100,30|ZR"
func NewJoyStickInput(inputLine string) (input *JoyStickInput) {
	input = &JoyStickInput{}
	line := strings.TrimRight(inputLine, "\r\n")
	line = strings.TrimRight(line, "\n")
	line = strings.ToLower(line)
	splits := strings.Split(line, "|")

	input.dPad = dPadTypeCenter
	for _, split := range splits {
		split = strings.Trim(split, " ")
		switch split {
		case "y":
			input.button.y = true
		case "b":
			input.button.b = true
		case "x":
			input.button.x = true
		case "a":
			input.button.a = true

		case "l":
			input.button.l = true
		case "r":
			input.button.r = true
		case "zl":
			input.button.zl = true
		case "zr":
			input.button.zr = true

		case "minus":
			input.button.minus = true
		case "plus":
			input.button.plus = true
		case "lpress":
			input.button.lPress = true
		case "rpress":
			input.button.rPress = true

		case "home":
			input.button.home = true
		case "capture":
			input.button.capture = true

		case "up":
			input.dPad = dPadTypeUp
		case "down":
			input.dPad = dPadTypeBottom
		case "left":
			input.dPad = dPadTypeLeft
		case "right":
			input.dPad = dPadTypeRight
		case "upleft":
			fallthrough
		case "leftup":
			input.dPad = dPadTypeTopLeft
		case "downleft":
			fallthrough
		case "leftdown":
			input.dPad = dPadTypeBottomLeft
		case "upright":
			fallthrough
		case "updown":
			input.dPad = dPadTypeTopRight
		case "downright":
			fallthrough
		case "rightdown":
			input.dPad = dPadTypeBottomRight
		case "center":
			input.dPad = dPadTypeCenter

		default:
			stick := strings.Split(split, "@")
			if len(stick) != 2 {
				break
			}
			coordinate := strings.Split(stick[1], ",")
			if len(coordinate) != 2 {
				break
			}
			x, err1 := strconv.Atoi(coordinate[0])
			if err1 != nil {
				break
			}
			if x > 127 {
				x = 127
			} else if x < -128 {
				x = -128
			}
			y, err1 := strconv.Atoi(coordinate[1])
			if err1 != nil {
				break
			}
			if y > 127 {
				y = 127
			} else if y < -128 {
				y = -128
			}
			switch stick[0] {
			case "lstick":
				input.stick.left.x = byte(x)
				input.stick.left.y = byte(y)
			case "rstick":
				input.stick.right.x = byte(x)
				input.stick.right.y = byte(y)
			}
		}
	}
	input.stick.left.x = input.stick.left.x + 128
	input.stick.left.y = input.stick.left.y + 128
	input.stick.right.x = input.stick.right.x + 128
	input.stick.right.y = input.stick.right.y + 128
	input.modify = true
	return
}

func (c *JoyStickInput) OutputBytes() (outputs [8]byte) {
	if c == nil {
		return
	}
	outputs[0] = c.orderSetBit(c.button.zr, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.zl, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.r, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.l, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.x, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.a, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.b, outputs[0], true)
	outputs[0] = c.orderSetBit(c.button.y, outputs[0], false)

	outputs[1] = c.orderSetBit(false, outputs[1], true)
	outputs[1] = c.orderSetBit(false, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.capture, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.home, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.rPress, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.lPress, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.plus, outputs[1], true)
	outputs[1] = c.orderSetBit(c.button.minus, outputs[1], false)

	outputs[2] = c.orderSetBit(false, outputs[2], true)
	outputs[2] = c.orderSetBit(false, outputs[2], true)
	outputs[2] = c.orderSetBit(false, outputs[2], true)
	outputs[2] = c.orderSetBit(false, outputs[2], false)
	outputs[2] = outputs[2] << 4
	outputs[2] = outputs[2] | byte(c.dPad)

	outputs[3] = c.stick.left.x
	outputs[4] = c.stick.left.y
	outputs[5] = c.stick.right.x
	outputs[6] = c.stick.right.y
	return
}

func (c *JoyStickInput) orderSetBit(flag bool, input byte, moveLeftFlg bool) (output byte) {
	if !c.modify {
		c.stick.left.x = 128
		c.stick.left.y = 128
		c.stick.right.x = 128
		c.stick.right.y = 128
		c.dPad = dPadTypeCenter
		c.modify = true
	}
	if flag {
		input = input | 0x01
	}
	if moveLeftFlg {
		input = input << 1
	}
	output = input
	return
}
