package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 60
const boxHeight = 20

// Styles för att rendera text i terminalen
// Dessa stilar används för att styla texten i terminalen
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

// Mode för att visa huvudmenyn
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

		// Om vi inte är i en övergång, rendera den aktuella vyn
		// beroende på vilket läge vi är i
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

// viewMainMenuWithOffset renderar huvudmenyn med en slide offset
func (m Model) viewMainMenuWithOffset() string {
	options := []string{"Render from example", "Some future option", "Quit"}
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

// Denna funktionen renderar en box med en border och centrerar den i terminalen
// och lägger till en slide offset
// och paddar eller trunkerar innehållet så att det passar i boxen.
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

	// Trimma eller padd höj
	content := slideAndBoxify(b.String(), boxWidth-4, boxHeight-4, m.slideOffset)

	// Rendera boxen med border
	box := borderStyle.Render(content)

	// Centrera boxen i terminalen och lägg till en slide offset
	// och lägg till en legend
	full := lipgloss.JoinVertical(lipgloss.Left, box, "", m.legend())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, full)

}

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

// slideAndBoxify tar en sträng och lägger till en slide offset
// och centrerar den i terminalen.
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
