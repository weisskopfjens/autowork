package main

import (
	"autowork/connection"
	"autowork/hid"
	"io"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

var mainapp fyne.App
var mainwindow fyne.Window
var combo1 *widget.Select
var recordingEntry *widget.Entry // Recordings
var entry5 *UintEntry            // Repeat
var entry6 *UintEntry            // Startat
var serialCom connection.SerialConnection
var com connection.Communicator
var hid1 *hid.HID
var status *widget.Label

var mousestep int

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	log.Debug("Loading config")
	appconfig := configuration{}
	appconfig.LoadConfig()

	log.Debug("Create HID")
	hid1 = hid.NewHID()

	log.Debug("Init gui")
	//mainapp = app.New()
	mainapp = app.NewWithID("jensweisskopf.autowork")
	mainwindow = mainapp.NewWindow("Autowork")

	mousestep = 1

	status = widget.NewLabel("Ready")
	StatusMsg = "Ready"

	//
	// Main menu
	//

	menuitem1 := fyne.NewMenuItem("Load", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mainwindow)
				return
			}
			if reader == nil {
				log.Debug("Cancel")
				return
			}
			var buffer []byte
			buffer, err = io.ReadAll(reader)
			if err != nil {
				log.Error(err)
			}
			recordingEntry.SetText(string(buffer))
			reader.Close()

			//loadData(reader)
		}, mainwindow)

		var err error
		var uri fyne.URI
		var luri fyne.ListableURI
		abspath, _ := filepath.Abs("./recordings")
		log.Debug(abspath)
		uri, err = storage.ParseURI(`file:` + abspath)
		if err != nil {
			log.Error(err)
		}
		luri, err = storage.ListerForURI(uri)
		if err != nil {
			log.Error(err)
		}
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.SetLocation(luri)
		fd.Show()
	})

	menuitem2 := fyne.NewMenuItem("Save", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, mainwindow)
				return
			}
			if writer == nil {
				log.Debug("Cancel")
				return
			}
			writer.Write([]byte(recordingEntry.Text))
			writer.Close()

			//saveData(writer, w)
		}, mainwindow)

		var err error
		var uri fyne.URI
		var luri fyne.ListableURI
		abspath, _ := filepath.Abs("./recordings")
		log.Debug(abspath)
		uri, err = storage.ParseURI(`file:` + abspath)
		if err != nil {
			log.Error(err)
		}
		luri, err = storage.ListerForURI(uri)
		if err != nil {
			log.Error(err)
		}

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.SetLocation(luri)
		fd.Show()
	})

	menuitem3 := fyne.NewMenuItem("begin", func() {
		if err := com.Begin(); err != nil {
			StatusMsg = "Error:" + err.Error()
		} else {
			if err := hid1.CheckConnection(); err != nil {
				StatusMsg = "Error:" + err.Error()
			}
			StatusMsg = "Connected"
		}
	})
	menuitem4 := fyne.NewMenuItem("end", func() {
		if err := com.End(); err != nil {
			StatusMsg = "Error:" + err.Error()
		} else {
			StatusMsg = "Disconnected"
		}
	})
	menuitem5 := fyne.NewMenuItem("start", HandleStartRecording)
	menuitem6 := fyne.NewMenuItem("stop", HandleStopRecording)
	menuitem7 := fyne.NewMenuItem("info", showinfo)
	menuitem8 := fyne.NewMenuItem("playback", HandlePlaybackRecording)
	menuitem9 := fyne.NewMenuItem("new", HandleNewRecording)

	menu1 := fyne.NewMenu("File", menuitem1, menuitem2)
	menu2 := fyne.NewMenu("Connection", menuitem3, menuitem4)
	menu3 := fyne.NewMenu("Record", menuitem9, menuitem5, menuitem6, menuitem8)
	menu4 := fyne.NewMenu("Help", menuitem7)

	mainmenu := fyne.NewMainMenu(menu1, menu2, menu3, menu4)

	mainwindow.SetMainMenu(mainmenu)
	mainwindow.SetMaster()

	//
	// Tab 1
	//

	button1 := widget.NewButton("L-SHIFT", HandleSpecialKey1)
	button2 := widget.NewButton("R-SHIFT", HandleSpecialKey2)
	button3 := widget.NewButton("L-META", HandleSpecialKey3)
	button4 := widget.NewButton("R-META", HandleSpecialKey4)
	button5 := widget.NewButton("ALT+TAB", HandleSpecialKey5)
	button6 := widget.NewButton("RELEASE", HandleSpecialKey6)

	commandGrid := container.New(layout.NewGridLayout(2),
		button1,
		button2,
		button3,
		button4,
		button5,
		button6,
	)

	catch1 := NewCatchEntry()
	catch1.SetPlaceHolder("Mouse")
	catch1.MultiLine = true
	catch1.OnTypedKey(HandleKeyInputMouse)

	catch2 := NewCatchEntry()
	catch2.SetPlaceHolder("Keyboard")
	catch2.MultiLine = true
	catch2.OnTypedKey(HandleKeyInputKeyboard)

	catchareas := container.NewGridWithColumns(2, catch1, catch2)

	commentEntry := widget.NewEntry()
	commentEntry.PlaceHolder = "Comment"
	commentButton := widget.NewButton("Insert comment", func() {
		recordingEntry.Append("// " + commentEntry.Text)
	})

	commentGrid := container.NewGridWithColumns(2, commentButton, commentEntry)

	tab1 := container.NewBorder(commandGrid, commentGrid, nil, nil, catchareas /*layout.NewSpacer()*/)

	//
	// Tab 2
	//

	label1 := widget.NewLabel("Serial port")
	combo1 := widget.NewSelect([]string{}, nil)
	combo1.SetOptions([]string{appconfig.Portname})
	combo1.SetSelected(appconfig.Portname)
	serialCom.SetPortName(appconfig.Portname)
	combo1.OnChanged = func(s string) {
		serialCom.SetPortName(s)
		appconfig.Portname = s
		appconfig.SaveConfig()
	}

	button10 := widget.NewButton("detect ports", func() {
		ports, err := serialCom.GetPorts()
		if err != nil {
			StatusMsg = "No serial port found"
		} else {
			StatusMsg = "Ports detected"
		}
		combo1.SetOptions(ports)
	})

	label2 := widget.NewLabel("Serial speed")
	entry2 := widget.NewEntry()
	entry2.SetText(appconfig.Speed)
	serialCom.SetSpeed(appconfig.Speed)
	entry2.OnChanged = func(s string) {
		serialCom.SetSpeed(s)
		appconfig.Speed = s
		appconfig.SaveConfig()
	}

	label3 := widget.NewLabel("IP")
	entry3 := widget.NewEntry()
	entry3.SetText(appconfig.IP)
	entry3.OnChanged = func(s string) { appconfig.IP = s; appconfig.SaveConfig() }

	label4 := widget.NewLabel("Com mode")
	combo2 := widget.NewSelect([]string{"Serial", "TCP", "UDP"}, nil)
	combo2.SetSelected(appconfig.Mode)
	HandleModeChange(appconfig.Mode)
	combo2.OnChanged = func(s string) {
		appconfig.Mode = s
		appconfig.SaveConfig()
		HandleModeChange(s)
	}

	tab2 := container.New(layout.NewFormLayout(),
		label1,
		combo1,
		layout.NewSpacer(),
		button10,
		label2,
		entry2,
		label3,
		entry3,
		label4,
		combo2,
	)

	//
	// Tab 3
	//

	button7 := widget.NewButton("Start", HandleStartRecording)
	button8 := widget.NewButton("Stop", HandleStopRecording)
	button9 := widget.NewButton("Playback", HandlePlaybackRecording)
	button11 := widget.NewButton("New", HandleNewRecording)
	label5 := widget.NewLabel("Repeat:")
	entry5 = NewUintEntry()
	label6 := widget.NewLabel("Start at line:")
	entry6 = NewUintEntry()
	entry5.SetText(strconv.Itoa(appconfig.Repeat))
	entry6.SetText(strconv.Itoa(appconfig.StartAtLine))
	entry5.OnChanged = func(s string) {
		var err error
		appconfig.Repeat, err = strconv.Atoi(s)
		if err != nil {
			log.Error("Configuration: Wrong repeat value")
			return
		}
		appconfig.SaveConfig()
	}
	entry6.OnChanged = func(s string) {
		var err error
		appconfig.StartAtLine, err = strconv.Atoi(s)
		if err != nil {
			log.Error("Configuration: Wrong start at line value")
			return
		}
		appconfig.SaveConfig()
	}

	label5.Alignment = fyne.TextAlignTrailing
	label6.Alignment = fyne.TextAlignTrailing
	playbackcontrol := container.NewGridWithColumns(5, button9, label5, entry5, label6, entry6)
	recordcontrol := container.NewGridWithColumns(3, button11, button7, button8)

	transportcontrol := container.NewVBox(recordcontrol, playbackcontrol)

	// Script / Recordings
	recordingEntry = widget.NewEntry()
	recordingEntry.MultiLine = true

	tab3 := container.NewBorder(transportcontrol, nil, nil, nil, recordingEntry)

	//
	// Tab
	//

	tabs := container.NewAppTabs(
		container.NewTabItem("Workspace", tab1),
		container.NewTabItem("Settings", tab2),
		container.NewTabItem("Recordings", tab3),
	)

	//
	// Main area
	//

	mainwindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		log.Debug(ev.Name)
		log.Debug(ev.Physical.ScanCode)
	})

	content := container.NewBorder(nil, status, nil, nil, tabs)
	mainwindow.SetContent(content)

	hid1.OnRecordLine(HandleRecordLine)
	go UpdateStatus(hid1, com, status)

	log.Debug("Show and run")

	mainwindow.Resize(fyne.NewSize(600, 300))
	mainwindow.ShowAndRun()
}
