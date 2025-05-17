package renderer

import "diagra/interpreter"

// PositionedNode is a struct that represents a node in the diagram with its position
type PositionedNode struct {
	Node interpreter.Node
	X, Y int
}

// PositionedEdge is a struct that represents an edge in the diagram with its start and end positions
type PositionedEdge struct {
	Edge         interpreter.Edge
	FromX, FromY int
	ToX, ToY     int
}

// Computelayout returns a layout for the diagram
// It uses a simple horizontal layout for nodes and edges
// The nodes are placed in a row with a fixed gap between them
func ComputeLayout(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	startX, startY := 150, 150
	gapX := 200

	// Place nodes horizontally (in a row)
	for i, n := range d.Nodes {
		x := startX + i*gapX
		y := startY
		pNodes = append(pNodes, PositionedNode{Node: n, X: x, Y: y})
	}

	// Create a position map for nodes
	// This map will be used to find the coordinates of each node
	posMap := map[string][2]int{}
	for _, pn := range pNodes {
		posMap[pn.Node.ID] = [2]int{pn.X, pn.Y}
	}

	// Connect nodes with directed edges
	// The edges are drawn from the center of the "from" node to the center of the "to" node
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

// ComputeVerticalLayout returns a layout for the diagram
// It uses a simple vertical layout for nodes and edges
func ComputeVerticalLayout(d interpreter.Diagram) ([]PositionedNode, []PositionedEdge) {
	var pNodes []PositionedNode
	var pEdges []PositionedEdge

	startX, startY := 400, 100 // middle of the screen
	gapY := 120

	// Place nodes vertically (in a column)
	for i, n := range d.Nodes {
		x := startX
		y := startY + i*gapY
		pNodes = append(pNodes, PositionedNode{Node: n, X: x, Y: y})
	}

	// Position map for nodes
	posMap := map[string][2]int{}
	for _, pn := range pNodes {
		posMap[pn.Node.ID] = [2]int{pn.X, pn.Y}
	}

	// Connect nodes with directed edges
	for _, e := range d.Edges {
		from := posMap[e.From]
		to := posMap[e.To]
		offset := 30 // half height of the node

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

// ComputeTreeLayout returns a layout for the diagram of type tree
// It uses a simple breadth-first search (BFS) to determine the levels of the nodes
// and assigns positions based on the levels
// The nodes are placed in a tree-like structure with a fixed gap between levels
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
		return pNodes, pEdges // no root found
	}

	// Build a BFS tree
	// BFS queue
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

	// Assign positions
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

	// Edges
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

// ComputeTreeLayoutRecursive returns a layout for the diagram of type tree
// It uses a recursive approach to determine the positions of the nodes
// and edges based on the tree structure
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

	// Find root: node that is not "To" in any edge
	// Create a map of targets
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

	// Recursive function to place nodes
	var place func(id string, depth int) int
	place = func(id string, depth int) int {
		c := children[id]
		if len(c) == 0 {
			// Leaf node: place itself
			x := currentX
			nodePositions[id] = [2]int{x, depth * yGap}
			currentX += xGap
			return x
		}

		// Inner node: place children first
		childXs := []int{}
		for _, cid := range c {
			cx := place(cid, depth+1)
			childXs = append(childXs, cx)
			edges = append(edges, PositionedEdge{
				Edge:  interpreter.Edge{From: id, To: cid},
				FromX: 0, FromY: 0, ToX: 0, ToY: 0, // fylls senare
			})
		}
		// Center the node under its children
		minX := childXs[0]
		maxX := childXs[len(childXs)-1]
		centerX := (minX + maxX) / 2
		nodePositions[id] = [2]int{centerX, depth * yGap}
		return centerX
	}

	place(rootID, 0)

	// Convert node positions to PositionedNode
	// and add a margin to the x and y coordinates
	// This is to avoid overlapping with the edges
	var pNodes []PositionedNode
	for id, pos := range nodePositions {
		node := nodeMap[id]
		pNodes = append(pNodes, PositionedNode{
			Node: node,
			X:    pos[0] + 100, // margin
			Y:    pos[1] + 50,
		})
	}

	// Uppdate edges with the correct positions
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
