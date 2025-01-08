package main

import (
	"fmt"
	"jensweisskopf/autowork/hid"
	"strconv"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
)

func HandleModeChange(s string) {
	switch s {
	case "Serial":
		com = &serialCom
		hid1.SetCom(com)
		log.Debugf("Set mode to %s", s)
	default:
		log.Debugf("Unsupported communication mode %s", s)
	}
}

func HandleKeyInputKeyboard(ke *fyne.KeyEvent) {
	log.Debugf("Keyboard: %s %s", ke.Name, strconv.Itoa(ke.Physical.ScanCode))
	if val, ok := hid.Keycode2HIDmod[ke.Physical.ScanCode]; ok {
		hid1.PressMod(val)
	} else {
		hid1.HitKey(hid.Keycode2HID[ke.Physical.ScanCode])
	}
}

func HandleKeyInputMouse(ke *fyne.KeyEvent) {
	log.Debug("Mouse: " + ke.Name)
	switch ke.Name {
	case "Left":
		hid1.MoveMouse(mousestep*-1, 0, 0, 0)

	case "Right":
		hid1.MoveMouse(mousestep, 0, 0, 0)

	case "Down":
		hid1.MoveMouse(0, mousestep, 0, 0)

	case "Up":
		hid1.MoveMouse(0, mousestep*-1, 0, 0)

	case "1":
		mousestep = 1

	case "2":
		mousestep = 5

	case "3":
		mousestep = 10

	case "4":
		mousestep = 20

	case "5":
		mousestep = 30

	case "6":
		mousestep = 50

	case "7":
		mousestep = 80

	case "8":
		mousestep = 100

	case "0":
		hid1.ResetMouse()

	case "Q":
		hid1.ClickMouse(1)

	case "W":
		hid1.ClickMouse(2)

	case "E":
		hid1.ClickMouse(3)

	case "R":
		hid1.ClickMouse(4)

	case "T":
		hid1.ClickMouse(5)
	}

}

func HandleRecordLine(l string) {
	recordingEntry.Append(l + "\n")
}

func HandlePlaybackRecording() {
	fmt.Println(entry5.Text)
	fmt.Println(entry6.Text)

	repeat, err := strconv.Atoi(entry5.Text)
	if err != nil {
		log.Error("Repeat must be a number")
		return
	}
	startat, err := strconv.Atoi(entry6.Text)
	if err != nil {
		log.Error("Start at line must be a number")
		return
	}
	status.SetText("Recording")
	hid1.PlaybackRecording(recordingEntry.Text, repeat, startat)
}

func HandleStartRecording() {
	hid1.StartRecording()
}

func HandleStopRecording() {
	hid1.StopRecording()
}

func HandleNewRecording() {
	defaulttext := `SetDelayPressMouse(1000)
SetDelayPressKey(1000)
SetDelayMoveMouse(100)
SetDelayResetMouse(1000)
SetDelayTransition(3000)
`
	recordingEntry.SetText(defaulttext)
}

func HandleSpecialKey1() {

}

func HandleSpecialKey2() {

}

func HandleSpecialKey3() {
	hid1.PressMod(0x08)
}

func HandleSpecialKey4() {
	hid1.PressMod(0x80)
}

func HandleSpecialKey5() {
	hid1.PressMod(0x04) //Alt_L
	hid1.PressKey(0x2B) //Tab
	hid1.ReleaseKeys()  //Release
}

func HandleSpecialKey6() {
	hid1.ReleaseKeys()
}
