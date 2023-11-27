package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/koniferous22/dot-user-git-util/utils"
)

type MultiSelectModel struct {
	HeaderText        string
	Options           []string
	Cursor            int
	Selected          []bool
	ShouldDisplayHelp bool
	ShouldExit        bool
}

const MultiSelectHelpText = "Press\n" +
	"* 'arrow-up'/'arrow-down'/'j'/'k' for Navigation\n" +
	"* 'space' for selection\n" +
	"* 'h'/'t' to toggle visiblity of help\n\n"

func CreateMultiSelectModel(headerText string, options []string, preselections []bool) MultiSelectModel {
	return MultiSelectModel{
		HeaderText:        headerText,
		Options:           options,
		Selected:          preselections,
		ShouldDisplayHelp: true,
	}
}

func (m MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (m MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Options)-1 {
				m.Cursor++
			}
		case " ":
			val := m.Selected[m.Cursor]
			if val {
				m.Selected[m.Cursor] = false
			} else {
				m.Selected[m.Cursor] = true
			}
		case "h", "t":
			m.ShouldDisplayHelp = !m.ShouldDisplayHelp
			return m, nil
		case "q", "esc", "ctrl+c":
			m.ShouldExit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MultiSelectModel) View() string {
	s := fmt.Sprintf("%s\n", m.HeaderText)
	if m.ShouldDisplayHelp {
		s += MultiSelectHelpText
	}
	for i, option := range m.Options {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		checked := "[ ]"
		var val bool
		if val = m.Selected[i]; val {
			checked = fmt.Sprintf("%s[x]%s", utils.ColorGreen, utils.Reset)
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, checked, option)
	}
	s += "\nPress ENTER to submit, q/esc/ctrl+c to quit.\n"
	return s
}
