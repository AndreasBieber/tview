package tview

import (
	"strings"

	"github.com/gdamore/tcell"
)

// DropDownOption is one option that can be selected in a drop-down primitive.
type DropDownOption struct {
	Text      string      // The text to be displayed in the drop-down.
	Selected  func()      // The (optional) callback for when this option was selected.
	Reference interface{} // An optional reference object.
}

func NewDropDownOption(text string) *DropDownOption {
	return &DropDownOption{Text: text}
}

// SetSelectedFunc sets the handler to be called when this option is selected.
func (d *DropDownOption) SetSelectedFunc(handler func()) *DropDownOption {
	d.Selected = handler
	return d
}

// SetReference allows you to store a reference of any type in this option.
func (d *DropDownOption) SetReference(reference interface{}) *DropDownOption {
	d.Reference = reference
	return d
}

// DropDown implements a selection widget whose options become visible in a
// drop-down list when activated.
//
// See https://github.com/rivo/tview/wiki/DropDown for an example.
type DropDown struct {
	*Box

	// The options from which the user can choose.
	options []*DropDownOption

	// Strings to be placed before and after each drop-down option.
	optionPrefix, optionSuffix string

	// The index of the currently selected option. Negative if no option is
	// currently selected.
	currentOption int

	// Strings to be placed before and after the current option.
	currentOptionPrefix, currentOptionSuffix string

	// The text to be displayed when no option has yet been selected.
	noSelection string

	// Set to true if the options are visible and selectable.
	open bool

	// The runes typed so far to directly access one of the list items.
	prefix string

	// The list element for the options.
	list *List

	// The text to be displayed before the input area.
	label string

	// The label color.
	labelColor tcell.Color

	// The label color when focused.
	labelColorActivated tcell.Color

	// The background color of the input area.
	fieldBackgroundColor tcell.Color

	// The background color of the input area when focused.
	fieldBackgroundColorActivated tcell.Color

	// The text color of the input area.
	fieldTextColor tcell.Color

	// The text color of the input area when focused.
	fieldTextColorActivated tcell.Color

	// The color for prefixes.
	prefixTextColor tcell.Color

	// The screen width of the label area. A value of 0 means use the width of
	// the label text.
	labelWidth int

	// The screen width of the input area. A value of 0 means extend as much as
	// possible.
	fieldWidth int

	// An optional function which is called when the user indicated that they
	// are done selecting options. The key which was pressed is provided (tab,
	// shift-tab, or escape).
	done func(tcell.Key)

	// A callback function set by the Form class and called when the user leaves
	// this form item.
	finished func(tcell.Key)

	// A callback function which is called when the user selected the drop-down's
	// selection.
	selected func(index int, option *DropDownOption)

	dragging bool // Set to true when mouse dragging is in progress.

	// The symbol to draw at the end of the field.
	dropDownSymbol rune

	// The chars to show when the option's text gets shortened.
	abbreviationChars string
}

// NewDropDown returns a new drop-down.
func NewDropDown() *DropDown {
	list := NewList()
	list.ShowSecondaryText(false).
		SetMainTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedBackgroundColor(Styles.PrimaryTextColor).
		SetHighlightFullLine(true).
		SetBackgroundColor(Styles.MoreContrastBackgroundColor)

	d := &DropDown{
		Box:                           NewBox(),
		currentOption:                 -1,
		list:                          list,
		labelColor:                    Styles.SecondaryTextColor,
		labelColorActivated:           Styles.SecondaryTextColor,
		fieldBackgroundColor:          Styles.ContrastBackgroundColor,
		fieldBackgroundColorActivated: Styles.PrimaryTextColor,
		fieldTextColor:                Styles.PrimaryTextColor,
		fieldTextColorActivated:       Styles.ContrastBackgroundColor,
		prefixTextColor:               Styles.ContrastSecondaryTextColor,
		dropDownSymbol:                '▼',
		abbreviationChars:             "...",
	}

	d.focus = d

	return d
}

// SetCurrentOption sets the index of the currently selected option. This may
// be a negative value to indicate that no option is currently selected. Calling
// this function will also trigger the "selected" callback (if there is one).
func (d *DropDown) SetCurrentOption(index int) *DropDown {
	if index >= 0 && index < len(d.options) {
		d.currentOption = index
		d.list.SetCurrentItem(index)
		if d.selected != nil {
			d.selected(index, d.options[index])
		}
		if d.options[index].Selected != nil {
			d.options[index].Selected()
		}
	} else {
		d.currentOption = -1
		d.list.SetCurrentItem(0) // Set to 0 because -1 means "last item".
		if d.selected != nil {
			d.selected(-1, nil)
		}
	}
	return d
}

// GetCurrentOption returns the index of the currently selected option as well
// as its text. If no option was selected, -1 and an empty string is returned.
func (d *DropDown) GetCurrentOption() (int, string) {
	var text string
	if d.currentOption >= 0 && d.currentOption < len(d.options) {
		text = d.options[d.currentOption].Text
	}
	return d.currentOption, text
}

// SetTextOptions sets the text to be placed before and after each drop-down
// option (prefix/suffix), the text placed before and after the currently
// selected option (currentPrefix/currentSuffix) as well as the text to be
// displayed when no option is currently selected. Per default, all of these
// strings are empty.
func (d *DropDown) SetTextOptions(prefix, suffix, currentPrefix, currentSuffix, noSelection string) *DropDown {
	d.currentOptionPrefix = currentPrefix
	d.currentOptionSuffix = currentSuffix
	d.noSelection = noSelection
	d.optionPrefix = prefix
	d.optionSuffix = suffix
	for index := 0; index < d.list.GetItemCount(); index++ {
		d.list.SetItemText(index, prefix+d.options[index].Text+suffix, "")
	}
	return d
}

// SetLabel sets the text to be displayed before the input area.
func (d *DropDown) SetLabel(label string) *DropDown {
	d.label = label
	return d
}

// GetLabel returns the text to be displayed before the input area.
func (d *DropDown) GetLabel() string {
	return d.label
}

// SetLabelWidth sets the screen width of the label. A value of 0 will cause the
// primitive to use the width of the label string.
func (d *DropDown) SetLabelWidth(width int) *DropDown {
	d.labelWidth = width
	return d
}

// SetLabelColor sets the color of the label.
func (d *DropDown) SetLabelColor(color tcell.Color) *DropDown {
	d.labelColor = color
	return d
}

// SetLabelColorActivated sets the color of the label when focused.
func (d *DropDown) SetLabelColorActivated(color tcell.Color) *DropDown {
	d.labelColorActivated = color
	return d
}

// SetFieldBackgroundColor sets the background color of the options area.
func (d *DropDown) SetFieldBackgroundColor(color tcell.Color) *DropDown {
	d.fieldBackgroundColor = color
	return d
}

// SetFieldBackgroundColorActivated sets the background color of the options area when focused.
func (d *DropDown) SetFieldBackgroundColorActivated(color tcell.Color) *DropDown {
	d.fieldBackgroundColorActivated = color
	return d
}

// SetFieldTextColor sets the text color of the options area.
func (d *DropDown) SetFieldTextColor(color tcell.Color) *DropDown {
	d.fieldTextColor = color
	return d
}

// SetFieldTextColorActivated sets the text color of the options area when focused.
func (d *DropDown) SetFieldTextColorActivated(color tcell.Color) *DropDown {
	d.fieldTextColorActivated = color
	return d
}

// SetDropDownTextColor sets text color of the drop down list.
func (d *DropDown) SetDropDownTextColor(color tcell.Color) *DropDown {
	d.list.SetMainTextColor(color)
	return d
}

// SetDropDownBackgroundColor sets the background color of the drop list.
func (d *DropDown) SetDropDownBackgroundColor(color tcell.Color) *DropDown {
	d.list.SetBackgroundColor(color)
	return d
}

// The text color of the selected option in the drop down list.
func (d *DropDown) SetDropDownSelectedTextColor(color tcell.Color) *DropDown {
	d.list.SetSelectedTextColor(color)
	return d
}

// The background color of the selected option in the drop down list.
func (d *DropDown) SetDropDownSelectedBackgroundColor(color tcell.Color) *DropDown {
	d.list.SetSelectedBackgroundColor(color)
	return d
}

// SetPrefixTextColor sets the color of the prefix string. The prefix string is
// shown when the user starts typing text, which directly selects the first
// option that starts with the typed string.
func (d *DropDown) SetPrefixTextColor(color tcell.Color) *DropDown {
	d.prefixTextColor = color
	return d
}

// SetFormAttributes sets attributes shared by all form items.
func (d *DropDown) SetFormAttributes(labelWidth int, bgColor, labelColor, labelColorActivated, fieldTextColor, fieldTextColorActivated, fieldBgColor, fieldBgColorActivated tcell.Color) FormItem {
	d.labelWidth = labelWidth
	d.backgroundColor = bgColor
	d.labelColor = labelColor
	d.labelColorActivated = labelColorActivated
	d.fieldTextColor = fieldTextColor
	d.fieldTextColorActivated = fieldTextColorActivated
	d.fieldBackgroundColor = fieldBgColor
	d.fieldBackgroundColorActivated = fieldBgColorActivated
	return d
}

// SetFieldWidth sets the screen width of the options area. A value of 0 means
// extend to as long as the longest option text.
func (d *DropDown) SetFieldWidth(width int) *DropDown {
	d.fieldWidth = width
	return d
}

// GetFieldWidth returns this primitive's field screen width.
func (d *DropDown) GetFieldWidth() int {
	if d.fieldWidth > 0 {
		return d.fieldWidth
	}
	fieldWidth := 0
	for _, option := range d.options {
		width := TaggedStringWidth(option.Text)
		if width > fieldWidth {
			fieldWidth = width
		}
	}
	fieldWidth += len(d.currentOptionPrefix) + len(d.currentOptionSuffix)
	fieldWidth += 3 // space + drop down symbol + space
	return fieldWidth
}

// AddOption adds a new selectable option to this drop-down.
func (d *DropDown) AddOption(option *DropDownOption) *DropDown {
	d.options = append(d.options, option)
	d.list.AddItem(d.optionPrefix+option.Text+d.optionSuffix, "", 0, nil)
	return d
}

// SetOptions replaces all current options with the ones provided and installs
// one callback function which is called when one of the options is selected.
// It will be called with the index and options. The "selected" parameter may be nil.
func (d *DropDown) SetOptions(texts []string, selected func(index int, option *DropDownOption)) *DropDown {
	d.list.Clear()
	d.options = nil
	for index, text := range texts {
		func(t string, i int) {
			d.AddOption(NewDropDownOption(text))
		}(text, index)
	}
	d.selected = selected
	return d
}

// SetSelectedFunc sets a handler which is called when the user selects the
// drop-down's option. This handler will be called in addition and prior to
// an option's optional individual handler. The handler is provided with the
// selected index and option. If "no option" was selected, these values
// are a -1 and nil.
func (d *DropDown) SetSelectedFunc(handler func(index int, option *DropDownOption)) *DropDown {
	d.selected = handler
	return d
}

// SetChangedFunc sets a handler which is called when the user changes the
// drop-down's option. This handler will be called in addition and prior to
// an option's optional individual handler. The handler is provided with the
// selected index and option. If "no option" was selected, these values
// are -1 and nil.
func (d *DropDown) SetChangedFunc(handler func(index int, option *DropDownOption)) *DropDown {
	d.list.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
		handler(index, d.options[index])
	})
	return d
}

// SetDoneFunc sets a handler which is called when the user is done selecting
// options. The callback function is provided with the key that was pressed,
// which is one of the following:
//
//   - KeyEscape: Abort selection.
//   - KeyTab: Move to the next field.
//   - KeyBacktab: Move to the previous field.
func (d *DropDown) SetDoneFunc(handler func(key tcell.Key)) *DropDown {
	d.done = handler
	return d
}

// SetFinishedFunc sets a callback invoked when the user leaves this form item.
func (d *DropDown) SetFinishedFunc(handler func(key tcell.Key)) FormItem {
	d.finished = handler
	return d
}

// Draw draws this primitive onto the screen.
func (d *DropDown) Draw(screen tcell.Screen) {
	d.Box.Draw(screen)

	labelColor := d.labelColor
	fieldBackgroundColor := d.fieldBackgroundColor
	fieldTextColor := d.fieldTextColor
	if d.GetFocusable().HasFocus() {
		labelColor = d.labelColorActivated
		fieldBackgroundColor = d.fieldBackgroundColorActivated
		fieldTextColor = d.fieldTextColorActivated
	}
	// Prepare.
	x, y, width, height := d.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	if d.labelWidth > 0 {
		labelWidth := d.labelWidth
		if labelWidth > rightLimit-x {
			labelWidth = rightLimit - x
		}
		Print(screen, d.label, x, y, labelWidth, AlignLeft, labelColor)
		x += labelWidth
	} else {
		_, drawnWidth := Print(screen, d.label, x, y, rightLimit-x, AlignLeft, labelColor)
		x += drawnWidth
	}

	// What's the longest option text?
	maxWidth := 0
	optionWrapWidth := TaggedStringWidth(d.optionPrefix + d.optionSuffix)
	for _, option := range d.options {
		strWidth := TaggedStringWidth(option.Text) + optionWrapWidth
		if strWidth > maxWidth {
			maxWidth = strWidth
		}
	}

	// Draw selection area.
	fieldWidth := d.fieldWidth
	if fieldWidth == 0 {
		fieldWidth = maxWidth
		if d.currentOption < 0 {
			noSelectionWidth := TaggedStringWidth(d.noSelection)
			if noSelectionWidth > fieldWidth {
				fieldWidth = noSelectionWidth
			}
		} else if d.currentOption < len(d.options) {
			currentOptionWidth := TaggedStringWidth(d.currentOptionPrefix + d.options[d.currentOption].Text + d.currentOptionSuffix)
			if currentOptionWidth > fieldWidth {
				fieldWidth = currentOptionWidth
			}
		}
	}
	if rightLimit-x < fieldWidth {
		fieldWidth = rightLimit - x
	}

	fieldStyle := tcell.StyleDefault.Background(fieldBackgroundColor)
	for index := 0; index < fieldWidth; index++ {
		screen.SetContent(x+index, y, ' ', nil, fieldStyle)
	}

	// Draw selected text.
	if d.open && len(d.prefix) > 0 {
		// Show the prefix.
		currentOptionPrefixWidth := TaggedStringWidth(d.currentOptionPrefix)
		prefixWidth := stringWidth(d.prefix)
		listItemText := d.options[d.list.GetCurrentItem()].Text
		Print(screen, d.currentOptionPrefix, x, y, fieldWidth, AlignLeft, fieldTextColor)
		Print(screen, d.prefix, x+currentOptionPrefixWidth, y, fieldWidth-currentOptionPrefixWidth, AlignLeft, d.prefixTextColor)
		if len(d.prefix) < len(listItemText) {
			Print(screen, listItemText[len(d.prefix):]+d.currentOptionSuffix, x+prefixWidth+currentOptionPrefixWidth, y, fieldWidth-prefixWidth-currentOptionPrefixWidth, AlignLeft, d.fieldTextColor)
		}
	} else {
		color := fieldTextColor
		text := d.noSelection
		if d.currentOption >= 0 && d.currentOption < len(d.options) {
			text = d.currentOptionPrefix + d.options[d.currentOption].Text + d.currentOptionSuffix
		}
		if fieldWidth > len(d.abbreviationChars)+3 && len(text) > fieldWidth {
			text = text[0:fieldWidth-3-len(d.abbreviationChars)] + d.abbreviationChars
		}
		Print(screen, text, x, y, fieldWidth, AlignLeft, color)
	}

	// Draw drop down symbol
	screen.SetContent(x+fieldWidth-2, y, d.dropDownSymbol, nil, new(tcell.Style).Foreground(fieldTextColor).Background(fieldBackgroundColor))

	// Draw options list.
	if d.HasFocus() && d.open {
		// We prefer to drop down but if there is no space, maybe drop up?
		lx := x
		ly := y + 1
		lwidth := maxWidth
		lheight := len(d.options)
		_, sheight := screen.Size()
		if ly+lheight >= sheight && ly-2 > lheight-ly {
			ly = y - lheight
			if ly < 0 {
				ly = 0
			}
		}
		if ly+lheight >= sheight {
			lheight = sheight - ly
		}
		d.list.SetRect(lx, ly, lwidth, lheight)
		d.list.Draw(screen)
	}
}

// InputHandler returns the handler for this primitive.
func (d *DropDown) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p Primitive)) {
		// Process key event.
		switch key := event.Key(); key {
		case tcell.KeyEnter, tcell.KeyRune, tcell.KeyDown:
			d.prefix = ""

			// If the first key was a letter already, it becomes part of the prefix.
			if r := event.Rune(); key == tcell.KeyRune && r != ' ' {
				d.prefix += string(r)
				d.evalPrefix()
			}

			d.openList(setFocus)
		case tcell.KeyEscape, tcell.KeyTab, tcell.KeyBacktab:
			if d.done != nil {
				d.done(key)
			}
			if d.finished != nil {
				d.finished(key)
			}
		}
	})
}

// evalPrefix selects an item in the drop-down list based on the current prefix.
func (d *DropDown) evalPrefix() {
	if len(d.prefix) > 0 {
		for index, option := range d.options {
			if strings.HasPrefix(strings.ToLower(option.Text), d.prefix) {
				d.list.SetCurrentItem(index)
				return
			}
		}

		// Prefix does not match any item. Remove last rune.
		r := []rune(d.prefix)
		d.prefix = string(r[:len(r)-1])
	}
}

// openList hands control over to the embedded List primitive.
func (d *DropDown) openList(setFocus func(Primitive)) {
	d.open = true
	optionBefore := d.currentOption

	d.list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if d.dragging {
			return // If we're dragging the mouse, we don't want to trigger any events.
		}

		// An option was selected. Close the list again.
		d.currentOption = index
		d.closeList(setFocus)

		// Trigger "selected" event.
		if d.selected != nil {
			d.selected(d.currentOption, d.options[d.currentOption])
		}
		if d.options[d.currentOption].Selected != nil {
			d.options[d.currentOption].Selected()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			d.prefix += string(event.Rune())
			d.evalPrefix()
		} else if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			if len(d.prefix) > 0 {
				r := []rune(d.prefix)
				d.prefix = string(r[:len(r)-1])
			}
			d.evalPrefix()
		} else if event.Key() == tcell.KeyEscape {
			d.currentOption = optionBefore
			d.list.SetCurrentItem(d.currentOption)
			d.closeList(setFocus)
			if d.selected != nil {
				if d.currentOption > -1 {
					d.selected(d.currentOption, d.options[d.currentOption])
				}
			}
		} else {
			d.prefix = ""
		}

		return event
	})

	setFocus(d.list)
}

// closeList closes the embedded List element by hiding it and removing focus
// from it.
func (d *DropDown) closeList(setFocus func(Primitive)) {
	d.open = false
	if d.list.HasFocus() {
		setFocus(d)
	}
}

// Focus is called by the application when the primitive receives focus.
func (d *DropDown) Focus(delegate func(p Primitive)) {
	d.Box.Focus(delegate)
	if d.open {
		delegate(d.list)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (d *DropDown) HasFocus() bool {
	if d.open {
		return d.list.HasFocus()
	}
	return d.hasFocus
}

// MouseHandler returns the mouse handler for this primitive.
func (d *DropDown) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return d.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		// Was the mouse event in the drop-down box itself (or on its label)?
		x, y := event.Position()
		_, rectY, _, _ := d.GetInnerRect()
		inRect := y == rectY
		if !d.open && !inRect {
			return d.InRect(x, y), nil // No, and it's not expanded either. Ignore.
		}

		// Handle dragging. Clicks are implicitly handled by this logic.
		switch action {
		case MouseLeftDown:
			consumed = d.open || inRect
			capture = d
			if !d.open {
				d.openList(setFocus)
				d.dragging = true
			} else if consumed, _ := d.list.MouseHandler()(MouseLeftClick, event, setFocus); !consumed {
				d.closeList(setFocus) // Close drop-down if clicked outside of it.
			}
		case MouseMove:
			if d.dragging {
				// We pretend it's a left click so we can see the selection during
				// dragging. Because we don't act upon it, it's not a problem.
				d.list.MouseHandler()(MouseLeftClick, event, setFocus)
				consumed = true
				capture = d
			}
		case MouseLeftUp:
			if d.dragging {
				d.dragging = false
				d.list.MouseHandler()(MouseLeftClick, event, setFocus)
				consumed = true
			}
		}

		return
	})
}
