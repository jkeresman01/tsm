package styles

import "github.com/charmbracelet/lipgloss"

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			DarkTheme creates a dark color theme.
//
//		@Return			Theme	Dark theme configuration
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			darkPalette returns the color palette for dark theme.
//
//		@Return			themePalette	Dark color palette
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func darkPalette() themePalette {
	return themePalette{
		border:    lipgloss.Color("15"),
		secondary: lipgloss.Color("241"),
		dimBG:     lipgloss.Color("236"),
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			darkSizes returns the size dimensions for dark theme.
//
//		@Return			themeSizes	Dark theme dimensions
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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
