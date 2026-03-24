package handlers

import (
	"container/heap"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

type AttackPathNode struct {
	NodeKey   string `json:"node_key,omitempty"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Type      string `json:"type"`
	Layer     int    `json:"layer"`
}

type AttackPathEdge struct {
	FromProductID   int    `json:"from_product_id"`
	ToProductID     int    `json:"to_product_id"`
	AttackTechnique string `json:"attack_technique"`
	RiskScore       int    `json:"risk_score"`
	Reason          string `json:"reason"`
}

type AttackPathResponse struct {
	TopologyID       uint                  `json:"topology_id"`
	TopologyName     string                `json:"topology_name"`
	SourceProductID  *int                  `json:"source_product_id,omitempty"`
	TargetProductID  *int                  `json:"target_product_id,omitempty"`
	BlockedProductID *int                  `json:"blocked_product_id,omitempty"`
	PathLength       int                   `json:"path_length"`
	OverallRisk      int                   `json:"overall_risk"`
	RiskLevel        string                `json:"risk_level"`
	Summary          AttackPathSummary     `json:"summary"`
	Nodes            []AttackPathNode      `json:"nodes"`
	Edges            []AttackPathEdge      `json:"edges"`
	KeyJumps         []AttackPathEdge      `json:"key_jumps"`
	Mitigations      []string              `json:"mitigations"`
	Simulation       *AttackPathSimulation `json:"simulation,omitempty"`
}

type AttackPathSummary struct {
	AvgJumpRisk       float64 `json:"avg_jump_risk"`
	MaxJumpRisk       int     `json:"max_jump_risk"`
	HighRiskJumpCount int     `json:"high_risk_jump_count"`
	WeakestNodeID     int     `json:"weakest_node_id"`
	WeakestNodeName   string  `json:"weakest_node_name"`
}

type AttackPathSimulation struct {
	BlockedProductID    int    `json:"blocked_product_id"`
	BlockedProductName  string `json:"blocked_product_name"`
	BaseRisk            int    `json:"base_risk"`
	SimulatedRisk       int    `json:"simulated_risk"`
	RiskReduction       int    `json:"risk_reduction"`
	BasePathLength      int    `json:"base_path_length"`
	SimulatedPathLength int    `json:"simulated_path_length"`
}

type graphArc struct {
	From string
	To   string
	Edge models.TopoEdge
}

type pqItem struct {
	node string
	cost int
	idx  int
}

type minHeap []*pqItem

func (h minHeap) Len() int { return len(h) }

func (h minHeap) Less(i, j int) bool {
	if h[i].cost == h[j].cost {
		return h[i].node < h[j].node
	}
	return h[i].cost < h[j].cost
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *minHeap) Push(x interface{}) {
	item := x.(*pqItem)
	item.idx = len(*h)
	*h = append(*h, item)
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.idx = -1
	*h = old[:n-1]
	return item
}

func (s *Service) GetAttackPathByTopoID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var topology models.NetworkTopo
	if err := s.DB.First(&topology, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}

	nodes, topoEdges, err := s.loadTopoGraph(topology.ID)
	if err != nil {
		s.internalError(c, "读取拓扑失败")
		return
	}
	if len(nodes) == 0 {
		c.JSON(http.StatusOK, AttackPathResponse{
			TopologyID:   topology.ID,
			TopologyName: topology.Name,
			PathLength:   0,
			OverallRisk:  0,
			RiskLevel:    "low",
			Summary:      AttackPathSummary{},
			Nodes:        []AttackPathNode{},
			Edges:        []AttackPathEdge{},
			KeyJumps:     []AttackPathEdge{},
			Mitigations:  []string{"拓扑中暂无节点，无法进行攻击路径推演"},
		})
		return
	}

	sourceProductID, hasSource, err := parseOptionalInt(c.Query("source_product_id"))
	if err != nil {
		s.badRequest(c, "source_product_id 参数无效")
		return
	}
	targetProductID, hasTarget, err := parseOptionalInt(c.Query("target_product_id"))
	if err != nil {
		s.badRequest(c, "target_product_id 参数无效")
		return
	}
	if hasSource != hasTarget {
		s.badRequest(c, "source_product_id 和 target_product_id 需同时提供")
		return
	}
	blockedProductID, hasBlocked, err := parseOptionalInt(c.Query("blocked_product_id"))
	if err != nil {
		s.badRequest(c, "blocked_product_id 参数无效")
		return
	}

	nodesByKey := make(map[string]models.TopoNode, len(nodes))
	productIDs := make([]int, 0)
	seenProduct := make(map[int]bool)
	for _, n := range nodes {
		nodesByKey[n.NodeKey] = n
		if n.ProductID == nil {
			continue
		}
		pid := int(*n.ProductID)
		if !seenProduct[pid] {
			productIDs = append(productIDs, pid)
			seenProduct[pid] = true
		}
	}
	topology.ProductIDs = productIDs

	var products []models.Product
	if len(topology.ProductIDs) > 0 {
		s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)
	}
	productMap := make(map[int]models.Product, len(products))
	for _, p := range products {
		p.FunctionIDs = utils.StringToIntSlice(p.FunctionIDsStr)
		productMap[int(p.ID)] = p
	}

	var productTypes []models.ProductType
	s.DB.Find(&productTypes)
	typeMap := make(map[uint]models.ProductType, len(productTypes))
	for _, t := range productTypes {
		typeMap[t.ID] = t
	}

	productToNodes := make(map[int][]string)
	for _, n := range nodes {
		if n.ProductID == nil {
			continue
		}
		pid := int(*n.ProductID)
		productToNodes[pid] = append(productToNodes[pid], n.NodeKey)
	}

	var sourcePtr *int
	var targetPtr *int
	var sourceCandidates []string
	var targetCandidates []string
	if hasSource && hasTarget {
		sourceCandidates = uniqueStrings(productToNodes[sourceProductID])
		targetCandidates = uniqueStrings(productToNodes[targetProductID])
		if len(sourceCandidates) == 0 || len(targetCandidates) == 0 {
			s.badRequest(c, "指定的 source/target 产品不在该拓扑中")
			return
		}
		sourcePtr = &sourceProductID
		targetPtr = &targetProductID
	} else {
		sourceCandidates, targetCandidates = autoSelectPathEndpoints(nodes, nil)
	}

	baseAdj := buildGraphArcs(topoEdges, nil)
	baseNodePath, baseArcs, ok := findBestPath(sourceCandidates, targetCandidates, baseAdj)
	if !ok || len(baseNodePath) < 2 {
		c.JSON(http.StatusOK, AttackPathResponse{
			TopologyID:       topology.ID,
			TopologyName:     topology.Name,
			SourceProductID:  sourcePtr,
			TargetProductID:  targetPtr,
			BlockedProductID: nil,
			PathLength:       0,
			OverallRisk:      0,
			RiskLevel:        "low",
			Summary:          AttackPathSummary{},
			Nodes:            []AttackPathNode{},
			Edges:            []AttackPathEdge{},
			KeyJumps:         []AttackPathEdge{},
			Mitigations:      []string{"当前拓扑中未发现可达攻击路径，建议核查关键边方向与分区隔离策略"},
		})
		return
	}

	pathNodes, pathEdges, overallRisk, keyJumps, mitigations, summary := buildPathArtifactsFromGraph(baseNodePath, baseArcs, nodesByKey, productMap, typeMap)
	baseSourceNodeKey := baseNodePath[0]
	baseTargetNodeKey := baseNodePath[len(baseNodePath)-1]

	var blockedPtr *int
	var simulation *AttackPathSimulation
	if hasBlocked {
		blockedNodeKeys, blockedName := collectBlockedNodeKeys(nodes, blockedProductID, productMap)
		if len(blockedNodeKeys) == 0 {
			s.badRequest(c, "blocked_product_id 不在该拓扑中")
			return
		}
		if hasSource && hasTarget && (sourceProductID == blockedProductID || targetProductID == blockedProductID) {
			s.badRequest(c, "source/target 不能与 blocked_product_id 相同")
			return
		}
		blockedPtr = &blockedProductID

		simulatedRisk, simulatedPathLength := simulateAgainstBaselineEndpoints(baseSourceNodeKey, baseTargetNodeKey, blockedNodeKeys, topoEdges, nodesByKey, productMap, typeMap)

		simulation = &AttackPathSimulation{
			BlockedProductID:    blockedProductID,
			BlockedProductName:  blockedName,
			BaseRisk:            overallRisk,
			SimulatedRisk:       simulatedRisk,
			RiskReduction:       maxInt(0, overallRisk-simulatedRisk),
			BasePathLength:      len(pathNodes),
			SimulatedPathLength: simulatedPathLength,
		}
	}

	c.JSON(http.StatusOK, AttackPathResponse{
		TopologyID:       topology.ID,
		TopologyName:     topology.Name,
		SourceProductID:  sourcePtr,
		TargetProductID:  targetPtr,
		BlockedProductID: blockedPtr,
		PathLength:       len(pathNodes),
		OverallRisk:      overallRisk,
		RiskLevel:        riskLevelByScore(overallRisk),
		Summary:          summary,
		Nodes:            pathNodes,
		Edges:            pathEdges,
		KeyJumps:         keyJumps,
		Mitigations:      mitigations,
		Simulation:       simulation,
	})
}

func buildPathArtifactsFromGraph(nodeKeys []string, arcs []graphArc, nodesByKey map[string]models.TopoNode, productMap map[int]models.Product, typeMap map[uint]models.ProductType) ([]AttackPathNode, []AttackPathEdge, int, []AttackPathEdge, []string, AttackPathSummary) {
	nodes := make([]AttackPathNode, 0, len(nodeKeys))
	for i, nodeKey := range nodeKeys {
		n := nodesByKey[nodeKey]
		respNode := AttackPathNode{
			NodeKey: nodeKey,
			Name:    n.Name,
			Brand:   n.Vendor,
			Type:    n.NodeType,
			Layer:   i + 1,
		}
		if n.ProductID != nil {
			pid := int(*n.ProductID)
			respNode.ProductID = pid
			if p, ok := productMap[pid]; ok {
				respNode.Name = p.Name
				respNode.Brand = p.Brand
				if pt, okType := typeMap[p.TypeID]; okType {
					respNode.Type = pt.Name
				}
			}
		}
		nodes = append(nodes, respNode)
	}

	edges := make([]AttackPathEdge, 0, len(arcs))
	jumpRisks := make([]int, 0, len(arcs))
	for _, arc := range arcs {
		fromNode := nodesByKey[arc.From]
		toNode := nodesByKey[arc.To]
		risk, technique, reason := estimateGraphEdgeRisk(arc, fromNode, toNode, productMap, typeMap)
		jumpRisks = append(jumpRisks, risk)

		edge := AttackPathEdge{
			AttackTechnique: technique,
			RiskScore:       risk,
			Reason:          reason,
		}
		if fromNode.ProductID != nil {
			edge.FromProductID = int(*fromNode.ProductID)
		}
		if toNode.ProductID != nil {
			edge.ToProductID = int(*toNode.ProductID)
		}
		edges = append(edges, edge)
	}

	overallRisk := calcPathRisk(edges)
	keyJumps := pickKeyJumps(edges)
	mitigations := buildGraphMitigations(nodeKeys, arcs, nodesByKey, productMap)
	summary := buildGraphPathSummary(nodeKeys, jumpRisks, nodesByKey, productMap)

	return nodes, edges, overallRisk, keyJumps, mitigations, summary
}

func parseOptionalInt(v string) (int, bool, error) {
	if v == "" {
		return 0, false, nil
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return 0, true, err
	}
	return parsed, true, nil
}

func estimateGraphEdgeRisk(arc graphArc, fromNode, toNode models.TopoNode, productMap map[int]models.Product, typeMap map[uint]models.ProductType) (int, string, string) {
	risk := arc.Edge.Risk
	if risk <= 0 {
		risk = 35
	}

	if arc.Edge.Weight > 1 {
		risk += minInt(15, (arc.Edge.Weight-1)*3)
	}

	if strings.EqualFold(fromNode.Criticality, "high") || strings.EqualFold(toNode.Criticality, "high") {
		risk += 10
	}

	reason := "存在可利用连通关系，具备横向移动条件"
	technique := mapEdgeTypeToTechnique(arc.Edge.EdgeType)

	if strings.EqualFold(arc.Edge.Direction, "bi") {
		risk += 8
		reason = "双向连通关系扩大了横向移动窗口"
	}

	if fromNode.ProductID != nil && toNode.ProductID != nil {
		src, srcOK := productMap[int(*fromNode.ProductID)]
		dst, dstOK := productMap[int(*toNode.ProductID)]
		if srcOK && dstOK {
			srcType := typeMap[src.TypeID].Name
			dstType := typeMap[dst.TypeID].Name
			if src.Brand != "" && src.Brand == dst.Brand {
				risk += 10
				reason = "相邻同厂商设备可能存在共性配置与漏洞利用链"
			}
			if srcType != "" && srcType == dstType {
				risk += 8
				technique = "same-tier-pivot"
				reason = "同类型节点之间更易形成同层横向扩散"
			}
			if len(src.FunctionIDs) >= 5 || len(dst.FunctionIDs) >= 5 {
				risk += 8
			}
		}
	}

	if risk > 100 {
		risk = 100
	}
	if risk < 0 {
		risk = 0
	}

	return risk, technique, reason
}

func calcPathRisk(edges []AttackPathEdge) int {
	if len(edges) == 0 {
		return 0
	}

	total := 0
	maxEdge := 0
	for _, edge := range edges {
		total += edge.RiskScore
		if edge.RiskScore > maxEdge {
			maxEdge = edge.RiskScore
		}
	}

	avg := total / len(edges)
	overall := int(float64(avg)*0.7 + float64(maxEdge)*0.3)
	if overall > 100 {
		return 100
	}
	if overall < 0 {
		return 0
	}
	return overall
}

func pickKeyJumps(edges []AttackPathEdge) []AttackPathEdge {
	if len(edges) == 0 {
		return []AttackPathEdge{}
	}

	sortedEdges := make([]AttackPathEdge, len(edges))
	copy(sortedEdges, edges)
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].RiskScore > sortedEdges[j].RiskScore
	})

	limit := 3
	if len(sortedEdges) < limit {
		limit = len(sortedEdges)
	}
	return sortedEdges[:limit]
}

func buildGraphPathSummary(nodeKeys []string, jumpRisks []int, nodesByKey map[string]models.TopoNode, productMap map[int]models.Product) AttackPathSummary {
	if len(jumpRisks) == 0 || len(nodeKeys) == 0 {
		return AttackPathSummary{}
	}

	total := 0
	maxRisk := 0
	highRiskCount := 0
	nodeRisk := make(map[string]int)

	for idx, r := range jumpRisks {
		total += r
		if r > maxRisk {
			maxRisk = r
		}
		if r >= 70 {
			highRiskCount++
		}

		from := nodeKeys[idx]
		to := nodeKeys[idx+1]
		if r > nodeRisk[from] {
			nodeRisk[from] = r
		}
		if r > nodeRisk[to] {
			nodeRisk[to] = r
		}
	}

	weakestNodeRisk := -1
	weakestNodeName := ""
	weakestNodeID := 0
	for _, key := range nodeKeys {
		r := nodeRisk[key]
		if r <= weakestNodeRisk {
			continue
		}
		n := nodesByKey[key]
		weakestNodeRisk = r
		weakestNodeName = n.Name
		if n.ProductID != nil {
			pid := int(*n.ProductID)
			weakestNodeID = pid
			if p, ok := productMap[pid]; ok {
				weakestNodeName = p.Name
			}
		}
	}

	avg := float64(total) / float64(len(jumpRisks))
	return AttackPathSummary{
		AvgJumpRisk:       attackPathRound2(avg),
		MaxJumpRisk:       maxRisk,
		HighRiskJumpCount: highRiskCount,
		WeakestNodeID:     weakestNodeID,
		WeakestNodeName:   weakestNodeName,
	}
}

func simulateAgainstBaselineEndpoints(
	baseSourceNodeKey string,
	baseTargetNodeKey string,
	blockedNodeKeys map[string]bool,
	topoEdges []models.TopoEdge,
	nodesByKey map[string]models.TopoNode,
	productMap map[int]models.Product,
	typeMap map[uint]models.ProductType,
) (int, int) {
	// 固定使用基线路径端点进行阻断对比，避免模拟前后目标不一致。
	if blockedNodeKeys[baseSourceNodeKey] || blockedNodeKeys[baseTargetNodeKey] {
		return 0, 0
	}

	simAdj := buildGraphArcs(topoEdges, blockedNodeKeys)
	simNodePath, simArcs, simOK, _ := shortestPath(baseSourceNodeKey, baseTargetNodeKey, simAdj)
	if !simOK || len(simNodePath) < 2 {
		return 0, 0
	}

	_, simEdges, simRisk, _, _, _ := buildPathArtifactsFromGraph(simNodePath, simArcs, nodesByKey, productMap, typeMap)
	return simRisk, len(simEdges) + 1
}

func buildGraphMitigations(nodeKeys []string, arcs []graphArc, nodesByKey map[string]models.TopoNode, productMap map[int]models.Product) []string {
	mitigations := []string{
		"对关键跳点之间实施最小权限访问控制与东西向微隔离策略",
		"在高风险跳点启用持续审计与关联告警，缩短横向移动发现时间",
	}

	hasBi := false
	hasHighRisk := false
	hasSameBrand := false
	for _, arc := range arcs {
		if strings.EqualFold(arc.Edge.Direction, "bi") {
			hasBi = true
		}
		if arc.Edge.Risk >= 75 {
			hasHighRisk = true
		}
		fromNode := nodesByKey[arc.From]
		toNode := nodesByKey[arc.To]
		if fromNode.ProductID != nil && toNode.ProductID != nil {
			src, okSrc := productMap[int(*fromNode.ProductID)]
			dst, okDst := productMap[int(*toNode.ProductID)]
			if okSrc && okDst && src.Brand != "" && src.Brand == dst.Brand {
				hasSameBrand = true
			}
		}
	}

	if hasBi {
		mitigations = append(mitigations, "优先收敛双向信任链路，改为按需单向授权并设置到期回收")
	}
	if hasHighRisk {
		mitigations = append(mitigations, "对高风险链路启用跳板审计与双因素认证，阻断凭证滥用扩散")
	}
	if hasSameBrand {
		mitigations = append(mitigations, "关键相邻节点优先采用异构厂商组合，降低同源漏洞联动风险")
	}
	if len(nodeKeys) >= 4 {
		mitigations = append(mitigations, "对长路径分段设置访问策略与分区边界，防止单点失陷后连续扩散")
	}

	return mitigations
}

func collectBlockedNodeKeys(nodes []models.TopoNode, blockedProductID int, productMap map[int]models.Product) (map[string]bool, string) {
	blocked := make(map[string]bool)
	blockedName := ""
	for _, n := range nodes {
		if n.ProductID == nil || int(*n.ProductID) != blockedProductID {
			continue
		}
		blocked[n.NodeKey] = true
		if blockedName == "" {
			if p, ok := productMap[blockedProductID]; ok {
				blockedName = p.Name
			} else {
				blockedName = n.Name
			}
		}
	}
	return blocked, blockedName
}

func filterBlockedCandidates(candidates []string, blocked map[string]bool) []string {
	if len(blocked) == 0 {
		return uniqueStrings(candidates)
	}
	result := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if blocked[c] {
			continue
		}
		result = append(result, c)
	}
	return uniqueStrings(result)
}

func buildGraphArcs(edges []models.TopoEdge, blockedNodeKeys map[string]bool) map[string][]graphArc {
	adj := make(map[string][]graphArc)
	for _, e := range edges {
		if blockedNodeKeys != nil && (blockedNodeKeys[e.FromNodeKey] || blockedNodeKeys[e.ToNodeKey]) {
			continue
		}
		adj[e.FromNodeKey] = append(adj[e.FromNodeKey], graphArc{From: e.FromNodeKey, To: e.ToNodeKey, Edge: e})
		if strings.EqualFold(e.Direction, "bi") {
			adj[e.ToNodeKey] = append(adj[e.ToNodeKey], graphArc{From: e.ToNodeKey, To: e.FromNodeKey, Edge: e})
		}
	}

	for key := range adj {
		sort.Slice(adj[key], func(i, j int) bool {
			if adj[key][i].To == adj[key][j].To {
				return adj[key][i].Edge.ID < adj[key][j].Edge.ID
			}
			return adj[key][i].To < adj[key][j].To
		})
	}
	return adj
}

func autoSelectPathEndpoints(nodes []models.TopoNode, blockedNodeKeys map[string]bool) ([]string, []string) {
	sources := make([]string, 0)
	targets := make([]string, 0)
	allHardware := make([]models.TopoNode, 0)

	const maxIntVal = int(^uint(0) >> 1)
	minLayer := maxIntVal
	maxLayer := -1

	for _, n := range nodes {
		if blockedNodeKeys != nil && blockedNodeKeys[n.NodeKey] {
			continue
		}
		isHardware := strings.EqualFold(n.NodeType, "hardware") || n.ProductID != nil
		if !isHardware {
			continue
		}
		allHardware = append(allHardware, n)
		zone := strings.ToLower(n.Zone)
		if zone == "edge" || zone == "internet" || zone == "dmz" {
			sources = append(sources, n.NodeKey)
		}
		if zone == "internal" || zone == "core" || zone == "data" {
			targets = append(targets, n.NodeKey)
		}
		if n.Layer < minLayer {
			minLayer = n.Layer
		}
		if n.Layer > maxLayer {
			maxLayer = n.Layer
		}
	}

	if len(sources) == 0 || len(targets) == 0 {
		sources = []string{}
		targets = []string{}
		for _, n := range allHardware {
			if n.Layer == minLayer {
				sources = append(sources, n.NodeKey)
			}
			if n.Layer == maxLayer {
				targets = append(targets, n.NodeKey)
			}
		}
	}

	if len(sources) == 0 || len(targets) == 0 {
		sources = []string{}
		targets = []string{}
		for _, n := range allHardware {
			sources = append(sources, n.NodeKey)
			targets = append(targets, n.NodeKey)
		}
	}

	return uniqueStrings(sources), uniqueStrings(targets)
}

func findBestPath(sources, targets []string, adj map[string][]graphArc) ([]string, []graphArc, bool) {
	if len(sources) == 0 || len(targets) == 0 {
		return nil, nil, false
	}

	bestCost := int(^uint(0) >> 1)
	bestNodes := []string(nil)
	bestArcs := []graphArc(nil)
	for _, src := range sources {
		for _, dst := range targets {
			if src == dst {
				continue
			}
			nodes, arcs, ok, cost := shortestPath(src, dst, adj)
			if !ok {
				continue
			}
			if cost < bestCost || (cost == bestCost && len(arcs) > len(bestArcs)) {
				bestCost = cost
				bestNodes = nodes
				bestArcs = arcs
			}
		}
	}

	if len(bestNodes) == 0 {
		return nil, nil, false
	}
	return bestNodes, bestArcs, true
}

func shortestPath(source, target string, adj map[string][]graphArc) ([]string, []graphArc, bool, int) {
	dist := map[string]int{source: 0}
	prevNode := make(map[string]string)
	prevArc := make(map[string]graphArc)

	h := &minHeap{}
	heap.Init(h)
	heap.Push(h, &pqItem{node: source, cost: 0})

	for h.Len() > 0 {
		item := heap.Pop(h).(*pqItem)
		if curCost, ok := dist[item.node]; ok && item.cost > curCost {
			continue
		}
		if item.node == target {
			break
		}

		for _, arc := range adj[item.node] {
			nextCost := item.cost + edgeTraversalCost(arc.Edge)
			if cur, ok := dist[arc.To]; !ok || nextCost < cur {
				dist[arc.To] = nextCost
				prevNode[arc.To] = item.node
				prevArc[arc.To] = arc
				heap.Push(h, &pqItem{node: arc.To, cost: nextCost})
			}
		}
	}

	best, ok := dist[target]
	if !ok {
		return nil, nil, false, 0
	}

	nodes := make([]string, 0)
	arcs := make([]graphArc, 0)
	for cur := target; cur != ""; {
		nodes = append(nodes, cur)
		if cur == source {
			break
		}
		arc, hasArc := prevArc[cur]
		if !hasArc {
			return nil, nil, false, 0
		}
		arcs = append(arcs, arc)
		cur = prevNode[cur]
	}

	reverseStrings(nodes)
	reverseArcs(arcs)
	return nodes, arcs, true, best
}

func edgeTraversalCost(edge models.TopoEdge) int {
	risk := edge.Risk
	if risk < 0 {
		risk = 0
	}
	if risk > 100 {
		risk = 100
	}

	weight := edge.Weight
	if weight <= 0 {
		weight = 1
	}

	// 风险越高，攻击者利用成本越低。
	cost := (101 - risk) + weight*5
	if strings.EqualFold(edge.Direction, "bi") {
		cost -= 6
	}
	if cost < 1 {
		return 1
	}
	return cost
}

func mapEdgeTypeToTechnique(edgeType string) string {
	switch strings.ToLower(edgeType) {
	case "trust":
		return "trust-abuse"
	case "depend":
		return "dependency-pivot"
	case "install":
		return "supply-chain-pivot"
	case "network":
		return "lateral-movement"
	default:
		return "lateral-movement"
	}
}

func uniqueStrings(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	seen := make(map[string]bool, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if seen[item] {
			continue
		}
		seen[item] = true
		result = append(result, item)
	}
	return result
}

func reverseStrings(items []string) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}

func reverseArcs(items []graphArc) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func attackPathRound2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func riskLevelByScore(score int) string {
	switch {
	case score >= 75:
		return "critical"
	case score >= 50:
		return "high"
	case score >= 25:
		return "medium"
	default:
		return "low"
	}
}
