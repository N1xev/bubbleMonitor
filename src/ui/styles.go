package ui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/N1xev/bubbleMonitor/src/config"
)

// ThemePalette holds colors for the application
type ThemePalette struct {
	Name       string
	Primary    compat.AdaptiveColor
	Secondary  compat.AdaptiveColor
	Success    compat.AdaptiveColor
	Warning    compat.AdaptiveColor
	Alert      compat.AdaptiveColor
	Text       compat.AdaptiveColor
	Muted      compat.AdaptiveColor
	Border     compat.AdaptiveColor
	Background compat.AdaptiveColor
}

// makeColor creates an AdaptiveColor with same value for Light/Dark
func makeColor(hex string) compat.AdaptiveColor {
	c := lipgloss.Color(hex)
	return compat.AdaptiveColor{Light: c, Dark: c}
}

// makeTTYColor creates an AdaptiveColor using ANSI color number (0-15)
func makeTTYColor(ansi int) compat.AdaptiveColor {
	c := lipgloss.ANSIColor(ansi)
	return compat.AdaptiveColor{Light: c, Dark: c}
}

// GetTheme returns the color palette for the given theme name
func GetTheme(name string) ThemePalette {
	switch name {
	case "light":
		return ThemePalette{
			Name:       "light",
			Primary:    makeColor("#2563EB"),
			Secondary:  makeColor("#7C3AED"),
			Success:    makeColor("#059669"),
			Warning:    makeColor("#D97706"),
			Alert:      makeColor("#DC2626"),
			Text:       makeColor("#1F2937"),
			Muted:      makeColor("#6B7280"),
			Border:     makeColor("#D1D5DB"),
			Background: makeColor("#FFFFFF"),
		}
	case "nord":
		return ThemePalette{
			Name:       "nord",
			Primary:    makeColor("#88C0D0"),
			Secondary:  makeColor("#B48EAD"),
			Success:    makeColor("#A3BE8C"),
			Warning:    makeColor("#EBCB8B"),
			Alert:      makeColor("#BF616A"),
			Text:       makeColor("#ECEFF4"),
			Muted:      makeColor("#D8DEE9"),
			Border:     makeColor("#4C566A"),
			Background: makeColor("#2E3440"),
		}
	case "dracula":
		return ThemePalette{
			Name:       "dracula",
			Primary:    makeColor("#BD93F9"),
			Secondary:  makeColor("#FF79C6"),
			Success:    makeColor("#50FA7B"),
			Warning:    makeColor("#FFB86C"),
			Alert:      makeColor("#FF5555"),
			Text:       makeColor("#F8F8F2"),
			Muted:      makeColor("#6272A4"),
			Border:     makeColor("#44475A"),
			Background: makeColor("#282A36"),
		}
	case "gruvbox":
		return ThemePalette{
			Name:       "gruvbox",
			Primary:    makeColor("#83A598"),
			Secondary:  makeColor("#D3869B"),
			Success:    makeColor("#B8BB26"),
			Warning:    makeColor("#FABD2F"),
			Alert:      makeColor("#FB4934"),
			Text:       makeColor("#EBDBB2"),
			Muted:      makeColor("#A89984"),
			Border:     makeColor("#504945"),
			Background: makeColor("#282828"),
		}
	case "rosepine":
		return ThemePalette{
			Name:       "rosepine",
			Primary:    makeColor("#C4A7E7"),
			Secondary:  makeColor("#EBBCBA"),
			Success:    makeColor("#9CCFD8"),
			Warning:    makeColor("#F6C177"),
			Alert:      makeColor("#EB6F92"),
			Text:       makeColor("#E0DEF4"),
			Muted:      makeColor("#908CAA"),
			Border:     makeColor("#403D52"),
			Background: makeColor("#191724"),
		}
	case "everforest":
		return ThemePalette{
			Name:       "everforest",
			Primary:    makeColor("#7FBBB3"),
			Secondary:  makeColor("#D699B6"),
			Success:    makeColor("#A7C080"),
			Warning:    makeColor("#DBBC7F"),
			Alert:      makeColor("#E67E80"),
			Text:       makeColor("#D3C6AA"),
			Muted:      makeColor("#859289"),
			Border:     makeColor("#3D484D"),
			Background: makeColor("#2D353B"),
		}
	case "nightowl":
		return ThemePalette{
			Name:       "nightowl",
			Primary:    makeColor("#82AAFF"),
			Secondary:  makeColor("#C792EA"),
			Success:    makeColor("#C5E478"),
			Warning:    makeColor("#FFD590"),
			Alert:      makeColor("#EF5350"),
			Text:       makeColor("#D6DEEB"),
			Muted:      makeColor("#5F7E97"),
			Border:     makeColor("#2A3F5F"),
			Background: makeColor("#011627"),
		}
	case "palenight":
		return ThemePalette{
			Name:       "palenight",
			Primary:    makeColor("#82AAFF"),
			Secondary:  makeColor("#C792EA"),
			Success:    makeColor("#C3E88D"),
			Warning:    makeColor("#FFCB6B"),
			Alert:      makeColor("#F07178"),
			Text:       makeColor("#BFCBDB"),
			Muted:      makeColor("#676E95"),
			Border:     makeColor("#3B4252"),
			Background: makeColor("#292D3E"),
		}
	case "material":
		return ThemePalette{
			Name:       "material",
			Primary:    makeColor("#89DDFF"),
			Secondary:  makeColor("#F78C6C"),
			Success:    makeColor("#C3E88D"),
			Warning:    makeColor("#FFCB6B"),
			Alert:      makeColor("#FF5370"),
			Text:       makeColor("#EEFFFF"),
			Muted:      makeColor("#546E7A"),
			Border:     makeColor("#2E3C43"),
			Background: makeColor("#263238"),
		}
	case "synthwave":
		return ThemePalette{
			Name:       "synthwave",
			Primary:    makeColor("#FF7EDB"),
			Secondary:  makeColor("#36F9F6"),
			Success:    makeColor("#72F1B8"),
			Warning:    makeColor("#FED800"),
			Alert:      makeColor("#FE4450"),
			Text:       makeColor("#FFFFFF"),
			Muted:      makeColor("#B6B1B1"),
			Border:     makeColor("#495495"),
			Background: makeColor("#262335"),
		}
	case "cobalt2":
		return ThemePalette{
			Name:       "cobalt2",
			Primary:    makeColor("#0088FF"),
			Secondary:  makeColor("#FF628C"),
			Success:    makeColor("#3AD900"),
			Warning:    makeColor("#FFC600"),
			Alert:      makeColor("#FF0000"),
			Text:       makeColor("#FFFFFF"),
			Muted:      makeColor("#8F9D6A"),
			Border:     makeColor("#0D3A58"),
			Background: makeColor("#193549"),
		}
	case "horizon":
		return ThemePalette{
			Name:       "horizon",
			Primary:    makeColor("#25B0BC"),
			Secondary:  makeColor("#E95678"),
			Success:    makeColor("#29D398"),
			Warning:    makeColor("#FAB795"),
			Alert:      makeColor("#EC6A88"),
			Text:       makeColor("#FDF0ED"),
			Muted:      makeColor("#6C6F93"),
			Border:     makeColor("#2E303E"),
			Background: makeColor("#1C1E26"),
		}
	case "oceanic":
		return ThemePalette{
			Name:       "oceanic",
			Primary:    makeColor("#6699CC"),
			Secondary:  makeColor("#C594C5"),
			Success:    makeColor("#99C794"),
			Warning:    makeColor("#FAC863"),
			Alert:      makeColor("#EC5f67"),
			Text:       makeColor("#CDD3DE"),
			Muted:      makeColor("#65737E"),
			Border:     makeColor("#343D46"),
			Background: makeColor("#1B2B34"),
		}
	case "palefire":
		return ThemePalette{
			Name:       "palefire",
			Primary:    makeColor("#95C4CE"),
			Secondary:  makeColor("#E1A3EE"),
			Success:    makeColor("#8DD4AA"),
			Warning:    makeColor("#E4C9AF"),
			Alert:      makeColor("#EA868F"),
			Text:       makeColor("#E2E0D7"),
			Muted:      makeColor("#7D7A68"),
			Border:     makeColor("#4A4845"),
			Background: makeColor("#242321"),
		}
	case "github":
		return ThemePalette{
			Name:       "github",
			Primary:    makeColor("#0969DA"),
			Secondary:  makeColor("#8250DF"),
			Success:    makeColor("#1A7F37"),
			Warning:    makeColor("#9A6700"),
			Alert:      makeColor("#D1242F"),
			Text:       makeColor("#1F2328"),
			Muted:      makeColor("#656D76"),
			Border:     makeColor("#D1D9E0"),
			Background: makeColor("#FFFFFF"),
		}
	case "moonlight":
		return ThemePalette{
			Name:       "moonlight",
			Primary:    makeColor("#82AAFF"),
			Secondary:  makeColor("#BAACFF"),
			Success:    makeColor("#C3E88D"),
			Warning:    makeColor("#FFCB8B"),
			Alert:      makeColor("#FF5572"),
			Text:       makeColor("#C8D3F5"),
			Muted:      makeColor("#7A88CF"),
			Border:     makeColor("#3E68D7"),
			Background: makeColor("#222436"),
		}
	case "shades":
		return ThemePalette{
			Name:       "shades",
			Primary:    makeColor("#00B4D8"),
			Secondary:  makeColor("#90E0EF"),
			Success:    makeColor("#06D6A0"),
			Warning:    makeColor("#FFD60A"),
			Alert:      makeColor("#EF476F"),
			Text:       makeColor("#F8F9FA"),
			Muted:      makeColor("#ADB5BD"),
			Border:     makeColor("#495057"),
			Background: makeColor("#212529"),
		}
	case "midnight":
		return ThemePalette{
			Name:       "midnight",
			Primary:    makeColor("#4D9DE0"),
			Secondary:  makeColor("#E15554"),
			Success:    makeColor("#7BC043"),
			Warning:    makeColor("#F1C40F"),
			Alert:      makeColor("#E74C3C"),
			Text:       makeColor("#ECF0F1"),
			Muted:      makeColor("#95A5A6"),
			Border:     makeColor("#34495E"),
			Background: makeColor("#0B132B"),
		}
	case "forest":
		return ThemePalette{
			Name:       "forest",
			Primary:    makeColor("#52B788"),
			Secondary:  makeColor("#95D5B2"),
			Success:    makeColor("#40916C"),
			Warning:    makeColor("#F4A261"),
			Alert:      makeColor("#E76F51"),
			Text:       makeColor("#E9F5DB"),
			Muted:      makeColor("#B7E4C7"),
			Border:     makeColor("#2D6A4F"),
			Background: makeColor("#1B4332"),
		}
	case "autumn":
		return ThemePalette{
			Name:       "autumn",
			Primary:    makeColor("#E07A5F"),
			Secondary:  makeColor("#F2CC8F"),
			Success:    makeColor("#81B29A"),
			Warning:    makeColor("#F4A261"),
			Alert:      makeColor("#CA6702"),
			Text:       makeColor("#3D405B"),
			Muted:      makeColor("#81B29A"),
			Border:     makeColor("#E07A5F"),
			Background: makeColor("#F4F1DE"),
		}
	case "cyberpunk":
		return ThemePalette{
			Name:       "cyberpunk",
			Primary:    makeColor("#00F0FF"),
			Secondary:  makeColor("#FF00FF"),
			Success:    makeColor("#00FF41"),
			Warning:    makeColor("#FFE400"),
			Alert:      makeColor("#FF0055"),
			Text:       makeColor("#FFFFFF"),
			Muted:      makeColor("#A0A0A0"),
			Border:     makeColor("#6600FF"),
			Background: makeColor("#0A0E27"),
		}
	case "sunset":
		return ThemePalette{
			Name:       "sunset",
			Primary:    makeColor("#FF6B6B"),
			Secondary:  makeColor("#FFE66D"),
			Success:    makeColor("#4ECDC4"),
			Warning:    makeColor("#F7B731"),
			Alert:      makeColor("#EE5A6F"),
			Text:       makeColor("#F8F9FA"),
			Muted:      makeColor("#CED4DA"),
			Border:     makeColor("#495057"),
			Background: makeColor("#2C3E50"),
		}
	case "ocean":
		return ThemePalette{
			Name:       "ocean",
			Primary:    makeColor("#0077B6"),
			Secondary:  makeColor("#00B4D8"),
			Success:    makeColor("#48CAE4"),
			Warning:    makeColor("#CAF0F8"),
			Alert:      makeColor("#E63946"),
			Text:       makeColor("#023E8A"),
			Muted:      makeColor("#90E0EF"),
			Border:     makeColor("#0096C7"),
			Background: makeColor("#F8F9FA"),
		}
	case "coffee":
		return ThemePalette{
			Name:       "coffee",
			Primary:    makeColor("#A97155"),
			Secondary:  makeColor("#D4A373"),
			Success:    makeColor("#8FA383"),
			Warning:    makeColor("#E2C799"),
			Alert:      makeColor("#B85042"),
			Text:       makeColor("#E8DCC4"),
			Muted:      makeColor("#9D8D8F"),
			Border:     makeColor("#6F5F5A"),
			Background: makeColor("#3E2723"),
		}
	case "solarized":
		return ThemePalette{
			Name:       "solarized",
			Primary:    makeColor("#268BD2"),
			Secondary:  makeColor("#6C71C4"),
			Success:    makeColor("#859900"),
			Warning:    makeColor("#B58900"),
			Alert:      makeColor("#DC322F"),
			Text:       makeColor("#839496"),
			Muted:      makeColor("#657B83"),
			Border:     makeColor("#073642"),
			Background: makeColor("#002B36"),
		}
	case "monokai":
		return ThemePalette{
			Name:       "monokai",
			Primary:    makeColor("#66D9EF"),
			Secondary:  makeColor("#AE81FF"),
			Success:    makeColor("#A6E22E"),
			Warning:    makeColor("#FD971F"),
			Alert:      makeColor("#F92672"),
			Text:       makeColor("#F8F8F2"),
			Muted:      makeColor("#75715E"),
			Border:     makeColor("#49483E"),
			Background: makeColor("#272822"),
		}
	case "catppuccin":
		return ThemePalette{
			Name:       "catppuccin",
			Primary:    makeColor("#89B4FA"),
			Secondary:  makeColor("#F5C2E7"),
			Success:    makeColor("#A6E3A1"),
			Warning:    makeColor("#F9E2AF"),
			Alert:      makeColor("#F38BA8"),
			Text:       makeColor("#CDD6F4"),
			Muted:      makeColor("#9399B2"),
			Border:     makeColor("#45475A"),
			Background: makeColor("#1E1E2E"),
		}
	case "tokyonight":
		return ThemePalette{
			Name:       "tokyonight",
			Primary:    makeColor("#7AA2F7"),
			Secondary:  makeColor("#BB9AF7"),
			Success:    makeColor("#9ECE6A"),
			Warning:    makeColor("#E0AF68"),
			Alert:      makeColor("#F7768E"),
			Text:       makeColor("#C0CAF5"),
			Muted:      makeColor("#565F89"),
			Border:     makeColor("#3B4261"),
			Background: makeColor("#1A1B26"),
		}
	case "onedark":
		return ThemePalette{
			Name:       "onedark",
			Primary:    makeColor("#61AFEF"),
			Secondary:  makeColor("#C678DD"),
			Success:    makeColor("#98C379"),
			Warning:    makeColor("#E5C07B"),
			Alert:      makeColor("#E06C75"),
			Text:       makeColor("#ABB2BF"),
			Muted:      makeColor("#5C6370"),
			Border:     makeColor("#3E4451"),
			Background: makeColor("#282C34"),
		}
	case "ayu":
		return ThemePalette{
			Name:       "ayu",
			Primary:    makeColor("#FF8F40"),
			Secondary:  makeColor("#5CCFEE"),
			Success:    makeColor("#B8CC52"),
			Warning:    makeColor("#F2AE49"),
			Alert:      makeColor("#FF3333"),
			Text:       makeColor("#CBCCC6"),
			Muted:      makeColor("#5C6773"),
			Border:     makeColor("#3B4254"),
			Background: makeColor("#1F2430"),
		}
	case "tty":
		// Uses native ANSI terminal colors (0-15)
		return ThemePalette{
			Name:       "tty",
			Primary:    makeTTYColor(14),
			Secondary:  makeTTYColor(13),
			Success:    makeTTYColor(10),
			Warning:    makeTTYColor(11),
			Alert:      makeTTYColor(9),
			Text:       makeTTYColor(15),
			Muted:      makeTTYColor(8),
			Border:     makeTTYColor(7),
			Background: makeTTYColor(0),
		}
	default: // dark
		return ThemePalette{
			Name:       "dark",
			Primary:    makeColor("#3B82F6"),
			Secondary:  makeColor("#8B5CF6"),
			Success:    makeColor("#10B981"),
			Warning:    makeColor("#F59E0B"),
			Alert:      makeColor("#EF4444"),
			Text:       makeColor("#F9FAFB"),
			Muted:      makeColor("#9CA3AF"),
			Border:     makeColor("#4B5563"),
			Background: makeColor("#111827"),
		}
	}
}

// GetThemeFromCustom creates a ThemePalette from CustomThemeConfig
func GetThemeFromCustom(custom *config.CustomThemeConfig) ThemePalette {
	if custom == nil {
		return GetTheme("dark")
	}
	return ThemePalette{
		Name:       "custom",
		Primary:    makeColor(custom.Primary),
		Secondary:  makeColor(custom.Secondary),
		Success:    makeColor(custom.Success),
		Warning:    makeColor(custom.Warning),
		Alert:      makeColor(custom.Alert),
		Text:       makeColor(custom.Text),
		Muted:      makeColor(custom.Muted),
		Border:     makeColor(custom.Border),
		Background: makeColor(custom.Background),
	}
}

// GetAppTheme helps resolve the theme from config
func GetAppTheme(theme string, custom *config.CustomThemeConfig) ThemePalette {
	if theme == "custom" {
		return GetThemeFromCustom(custom)
	}
	return GetTheme(theme)
}
