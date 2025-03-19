package main

import (
	"fmt"
	"jensweisskopf/autowork/connection"
	"jensweisskopf/autowork/hid"
	"time"

	"fyne.io/fyne/v2/widget"
)

func UpdateStatus(hid *hid.HID, con connection.Communicator, w *widget.Label) {
	for {
		var PlaybackInfo string
		if hid.IsPlaying {
			PlaybackInfo = fmt.Sprintf("Iteration %d of %d | Current line %d | Duration %ds | Estimate %ds",
				hid.CurrentIteration,
				hid.NumberOfIterations,
				hid.CurrentLine,
				uint(hid.DurationPerIteration.Abs().Seconds()),
				uint(hid.EstimatedDuration.Abs().Seconds()))
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
