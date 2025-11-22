package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var blockFont = map[rune][]string{
	'A': {" ### ", "#   #", "#####", "#   #", "#   #"},
	'B': {"#### ", "#   #", "#### ", "#   #", "#### "},
	'C': {" ####", "#    ", "#    ", "#    ", " ####"},
	'D': {"#### ", "#   #", "#   #", "#   #", "#### "},
	'E': {"#####", "#    ", "###  ", "#    ", "#####"},
	'F': {"#####", "#    ", "###  ", "#    ", "#    "},
	'G': {" ### ", "#    ", "#  ##", "#   #", " ### "},
	'H': {"#   #", "#   #", "#####", "#   #", "#   #"},
	'I': {" ### ", "  #  ", "  #  ", "  #  ", " ### "},
	'J': {"   ##", "    #", "    #", "#   #", " ### "},
	'K': {"#   #", "#  # ", "###  ", "#  # ", "#   #"},
	'L': {"#    ", "#    ", "#    ", "#    ", "#####"},
	'M': {"#   #", "## ##", "# # #", "#   #", "#   #"},
	'N': {"#   #", "##  #", "# # #", "#  ##", "#   #"},
	'O': {" ### ", "#   #", "#   #", "#   #", " ### "},
	'P': {"#### ", "#   #", "#### ", "#    ", "#    "},
	'Q': {" ### ", "#   #", "#   #", "#  ##", " ####"},
	'R': {"#### ", "#   #", "#### ", "#  # ", "#   #"},
	'S': {" ####", "#    ", " ### ", "    #", "#### "},
	'T': {"#####", "  #  ", "  #  ", "  #  ", "  #  "},
	'U': {"#   #", "#   #", "#   #", "#   #", " ### "},
	'V': {"#   #", "#   #", "#   #", " # # ", "  #  "},
	'W': {"#   #", "#   #", "# # #", "## ##", "#   #"},
	'X': {"#   #", " # # ", "  #  ", " # # ", "#   #"},
	'Y': {"#   #", " # # ", "  #  ", "  #  ", "  #  "},
	'Z': {"#####", "   # ", "  #  ", " #   ", "#####"},
	'0': {" ### ", "#  ##", "# # #", "##  #", " ### "},
	'1': {"  #  ", " ##  ", "  #  ", "  #  ", " ### "},
	'2': {" ### ", "#   #", "   # ", "  #  ", "#####"},
	'3': {" ### ", "#   #", "   # ", "#   #", " ### "},
	'4': {"#  # ", "#  # ", "#####", "   # ", "   # "},
	'5': {"#####", "#    ", "#### ", "    #", "#### "},
	'6': {" ### ", "#    ", "#### ", "#   #", " ### "},
	'7': {"#####", "   # ", "  #  ", " #   ", " #   "},
	'8': {" ### ", "#   #", " ### ", "#   #", " ### "},
	'9': {" ### ", "#   #", " ####", "    #", " ### "},
	' ': {"     ", "     ", "     ", "     ", "     "},
}

func renderASCII(text string, color lipgloss.Color) string {
	text = strings.ToUpper(text)
	lines := make([]string, 5)

	for _, ch := range text {
		pat, ok := blockFont[ch]
		if !ok {
			pat = blockFont[' ']
		}
		for i := 0; i < 5; i++ {
			line := strings.ReplaceAll(pat[i], "#", "█")
			lines[i] += " " + line
		}
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	return style.Render(strings.Join(lines, "\n"))
}

type state int

const (
	pickColor state = iota
	enterText
	showOutput
)

type model struct {
	step      state
	colors    []string
	cursor    int
	colorCode string
	text      string
	output    string
	blink     bool
}

func initialModel() model {
	return model{
		step:   pickColor,
		colors: []string{"Red", "Green", "Blue", "Yellow", "Pink", "Cyan", "White"},
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return "blink"
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case string:
		if msg == "blink" {
			m.blink = !m.blink
			return m, tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg { return "blink" })
		}

	case tea.KeyMsg:
		if m.step == pickColor {
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.colors)-1 {
					m.cursor++
				}
			case "enter":
				switch m.colors[m.cursor] {
				case "Red":
					m.colorCode = "196"
				case "Green":
					m.colorCode = "46"
				case "Blue":
					m.colorCode = "27"
				case "Yellow":
					m.colorCode = "220"
				case "Pink":
					m.colorCode = "205"
				case "Cyan":
					m.colorCode = "51"
				case "White":
					m.colorCode = "15"
				}

				m.step = enterText
			case "q":
				return m, tea.Quit
			}
			return m, nil
		}
		if m.step == enterText {
			switch msg.String() {

			case "enter":
				m.output = renderASCII(m.text, lipgloss.Color(m.colorCode))
				m.step = showOutput

			case "backspace":
				if len(m.text) > 0 {
					m.text = m.text[:len(m.text)-1]
				}

			case "q":
				return m, tea.Quit

			default:
				if len(msg.Runes) > 0 {
					m.text += string(msg.Runes)
				}
			}
			return m, nil
		}
		if m.step == showOutput && msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.step == pickColor {
		header := lipgloss.NewStyle().
			Foreground(lipgloss.Color("45")).
			Bold(true).
			Render("🎨 Choose a color")

		list := ""
		for i, c := range m.colors {
			colorCode := "15" // default white

			switch c {
			case "Red":
				colorCode = "196"
			case "Green":
				colorCode = "46"
			case "Blue":
				colorCode = "27"
			case "Yellow":
				colorCode = "220"
			case "Pink":
				colorCode = "205"
			case "Cyan":
				colorCode = "51"
			case "White":
				colorCode = "15"
			}

			if i == m.cursor {
				list += lipgloss.NewStyle().
					Foreground(lipgloss.Color(colorCode)).
					Bold(true).
					Render("→ "+c) + "\n"
			} else {
				list += "  " + c + "\n"
			}
		}

		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("60")).
			Render("[ ↑↓ select ] [ Enter confirm ] [ q quit ]")

		return "\n" + header + "\n\n" + list + "\n" + footer + "\n"
	}

	if m.step == enterText {
		cursor := " "
		if m.blink {
			cursor = "█"
		}

		header := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.colorCode)).
			Bold(true).
			Render("💬 Type your text")

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(m.colorCode)).
			Padding(1, 2)

		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("60")).
			Render("[ Enter print ] [ q quit ]")

		return "\n" + header + "\n" +
			box.Render(m.text+cursor) +
			"\n\n" + footer + "\n"
	}
	if m.step == showOutput {

		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("60")).
			Render("[ q quit ]")

		return "\n\n" + m.output + "\n\n" + footer + "\n"
	}

	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
