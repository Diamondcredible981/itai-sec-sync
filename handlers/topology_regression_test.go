package handlers

import (
	"testing"

	"github.com/iMayday-Yee/XinchuangAnalyze/models"
)

func TestMissingProductIDs(t *testing.T) {
	requested := []int{1, 2, 3, 4}
	productMap := map[int]models.Product{
		1: {ID: 1, Name: "p1"},
		3: {ID: 3, Name: "p3"},
	}

	missing := missingProductIDs(requested, productMap)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing products, got %d", len(missing))
	}
	if missing[0] != 2 || missing[1] != 4 {
		t.Fatalf("unexpected missing products: %#v", missing)
	}
}

func TestIsClearGraphRequest(t *testing.T) {
	updateData := models.NetworkTopo{
		ProductIDs: []int{},
		Nodes:      []models.TopoNode{},
		Edges:      []models.TopoEdge{},
	}

	if !isClearGraphRequest(false, true, true, updateData) {
		t.Fatalf("expected explicit empty nodes+edges to be treated as clear-graph request")
	}

	nonEmptyProducts := models.NetworkTopo{
		ProductIDs: []int{1},
		Nodes:      []models.TopoNode{},
		Edges:      []models.TopoEdge{},
	}
	if isClearGraphRequest(true, true, true, nonEmptyProducts) {
		t.Fatalf("expected non-empty product_ids to disable clear-graph shortcut")
	}
}

func TestSimulateAgainstBaselineEndpoints_WhenEndpointBlocked_ReturnsZeroPath(t *testing.T) {
	topoEdges := []models.TopoEdge{
		{FromNodeKey: "a", ToNodeKey: "b", EdgeType: "network", Direction: "uni", Weight: 1, Risk: 20},
		{FromNodeKey: "c", ToNodeKey: "d", EdgeType: "network", Direction: "uni", Weight: 1, Risk: 20},
	}

	blocked := map[string]bool{"a": true}
	risk, length := simulateAgainstBaselineEndpoints("a", "b", blocked, topoEdges, nil, nil, nil)
	if risk != 0 || length != 0 {
		t.Fatalf("expected blocked baseline endpoint to produce zero simulation, got risk=%d length=%d", risk, length)
	}
}
