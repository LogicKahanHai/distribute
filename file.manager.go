package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type model struct {
	files    []os.DirEntry
	cursor   int
	selected map[string]struct{}
	cwd      string
	err      error
	final    []string
}

func joinFilePath(segs ...string) string {
	return filepath.Join(segs...)
}

func initialiseModel(config ConfigFile) model {
	m := model{
		err: nil,
	}
	cwd, err := os.Getwd()
	if err != nil {
		m.err = err
	}
	ent, err := os.ReadDir(cwd)
	if err != nil {
		m.err = err
	}

	m.files = ent
	m.selected = make(map[string]struct{})
	m.cwd = cwd

	for _, files := range config.Build.Files {
		m.selected[files] = struct{}{}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) HandleSelect() error {

	if m.files[m.cursor].IsDir() == true {
		ent, err := os.ReadDir(m.files[m.cursor].Name())
		if err != nil {
			return err
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		m.cwd = joinFilePath(cwd, m.files[m.cursor].Name())
		m.files = ent
		m.cursor = 0
	} else {
		_, ok := m.selected[joinFilePath(m.cwd, m.files[m.cursor].Name())]
		if ok {
			delete(m.selected, joinFilePath(m.cwd, m.files[m.cursor].Name()))
		} else {
			m.selected[joinFilePath(m.cwd, m.files[m.cursor].Name())] = struct{}{}
		}
	}
	return nil
}

func (m *model) HandleGoBackDir() error {
	prevDir := filepath.Dir(m.cwd)
	ent, err := os.ReadDir(prevDir)
	if err != nil {
		return err
	}
	m.files = ent
	m.cwd = prevDir
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case " ":
			err := m.HandleSelect()
			if err != nil {
				m.err = err
				return m, nil
			}

		case "h", "backspace", "left":
			err := m.HandleGoBackDir()
			if err != nil {
				m.err = err
				return m, nil
			}
		case "enter":
			m.final = []string{}
			for key := range m.selected {
				m.final = append(m.final, key)
			}
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) ShowError() string {
	var style = lip.NewStyle().
		Bold(true).
		Foreground(ErrorText).
		Align(lip.Center)

	s := style.Render(m.err.Error())

	return s
}

func (m model) ShowFooter() string {
	var style = lip.NewStyle().
		Foreground(FooterGrey)
	var s string

	if m.err != nil {
		s = style.Render("Press q to quit.")
	} else {
		s = style.Render("Press q to quit.\tj/k or up/down for movement.\th or backspace to go up a directory.\n")
	}

	return s

}

func (m model) ShowTitle() string {
	var style = lip.NewStyle().
		Bold(true).
		Background(TitleBackground).
		Foreground(TitleForeground).
		Align(lip.Center).
		PaddingLeft(5).
		PaddingRight(5)

	return style.Render("Distribute")
}

func (m model) ShowHighlightedText(s string) string {
	var style = lip.NewStyle().
		Bold(true).
		Foreground(TextHighlight)

	return style.Render(s)
}

func (m model) ShowSubTitle(s string) string {
	var style = lip.NewStyle().
		Bold(true).
		Background(SubtitleBackground).
		Foreground(SubtitleForeground).
		Align(lip.Center).
		PaddingLeft(5).
		PaddingRight(5)

	return style.Render(s)
}

func (m model) ShowCursor(s string) string {
	var style = lip.NewStyle().
		Foreground(Cursor)

	return style.Render(s)
}

func (m model) ShowSelected(s string) string {
	var style = lip.NewStyle().
		Foreground(ListSelected)

	return style.Render(s)
}

func (m model) View() string {
	sectionBreak := "\n\n"
	if m.err != nil {
		s := m.ShowTitle()
		s += m.ShowError()
		s += "\n\n"
		s += m.ShowFooter()
		return s
	}
	var style = lip.NewStyle().
		Foreground(GenericText)

	s := m.ShowTitle()
	s += sectionBreak
	s += m.ShowSubTitle("What files to select?")

	s += sectionBreak

	// Iterate over our choices
	for i, choice := range m.files {
		filename := choice.Name()

		// Is the cursor pointing at this choice?
		cursor := style.Render(" ") // no cursor
		if m.cursor == i {
			cursor = m.ShowCursor(">") // cursor!
		}

		// Is this choice selected?
		checked := style.Render("[ ]") // not selected
		if _, ok := m.selected[joinFilePath(m.cwd, m.files[i].Name())]; ok {
			checked = m.ShowSelected("[x]") // selected!
			filename = m.ShowSelected(filename)
		}

		// Render the row
		if choice.IsDir() {
			filename = m.ShowHighlightedText(filename)
			s += fmt.Sprintf("%s 📁 %s\n", cursor, filename)
		} else {
			s += fmt.Sprintf("%s %s %s\n", cursor, checked, filename)
		}
	}

	// The footer
	s += m.ShowFooter()

	// Send the UI for rendering
	return s
}

func SelectFiles(config ConfigFile) []string {

	p := tea.NewProgram(initialiseModel(config), tea.WithAltScreen())
	output, err := p.Run()
	if err != nil {
		os.Exit(1)
	}

	if m, ok := output.(model); ok {
		return m.final
	}
	return []string{}
}
