package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pageton/dbview/internal/theme"
)

// Components provides reusable styled UI components.
type Components struct {
	cl theme.Colors
}

// NewComponents creates a new Components instance with the given theme.
func NewComponents(cl theme.Colors) Components {
	return Components{cl: cl}
}

// Header renders a styled header section.
func (c Components) Header(title string, subtitle string) string {
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(c.cl.HeaderBg)).
		Foreground(lipgloss.Color(c.cl.HeaderFg)).
		Bold(true).
		Padding(0, 1)

	subtitleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(c.cl.HeaderBg)).
		Foreground(lipgloss.Color(c.cl.Dim)).
		Padding(0, 1)

	if subtitle != "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(title),
			subtitleStyle.Render(subtitle),
		)
	}
	return titleStyle.Render(title)
}

// Card renders a card-like container with content.
func (c Components) Card(title string, content string) string {
	titleStyle := theme.SectionTitleStyle(c.cl)
	contentStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Background(lipgloss.Color(c.cl.Bg))

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.cl.Border)).
		Padding(0, 1)

	if title != "" {
		return card.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render(title),
				contentStyle.Render(content),
			))
	}
	return card.Render(contentStyle.Render(content))
}

// Section renders a section with a title and content.
func (c Components) Section(title string, content string) string {
	titleStyle := theme.SectionTitleStyle(c.cl)
	titleBar := titleStyle.Render(" " + title + " ")

	return lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		content,
	)
}

// Breadcrumb renders breadcrumb navigation.
func (c Components) Breadcrumb(items []string) string {
	for i := range items {
		if i > 0 {
			items[i] = "› " + items[i]
		}
	}
	breadcrumb := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Dim)).
		Render(strings.Join(items, " "))

	return lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color(c.cl.HeaderBg)).
		Render(breadcrumb)
}

// StatusBadge renders a status badge.
func (c Components) StatusBadge(text string, status string) string {
	var style lipgloss.Style
	switch status {
	case "success":
		style = lipgloss.NewStyle().
			Background(lipgloss.Color(c.cl.Ok)).
			Foreground(lipgloss.Color(c.cl.Bg))
	case "error":
		style = lipgloss.NewStyle().
			Background(lipgloss.Color(c.cl.Err)).
			Foreground(lipgloss.Color(c.cl.White))
	case "warning":
		style = lipgloss.NewStyle().
			Background(lipgloss.Color(c.cl.Warn)).
			Foreground(lipgloss.Color(c.cl.Bg))
	default:
		style = lipgloss.NewStyle().
			Background(lipgloss.Color(c.cl.AccentDim)).
			Foreground(lipgloss.Color(c.cl.White))
	}

	return style.Render(" " + text + " ")
}

// InfoBox renders an informational box.
func (c Components) InfoBox(title string, content string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Accent)).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.cl.Accent)).
		Padding(1, 2)

	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("ℹ "+title),
			content,
		))
}

// WarningBox renders a warning box.
func (c Components) WarningBox(title string, content string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Warn)).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.cl.Warn)).
		Padding(1, 2)

	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("⚠ "+title),
			content,
		))
}

// ErrorBox renders an error box.
func (c Components) ErrorBox(title string, content string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Err)).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(c.cl.Err)).
		Padding(1, 2)

	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("✗ "+title),
			content,
		))
}

// HelpBar renders a contextual help bar with key hints.
func (c Components) HelpBar(hints []string) string {
	helpStyle := theme.HelpStyle(c.cl)

	// Create sections of hints
	var sections [][]string
	if len(hints) > 0 {
		// Navigation section (first 3 hints)
		navEnd := 3
		if navEnd > len(hints) {
			navEnd = len(hints)
		}
		sections = append(sections, hints[0:navEnd])

		// Actions section (next 3 hints)
		actionEnd := 6
		if actionEnd > len(hints) {
			actionEnd = len(hints)
		}
		if actionEnd > 3 {
			sections = append(sections, hints[3:actionEnd])
		}

		// Other section (remaining hints)
		if len(hints) > 6 {
			sections = append(sections, hints[6:])
		}
	}

	var renderedSections []string
	for _, section := range sections {
		renderedSections = append(renderedSections, helpStyle.Render(" "+strings.Join(section, " • ")))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedSections...)
}

// LoadingSpinner renders a loading spinner with text.
func (c Components) LoadingSpinner(text string) string {
	spinnerFrame := theme.SpinnerFrames[0] // This will be updated by the model
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Accent)).
		Render(fmt.Sprintf("%s %s", spinnerFrame, text))
}

// EmptyState renders an empty state message.
func (c Components) EmptyState(message string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Dim)).
		Align(lipgloss.Center).
		Padding(2, 4).
		Render(message)
}

// Divider renders a horizontal divider.
func (c Components) Divider() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.cl.Border)).
		Render(strings.Repeat("─", 50))
}

// FocusIndicator renders a focus indicator for active elements.
func (c Components) FocusIndicator(text string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(c.cl.Accent)).
		Foreground(lipgloss.Color(c.cl.White)).
		Bold(true).
		Padding(0, 1).
		Render("▶ " + text)
}

// ShortText truncates text with ellipsis if too long.
func (c Components) ShortText(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max-3] + "..."
}
