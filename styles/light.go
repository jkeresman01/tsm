package styles

import "github.com/charmbracelet/lipgloss"

func LightTheme() Theme {
	p := lightPalette()
	s := lightSizes()

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

func lightPalette() themePalette {
	return themePalette{
		border:    lipgloss.Color("0"),
		secondary: lipgloss.Color("8"),
		dimBG:     lipgloss.Color("252"),
	}
}

func lightSizes() themeSizes {
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
