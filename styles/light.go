package styles

import "github.com/charmbracelet/lipgloss"

func LightTheme() Theme {
	p := lightPalette()
	s := lightSizes()

	return Theme{
		BorderColor:      p.border,
		SecondaryColor:   p.secondary,
		AccentColor:      p.accent,
		HighlightColor:   p.highlight,
		LeftPanelWidth:   s.leftPanel,
		RightPanelWidth:  s.rightPanel,
		ContainerHeight:  s.containerHeight,
		OuterStyle:       outerStyle(p, s),
		HeaderStyle:      headerStyle(p),
		FooterStyle:      footerStyle(p),
		DimmedBackground: dimStyle(p),
		ListStyle:        listStyle(s),
		PreviewStyle:     previewStyle(s),
	}
}

func lightPalette() themePalette {
	return themePalette{
		border:    lipgloss.Color("#7C3AED"),
		secondary: lipgloss.Color("#64748B"),
		dimBG:     lipgloss.Color("#F1F5F9"),
		accent:    lipgloss.Color("#8B5CF6"),
		highlight: lipgloss.Color("#A855F7"),
	}
}

func lightSizes() themeSizes {
	return themeSizes{
		leftPanel:       45,
		rightPanel:      55,
		containerHeight: 28,
		paddingX:        2,
		paddingY:        1,
		marginX:         1,
		marginY:         1,
	}
}
