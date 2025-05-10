package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 60
const boxHeight = 20

// Styles for the TUI
// These styles are used to render the TUI elements
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Bold(true).
			Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFA3")).
			Bold(true)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#007AFF")).
				Bold(true).
				Padding(0, 1)
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FFA3")).
			Padding(1, 2)
	legendStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			Padding(0, 2) // Padding för legend
)

// View renders the model to a string.
// It handles the sliding transition and centers the content in the terminal.
// It returns the rendered string to be displayed in the terminal.
func (m Model) View() string {

	var main string

	if m.transitioning {
		switch {
		case m.slidingIn:
			switch m.slideTarget {
			case modeFilePicker:
				main = m.viewFilePickerWithOffset()
			}
		case m.slideOut:
			switch m.slideTarget {
			case modeMenu:
				main = m.viewFilePickerWithOffset() // Slide *out* of file picker
			}
		}
	} else {
		switch m.mode {
		case modeMenu:
			main = m.viewMainMenuWithOffset()
		case modeFilePicker:
			main = m.viewFilePickerWithOffset()
		default:
			main = "Unknown mode"
		}
	}
	return main
}

// viewMainMenuWithOffset renders the main menu with a slide offset
// and centers it in the terminal.
// It shows the list of options and allows the user to select one.
func (m Model) viewMainMenuWithOffset() string {
	options := []string{"Render from example", "Render all diagrams", "Quit"}
	var b strings.Builder

	b.WriteString(titleStyle.Render("📊 Diagra") + "\n\n")
	b.WriteString("Select Mode:\n\n")

	for i, option := range options {
		if m.cursor == i {
			// b.WriteString(cursorStyle.Render("➜ ") + selectedItemStyle.Render(option) + "\n")
			line := selectedItemStyle.Render("➜ " + option)
			b.WriteString(line + "\n")
		} else {
			line := "  " + itemStyle.Render(option)
			b.WriteString(line + "\n")
		}
	}

	lines := strings.Split(b.String(), "\n")
	prefix := strings.Repeat(" ", m.slideOffset)
	for i := range lines {
		lines[i] = prefix + lines[i]
	}

	status := " " + m.output
	b.WriteString("\n" + itemStyle.Render(status))

	content := slideAndBoxify(b.String(), boxWidth-4, boxHeight-4, m.slideOffset)

	full := lipgloss.JoinVertical(lipgloss.Left,
		borderStyle.Render(content),
		"",
		m.legend(),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, full)
}

// viewFilePickerWithOffset renders the file picker with a slide offset
// and centers it in the terminal.
// It shows the list of files and allows the user to select one.
func (m Model) viewFilePickerWithOffset() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📂 Choose a Diagram") + "\n\n")

	for i, f := range m.files {
		if m.cursor == i {
			b.WriteString(cursorStyle.Render("➜ ") + selectedItemStyle.Render(f) + "\n")
		} else {
			b.WriteString("  " + itemStyle.Render(f) + "\n")
		}
	}

	// Status line (always added)
	status := m.spinner.View() + " " + m.output
	b.WriteString("\n" + itemStyle.Render(status))

	content := slideAndBoxify(b.String(), boxWidth-4, boxHeight-4, m.slideOffset)

	box := borderStyle.Render(content)

	full := lipgloss.JoinVertical(lipgloss.Left, box, "", m.legend())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, full)

}

// legend returns the legend for the current mode.
// It shows the available key bindings for navigation and actions.
// The legend is displayed at the bottom of the screen.
func (m Model) legend() string {

	switch m.mode {
	case modeFilePicker:
		return legendStyle.Render("↑/↓ or j/k: Navigate|Enter: Select|Esc: Back|q: Quit")
	case modeMenu:
		return legendStyle.Render("↑/↓ or j/k: Navigate|Enter: Choose|Esc/q: Quit")
	default:
		return ""
	}
}

// slideAndBoxify adds a slide effect to the content and boxes it.
// It takes the content, content width, content height, and slide offset as parameters.
// It returns the formatted content as a string.
func slideAndBoxify(content string, contentWidth, contentHeight, slideOffset int) string {
	lines := strings.Split(content, "\n")

	prefix := strings.Repeat(" ", slideOffset)
	for i := range lines {
		lines[i] = prefix + lines[i]
	}

	for len(lines) < contentHeight {
		lines = append(lines, "")
	}
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	for i, line := range lines {
		// r := []rune(line)
		// if len(r) > contentWidth {
		// 	lines[i] = string(r[:contentWidth])
		// } else {
		// 	lines[i] = line + strings.Repeat(" ", contentWidth-len(r))
		// }
		lines[i] = padRight(line, contentWidth)
	}

	// Säkerställ att varje rad är exakt contentWidth
	for i, line := range lines {
		lines[i] = padRight(line, contentWidth)
	}

	return strings.Join(lines, "\n")
}

func padRight(s string, width int) string {
	r := []rune(s)
	if len(r) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(r))
}
