package routes

import (
	"fmt"
	"net/http"
	"sample_server/db"
	"sample_server/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func getBlastRadius(c *gin.Context) {
	id := c.Param("id")
	maxHops := c.Param("hop_length")
	label := c.Param("label")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param: txn_id"})
		return
	}

	session := db.Driver.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c)

	query := fmt.Sprintf(`
	// Step 1: Traverse nodes + relationships
	MATCH p = (n:Transaction {id: $id})-[*1..%s]-(m)
	WITH collect(DISTINCT nodes(p)) AS nodePaths, collect(DISTINCT relationships(p)) AS relPaths
	UNWIND nodePaths AS np
	UNWIND np AS node
	WITH collect(DISTINCT node) AS nodes, relPaths

	UNWIND relPaths AS rp
	UNWIND rp AS rel
	WITH nodes, collect(DISTINCT rel) AS relationships

	// Step 2: Drop old projection
	CALL gds.graph.exists('subgraph') YIELD exists
	WITH nodes, relationships, exists
	CALL apoc.do.when(
		exists,
		'CALL gds.graph.drop("subgraph") YIELD graphName RETURN graphName',
		'RETURN null AS graphName',
		{}
	) YIELD value
	WITH nodes, relationships

	// Step 3: Create new projection
	CALL gds.graph.project.cypher(
		'subgraph',
		'UNWIND $nodes AS n RETURN id(n) AS id',
		'
		  UNWIND $rels AS r
		  RETURN id(startNode(r)) AS source, id(endNode(r)) AS target, type(r) AS type
		',
		{ parameters: { nodes: nodes, rels: relationships } }
	)
	YIELD graphName

	// Step 4: Return nodes + relationships
	UNWIND nodes AS n
	OPTIONAL MATCH (n)-[r]->(m)
	WHERE r IN relationships
	RETURN 
		collect(DISTINCT {id: id(n), label: head(labels(n))}) AS nodes,
		collect(DISTINCT {source: id(startNode(r)), target: id(endNode(r)), type: type(r)}) AS rels
`, maxHops)

	params := map[string]interface{}{
		"id":         id,
		"hop_length": maxHops,
		"label":      label,
	}

	result, err := session.Run(c, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blast models.BlastRad
	if result.Next(c) {
		rec := result.Record()

		// nodes
		rawNodes := rec.Values[0].([]interface{})
		for _, r := range rawNodes {
			n := r.(map[string]interface{})
			blast.Nodes = append(blast.Nodes, models.Node{
				ID:    fmt.Sprintf("%v", n["id"]),
				Label: n["label"].(string),
			})
		}

		// relationships
		rawRels := rec.Values[1].([]interface{})
		for _, r := range rawRels {
			rel := r.(map[string]interface{})
			blast.Relationships = append(blast.Relationships, models.Relationship{
				Source: fmt.Sprintf("%v", rel["source"]),
				Target: fmt.Sprintf("%v", rel["target"]),
				Type:   rel["type"].(string),
			})
		}
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blast)
}

func getLouvianCommunities(c *gin.Context) {
	session := db.Driver.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(c)

	query := `
	CALL gds.louvain.stream('subgraph')
	YIELD nodeId, communityId
	WITH gds.util.asNode(nodeId) AS node, communityId
	RETURN communityId, node.id AS node_id
	ORDER BY communityId, node_id
	`

	result, err := session.Run(c, query, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commMap := make(map[int][]string)

	for result.Next(c) {
		rec := result.Record()
		communityID := int(rec.Values[0].(int64))
		nodeID := fmt.Sprintf("%v", rec.Values[1])

		commMap[communityID] = append(commMap[communityID], nodeID)
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var communities []models.Community
	for cid, nodes := range commMap {
		communities = append(communities, models.Community{
			ID:    fmt.Sprintf("%d", cid),
			Nodes: nodes,
		})
	}

	c.JSON(http.StatusOK, models.Communities{Communities: communities})
}
