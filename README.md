# osmgraph

[![Release](https://img.shields.io/github/v/release/ize-302/osmgraph?style=for-the-badge&logo=github&color=blue)](https://github.com/ize-302/osmgraph/releases)
[![Go Reference](https://img.shields.io/badge/go.dev-reference-00ADD8?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/ize-302/osmgraph)

Builds a road graph from an OpenStreetMap PBF file (osm.pbf). Returns a node map (coordinates) and an adjacency map (edges) suitable for pathfinding.

## Install

```sh
go get github.com/ize-302/osmgraph
```

## Usage

```go
import "github.com/ize-302/osmgraph/osmgraph"

f, err := os.Open("region.osm.pbf")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

nodes, edges, err := osmgraph.GraphBuilder(f, osmgraph.DefaultRoadFilter, osmgraph.DefaultOneway)
if err != nil {
    log.Fatal(err)
}

// nodes: map[int64]osm.Node  — road-graph vertices with lat/lon
// edges: map[int64][]int64   — adjacency list (node ID → reachable node IDs)
```
