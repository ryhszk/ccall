package ui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	te "github.com/muesli/termenv"
	"golang.org/x/term"
)

const unfocusedColor = "8"
const mainColor = "10"
const registColor = "13"
const searchColor = "14"

var (
	color                 = te.ColorProfile().Color
	focusedPrompt         = blinkString(">")
	blurredPrompt         = " "
	focusedRegisterButton = "[ " + colorString("Register", registColor) + " ]"
	blurredRegisterButton = "[ " + colorString("Register", unfocusedColor) + " ]"
	focusedSearchButton   = "[ " + colorString("Search", searchColor) + " ]"
	blurredSearchButton   = "[ " + colorString("Search", unfocusedColor) + " ]"
)

func colorString(str, ccode string) string {
	return te.String(str).Foreground(color(ccode)).String()
}

func blinkString(str string) string {
	return te.String(str).Blink().String()
}

func NewApp() error {
	return tea.NewProgram(initModel()).Start()
}

type Mode int

const (
	_ Mode = iota
	Main
	Regist
	Search
)

type model struct {
	cursor int
	mode   Mode
	main   []mainModel
	regist registModel
	search searchModel
	tsize  termSize
}

type mainModel struct {
	name string
	cmd  string
}

type registModel struct {
	nameInput  textinput.Model
	cmdInput   textinput.Model
	doneButton string
}

type searchModel struct {
	searchInput  textinput.Model
	searchButton string
}

type termSize struct {
	width  int
	height int
}

func initMainModel() []mainModel {
	// names := []textinput.Model{}
	// cmds := []textinput.Model{}

	// w, _, _ := term.GetSize(syscall.Stdin)
	// adjw := w - 27
	ems := []mainModel{}
	for i := 0; i < 3; i++ {
		var em mainModel
		em.cmd = fmt.Sprintf("ps aux | awk '{ if(NR>1){p[$1] += $3; n[$1]++ }}END{for(i in p) print p[i], n[i], i}' : %d", i)
		// str := "Command nameaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbccccccccccccccccccccccccccccccccccccccccccccc"
		str := "Command nameaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		// em.name = str[:adjw]
		em.name = str
		// em.nameInput.CharLimit = 50
		// em.nameInput.Width = 50
		// if i != 0 {
		// 	em.nameInput.Prompt = blurredPrompt
		// } else {
		// 	em.nameInput.Focus()
		// 	em.nameInput.Prompt = focusedPrompt
		// 	em.nameInput.TextColor = focusedColor
		// }
		ems = append(ems, em)
	}

	return ems
}

func initRegistModel() registModel {
	name := textinput.NewModel()
	name.Placeholder = "Input the distinguished name of the following command."
	name.Focus()
	name.Prompt = focusedPrompt
	name.TextColor = registColor
	name.CharLimit = 50
	name.Width = 50

	cmd := textinput.NewModel()
	cmd.Placeholder = "Input the command you want to run in shell."
	cmd.Prompt = blurredPrompt
	cmd.CharLimit = 500
	cmd.Width = 50

	return registModel{
		name,
		cmd,
		blurredRegisterButton,
	}
}

func initSearchModel() searchModel {
	search := textinput.NewModel()
	search.Placeholder = "Input the string you wish to search for."
	search.Focus()
	search.Prompt = focusedPrompt
	search.TextColor = searchColor
	search.CharLimit = 500
	search.Width = 10

	return searchModel{
		search,
		blurredSearchButton,
	}
}

func initModel() model {
	ts := termSize{
		width:  0,
		height: 0,
	}

	return model{
		cursor: 0,
		mode:   Main,
		main:   initMainModel(),
		regist: initRegistModel(),
		search: initSearchModel(),
		tsize:  ts,
	}

}
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	offset := 20
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tsize.width = msg.Width - offset
		m.tsize.height = msg.Height - offset
	}

	switch m.mode {
	case Main:
		return m.mainUpdate(msg)
	case Regist:
		return m.RegistUpdate(msg)
	case Search:
		return m.searchUpdate(msg)
	}

	return m, nil
}

// func (m *model) cycleCursor() {

// 	if m.cursor > len(m.txtModels)-1 {
// 		m.cursor = 0
// 	} else if m.cursor < 0 {
// 		m.cursor = len(m.txtModels) - 1
// 	}

// 	for i := 0; i <= len(m.txtModels)-1; i++ {
// 		if i == m.cursor {
// 			// Set focused state
// 			m.txtModels[i].Focus()
// 			m.txtModels[i].Prompt = focusedPrompt
// 			m.txtModels[i].TextColor = focusedColor
// 			continue
// 		}
// 		// Remove focused state
// 		m.txtModels[i].Blur()
// 		m.txtModels[i].Prompt = blurredPrompt
// 		m.txtModels[i].TextColor = ""
// 	}
// }

// func (m model) cycleCursor(t []textinput.Model) {

// 	if m.cursor > len(t) {
// 		m.cursor = 0
// 	} else if m.cursor < -0 {
// 		m.cursor = len(t) - 1
// 	}

// 	for i := 0; i <= len(t)-1; i++ {
// 		if i == m.cursor {
// 			// Set focused state
// 			t.txtModels[i].Focus()
// 			m.txtModels[i].Prompt = focusedPrompt
// 			m.txtModels[i].TextColor = focusedColor
// 			continue
// 		}
// 		// Remove focused state
// 		m.txtModels[i].Blur()
// 		m.txtModels[i].Prompt = blurredPrompt
// 		m.txtModels[i].TextColor = ""
// 	}
// }

func (m model) mainUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "r":
			m.mode = Regist
			return m, nil

		case "s":
			m.mode = Search
			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		// Cycle between inputs
		case "tab", "shift+tab", "enter", "up", "down":

			inputs := m.main

			// []textinput.Model{
			// 	m.nameInput,
			// 	m.emailInput,
			// 	m.passwordInput,
			// }

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			// if s == "enter" && m.cursor == len(inputs) {
			// 	return m, tea.Quit
			// }

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor > len(inputs)-1 {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(inputs) - 1
			}

			// for i := 0; i <= len(inputs)-1; i++ {
			// 	if i == m.cursor {
			// 		// Set focused state
			// 		inputs[i].nameInput.Focus()
			// 		inputs[i].nameInput.Prompt = focusedPrompt
			// 		inputs[i].nameInput.TextColor = focusedColor
			// 		continue
			// 	}
			// 	// Remove focused state
			// 	inputs[i].nameInput.Blur()
			// 	inputs[i].nameInput.Prompt = blurredPrompt
			// 	inputs[i].nameInput.TextColor = ""
			// }

			// m.nameInput = inputs[0]
			// m.emailInput = inputs[1]
			// m.passwordInput = inputs[2]

			// if m.cursor == len(inputs) {
			// 	m.submitButton = focusedSubmitButton
			// } else {
			// 	m.submitButton = blurredSubmitButton
			// }

			return m, nil
		}
	}

	// Handle character input and blinks
	// if m.mode != main {
	// 	m, cmd = updateExecInputs(msg, m)
	// }
	return m, cmd
}

func (m model) RegistUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		// case "enter":
		// 	m.mode = main

		// case "up":
		// 	m.cursor--

		// case "down":
		// 	m.cursor++

		case "tab", "shift+tab", "enter", "up", "down":

			inputs := []textinput.Model{
				m.regist.nameInput,
				m.regist.cmdInput,
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.cursor == len(inputs) {
				m.mode = Main
				// return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor > len(inputs) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.cursor {
					// Set focused state
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					inputs[i].TextColor = registColor
					continue
				}
				// Remove focused state
				inputs[i].Blur()
				inputs[i].Prompt = blurredPrompt
				inputs[i].TextColor = ""
			}

			m.regist.nameInput = inputs[0]
			m.regist.cmdInput = inputs[1]

			if m.cursor == len(inputs) {
				m.regist.doneButton = focusedRegisterButton
			} else {
				m.regist.doneButton = blurredRegisterButton
			}

			return m, nil
		}
	}

	m, cmd = updateRegistInputs(msg, m)
	return m, cmd
}

func updateRegistInputs(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.regist.nameInput, cmd = m.regist.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.regist.cmdInput, cmd = m.regist.cmdInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) searchUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		// case "enter":
		// 	m.mode = Main

		// case "up":
		// 	m.cursor--

		// case "down":
		// 	m.cursor++

		case "tab", "shift+tab", "enter", "up", "down":

			inputs := []textinput.Model{
				m.search.searchInput,
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.cursor == len(inputs) {
				m.mode = Main
				// return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor > len(inputs) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.cursor {
					// Set focused state
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					inputs[i].TextColor = searchColor
					continue
				}
				// Remove focused state
				inputs[i].Blur()
				inputs[i].Prompt = blurredPrompt
				inputs[i].TextColor = ""
			}

			m.search.searchInput = inputs[0]

			if m.cursor == len(inputs) {
				m.search.searchButton = focusedSearchButton
			} else {
				m.search.searchButton = blurredSearchButton
			}

			return m, nil
		}
	}

	m, cmd = updateSearchInputs(msg, m)
	return m, cmd
}

func updateSearchInputs(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.search.searchInput, cmd = m.search.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func horizontalBar(w int) string {
	s := strings.Builder{}
	for i := 0; i < w; i++ {
		s.WriteString("-")
	}
	return s.String()
}

func modeName(mode Mode) string {
	var s string
	switch mode {
	case Main:
		s = colorString("Main", mainColor)
	case Regist:
		s = colorString("Regist", registColor)
	case Search:
		s = colorString("Search", searchColor)
	}
	return s
}

func (m model) View() string {
	adjw := m.tsize.width
	s := strings.Builder{}
	// s.WriteString(displayMode(m.mode))

	s.WriteString(displayMode(m.mode))
	s.WriteString("+" + horizontalBar(adjw) + "+\n")
	// modes := fmt.Sprintf("+--- MODE: %8s ", modeName(m.mode))
	// s.WriteString(modes + horizontalBar(adjw-17) + "+\n")
	// displayMode(m.mode)
	switch m.mode {
	case Main:
		s.WriteString(m.mainView())
	case Regist:
		s.WriteString(m.RegistView())
	case Search:
		s.WriteString(m.searchView())
	}
	s.WriteString(colorString(fmt.Sprintf(" 'esc/ctrl+c/ctrl+q': Quit this app.\n"), unfocusedColor))

	return s.String()
}

func displayMode(mode Mode) string {
	s := strings.Builder{}
	s.WriteString("+----------------+\n")
	s.WriteString("| MODE:")
	switch mode {
	case Main:
		s.WriteString(colorString(" Main     ", mainColor))
	case Regist:
		s.WriteString(colorString(" Register ", registColor))
	case Search:
		s.WriteString(colorString(" Search   ", searchColor))
	}
	s.WriteString("|\n")
	// s.WriteString("+----------------+\n")

	return s.String()
}

func (m model) mainView() string {
	w, _, _ := term.GetSize(syscall.Stdin)
	adjw := w - 23
	// var adjw int = m.tsize.width
	s := strings.Builder{}
	for i := 0; i < len(m.main); i++ {

		if m.cursor == i {
			s.WriteString("| " + focusedPrompt)
		} else {
			s.WriteString("| " + blurredPrompt)
		}

		fmtStr := fmt.Sprintf("%2d. %s", i, m.main[i].name)
		if len(fmtStr) > adjw {
			fmtStr = fmtStr[:adjw-3] + "..."
		}
		if m.cursor == i {
			s.WriteString(colorString(fmtStr, mainColor))
		} else {
			s.WriteString(fmtStr)

		}
		gapSize := (adjw + 1) - len(fmtStr)
		s.WriteString(horizontalSpace(gapSize) + "|\n")
		// marginSize := (adjw + 2) - len(slicestr)
		// tmps.WriteString(horizontalSpace(marginSize) + "|\n")

		// s.WriteString(fmt.Sprintf("%v %v", len(fmtStr), adjw))
		// for i := 0; i < len(fmtStr); i++ {

		// }
		// s.WriteString(" |\n")
	}

	s.WriteString("|" + horizontalSpace(adjw+3) + "|\n")
	cmdstr := m.main[m.cursor].cmd
	if len(cmdstr) > adjw && adjw >= 0 {
		retCnt := len(cmdstr) / adjw
		tmps := strings.Builder{}
		for i := 0; i <= retCnt; i++ {
			si := adjw * i
			ei := adjw * (i + 1)
			// fmt.Printf("i:%d s:%d e:%d w:%d adjw:%d, retCnt:%d len:%d\n", i, s, e, w, adjw, retCnt, len(cmdstr))
			tmps.WriteString("| ")
			sliceStr := ""
			if ei > len(cmdstr) {
				sliceStr = cmdstr[si:]
			} else {
				sliceStr = cmdstr[si:ei]
			}
			tmps.WriteString(colorString(sliceStr, mainColor))
			gapSize := (adjw + 2) - len(sliceStr)
			tmps.WriteString(horizontalSpace(gapSize) + "|\n")
		}
		cmdstr = tmps.String()
	}

	s.WriteString(fmt.Sprintf("+-- select command (No.%2d) ", m.cursor) + horizontalBar(adjw-23) + "+\n")
	s.WriteString(cmdstr)
	s.WriteString("+" + horizontalBar(adjw+3) + "+\n")
	s.WriteString(colorString(fmt.Sprintf(" '↑/↓': Up/Down, 'enter': Execute command.\n"), unfocusedColor))
	s.WriteString(colorString(fmt.Sprintf(" 'r/e/s': Registration/Edit/Search mode.\n"), unfocusedColor))

	return s.String()
}

func horizontalSpace(w int) string {
	s := strings.Builder{}
	for i := 0; i < w; i++ {
		s.WriteString(" ")
	}
	return s.String()
}

func (m model) RegistView() string {
	w, _, _ := term.GetSize(syscall.Stdin)
	adjw := w - 20
	s := strings.Builder{}

	s.WriteString("\n")
	inputs := []string{
		m.regist.nameInput.View(),
		m.regist.cmdInput.View(),
	}

	for i := 0; i < len(inputs); i++ {
		s.WriteString(inputs[i])
		if i < len(inputs)-1 {
			s.WriteString("\n")
		}
	}

	s.WriteString("\n\n" + m.regist.doneButton + "\n")
	s.WriteString("+" + horizontalBar(adjw) + "+\n")

	return s.String()
}

func (m *model) searchView() string {
	w, _, _ := term.GetSize(syscall.Stdin)
	adjw := w - 20

	s := strings.Builder{}
	inputs := []string{
		m.search.searchInput.View(),
	}

	m.search.searchInput.Width = 10
	for i := 0; i < len(inputs); i++ {
		// if len(inputs[i]) > adjw {
		// 	m.search.searchInput.Width = adjw
		// }
		s.WriteString("| " + inputs[i])
		if i < len(inputs)-1 {
			s.WriteString("\n")
		}
	}
	s.WriteString("\n")
	s.WriteString("|\n")
	s.WriteString("| " + m.search.searchButton)
	gapSize := (adjw + 8) - len(m.search.searchButton)
	s.WriteString(horizontalSpace(gapSize) + "|\n")
	s.WriteString("+" + horizontalBar(adjw) + "+\n")
	s.WriteString(fmt.Sprintf("%v", m.search.searchInput.Width))

	return s.String()
}
