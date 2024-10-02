package timeseries

import (
	"time"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New(color lipgloss.Color) Model {
	// create new time series chart
	width := 140
	height := 10
	chart := tslc.New(width, height)

	// additional chart code goes here
	// add default data set
	dataSet := []float64{0, 2, 4, 6, 8, 10, 8, 6, 4, 2, 0}
	for i, v := range dataSet {
		date := time.Now().Add(time.Minute * time.Duration(i))
		chart.Push(tslc.TimePoint{date, v})
	}

	chart.SetStyle(lipgloss.NewStyle().Foreground(color))
	chart.XLabelFormatter = func(i int, v float64) string {
		t := time.Unix(int64(v), 0).Local()
		return t.Format("15:04:05")
	}

	return Model{chart: chart, color: color}
}

type Model struct {
	chart tslc.Model
	title string
	color lipgloss.Color
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.chart.Resize(msg.Width, msg.Height)
	}

	// forward Bubble Tea Msg to time series chart
	// and draw all data sets using braille runes
	m.chart, _ = m.chart.Update(msg)
	m.chart.DrawBrailleAll()
	return m, nil
}

func (m Model) View() string {
	return m.chart.View()
}
