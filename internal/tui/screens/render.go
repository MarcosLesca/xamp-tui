package screens

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Theme define los colores del estilo Catppuccin.
type Theme struct {
	// Base colors
	Base    lipgloss.Color
	Surface lipgloss.Color
	Overlay lipgloss.Color
	Text    lipgloss.Color
	Subtext lipgloss.Color

	// Accent colors
	Blue   lipgloss.Color
	Mauve  lipgloss.Color
	Green  lipgloss.Color
	Yellow lipgloss.Color
	Sky    lipgloss.Color
	Peach  lipgloss.Color
	Red    lipgloss.Color
	Pink   lipgloss.Color
}

// newTheme crea el tema basado en el fondo del terminal.
func newTheme() Theme {
	if termenv.HasDarkBackground() {
		return newCatppuccinMocha()
	}
	return newCatppuccinLatte()
}

func newCatppuccinMocha() Theme {
	return Theme{
		Base:    lipgloss.Color("#1e1e2e"),
		Surface: lipgloss.Color("#313244"),
		Overlay: lipgloss.Color("#45475a"),
		Text:    lipgloss.Color("#cdd6f4"),
		Subtext: lipgloss.Color("#a6adc8"),
		Blue:    lipgloss.Color("#89b4fa"),
		Mauve:   lipgloss.Color("#cba6f7"),
		Green:   lipgloss.Color("#a6e3a1"),
		Yellow:  lipgloss.Color("#f9e2af"),
		Sky:     lipgloss.Color("#89dceb"),
		Peach:   lipgloss.Color("#fab387"),
		Red:     lipgloss.Color("#f38ba8"),
		Pink:    lipgloss.Color("#f5c2e7"),
	}
}

func newCatppuccinLatte() Theme {
	return Theme{
		Base:    lipgloss.Color("#eff1f5"),
		Surface: lipgloss.Color("#e6e9ef"),
		Overlay: lipgloss.Color("#ccd0da"),
		Text:    lipgloss.Color("#4c4f69"),
		Subtext: lipgloss.Color("#6c6f85"),
		Blue:    lipgloss.Color("#04a5e5"),
		Mauve:   lipgloss.Color("#8839ef"),
		Green:   lipgloss.Color("#40a02b"),
		Yellow:  lipgloss.Color("#df8e1d"),
		Sky:     lipgloss.Color("#209fb5"),
		Peach:   lipgloss.Color("#fe640b"),
		Red:     lipgloss.Color("#d20f39"),
		Pink:    lipgloss.Color("#ea76cb"),
	}
}

// Styles define los estilos visuales de la TUI.
type Styles struct {
	Theme *Theme
}

// NewStyles crea los estilos con tema Catppuccin.
func NewStyles() *Styles {
	return &Styles{
		Theme: func() *Theme {
			t := newTheme()
			return &t
		}(),
	}
}

// BaseLayout devuelve el layout base al estilo gentle-ai/career-ops (full width).
func BaseLayout(styles *Styles, width int, content string, hints []string) string {
	t := styles.Theme

	if width < 80 {
		width = 80
	}

	// === HEADER ===
	headerStyle := lipgloss.NewStyle().
		Foreground(t.Blue).
		Background(t.Base).
		Bold(true).
		Width(width).
		Padding(0, 1)

	// === CONTENT (sin centrado, full width) ===
	contentRendered := content

	// === FOOTER ===
	keyStyle := lipgloss.NewStyle().Foreground(t.Green).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(t.Subtext)
	brandStyle := lipgloss.NewStyle().Foreground(t.Overlay)

	var keyParts []string
	for i := 0; i < len(hints); i += 2 {
		key := hints[i]
		desc := ""
		if i+1 < len(hints) {
			desc = hints[i+1]
		}
		if desc != "" {
			keyParts = append(keyParts, keyStyle.Render(key)+descStyle.Render(" "+desc))
		} else {
			keyParts = append(keyParts, keyStyle.Render(key))
		}
	}
	keysRendered := strings.Join(keyParts, descStyle.Render("  "))

	brand := brandStyle.Render("xampp-tui")

	gap := width - lipgloss.Width(keysRendered) - lipgloss.Width(brand) - 4
	if gap < 1 {
		gap = 1
	}

	footer := lipgloss.NewStyle().
		Foreground(t.Subtext).
		Background(t.Surface).
		Width(width).
		Padding(0, 1).
		Render(keysRendered + strings.Repeat(" ", gap) + brand)

	border := lipgloss.NewStyle().
		Foreground(t.Overlay).
		Width(width).
		Render(strings.Repeat("─", width))

	return headerStyle.Render(" ✦ XAMPP-TUI ") + "\n" +
		border + "\n" +
		contentRendered + "\n" +
		border + "\n" +
		footer
}

// JoinVertical une líneas verticalmente.
func JoinVertical(lines []string) string {
	return strings.Join(lines, "\n")
}