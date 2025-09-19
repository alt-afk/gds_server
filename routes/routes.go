package routes

import (
	"fmt"
	"net/http"
	"sample_server/db"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func RegisterRoutes(r *gin.Engine) {
	gdsGroupOld := r.Group("/gds")
	{
		gdsGroupOld.GET("/louvian/:id/:hop_length", getLouvianCluster)
		gdsGroupOld.GET("/labelprop/:id/:hop_length", getLabelpropCluster)
	}

	gdsGroupNew := r.Group("/gds/v2")
	{
		gdsGroupNew.GET("/:algo", func(c *gin.Context) {
			algo := c.Param("algo")

			switch algo {
			case "louvain":
				getLouvianCommunities(c)
			case "labelprop":
				getLabelpropCluster(c)
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("unsupported algorithm: %s", algo),
				})
			}
		})
	}

	r.GET("/blastrad/:label/:id/:hop_length", getBlastRadius)

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

	// label := c.Param("label")
	id := c.Param("id")
	hop := c.Param("hop_length")

	params := map[string]interface{}{
		"id":         id,
		"hop_length": hop,
	}

	result, err := session.ExecuteWrite(c, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(c, db.QueryLouvian, params)
		if err != nil {
			return nil, err
		}

		// Map: communityId => list of {label, id}
		clusters := make(map[int][]map[string]interface{})

		for records.Next(c) {
			record := records.Record()
			communityID := int(record.Values[0].(int64))
			label := record.Values[1].(string)
			id := record.Values[2]

			nodeInfo := map[string]interface{}{
				"label":   label,
				"node_id": id,
			}

			clusters[communityID] = append(clusters[communityID], nodeInfo)

			fmt.Printf("Community %d: %+v\n", communityID, nodeInfo)
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

	// label := c.Param("label")
	id := c.Param("id")
	hop := c.Param("hop_length")

	params := map[string]interface{}{
		"id":         id,
		"hop_length": hop,
	}

	result, err := session.ExecuteWrite(c, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(c, db.QueryLabelprop, params)
		if err != nil {
			return nil, err
		}

		// Map: communityId => list of {label, id}
		clusters := make(map[int][]map[string]interface{})

		for records.Next(c) {
			record := records.Record()
			communityID := int(record.Values[0].(int64))
			label := record.Values[1].(string)
			id := record.Values[2]

			nodeInfo := map[string]interface{}{
				"label":   label,
				"node_id": id,
			}

			clusters[communityID] = append(clusters[communityID], nodeInfo)

			fmt.Printf("Community %d: %+v\n", communityID, nodeInfo)
		}

		return clusters, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clusters": result})
}
