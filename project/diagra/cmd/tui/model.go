package tui

import (
	"fmt"
	"path"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	modeMenu Mode = iota
	modeFilePicker
)

type Model struct {
	mode          Mode
	cursor        int
	files         []string
	output        string
	spinner       spinner.Model
	loading       bool
	width         int
	height        int
	transitioning bool
	slideOffset   int
	slideTarget   Mode
	slidingIn     bool
	slideOut      bool
	renderStart   time.Time
}

// InitialModel skapar en ny instans av Model med angivna diagFiles
func InitialModel(diagFiles []string) Model {
	return Model{
		files:   diagFiles,
		cursor:  0,
		output:  "",
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		loading: false,
	}
}

// init() initierar modellen och returnerar en cmd för att starta
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update hanterar meddelanden och uppdaterar modellen
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			switch m.mode {
			case modeMenu:
				switch m.cursor {
				case 0:
					m.slidingIn = true
					m.transitioning = true
					m.slideOffset = 10
					m.slideTarget = modeFilePicker
					return m, slideTick()
				case 1:
					m.output = "🔧 This option is not yet implemented."
					return m, clearOutputAfter(2 * time.Second)
				case 2:
					return m, tea.Quit
				}
			case modeFilePicker:

				m.output = "Rendering..."
				m.loading = true
				m.renderStart = time.Now()
				m.spinner = spinner.New(spinner.WithSpinner(spinner.Dot))
				cmd := renderDiagCmd(m.files[m.cursor])

				return m, tea.Batch(cmd, m.spinner.Tick)
			}

		case "esc":
			switch m.mode {
			case modeFilePicker:
				m.transitioning = true
				m.slideOut = true
				m.slideOffset = 0
				m.slideTarget = modeMenu
				return m, slideTick()

			case modeMenu:
				return m, tea.Quit
			}
		}

	case renderFinishedMsg:
		m.loading = false
		duration := time.Since(m.renderStart).Milliseconds()
		m.output = fmt.Sprintf("✅ Rendering finished in %dms", duration)
		return m, clearOutputAfter(2 * time.Second)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case slideMsg:
		if m.slidingIn && m.slideOffset > 0 {
			m.slideOffset--
			return m, slideTick()
		}
		if m.slidingIn {
			m.slidingIn = false
			m.transitioning = false
			m.mode = m.slideTarget
			m.cursor = 0
			m.slideOffset = 0
			return m, nil
		}

		// Slide out (file picker sliding out to right)
		if m.slideOut && m.slideOffset < 10 {
			m.slideOffset++
			return m, slideTick()
		}
		if m.slideOut {
			m.slideOut = false
			m.transitioning = false
			m.mode = m.slideTarget
			m.cursor = 0
			m.slideOffset = 0
			return m, nil
		}
	case clearOutputMsg:
		m.output = ""
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// Custom message for async rendering
type renderFinishedMsg string

func renderDiagCmd(filename string) tea.Cmd {

	return func() tea.Msg {
		path := path.Join("example", filename)
		renderDiagToSVG(path)
		return renderFinishedMsg("Rendering finished")
	}
}

// func fakeFadeTo(nextMode Mode) tea.Cmd {
// 	return func() tea.Msg {
// 		time.Sleep(500 * time.Millisecond)
// 		return modeSwitchMsg(nextMode)
// 	}
// }

// type modeSwitchMsg Mode

type slideMsg struct{}

func slideTick() tea.Cmd {
	return tea.Tick(40*time.Millisecond, func(t time.Time) tea.Msg {
		return slideMsg{}
	})
}

type clearOutputMsg struct{}

func clearOutputAfter(delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay)
		return clearOutputMsg{}
	}
}
