package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	// lip "github.com/charmbracelet/lipgloss"
)

type model struct {
	files    []os.DirEntry
	cursor   int
	selected map[string]struct{}
	cwd      string
	err      error
	final    []string
}

func initialiseModel() model {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ent, err := os.ReadDir(cwd)
	if err != nil {
		panic(err)
	}

	return model{
		files:    ent,
		selected: make(map[string]struct{}),
		cwd:      cwd,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
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
			if m.files[m.cursor].IsDir() != true {
				_, ok := m.selected[m.cwd+"/"+m.files[m.cursor].Name()]
				if ok {
					delete(m.selected, m.cwd+"/"+m.files[m.cursor].Name())
				} else {
					m.selected[m.cwd+"/"+m.files[m.cursor].Name()] = struct{}{}
				}
			} else {
				ent, err := os.ReadDir(m.files[m.cursor].Name())
				if err != nil {
					panic(err)
				}
				cwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				m.cwd = cwd + "/" + m.files[m.cursor].Name()
				m.files = ent
				m.cursor = 0
			}
		case "h", "backspace", "left":
			prevDir := filepath.Dir(m.cwd)
			ent, err := os.ReadDir(prevDir)
			if err != nil {
				panic(err)
			}
			m.files = ent
			m.cwd = prevDir
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

func (m model) View() string {
	s := "What files to select?\n\n"

	// Iterate over our choices
	for i, choice := range m.files {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[m.cwd+"/"+m.files[i].Name()]; ok {
			checked = "x" // selected!
		}

		// Render the row
		if choice.IsDir() {
			s += fmt.Sprintf("%s 📁 %s\n", cursor, choice)
		} else {
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	// The footer
	s += "\nPress q to quit.\tj/k or up/down for movement.\th or backspace to go up a directory.\n"

	if m.err != nil {
		s += fmt.Sprintf("\n%v\n", m.err)
	}

	// Send the UI for rendering
	return s
}

func output() {

	p := tea.NewProgram(initialiseModel())
	output, err := p.Run()
	if err != nil {
		os.Exit(1)
	}

	if m, ok := output.(model); ok {
		fmt.Printf("Final Output: %s\n", m.final)
	}
}
