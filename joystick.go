package joystick

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type JoyStickCommandType int

const (
	JoyStickCommandInstant JoyStickCommandType = iota
	JoyStickCommandMacro

	SendKeyMinDuration time.Duration = time.Millisecond * 5
)

type JoyStick struct {
	DeviceName       string
	CommandType      JoyStickCommandType
	isOpened         bool
	fp               *os.File
	inputChan        chan *JoyStickInput
	sentKeyChan      chan bool
	abortSendKeyChan chan bool
	closeChan        chan bool
	wg               sync.WaitGroup
}

func NewJoyStick(options ...func(*JoyStick)) (joyStick *JoyStick) {
	joyStick = &JoyStick{
		DeviceName:  "/dev/hidg0",
		CommandType: JoyStickCommandInstant,
	}
	for _, option := range options {
		option(joyStick)
	}
	return
}

func (c *JoyStick) Open() (err error) {
	if c.isOpened {
		return
	}
	c.fp, err = os.OpenFile(c.DeviceName, os.O_RDWR|os.O_SYNC, os.ModeDevice)
	if err != nil {
		return
	}
	c.wg.Add(1)
	c.inputChan = make(chan *JoyStickInput)
	c.sentKeyChan = make(chan bool, 1)
	c.abortSendKeyChan = make(chan bool, 1)
	c.closeChan = make(chan bool)
	c.isOpened = true
	c.sentKeyChan <- true
	go func() {
		jumpLoop := false
		t := time.NewTicker(time.Millisecond * 5)
		input := NewJoyStickInput("")
		lastStr := ""
		lastTime := time.Time{}
		for {
			select {
			case i, ok := <-c.inputChan:
				if !ok {
					jumpLoop = true
					break
				}
				input = i
			case <-c.closeChan:
				jumpLoop = true
			case <-t.C:
			}
			t.Stop()
			if jumpLoop {
				break
			}

			output := input.OutputBytes()
			str := fmt.Sprintf("%X", output)
			if str == lastStr && time.Since(lastTime) < time.Second {
				continue
			}
			lastStr = str
			lastTime = time.Now()
			_, err1 := c.fp.Write(output[:])
			// log.Printf("ts %d:\t%X\n", time.Now().UnixNano(), output[:])
			if err1 != nil {
				log.Println(err1)
				break
			}
		}

		output := NewJoyStickInput("").OutputBytes()
		c.fp.Write(output[:])

		c.isOpened = false
		err = c.fp.Close()
		c.fp = nil
		close(c.closeChan)
		close(c.inputChan)
		close(c.sentKeyChan)
		close(c.abortSendKeyChan)
		c.wg.Done()
	}()
	return
}

func (c *JoyStick) Close() (err error) {
	if !c.isOpened {
		return
	}
	c.closeChan <- true
	c.wg.Wait()
	return
}

func (c *JoyStick) SyncSendKey(input *JoyStickInput, duration time.Duration) (err error) {
	if !c.isOpened {
		err = fmt.Errorf("device was closed")
		return
	}
	select {
	case <-time.After(time.Millisecond):
		err = fmt.Errorf("previous task not completed")
		return
	case <-c.sentKeyChan:
	}
	c.clearAbortChannel()
	defer func() {
		c.sentKeyChan <- true
	}()

	c.inputChan <- input
	if duration <= 0 {
		return
	}

	if duration < SendKeyMinDuration {
		duration = SendKeyMinDuration
	}
	t := time.NewTicker(duration)
	defer t.Stop()
	select {
	case <-t.C:
	case <-c.abortSendKeyChan:
	}
	c.inputChan <- NewJoyStickInput("")
	return
}

func (c *JoyStick) clearAbortChannel() (err error) {
	if !c.isOpened {
		err = fmt.Errorf("device was closed")
		return
	}
	cleared := false
	for !cleared {
		select {
		case <-c.abortSendKeyChan:
		case <-time.After(time.Millisecond):
			cleared = true
		}
	}
	return
}

func (c *JoyStick) AbortSendKey() (err error) {
	if !c.isOpened {
		err = fmt.Errorf("device was closed")
		return
	}
	go func() {
		c.abortSendKeyChan <- true
	}()
	return
}
