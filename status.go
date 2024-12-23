package main

import (
	"autowork/connection"
	"autowork/hid"
	"fmt"
	"time"

	"fyne.io/fyne/v2/widget"
)

var StatusMsg string

func UpdateStatus(hid *hid.HID, con connection.Communicator, w *widget.Label) {
	for {
		var PlaybackInfo string
		if hid.IsPlaying {
			PlaybackInfo = fmt.Sprintf("Iteration %d of %d | Current line %d",
				hid.CurrentIteration,
				hid.NumberOfIterations,
				hid.CurrentLine)
		}
		status := fmt.Sprintf("%s | Connection: %t | Record: %t | Play: %t | %s",
			StatusMsg,
			con.IsConnected(),
			hid.IsRecording,
			hid.IsPlaying,
			PlaybackInfo)

		w.SetText(status)
		time.Sleep(100 * time.Millisecond)
	}
}
