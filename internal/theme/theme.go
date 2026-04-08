package theme

import "github.com/charmbracelet/lipgloss"

// Colors holds the color palette for a single theme.
// Each field is a hex color string (e.g. "#7D56F4").
type Colors struct {
	Accent    string
	AccentDim string
	Err       string
	Ok        string
	Warn      string
	Dim       string
	White     string
	Bg        string
	Border    string
	HeaderBg  string
	HeaderFg  string
}

// Themes maps theme names to their color palettes.
var Themes = map[string]Colors{
	"dark": {
		Accent: "#7D56F4", AccentDim: "#5A3EBF", Err: "#FF4444",
		Ok: "#44FF88", Warn: "#FFD700", Dim: "#888888",
		White: "#FFFFFF", Bg: "#1A1A2E", Border: "#7D56F4",
		HeaderBg: "#2D2A3E", HeaderFg: "#F8F8F2",
	},
	"light": {
		Accent: "#5B4BC4", AccentDim: "#8B7FD4", Err: "#CC0000",
		Ok: "#228B22", Warn: "#B8860B", Dim: "#666666",
		White: "#000000", Bg: "#F5F5F5", Border: "#5B4BC4",
		HeaderBg: "#FFFFFF", HeaderFg: "#333333",
	},
	"gruvbox-dark": {
		Accent: "#FE8019", AccentDim: "#D65D0E", Err: "#FB4934",
		Ok: "#B8BB26", Warn: "#FABD2F", Dim: "#665C54",
		White: "#EBDBB2", Bg: "#282828", Border: "#FE8019",
		HeaderBg: "#3C3836", HeaderFg: "#FBF1C7",
	},
	"gruvbox-light": {
		Accent: "#AF3A03", AccentDim: "#D65D0E", Err: "#9D0006",
		Ok: "#79740E", Warn: "#B57614", Dim: "#928374",
		White: "#3C3836", Bg: "#FBF1C7", Border: "#AF3A03",
		HeaderBg: "#F2E5BC", HeaderFg: "#3C3836",
	},
	"catppuccin-latte": {
		Accent: "#8839EF", AccentDim: "#7287FD", Err: "#D20F39",
		Ok: "#40A02B", Warn: "#DF8E1D", Dim: "#9CA0B0",
		White: "#4C4F69", Bg: "#EFF1F5", Border: "#8839EF",
		HeaderBg: "#E6E9EF", HeaderFg: "#4C4F69",
	},
	"catppuccin-frappe": {
		Accent: "#CA9EE6", AccentDim: "#BABBF1", Err: "#E78284",
		Ok: "#A6D189", Warn: "#E5C890", Dim: "#626880",
		White: "#C6D0F5", Bg: "#303446", Border: "#CA9EE6",
		HeaderBg: "#37414D", HeaderFg: "#C6D0F5",
	},
	"catppuccin-macchiato": {
		Accent: "#C6A0F6", AccentDim: "#B4BEFE", Err: "#ED8796",
		Ok: "#A6DA95", Warn: "#EED49F", Dim: "#5B6078",
		White: "#CAD3F8", Bg: "#24273A", Border: "#C6A0F6",
		HeaderBg: "#302D41", HeaderFg: "#CAD3F8",
	},
	"catppuccin-mocha": {
		Accent: "#CBA6F7", AccentDim: "#B4BEFE", Err: "#F38BA8",
		Ok: "#A6E3A1", Warn: "#F9E2AF", Dim: "#585B70",
		White: "#CDD6F4", Bg: "#1E1E2E", Border: "#CBA6F7",
		HeaderBg: "#181825", HeaderFg: "#CDD6F4",
	},
}

// ThemeOrder defines the cycle order when pressing T to switch themes.
var ThemeOrder = []string{
	"dark", "light",
	"gruvbox-dark", "gruvbox-light",
	"catppuccin-latte", "catppuccin-frappe", "catppuccin-macchiato", "catppuccin-mocha",
}

// SpinnerFrames are the braille spinner animation frames.
var SpinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// Resolve returns the color palette for the given theme name,
// falling back to "dark" if the theme is unknown.
func Resolve(name string) Colors {
	if cl, ok := Themes[name]; ok {
		return cl
	}
	return Themes["dark"]
}

// --- Style helpers ---

// Styled renders s with the given foreground color.
func Styled(s, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(s)
}

// StyledBold renders s bold with the given foreground color.
func StyledBold(s, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).Render(s)
}

// HelpStyle returns a dim-foreground style.
func HelpStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cl.Dim))
}

// WarnStyle returns a warn-foreground style.
func WarnStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cl.Warn))
}

// ErrStyle returns an error-foreground style.
func ErrStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cl.Err))
}

// OkStyle returns an ok-foreground style.
func OkStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cl.Ok))
}

// DimStyle returns a dim-foreground style.
func DimStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(cl.Dim))
}

// BoxError returns a rounded-border box styled with the error color.
func BoxError(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cl.Err)).
		Padding(1, 3)
}

// BoxAccent returns a rounded-border box styled with the accent color.
func BoxAccent(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cl.Accent)).
		Padding(1, 3)
}

// BoxOk returns a rounded-border box styled with the ok color.
func BoxOk(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cl.Ok)).
		Padding(1, 3)
}

// BorderedTable returns a rounded-border style for table wrapping,
// using the theme's border color (fixes the previous bug where only
// the "light" theme got colored borders).
func BorderedTable(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cl.Border))
}

// HeaderStyle returns a style for section headers.
func HeaderStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(cl.HeaderBg)).
		Foreground(lipgloss.Color(cl.HeaderFg)).
		Bold(true)
}

// CardStyle returns a card-like container with borders.
func CardStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(cl.Border)).
		Padding(0, 1)
}

// SectionTitleStyle returns a style for section titles.
func SectionTitleStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(cl.Accent)).
		Bold(true)
}

// FocusStyle returns a style for focused elements.
func FocusStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(cl.Accent)).
		Foreground(lipgloss.Color(cl.White)).
		Bold(true)
}

// SubtleStyle returns a subtle, dimmed style for secondary elements.
func SubtleStyle(cl Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(cl.Dim))
}
