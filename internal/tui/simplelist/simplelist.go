package simplelist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct {
	selectedPrefix string
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(d.selectedPrefix + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

func New(title string, names ...string) Model {
	items := make([]list.Item, len(names))
	for i, name := range names {
		items[i] = item(name)
	}

	l := list.New(items, itemDelegate{}, 18, 10)
	l.Title = title
	l.Styles.Title = titleStyle

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)

	l.DisableQuitKeybindings()

	return Model{lm: l}
}

type Model struct {
	lm    list.Model
	Focus bool
}

func (l Model) Init() tea.Cmd {
	return nil
}

func (l Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	l.lm, cmd = l.lm.Update(msg)
	return l, cmd
}

func (l Model) View() string {
	if l.Focus {
		l.lm.SetDelegate(itemDelegate{selectedPrefix: "> "})
	} else {
		l.lm.SetDelegate(itemDelegate{selectedPrefix: "  "})
	}
	return l.lm.View()
}
