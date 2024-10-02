package main

import (
	"bubbletea-eds/internal/tui/simplelist"
	"bubbletea-eds/internal/tui/timeseries"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	boxer "github.com/treilik/bubbleboxer"
)

func main() {
	// layout-tree definition
	m := model{
		tui:            boxer.Boxer{},
		viewsChoice:    simplelist.New("views", "pods", "weights"),
		intervalChoice: simplelist.New("intervals", "5m", "15m", "30m", "1h", "3h"),
	}

	m.viewsChoice.Focus = true

	rowSizeFunc := func(_ boxer.Node, widthOrHeight int) []int {
		w := (widthOrHeight - 3) / 2
		return []int{3, w, w}
	}

	headerRow := boxer.CreateNoBorderNode()
	headerRow.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("topleft", NewStringer("   ", lipgloss.NewStyle()))),
		stripErr(m.tui.CreateLeaf("callers", NewStringer("callers", lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Align(lipgloss.Center)))),
		stripErr(m.tui.CreateLeaf("targets", NewStringer("targets", lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Align(lipgloss.Center)))),
	}
	headerRow.SizeFunc = rowSizeFunc

	color := lipgloss.Color("1")
	zoneARow := boxer.CreateNoBorderNode()
	zoneARow.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("zonea", NewStringer("a", lipgloss.NewStyle().Foreground(color)))),
		stripErr(m.tui.CreateLeaf("zonea-caller", timeseries.New(color))),
		stripErr(m.tui.CreateLeaf("zonea-target", timeseries.New(color))),
	}
	zoneARow.SizeFunc = rowSizeFunc

	color = lipgloss.Color("2")
	zoneBRow := boxer.CreateNoBorderNode()
	zoneBRow.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("zoneb", NewStringer("b", lipgloss.NewStyle().Foreground(color)))),
		stripErr(m.tui.CreateLeaf("zoneb-caller", timeseries.New(color))),
		stripErr(m.tui.CreateLeaf("zoneb-target", timeseries.New(color))),
	}
	zoneBRow.SizeFunc = rowSizeFunc

	color = lipgloss.Color("3")
	zoneCRow := boxer.CreateNoBorderNode()
	zoneCRow.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("zonec", NewStringer("c", lipgloss.NewStyle().Foreground(color)))),
		stripErr(m.tui.CreateLeaf("zonec-caller", timeseries.New(color))),
		stripErr(m.tui.CreateLeaf("zonec-target", timeseries.New(color))),
	}
	zoneCRow.SizeFunc = rowSizeFunc

	displayNode := boxer.CreateNoBorderNode()
	displayNode.VerticalStacked = true
	displayNode.Children = []boxer.Node{
		headerRow,
		zoneARow,
		zoneBRow,
		zoneCRow,
	}
	displayNode.SizeFunc = func(_ boxer.Node, widthOrHeight int) []int {
		h := (widthOrHeight - 1) / 3
		return []int{1, h, h, h}
	}

	screen := boxer.CreateNoBorderNode()
	screen.VerticalStacked = true
	screen.SizeFunc = func(_ boxer.Node, widthOrHeight int) []int {
		return []int{
			widthOrHeight - 1,
			1,
		}
	}
	screen.Children = []boxer.Node{
		{ // Main part
			SizeFunc: func(_ boxer.Node, widthOrHeight int) []int {
				return []int{
					15,
					widthOrHeight - 15,
				}
			},
			Children: []boxer.Node{
				{
					// left sidebar
					VerticalStacked: true,
					Children: []boxer.Node{
						stripErr(m.tui.CreateLeaf("viewsChoice", m.viewsChoice)),
						stripErr(m.tui.CreateLeaf("intervalChoice", m.intervalChoice)),
					},
				},
				displayNode,
			},
		},

		stripErr(m.tui.CreateLeaf("statusBar", NewStringer("status", lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("15"))))),
	}

	m.tui.LayoutTree = screen
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}

func stripErr(n boxer.Node, _ error) boxer.Node {
	return n
}

type model struct {
	tui            boxer.Boxer
	viewsChoice    simplelist.Model
	intervalChoice simplelist.Model
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		_ = m.tui.UpdateSize(msg)
	}
	return m, nil
}
func (m model) View() string {
	return m.tui.View()
}

func NewStringer(text string, style lipgloss.Style) tea.Model {
	return stringer{
		text:  text,
		style: style,
	}
}

type stringer struct {
	text  string
	style lipgloss.Style

	width int
}

func (s stringer) Init() tea.Cmd { return nil }
func (s stringer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	}
	return s, nil
}
func (s stringer) View() string {
	return s.style.Width(s.width).Render(s.text)
}
