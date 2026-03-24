package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
	"gorm.io/gorm"
)

// 获取所有网络拓扑
func (s *Service) ListTopos(c *gin.Context) {
	var topologies []models.NetworkTopo
	s.DB.Find(&topologies)

	var allProducts []models.Product
	s.DB.Find(&allProducts)
	productMap := make(map[int]models.Product)
	for _, p := range allProducts {
		p.FunctionIDs = utils.StringToIntSlice(p.FunctionIDsStr)
		productMap[int(p.ID)] = p
	}

	for i := range topologies {
		ids, err := s.deriveProductIDsFromTopo(topologies[i].ID)
		if err != nil {
			s.internalError(c, fmt.Sprintf("读取拓扑 %d 失败", topologies[i].ID))
			return
		}
		topologies[i].ProductIDs = ids

		nodes, edges, err := s.loadTopoGraph(topologies[i].ID)
		if err != nil {
			s.internalError(c, fmt.Sprintf("读取拓扑图 %d 失败", topologies[i].ID))
			return
		}
		topologies[i].Nodes = nodes
		topologies[i].Edges = edges

		if len(topologies[i].ProductIDs) > 0 {
			products := make([]models.Product, 0, len(topologies[i].ProductIDs))
			for _, pid := range topologies[i].ProductIDs {
				if p, ok := productMap[pid]; ok {
					products = append(products, p)
				}
			}
			topologies[i].Products = products
		}
	}

	c.JSON(http.StatusOK, topologies)
}

// 获取单个网络拓扑
func (s *Service) GetTopo(c *gin.Context) {
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

	ids, err := s.deriveProductIDsFromTopo(topology.ID)
	if err != nil {
		s.internalError(c, "读取拓扑失败")
		return
	}
	topology.ProductIDs = ids

	nodes, edges, err := s.loadTopoGraph(topology.ID)
	if err != nil {
		s.internalError(c, "读取拓扑图失败")
		return
	}
	topology.Nodes = nodes
	topology.Edges = edges

	if len(topology.ProductIDs) > 0 {
		var products []models.Product
		s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)
		for i := range products {
			products[i].FunctionIDs = utils.StringToIntSlice(products[i].FunctionIDsStr)
		}
		topology.Products = products
	}

	c.JSON(http.StatusOK, topology)
}

// 添加网络拓扑
func (s *Service) AddTopo(c *gin.Context) {
	var topology models.NetworkTopo
	if err := c.ShouldBindJSON(&topology); err != nil {
		s.badRequest(c, err.Error())
		return
	}
	if topology.Name == "" {
		s.badRequest(c, "拓扑名称不能为空")
		return
	}

	productIDs, nodes, edges, err := s.normalizeTopoInput(topology.ProductIDs, topology.Nodes, topology.Edges)
	if err != nil {
		s.badRequest(c, err.Error())
		return
	}

	if err := s.validateProductIDs(productIDs); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	topology.ProductIDs = productIDs

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&topology).Error; err != nil {
			return err
		}

		for i := range nodes {
			nodes[i].ID = 0
			nodes[i].TopoID = topology.ID
		}
		for i := range edges {
			edges[i].ID = 0
			edges[i].TopoID = topology.ID
		}

		if len(nodes) > 0 {
			if err := tx.Create(&nodes).Error; err != nil {
				return err
			}
		}
		if len(edges) > 0 {
			if err := tx.Create(&edges).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.internalError(c, "创建失败")
		return
	}
	topology.Nodes = nodes
	topology.Edges = edges

	if len(productIDs) > 0 {
		var products []models.Product
		s.DB.Where("id IN ?", productIDs).Find(&products)
		for i := range products {
			products[i].FunctionIDs = utils.StringToIntSlice(products[i].FunctionIDsStr)
		}
		topology.Products = products
	}

	c.JSON(http.StatusCreated, topology)

	return

}

// 更新网络拓扑
func (s *Service) UpdateTopo(c *gin.Context) {
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

	rawBody, err := c.GetRawData()
	if err != nil {
		s.badRequest(c, "无效的请求体")
		return
	}

	var rawPayload map[string]json.RawMessage
	if err := json.Unmarshal(rawBody, &rawPayload); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	var updateData models.NetworkTopo
	if err := json.Unmarshal(rawBody, &updateData); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	_, hasProductIDsField := rawPayload["product_ids"]
	_, hasNodesField := rawPayload["nodes"]
	_, hasEdgesField := rawPayload["edges"]

	updateGraph := hasProductIDsField || hasNodesField || hasEdgesField
	if updateGraph {
		clearGraph := hasNodesField && hasEdgesField && len(updateData.Nodes) == 0 && len(updateData.Edges) == 0 && (!hasProductIDsField || len(updateData.ProductIDs) == 0)
		if clearGraph {
			topology.ProductIDs = []int{}
			topology.Nodes = []models.TopoNode{}
			topology.Edges = []models.TopoEdge{}
		} else {
			productIDs, nodes, edges, err := s.normalizeTopoInput(updateData.ProductIDs, updateData.Nodes, updateData.Edges)
			if err != nil {
				s.badRequest(c, err.Error())
				return
			}
			if err := s.validateProductIDs(productIDs); err != nil {
				s.badRequest(c, err.Error())
				return
			}

			topology.ProductIDs = productIDs
			topology.Nodes = nodes
			topology.Edges = edges
		}
	}

	// 更新其他字段
	if updateData.Name != "" {
		topology.Name = updateData.Name
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&topology).Error; err != nil {
			return err
		}

		if updateGraph {
			if err := tx.Where("topo_id = ?", topology.ID).Delete(&models.TopoEdge{}).Error; err != nil {
				return err
			}
			if err := tx.Where("topo_id = ?", topology.ID).Delete(&models.TopoNode{}).Error; err != nil {
				return err
			}

			for i := range topology.Nodes {
				topology.Nodes[i].ID = 0
				topology.Nodes[i].TopoID = topology.ID
			}
			for i := range topology.Edges {
				topology.Edges[i].ID = 0
				topology.Edges[i].TopoID = topology.ID
			}

			if len(topology.Nodes) > 0 {
				if err := tx.Create(&topology.Nodes).Error; err != nil {
					return err
				}
			}
			if len(topology.Edges) > 0 {
				if err := tx.Create(&topology.Edges).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		s.internalError(c, "更新失败")
		return
	}

	c.JSON(http.StatusOK, topology)
}

// 删除网络拓扑
func (s *Service) DeleteTopo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("topo_id = ?", id).Delete(&models.TopoEdge{}).Error; err != nil {
			return err
		}
		if err := tx.Where("topo_id = ?", id).Delete(&models.TopoNode{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.NetworkTopo{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.internalError(c, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 复制网络拓扑
func (s *Service) CopyTopo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var originalTopo models.NetworkTopo
	if err := s.DB.First(&originalTopo, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}

	newTopo := models.NetworkTopo{Name: originalTopo.Name + "_副本"}
	nodes, edges, err := s.loadTopoGraph(originalTopo.ID)
	if err != nil {
		s.internalError(c, "复制失败")
		return
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newTopo).Error; err != nil {
			return err
		}

		for i := range nodes {
			nodes[i].ID = 0
			nodes[i].TopoID = newTopo.ID
		}
		for i := range edges {
			edges[i].ID = 0
			edges[i].TopoID = newTopo.ID
		}

		if len(nodes) > 0 {
			if err := tx.Create(&nodes).Error; err != nil {
				return err
			}
		}
		if len(edges) > 0 {
			if err := tx.Create(&edges).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.internalError(c, "复制失败")
		return
	}

	newTopo.ProductIDs = extractProductIDsFromNodes(nodes)
	newTopo.Nodes = nodes
	newTopo.Edges = edges

	c.JSON(http.StatusCreated, newTopo)
}

// 获取拓扑的可视化数据
func (s *Service) GetTopoVisualization(c *gin.Context) {
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

	topology.ProductIDs, err = s.deriveProductIDsFromTopo(topology.ID)
	if err != nil {
		s.internalError(c, "读取拓扑失败")
		return
	}
	topology.Nodes, topology.Edges, err = s.loadTopoGraph(topology.ID)
	if err != nil {
		s.internalError(c, "读取拓扑图失败")
		return
	}

	var nodes []map[string]interface{}
	var edges []map[string]interface{}

	nodeMap := make(map[string]models.TopoNode)

	var productTypes []models.ProductType
	s.DB.Find(&productTypes)
	ptMap := make(map[uint]models.ProductType)
	for _, pt := range productTypes {
		ptMap[pt.ID] = pt
	}

	var products []models.Product
	s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)
	productMap := make(map[int]models.Product)
	for _, p := range products {
		productMap[int(p.ID)] = p
	}

	sort.Slice(topology.Nodes, func(i, j int) bool {
		if topology.Nodes[i].Layer == topology.Nodes[j].Layer {
			return topology.Nodes[i].ID < topology.Nodes[j].ID
		}
		return topology.Nodes[i].Layer < topology.Nodes[j].Layer
	})

	for _, n := range topology.Nodes {
		nodeMap[n.NodeKey] = n
		item := map[string]interface{}{
			"id":          n.NodeKey,
			"label":       n.Name,
			"node_type":   n.NodeType,
			"zone":        n.Zone,
			"criticality": n.Criticality,
			"x":           100 + (n.Layer%5)*150,
			"y":           100 + (n.Layer/5)*100,
		}

		if n.ProductID != nil {
			if p, ok := productMap[int(*n.ProductID)]; ok {
				item["product"] = p
				if pt, ok := ptMap[p.TypeID]; ok {
					item["type"] = pt.Name
					item["category"] = pt.Name
					item["icon"] = pt.Icon
					item["color"] = pt.Color
				}
			}
		}

		nodes = append(nodes, item)
	}

	for _, e := range topology.Edges {
		if _, ok := nodeMap[e.FromNodeKey]; !ok {
			continue
		}
		if _, ok := nodeMap[e.ToNodeKey]; !ok {
			continue
		}
		edges = append(edges, map[string]interface{}{
			"source":    e.FromNodeKey,
			"target":    e.ToNodeKey,
			"type":      e.EdgeType,
			"direction": e.Direction,
			"risk":      e.Risk,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes":    nodes,
		"edges":    edges,
		"topology": topology,
	})
}

func (s *Service) normalizeTopoInput(productIDs []int, nodes []models.TopoNode, edges []models.TopoEdge) ([]int, []models.TopoNode, []models.TopoEdge, error) {
	if len(productIDs) == 0 && len(nodes) == 0 {
		return nil, nil, nil, fmt.Errorf("请提供 product_ids 或 nodes")
	}

	productIDs = utils.UniqueIntSlice(productIDs)

	if len(nodes) == 0 {
		return nil, nil, nil, fmt.Errorf("新模式必须提供 nodes")
	}

	allowedNodeTypes := map[string]bool{
		"hardware": true,
		"software": true,
		"os":       true,
		"service":  true,
	}
	allowedCriticality := map[string]bool{
		"low":      true,
		"normal":   true,
		"high":     true,
		"critical": true,
	}
	allowedEdgeTypes := map[string]bool{
		"network": true,
		"install": true,
		"depend":  true,
		"trust":   true,
	}
	allowedDirection := map[string]bool{
		"uni": true,
		"bi":  true,
	}

	hardwareCount := 0

	for i := range nodes {
		if nodes[i].NodeKey == "" {
			nodes[i].NodeKey = fmt.Sprintf("node-%d", i+1)
		}
		if nodes[i].NodeType == "" {
			if nodes[i].ProductID != nil {
				nodes[i].NodeType = "hardware"
			} else {
				nodes[i].NodeType = "service"
			}
		}
		if nodes[i].Layer == 0 {
			nodes[i].Layer = i + 1
		}
		if nodes[i].Criticality == "" {
			nodes[i].Criticality = "normal"
		}
		if nodes[i].Zone == "" {
			nodes[i].Zone = "default"
		}

		if !allowedNodeTypes[nodes[i].NodeType] {
			return nil, nil, nil, fmt.Errorf("无效的 node_type: %s", nodes[i].NodeType)
		}
		if !allowedCriticality[nodes[i].Criticality] {
			return nil, nil, nil, fmt.Errorf("无效的 criticality: %s", nodes[i].Criticality)
		}

		if nodes[i].NodeType == "hardware" {
			hardwareCount++
			if nodes[i].ProductID == nil {
				return nil, nil, nil, fmt.Errorf("硬件节点必须绑定 product_id: %s", nodes[i].NodeKey)
			}
		} else {
			if nodes[i].ProductID != nil {
				return nil, nil, nil, fmt.Errorf("非硬件节点不能绑定 product_id: %s", nodes[i].NodeKey)
			}
		}
	}

	if hardwareCount == 0 {
		return nil, nil, nil, fmt.Errorf("拓扑至少需要一个硬件节点")
	}

	nodeKeyMap := make(map[string]bool)
	for _, n := range nodes {
		if nodeKeyMap[n.NodeKey] {
			return nil, nil, nil, fmt.Errorf("node_key 不能重复: %s", n.NodeKey)
		}
		nodeKeyMap[n.NodeKey] = true
	}

	if len(edges) == 0 {
		edges = buildDefaultEdgesFromNodes(nodes)
	}
	if len(nodes) > 1 && len(edges) == 0 {
		return nil, nil, nil, fmt.Errorf("多节点拓扑必须包含边")
	}

	edgeSet := make(map[string]bool)
	adj := make(map[string][]string)

	for i := range edges {
		if !nodeKeyMap[edges[i].FromNodeKey] || !nodeKeyMap[edges[i].ToNodeKey] {
			return nil, nil, nil, fmt.Errorf("edge 引用了不存在的节点")
		}
		if edges[i].FromNodeKey == edges[i].ToNodeKey {
			return nil, nil, nil, fmt.Errorf("edge 不能自环: %s", edges[i].FromNodeKey)
		}
		if edges[i].EdgeType == "" {
			edges[i].EdgeType = "network"
		}
		if edges[i].Direction == "" {
			edges[i].Direction = "uni"
		}
		if !allowedEdgeTypes[edges[i].EdgeType] {
			return nil, nil, nil, fmt.Errorf("无效的 edge_type: %s", edges[i].EdgeType)
		}
		if !allowedDirection[edges[i].Direction] {
			return nil, nil, nil, fmt.Errorf("无效的 direction: %s", edges[i].Direction)
		}
		if edges[i].Weight == 0 {
			edges[i].Weight = 1
		}
		if edges[i].Weight < 0 {
			return nil, nil, nil, fmt.Errorf("edge weight 不能为负数")
		}
		if edges[i].Risk < 0 || edges[i].Risk > 100 {
			return nil, nil, nil, fmt.Errorf("edge risk 必须在 0~100 之间")
		}

		edgeKey := fmt.Sprintf("%s|%s|%s|%s", edges[i].FromNodeKey, edges[i].ToNodeKey, edges[i].EdgeType, edges[i].Direction)
		if edgeSet[edgeKey] {
			return nil, nil, nil, fmt.Errorf("重复边定义: %s -> %s", edges[i].FromNodeKey, edges[i].ToNodeKey)
		}
		edgeSet[edgeKey] = true

		adj[edges[i].FromNodeKey] = append(adj[edges[i].FromNodeKey], edges[i].ToNodeKey)
		adj[edges[i].ToNodeKey] = append(adj[edges[i].ToNodeKey], edges[i].FromNodeKey)
	}

	if len(nodes) > 1 && !isGraphConnected(nodes, adj) {
		return nil, nil, nil, fmt.Errorf("拓扑图必须连通")
	}

	derived := extractProductIDsFromNodes(nodes)
	if len(derived) == 0 {
		derived = productIDs
	}
	if len(derived) == 0 {
		return nil, nil, nil, fmt.Errorf("拓扑中未找到可用设备节点")
	}

	return derived, nodes, edges, nil
}

func isGraphConnected(nodes []models.TopoNode, adj map[string][]string) bool {
	if len(nodes) == 0 {
		return true
	}

	visited := make(map[string]bool)
	queue := []string{nodes[0].NodeKey}
	visited[nodes[0].NodeKey] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, nxt := range adj[cur] {
			if !visited[nxt] {
				visited[nxt] = true
				queue = append(queue, nxt)
			}
		}
	}

	return len(visited) == len(nodes)
}

func buildDefaultGraphFromProducts(productIDs []int) ([]models.TopoNode, []models.TopoEdge) {
	nodes := make([]models.TopoNode, 0, len(productIDs))
	edges := make([]models.TopoEdge, 0)

	for i, pid := range productIDs {
		pidCopy := uint(pid)
		key := fmt.Sprintf("product-%d-%d", pid, i+1)
		nodes = append(nodes, models.TopoNode{
			NodeKey:     key,
			NodeType:    "hardware",
			Name:        fmt.Sprintf("设备-%d", pid),
			ProductID:   &pidCopy,
			Criticality: "normal",
			Zone:        "default",
			Layer:       i + 1,
		})

		if i > 0 {
			edges = append(edges, models.TopoEdge{
				FromNodeKey: nodes[i-1].NodeKey,
				ToNodeKey:   key,
				EdgeType:    "network",
				Direction:   "uni",
				Weight:      1,
			})
		}
	}

	return nodes, edges
}

func buildDefaultEdgesFromNodes(nodes []models.TopoNode) []models.TopoEdge {
	if len(nodes) < 2 {
		return []models.TopoEdge{}
	}

	sortedNodes := make([]models.TopoNode, len(nodes))
	copy(sortedNodes, nodes)
	sort.Slice(sortedNodes, func(i, j int) bool {
		if sortedNodes[i].Layer == sortedNodes[j].Layer {
			return sortedNodes[i].NodeKey < sortedNodes[j].NodeKey
		}
		return sortedNodes[i].Layer < sortedNodes[j].Layer
	})

	edges := make([]models.TopoEdge, 0, len(sortedNodes)-1)
	for i := 0; i < len(sortedNodes)-1; i++ {
		edges = append(edges, models.TopoEdge{
			FromNodeKey: sortedNodes[i].NodeKey,
			ToNodeKey:   sortedNodes[i+1].NodeKey,
			EdgeType:    "network",
			Direction:   "uni",
			Weight:      1,
		})
	}
	return edges
}

func extractProductIDsFromNodes(nodes []models.TopoNode) []int {
	sortedNodes := make([]models.TopoNode, len(nodes))
	copy(sortedNodes, nodes)
	sort.Slice(sortedNodes, func(i, j int) bool {
		if sortedNodes[i].Layer == sortedNodes[j].Layer {
			return sortedNodes[i].ID < sortedNodes[j].ID
		}
		return sortedNodes[i].Layer < sortedNodes[j].Layer
	})

	ids := make([]int, 0)
	seen := make(map[int]bool)
	for _, n := range sortedNodes {
		if n.ProductID == nil {
			continue
		}
		id := int(*n.ProductID)
		if !seen[id] {
			ids = append(ids, id)
			seen[id] = true
		}
	}
	return ids
}

func (s *Service) validateProductIDs(productIDs []int) error {
	if len(productIDs) == 0 {
		return nil
	}

	var existingProducts []models.Product
	s.DB.Find(&existingProducts)

	validProductIDs := make([]int, 0, len(existingProducts))
	for _, p := range existingProducts {
		validProductIDs = append(validProductIDs, int(p.ID))
	}

	if !utils.ValidateIntSlice(productIDs, validProductIDs) {
		return fmt.Errorf("部分产品不存在")
	}

	return nil
}
