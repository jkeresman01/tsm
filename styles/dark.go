package styles

import "github.com/charmbracelet/lipgloss"

func DarkTheme() Theme {
	p := darkPalette()
	s := darkSizes()

	return Theme{
		BorderColor:      p.border,
		SecondaryColor:   p.secondary,
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

func darkPalette() themePalette {
	return themePalette{
		border:    lipgloss.Color("15"),
		secondary: lipgloss.Color("241"),
		dimBG:     lipgloss.Color("236"),
	}
}

func darkSizes() themeSizes {
	return themeSizes{
		leftPanel:       40,
		rightPanel:      60,
		containerHeight: 25,
		paddingX:        2,
		paddingY:        1,
		marginX:         0,
		marginY:         1,
	}
}
