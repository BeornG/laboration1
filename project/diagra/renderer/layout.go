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

// ComputeTreeLayout returnerar en layout för diagram av typen tree
func ComputeTreeLayout(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	// Mappa noder
	nodeMap := map[string]interpreter.Node{}
	for _, n := range d.Nodes {
		nodeMap[n.ID] = n
	}

	// Hitta root: nod som inte är "To" i någon kant
	targets := map[string]bool{}
	for _, e := range d.Edges {
		targets[e.To] = true
	}
	var rootID string
	for _, n := range d.Nodes {
		if !targets[n.ID] {
			rootID = n.ID
			break
		}
	}
	if rootID == "" {
		return pNodes, pEdges // ingen root hittad
	}

	// Bygg nivåer (enkelt BFS)
	type queueItem struct {
		ID    string
		Level int
	}
	queue := []queueItem{{rootID, 0}}
	levels := map[int][]string{}
	visited := map[string]bool{}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if visited[item.ID] {
			continue
		}
		visited[item.ID] = true

		levels[item.Level] = append(levels[item.Level], item.ID)

		for _, e := range d.Edges {
			if e.From == item.ID {
				queue = append(queue, queueItem{e.To, item.Level + 1})
			}
		}
	}

	// Tilldela positioner
	yStart := 100
	xStart := 100
	xGap := 160
	yGap := 120

	posMap := map[string][2]int{}

	for level, ids := range levels {
		for i, id := range ids {
			x := xStart + i*xGap
			y := yStart + level*yGap
			pNodes = append(pNodes, PositionedNode{
				Node: nodeMap[id],
				X:    x,
				Y:    y,
			})
			posMap[id] = [2]int{x, y}
		}
	}

	// Kanter
	for _, e := range d.Edges {
		from := posMap[e.From]
		to := posMap[e.To]
		pEdges = append(pEdges, PositionedEdge{
			Edge:  e,
			FromX: from[0],
			FromY: from[1] + 25, // halva höjden ner
			ToX:   to[0],
			ToY:   to[1] - 25, // halva höjden upp
		})
	}

	return pNodes, pEdges
}

// ComputeTreeLayoutRecursive returnerar en layout för diagram av typen tree
func ComputeTreeLayoutRecursive(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	nodeMap := map[string]interpreter.Node{}
	children := map[string][]string{}
	edges := []PositionedEdge{}

	for _, n := range d.Nodes {
		nodeMap[n.ID] = n
	}
	for _, e := range d.Edges {
		children[e.From] = append(children[e.From], e.To)
	}

	// Hitta root: nod som inte är måltavla
	isTarget := map[string]bool{}
	for _, e := range d.Edges {
		isTarget[e.To] = true
	}
	var rootID string
	for _, n := range d.Nodes {
		if !isTarget[n.ID] {
			rootID = n.ID
			break
		}
	}
	if rootID == "" {
		return nil, nil
	}

	nodePositions := map[string][2]int{}
	currentX := 0
	yGap := 120
	xGap := 80

	// Rekursiv placering
	var place func(id string, depth int) int
	place = func(id string, depth int) int {
		c := children[id]
		if len(c) == 0 {
			// Löv: placera och returnera x
			x := currentX
			nodePositions[id] = [2]int{x, depth * yGap}
			currentX += xGap
			return x
		}

		// Inre nod: placera barn först
		childXs := []int{}
		for _, cid := range c {
			cx := place(cid, depth+1)
			childXs = append(childXs, cx)
			edges = append(edges, PositionedEdge{
				Edge:  interpreter.Edge{From: id, To: cid},
				FromX: 0, FromY: 0, ToX: 0, ToY: 0, // fylls senare
			})
		}
		// Centera denna nod över barnen
		minX := childXs[0]
		maxX := childXs[len(childXs)-1]
		centerX := (minX + maxX) / 2
		nodePositions[id] = [2]int{centerX, depth * yGap}
		return centerX
	}

	place(rootID, 0)

	// Konvertera till []PositionedNode
	var pNodes []PositionedNode
	for id, pos := range nodePositions {
		node := nodeMap[id]
		pNodes = append(pNodes, PositionedNode{
			Node: node,
			X:    pos[0] + 100, // margin
			Y:    pos[1] + 50,
		})
	}

	// Uppdatera edges med riktiga koordinater
	posMap := map[string][2]int{}
	for _, pn := range pNodes {
		posMap[pn.Node.ID] = [2]int{pn.X, pn.Y}
	}
	for i, e := range edges {
		from := posMap[e.Edge.From]
		to := posMap[e.Edge.To]
		edges[i].FromX = from[0]
		edges[i].FromY = from[1] + 25
		edges[i].ToX = to[0]
		edges[i].ToY = to[1] - 25
	}

	return pNodes, edges
}
