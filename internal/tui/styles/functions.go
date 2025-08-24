package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func gradientText(text string, startColor, endColor lipgloss.Color, base, additional lipgloss.Style) string {
	n := len(text)
	result := ""

	for i := range n {
		ratio := float64(i) / float64(n-1)
		color := interpolateColor(startColor, endColor, ratio)

		style := base.Foreground(lipgloss.Color(color))
		result += style.Render(string(text[i]))
	}

	return additional.Render(result)
}

func interpolateColor(start, end lipgloss.Color, ratio float64) string {
	r1, g1, b1 := hexToRGB(string(start))
	r2, g2, b2 := hexToRGB(string(end))

	r := int(float64(r1) + (float64(r2)-float64(r1))*ratio)
	g := int(float64(g1) + (float64(g2)-float64(g1))*ratio)
	b := int(float64(b1) + (float64(b2)-float64(b1))*ratio)

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func ApplyGradientToFooter(text string) string {
	return gradientText("REQ", footerNameFGFrom, footerNameFGTo, footerNameStyle, footerNameBGStyle)
}
