package prompts

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/koniferous22/dot-user-git-util/utils"
)

type YesNoModel struct {
	Question          string
	InvalidInput      bool
	Entered           bool
	Result            bool
	ShouldExit        bool
	ShouldDisplayHelp bool
	ShouldQuitOnNo    bool
}

const YesNoHelpText = "Press\n" +
	"* 'y' for Yes\n" +
	"* 'n' for No\n" +
	"* 'q'/'esc'/'ctrl+c' to Quit\n" +
	"* 'h'/'t' to toggle visiblity of help\n"

func CreateYesNoModel(question string, shouldQuitOnNo bool) YesNoModel {
	return YesNoModel{
		Question:          question,
		ShouldDisplayHelp: true,
		ShouldQuitOnNo:    shouldQuitOnNo,
	}
}

func (m YesNoModel) Init() tea.Cmd {
	return nil
}

func (m YesNoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.InvalidInput = false
			m.Result = true
			m.Entered = true
			return m, tea.Quit
		case "n", "N":
			m.InvalidInput = false
			m.Result = false
			m.Entered = true
			m.ShouldExit = m.ShouldQuitOnNo
			return m, tea.Quit
		case "h", "t":
			m.ShouldDisplayHelp = !m.ShouldDisplayHelp
			return m, nil
		case "q", "esc", "ctrl+c":
			m.ShouldExit = true
			return m, tea.Quit
		default:
			m.InvalidInput = true
			return m, nil
		}
	}
	return m, nil
}
func (m YesNoModel) View() string {
	s := fmt.Sprintf("%s%s%s\n", utils.FontBold, m.Question, utils.Reset)
	if m.ShouldDisplayHelp {
		s += YesNoHelpText
	}

	if m.InvalidInput {
		s += fmt.Sprintf("%sInvalid Input%s\n", utils.ColorRed, utils.Reset)
	}
	if m.Entered {
		if m.Result {
			s += fmt.Sprintf("%sEntered Yes%s\n", utils.ColorGreen, utils.Reset)
		} else {
			s += fmt.Sprintf("%sEntered No%s\n", utils.ColorRed, utils.Reset)
		}
	}
	return s
}
