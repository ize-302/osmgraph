// Package osmgraph
package osmgraph

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// DefaultRoadFilter includes common driveable highway types
var DefaultRoadFilter = func(tags map[string]string) bool {
	switch tags["highway"] {
	case "motorway", "trunk", "primary", "secondary", "tertiary", "residential", "service":
		return true
	}
	return false
}

// DefaultOneway interprets standard OSM oneway tags.
var DefaultOneway = func(tags map[string]string) (isOneway bool, isReversed bool) {
	o := tags["oneway"]
	h := tags["highway"]
	isOneway = o == "yes" || o == "1" || o == "true" || (h == "motorway" && o == "")
	isReversed = o == "-1"
	return
}

// GraphBuilder builds and returns a graph consisting of two maps: vertices & edges
func GraphBuilder(
	r io.Reader,
	filter func(tags map[string]string) bool,
	oneway func(tags map[string]string) (isOneway bool, isReversed bool),
) (map[int64]osm.Node, map[int64][]int64, error) {
	if filter == nil || oneway == nil {
		return nil, nil, fmt.Errorf("osmgraph: filter and oneway must not be nil")
	}

	nodes := make(map[int64]osm.Node)
	edges := make(map[int64][]int64)

	scanner := osmpbf.New(context.Background(), r, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Node:
			nodes[int64(o.ID)] = *o

		case *osm.Way:
			tags := o.Tags.Map()
			if !filter(tags) {
				continue
			}

			if len(o.Nodes) < 2 {
				continue
			}

			isOneway, isReversed := oneway(tags)

			for i := 0; i < len(o.Nodes)-1; i++ {
				from := int64(o.Nodes[i].ID)
				to := int64(o.Nodes[i+1].ID)

				if isReversed {
					edges[to] = append(edges[to], from)
					if !isOneway {
						edges[from] = append(edges[from], to)
					}
				} else {
					edges[from] = append(edges[from], to)
					if !isOneway {
						edges[to] = append(edges[to], from)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("osmgraph: scan error: %w", err)
	}

	referenced := make(map[int64]struct{})
	for from, tos := range edges {
		referenced[from] = struct{}{}
		for _, to := range tos {
			referenced[to] = struct{}{}
		}
	}
	for id := range nodes {
		if _, ok := referenced[id]; !ok {
			delete(nodes, id)
		}
	}

	return nodes, edges, nil
}
