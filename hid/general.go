package hid

import (
	"errors"
	"fmt"
	"jensweisskopf/autowork/connection"
	"strings"
	"time"

	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
	log "github.com/sirupsen/logrus"
)

/*

pri
wri write a character
raw press a key
mod press a modification key
pre press a key
rel release all keys
rer release a key (rawcode)
mov move mouse
cli click mouse
mpr press mouse button
mre release mouse button
sta return status of the bl connection
eco echo serial input
deb

*/

type HID struct {
	com                connection.Communicator
	IsRecording        bool
	IsPlaying          bool
	onrecordline       func(string)
	scriptEnv          *env.Env
	StartTimestamp     time.Time
	LastCmdTimestamp   time.Time
	delayClickMouse    uint
	delayPressKey      uint
	delayMoveMouse     uint
	delayResetMouse    uint
	delayTransition    uint
	forcedDelay        uint
	lastCommand        uint
	CurrentLine        uint
	CurrentIteration   uint
	NumberOfIterations uint
}

func NewHID() *HID {
	hid := &HID{}
	hid.com = nil
	hid.IsPlaying = false
	hid.IsRecording = false

	log.Print("Create new global scope")
	hid.scriptEnv = env.NewEnv()
	err := hid.scriptEnv.Define("println", fmt.Println)
	if err != nil {
		log.Error(err)
	}
	hid.scriptEnv.Define("ResetMouse", hid.ResetMouse)
	hid.scriptEnv.Define("MoveMouse", hid.MoveMouse)
	hid.scriptEnv.Define("ClickMouse", hid.ClickMouse)
	hid.scriptEnv.Define("HitKey", hid.HitKey)
	hid.scriptEnv.Define("PressMod", hid.PressMod)
	hid.scriptEnv.Define("PressKey", hid.PressKey)
	hid.scriptEnv.Define("ReleaseKey", hid.ReleaseKey)
	hid.scriptEnv.Define("ReleaseKeys", hid.ReleaseKeys)
	hid.scriptEnv.Define("Delay", hid.Delay)
	hid.scriptEnv.Define("SetDelayPressMouse", hid.SetDelayPressMouse)
	hid.scriptEnv.Define("SetDelayPressKey", hid.SetDelayPressKey)
	hid.scriptEnv.Define("SetDelayMoveMouse", hid.SetDelayMoveMouse)
	hid.scriptEnv.Define("SetDelayResetMouse", hid.SetDelayResetMouse)
	hid.scriptEnv.Define("SetDelayTransition", hid.SetDelayTransition)

	hid.lastCommand = 0
	hid.forcedDelay = 0

	return hid
}

func (h *HID) SetCom(com connection.Communicator) {
	h.com = com
}

func (h *HID) CheckConnection() error {
	cmd := "sta"
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	response, err := h.readFromCommunicator()
	if err != nil {
		return err
	}
	log.Printf("Response = %s", response)
	if response != "ok" {
		return errors.New(response)
	}
	return nil
}

func (h *HID) StartRecording() {
	if !h.IsPlaying {
		h.IsRecording = true
		h.StartTimestamp = time.Now()
		h.LastCmdTimestamp = h.StartTimestamp
		log.Debug("Recording started")
	} else {
		log.Error("Can not record while playback")
	}

}

func (h *HID) StopRecording() {
	h.IsRecording = false
	h.IsPlaying = false
	log.Debug("Stop recording and playback")
}

// ResetMouse move the mouse to absolut zero
// Command id = 1
func (h *HID) ResetMouse() error {
	h.processingDelay(1, h.delayResetMouse)

	for i := 0; i < 10; i++ {
		cmd := "mov -127 -127 0 0"
		if err := h.writeToCommunicator(cmd); err != nil {
			return err
		}
	}
	if err := h.recordLine("ResetMouse()"); err != nil {
		return err
	}
	return nil
}

func (h *HID) MoveMouse(x int, y int, w int, sw int) error {
	h.processingDelay(2, h.delayMoveMouse)

	cmd := fmt.Sprintf("mov %d %d %d %d", x, y, w, sw)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("MoveMouse(%d, %d, %d, %d)", x, y, w, sw)
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) ClickMouse(b int) error {
	h.processingDelay(3, h.delayClickMouse)

	cmd := fmt.Sprintf("cli %d", b)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("ClickMouse(%d)", b)
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) PressKey(k int) error {
	h.processingDelay(4, h.delayPressKey)

	cmd := fmt.Sprintf("raw %d", k)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("PressKey(%d)", k)
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) HitKey(k int) error {
	h.processingDelay(4, h.delayPressKey)

	cmd := fmt.Sprintf("raw %d", k)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("rer %d", k)
	//cmd = fmt.Sprint("rel")
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("HitKey(%d) // %s", k, HID2str[k])
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) PressMod(k int) error {
	cmd := fmt.Sprintf("mod %d", k)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("PressMod(%d)", k)
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) ReleaseKey(k int) error {
	cmd := fmt.Sprintf("rer %d", k)
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = fmt.Sprintf("ReleaseKey(%d)", k)
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) ReleaseKeys() error {
	cmd := "rel"
	if err := h.writeToCommunicator(cmd); err != nil {
		return err
	}
	cmd = "ReleaseKeys()"
	if err := h.recordLine(cmd); err != nil {
		return err
	}
	return nil
}

func (h *HID) OnRecordLine(f func(l string)) {
	h.onrecordline = f
}

func (h *HID) writeToCommunicator(s string) error {
	s = s + "\n"
	if err := h.com.Write(s); err != nil {
		return err
	}
	log.Debug("Write:", s)
	return nil
}

func (h *HID) readFromCommunicator() (string, error) {
	s, err := h.com.Read()
	if err != nil {
		return "", err
	}
	log.Debug("Read:", s)
	return s, nil
}

func (h *HID) recordLine(s string) error {
	if h.IsRecording {
		if h.onrecordline != nil {
			h.onrecordline(s)
		} else {
			return errors.New("no record handler available")
		}
	}
	return nil
}

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func (h *HID) PlaybackRecording(script string, repeat int, startAtLine int) {
	if h.IsRecording {
		log.Error("Can't playback while recording")
		return
	}
	if h.IsPlaying {
		log.Error("Playback already running")
		return
	}
	log.Debug("Execute script")
	h.IsPlaying = true
	go h.processPlaybackRecording(script, repeat, startAtLine)
}

func (h *HID) processPlaybackRecording(script string, repeat int, startAtLine int) {
	var err error
	lines := strings.Split(strings.ReplaceAll(script, "\r\n", "\n"), "\n")

	for i := 1; i < repeat+1; i++ {
		log.Debugf("Iteration %d from %d", i, repeat)
		h.CurrentIteration = uint(i)
		h.NumberOfIterations = uint(repeat)
		for k := startAtLine; k < len(lines); k = k + 1 {
			log.Debugf("Execute line %d: %s", k, lines[k])
			h.CurrentLine = uint(k)
			_, err = vm.Execute(h.scriptEnv, nil, lines[k])
			if err != nil {
				log.Errorf("Error at line %d,%e", k, err)
			}
			if !h.IsPlaying {
				log.Errorf("Playback aborted after line %d", k)
				break
			}
		}
	}
	h.IsPlaying = false
}

func (h *HID) Delay(d uint) {
	h.forcedDelay = d
}

func (h *HID) SetDelayPressMouse(m uint) {
	h.delayClickMouse = m
}

func (h *HID) SetDelayPressKey(m uint) {
	h.delayPressKey = m
}

func (h *HID) SetDelayMoveMouse(m uint) {
	h.delayMoveMouse = m
}

func (h *HID) SetDelayResetMouse(m uint) {
	h.delayResetMouse = m
}

func (h *HID) SetDelayTransition(m uint) {
	h.delayTransition = m
}

func (h *HID) processingDelay(commandID uint, commandDelay uint) error {
	if h.lastCommand != commandID {
		if h.delayTransition > 0 {
			h.forcedDelay = h.delayTransition
		}
	} else {
		if commandDelay > 0 {
			h.forcedDelay = commandDelay
		}
	}
	if h.IsRecording {
		deltaInMillis := time.Since(h.LastCmdTimestamp) / 1000000
		h.LastCmdTimestamp = time.Now()
		cmd := fmt.Sprintf("Delay(%d)", uint(deltaInMillis))
		if err := h.recordLine(cmd); err != nil {
			return err
		}
	} else {
		time.Sleep(time.Duration(h.forcedDelay) * time.Millisecond)
	}
	h.lastCommand = commandID
	return nil
}
