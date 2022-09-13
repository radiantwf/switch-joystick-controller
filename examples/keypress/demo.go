package main

import (
	"log"
	"time"

	joystick "github.com/radiantwf/switch-joystick-controller"
)

func main() {
	stick := joystick.NewJoyStick()

	err := stick.Open()
	if err != nil {
		panic(err)
	}
	defer stick.Close()

	log.Println("首次发送Button A（唤醒手柄时无实际按键操作）")
	input := joystick.NewJoyStickInput("A")
	stick.SyncSendKey(input, time.Millisecond*100)
	log.Println("延时15s，等待手柄唤醒")
	time.Sleep(time.Second * 15)

	input = joystick.NewJoyStickInput("A")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("B")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("X")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("Y")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)

	input = joystick.NewJoyStickInput("L")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("R")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("ZL")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("ZR")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)

	input = joystick.NewJoyStickInput("Minus")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("Plus")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("LPress")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("RPress")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)

	input = joystick.NewJoyStickInput("Up")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("Down")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("Left")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
	input = joystick.NewJoyStickInput("Right")
	stick.SyncSendKey(input, time.Millisecond*100)
	time.Sleep(time.Second)
}

// echo 0400088080808000 | xxd -r -ps > /dev/hidg0
