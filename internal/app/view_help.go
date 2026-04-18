package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderHelpNew() string {
	cl := m.c()

	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(cl.Accent)).
		Padding(1, 3).
		Width(m.width - 6).
		Background(lipgloss.Color(cl.Bg))

	var sections [][]string
	sections = append(sections, []string{
		"TABLES VIEW",
		"  ↑↓ / jk      Navigate tables",
		"  enter        Open table data",
		"  s            View schema",
		"  x            Drop table (confirm)",
		"  D            Database stats",
		"  F            Flush table",
		"  /            SQL query",
		"  r            Reload",
	})
	sections = append(sections, []string{
		"DATA VIEW",
		"  ←→ / hl      Select column",
		"  ↑↓           Scroll rows",
		"  enter        Row detail (full fields)",
		"  1-9          Sort by column N",
		"  e            Edit cell",
		"  x            Delete row",
		"  d            Duplicate row",
		"  a            Add row",
		"  I            Import CSV/JSON",
		"  E            Export",
		"  c            Copy cell",
		"  C            Copy row",
		"  [ ]          Previous/next page",
		"  { }          First/last page",
		"  ctrl+f       Live filter",
		"  s            View schema",
		"  /            SQL query",
		"  r            Reload",
	})
	sections = append(sections, []string{
		"QUERY VIEW",
		"  ↑/↓          Query history",
		"  enter        Execute query",
		"  esc          Back",
	})
	sections = append(sections, []string{
		"GLOBAL",
		"  T            Cycle theme",
		"  M            Toggle mouse",
		"  Q            Query log",
		"  ?            This help",
		"  q            Quit",
		"  esc          Go back",
		"  ctrl+c       Force quit",
	})

	var content strings.Builder

	// Header
	headerTitle := lipgloss.NewStyle().
		Background(lipgloss.Color(cl.HeaderBg)).
		Foreground(lipgloss.Color(cl.HeaderFg)).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf(" %s Viewer — Help ", m.viewerTitle()))
	content.WriteString(headerTitle)
	content.WriteString("\n\n")

	for _, section := range sections {
		sectionTitle := section[0]
		sectionContent := section[1:]

		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(cl.Accent)).
			Bold(true).
			Padding(0, 1)
		content.WriteString(titleStyle.Render("▌ " + sectionTitle))
		content.WriteString("\n")

		for _, line := range sectionContent {
			content.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color(cl.White)).
				Padding(0, 2).
				Render(line))
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(cl.Dim)).
		Render(" Press any key to close"))

	return box.Render(content.String())
}
