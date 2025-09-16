package routes

import (
	"sample_server/db"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func RegisterRoutes(r *gin.Engine) {
	gdsGroup := r.Group("/gds")
	{
		gdsGroup.GET("/louvian/:label/:id/:intensity", getLouvianCluster)
		gdsGroup.GET("/labelprop/:label/:id/:intensity", getLabelpropCluster)
	}
	r.DELETE("/reset", resetDatabase)
}

func resetDatabase(c *gin.Context) {

	session := db.Driver.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c)

	// delete database
	query := "MATCH (n) DETACH DELETE n"
	_, err := session.ExecuteWrite(c, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(c, query, nil)
		return nil, err
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "Database reset successful"})
}

func getLouvianCluster(c *gin.Context) {
	session := db.Driver.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c)

	label := c.Param("label")
	id := c.Param("id")
	hop := c.Param("intensity")

query := fmt.Sprintf(`
	MATCH (start:%s {id: $id})
	CALL apoc.path.subgraphNodes(start, {
		maxLevel: toInteger($hop_length),
		relationshipFilter: "CONNECTED_TO|WORKS_ON|USED_MACHINE|INTRODUCED_BY",
		labelFilter: ">%s"
	})
	YIELD node
	WITH collect(node) AS nodes

	CALL gds.graph.exists('subgraph') YIELD exists
	WITH exists, nodes
	CALL apoc.do.when(
		exists,
		'CALL gds.graph.drop("subgraph", false)',
		'RETURN null',
		{}
	) YIELD value
	WITH nodes

	CALL gds.graph.project.cypher(
		'subgraph',
		'UNWIND $nodes AS n RETURN id(n) AS id',
		'
		UNWIND $nodes AS n
		MATCH (n)-[r]->(m)
		WHERE m IN $nodes
		RETURN id(n) AS source, id(m) AS target
		',
		{parameters: {nodes: nodes}}
	)
	YIELD graphName AS createdGraph

	CALL {
		WITH createdGraph
		CALL gds.louvain.stream(createdGraph)
		YIELD nodeId, communityId
		RETURN gds.util.asNode(nodeId) AS node, communityId, createdGraph AS gName
	}
	WITH node, communityId, gName
	CALL gds.graph.drop(gName, false) YIELD graphName AS droppedGraphName
	RETURN node, communityId
`, label, label)




	params := map[string]interface{}{
		"id":         id,
		"hop_length": hop,
	}

	result, err := session.ExecuteWrite(c, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(c, query, params)
		if err != nil {
			return nil, err
		}

		// Map: communityId => list of node properties
		clusters := make(map[int][]map[string]interface{})

		for records.Next(c) {
			record := records.Record()
			node := record.Values[0].(neo4j.Node)
			communityID := int(record.Values[1].(int64))

			clusters[communityID] = append(clusters[communityID], node.Props)

			fmt.Println(communityID, node.Props)
		}

		return clusters, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clusters": result})
}


func getLabelpropCluster(c *gin.Context) {
	session := db.Driver.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c)

	label := c.Param("label")
	id := c.Param("id")
	hop := c.Param("intensity")

	query := fmt.Sprintf(`
	MATCH (start:%s {id: $id})
	CALL apoc.path.subgraphNodes(start, {
		maxLevel: toInteger($hop_length),
		relationshipFilter: "CONNECTED_TO|WORKS_ON|USED_MACHINE|FOR_PERSON|PERFORMED_BY",
		labelFilter: ">%s"
	})
	YIELD node
	WITH collect(node) AS nodes

	CALL gds.graph.project.cypher(
		'subgraph',
		'UNWIND $nodes AS n RETURN id(n) AS id',
		'
		UNWIND $nodes AS n
		MATCH (n)-[r]->(m)
		WHERE m IN $nodes
		RETURN id(n) AS source, id(m) AS target
		',
		{parameters: {nodes: nodes}}
	)
	YIELD graphName

	CALL {
		WITH graphName
		CALL gds.labelPropagation.stream(graphName)
		YIELD nodeId, communityId
		RETURN gds.util.asNode(nodeId) AS node, communityId, graphName
	}
	WITH node, communityId, graphName

	CALL gds.graph.drop(graphName, false) YIELD graphName AS droppedGraphName

	RETURN node, communityId
`, label, label)


	params := map[string]interface{}{
		"id":         id,
		"hop_length": hop,
	}

	result, err := session.ExecuteWrite(c, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(c, query, params)
		if err != nil {
			return nil, err
		}

		// Map: communityId => list of node properties
		clusters := make(map[int][]map[string]interface{})

		for records.Next(c) {
			record := records.Record()
			node := record.Values[0].(neo4j.Node)
			communityID := int(record.Values[1].(int64))

			clusters[communityID] = append(clusters[communityID], node.Props)
		}

		return clusters, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clusters": result})
}
