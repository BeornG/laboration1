package tui

import (
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
				case 2:
					return m, tea.Quit
				}
			case modeFilePicker:
				renderDiagCmd(m.files[m.cursor])
				m.loading = true
				m.spinner, _ = m.spinner.Update(msg)
				m.output = "Rendering..."
			}

			// switch m.cursor { // 0 = Render from example, 1 = Some future option, 2 = Quit
			// case 0:
			// 	m.slidingIn = true
			// 	m.transitioning = true
			// 	m.slideOffset = 10
			// 	m.slideTarget = modeFilePicker
			// 	return m, slideTick()
			// case 1:
			// 	m.output = "🔧 This option is not yet implemented."
			// case 2:
			// 	return m, tea.Quit
			// }
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
		m.output = string(msg)

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

	}

	return m, tea.Batch(cmds...)
}

// Custom message for async rendering
type renderFinishedMsg string

func renderDiagCmd(filename string) tea.Cmd {
	path := path.Join("example", filename)
	renderDiagToSVG(path)
	return func() tea.Msg {
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
