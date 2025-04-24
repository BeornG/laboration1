package renderer

import (
	"diagra/interpreter"
	"fmt"
	"strings"
)

// Konstanter för att definiera layouten av noder och kanter i SVG-diagrammet.
const (
	nodeSpacingX = 200 // Avstånd mellan noder i X-led
	nodeSpacingY = 100 // Avstånd mellan noder i Y-led
	margin       = 100 // Marginal runt diagrammet
	nodeHeight   = 100 // Höjd på varje nod
)

// RenderSVG tar ett diagram och returnerar en SVG-sträng som representerar diagrammet.
func RenderSVG(d interpreter.Diagram) string {
	var sb strings.Builder

	defaultHeight := 600
	defaultWidth := 800

	height := defaultHeight
	width := defaultWidth

	// räkna total noder och kanter
	switch d.Layout {
	case "vertical":
		height = len(d.Nodes)*nodeSpacingY + margin + nodeHeight
		width = len(d.Nodes)*nodeSpacingX + margin
	default:
		height = len(d.Nodes)*nodeSpacingY + margin
		width = len(d.Nodes)*nodeSpacingX + margin
	}

	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+"\n",
		width, height, width, height,
	))

	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	// switch d.Layout {
	// case "vertical":
	// 	pNodes, pEdges = ComputeVerticalLayout(d)
	// default:
	// 	pNodes, pEdges = ComputeLayout(d) // horisontell som fallback
	// }

	switch d.Name {
	case "flowchart":
		if d.Layout == "vertical" {
			pNodes, pEdges = ComputeVerticalLayout(d)
		} else {
			pNodes, pEdges = ComputeLayout(d) // horisontell som fallback
		}
	case "tree":
		pNodes, pEdges = ComputeTreeLayout(d)
	}

	// Noder
	for _, n := range pNodes {
		x, y := n.X, n.Y
		if n.Node.Shape == "ellipse" {
			sb.WriteString(fmt.Sprintf(
				`  <ellipse cx="%d" cy="%d" rx="50" ry="25" fill="%s" stroke="%s" stroke-width="2"/>`+"\n",
				x, y, n.Node.Color, n.Node.Border,
			))
		} else {
			sb.WriteString(fmt.Sprintf(
				`  <rect x="%d" y="%d" width="100" height="50" rx="10" ry="10" fill="%s" stroke="%s" stroke-width="2"/>`+"\n",
				x-50, y-25, n.Node.Color, n.Node.Border,
			))
		}

		sb.WriteString(fmt.Sprintf(
			`  <text x="%d" y="%d" font-size="14" text-anchor="middle" fill="%s">%s</text>`+"\n",
			x, y+5, n.Node.Text, n.Node.Label,
		))
	}

	// Kanter
	for _, e := range pEdges {
		sb.WriteString(fmt.Sprintf(
			`  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="%s" marker-end="url(#arrow)"/>`+"\n",
			e.FromX, e.FromY, e.ToX, e.ToY, e.Edge.Color, e.Edge.Width,
		))
		midX := (e.FromX + e.ToX) / 2
		midY := (e.FromY + e.ToY) / 2

		labelX := midX + 10 // flytta etiketten åt sidan
		labelY := midY - 5  // lite ovanför linjen

		if d.Layout == "vertical" {
			labelX += 10
			labelY -= 5
		} else { // default är horisontell
			labelY -= 10
			labelX -= 30
		}

		sb.WriteString(fmt.Sprintf(
			`  <text x="%d" y="%d" font-size="12" text-anchor="start" fill="#37474f">%s</text>`+"\n",
			labelX, labelY, e.Edge.Label,
		))

	}

	// Pilar
	sb.WriteString(`
  <defs>
    <marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5"
            markerWidth="6" markerHeight="6"
            orient="auto-start-reverse">
      <path d="M 0 0 L 10 5 L 0 10 z" fill="#37474f"/>
    </marker>
  </defs>
`)

	sb.WriteString(`</svg>`)
	return sb.String()
}
