package renderer

import "diagra/interpreter"

// Resultatstruktur
type PositionedNode struct {
	Node interpreter.Node
	X, Y int
}

type PositionedEdge struct {
	Edge         interpreter.Edge
	FromX, FromY int
	ToX, ToY     int
}

// Returnerar layout av noder och kanter i en diagram
func ComputeLayout(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	startX, startY := 150, 150
	gapX := 200

	// Placera noder horisontellt på samma rad
	for i, n := range d.Nodes {
		x := startX + i*gapX
		y := startY
		pNodes = append(pNodes, PositionedNode{Node: n, X: x, Y: y})
	}

	// Skapa en ID → position-karta
	posMap := map[string][2]int{}
	for _, pn := range pNodes {
		posMap[pn.Node.ID] = [2]int{pn.X, pn.Y}
	}

	// Koppla kanter till koordinater
	for _, e := range d.Edges {
		from := posMap[e.From]
		to := posMap[e.To]
		offset := 60 // för att inte peka rakt in i boxarna

		pEdges = append(pEdges, PositionedEdge{
			Edge:  e,
			FromX: from[0] + offset,
			FromY: from[1],
			ToX:   to[0] - offset,
			ToY:   to[1],
		})
	}

	return pNodes, pEdges
}

// ComputeVerticalLayout returnerar en layout för diagrammet i vertikal stil
func ComputeVerticalLayout(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	startX, startY := 400, 100 // mitten på canvas
	gapY := 120

	// Placera noder vertikalt (på rad)
	for i, n := range d.Nodes {
		x := startX
		y := startY + i*gapY
		pNodes = append(pNodes, PositionedNode{Node: n, X: x, Y: y})
	}

	// Positionskarta
	posMap := map[string][2]int{}
	for _, pn := range pNodes {
		posMap[pn.Node.ID] = [2]int{pn.X, pn.Y}
	}

	// Lägg till riktade pilar neråt
	for _, e := range d.Edges {
		from := posMap[e.From]
		to := posMap[e.To]
		offset := 30 // halv höjd på boxen

		pEdges = append(pEdges, PositionedEdge{
			Edge:  e,
			FromX: from[0],
			FromY: from[1] + offset,
			ToX:   to[0],
			ToY:   to[1] - offset,
		})
	}

	return pNodes, pEdges
}
