package main

import (
	"strconv"

	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// FloatEntry
type FloatEntry struct {
	widget.Entry
}

func NewFloatEntry() *FloatEntry {
	entry := &FloatEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *FloatEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

/*func (e *FloatEntry) TypedKey(key *fyne.KeyEvent) {
	fmt.Printf("FloatEntry: key %s", key.Name)
}*/

func (e *FloatEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *FloatEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// UintEntry
type UintEntry struct {
	widget.Entry
}

func NewUintEntry() *UintEntry {
	entry := &UintEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *UintEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

/*func (e *UintEntry) TypedKey(key *fyne.KeyEvent) {
	fmt.Printf("UintEntry: key %s", key.Name)
}*/

func (e *UintEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseUint(content, 10, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *UintEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// IntEntry
type IntEntry struct {
	widget.Entry
}

func NewIntEntry() *IntEntry {
	entry := &IntEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *IntEntry) TypedRune(r rune) {
	if len(e.Text) == 0 && r == '-' {
		e.Entry.TypedRune(r)
	}
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

/*func (e *IntEntry) TypedKey(key *fyne.KeyEvent) {
	fmt.Printf("IntEntry: key %s", key.Name)
}*/

func (e *IntEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseInt(content, 10, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *IntEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

//
// CatchEntry
//

type CatchEntry struct {
	widget.Entry
	handleTypedKey  func(*fyne.KeyEvent)
	handleTypedRune func(rune)
}

func NewCatchEntry() *CatchEntry {
	entry := &CatchEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *CatchEntry) TypedRune(r rune) {
	if e.handleTypedRune != nil {
		e.handleTypedRune(r)
	}
}

func (e *CatchEntry) OnTypedRune(f func(rune)) {
	e.handleTypedRune = f
}

func (e *CatchEntry) TypedKey(key *fyne.KeyEvent) {
	if e.handleTypedKey != nil {
		e.handleTypedKey(key)
	}
}

func (e *CatchEntry) OnTypedKey(f func(*fyne.KeyEvent)) {
	e.handleTypedKey = f
}

func (e *CatchEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *CatchEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// TimeEntry
type TimeEntry struct {
	widget.Entry
}

func NewTimeEntry() *TimeEntry {
	entry := &TimeEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *TimeEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		if len(e.Entry.Text) == 2 {
			e.Entry.TypedRune(':')
		}
		if len(e.Entry.Text) == 5 {
			e.Entry.TypedRune(':')
		}
		e.Entry.TypedRune(r)
	}
}

/*func (e *FloatEntry) TypedKey(key *fyne.KeyEvent) {
	fmt.Printf("FloatEntry: key %s", key.Name)
}*/

func (e *TimeEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *TimeEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
