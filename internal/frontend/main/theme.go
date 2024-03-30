package yaac_frontend_main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type yaacTheme struct {
	bg      color.RGBA
	accent  color.RGBA
	warn    color.RGBA
	page_bg color.RGBA
	elem    color.RGBA
	hover   color.RGBA
	green   color.RGBA
	yellow  color.RGBA
	text    color.RGBA
	white   color.RGBA
}

// Assert that the Theme implements the theme interface
var _ fyne.Theme = (*yaacTheme)(nil)
var ytheme fyne.Theme = yaacTheme{
	bg:      color.RGBA{230, 233, 235, 127},
	accent:  color.RGBA{227, 0, 27, 255},
	warn:    color.RGBA{227, 0, 27, 255},
	page_bg: color.RGBA{230, 233, 235, 50},
	elem:    color.RGBA{255, 255, 255, 0},
	hover:   color.RGBA{209, 209, 209, 255},
	green:   color.RGBA{51, 255, 0, 255},
	yellow:  color.RGBA{255, 229, 0, 255},
	text:    color.RGBA{0, 0, 0, 255},
	white:   color.RGBA{243, 243, 243, 255},
}

func (m yaacTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		// Handle background color
		return m.bg
	case theme.ColorNameButton:
		// Handle button color
		return m.white
	case theme.ColorNameDisabledButton:
		// Handle disabled button color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameDisabled:
		// Handle disabled foreground color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameError:
		// Handle error foreground color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameFocus:
		// Handle focus color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameForeground:
		// Handle foreground color
		return m.text
	case theme.ColorNameHeaderBackground:
		// Handle header background color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameHover:
		// Handle hover color
		return m.hover
	case theme.ColorNameHyperlink:
		// Handle hyperlink color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameInputBackground:
		// Handle input field background color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameInputBorder:
		// Handle input field border color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameMenuBackground:
		// Handle menu background color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameOverlayBackground:
		// Handle overlay background color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNamePlaceHolder:
		// Handle placeholder text color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNamePressed:
		// Handle tap overlay color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNamePrimary:
		// Handle primary color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameScrollBar:
		// Handle scrollbar color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameSelection:
		// Handle selection color
		return m.warn
	case theme.ColorNameSeparator:
		// Handle separator bar color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameShadow:
		// Handle shadow color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameSuccess:
		// Handle success foreground color
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	case theme.ColorNameWarning:
		// Handle warning foreground color
		return m.warn
	default:
		// Handle default case (if needed)
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

func (m yaacTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m yaacTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m yaacTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
