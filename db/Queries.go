package db
const(
	QueryLouvian = `
	WITH $id AS id, toInteger($hop_length) AS hop_length

// Step 1: Traverse the subgraph
MATCH (start:Operator {id: id})
CALL apoc.path.subgraphNodes(start, {
  maxLevel: hop_length,
  relationshipFilter: "CONNECTED_TO|PERFORMED_BY|WORKS_ON|FOR_PERSON",
  bfs: true
}) YIELD node
WITH collect(DISTINCT node) AS nodes

// Step 2: Drop existing projection if it exists
CALL gds.graph.exists('subgraph') YIELD exists
WITH nodes, exists
CALL apoc.do.when(
  exists,
  'CALL gds.graph.drop("subgraph") YIELD graphName RETURN graphName',
  'RETURN null AS graphName',
  {}
) YIELD value
WITH nodes

// Step 3: Project subgraph into GDS
CALL gds.graph.project.cypher(
  'subgraph',
  'UNWIND $nodes AS n RETURN id(n) AS id',
  '
    UNWIND $nodes AS n
    MATCH (n)-[r]-(m)
    WHERE m IN $nodes
    RETURN id(n) AS source, id(m) AS target
  ',
  { parameters: { nodes: nodes } }
)
YIELD graphName AS createdGraph
WITH createdGraph

// Step 4: Run Louvain algorithm
CALL gds.louvain.stream(createdGraph)
YIELD nodeId, communityId
WITH gds.util.asNode(nodeId) AS node, communityId

// Step 5: Return results
RETURN communityId, labels(node)[0] AS label, node.id AS node_id
ORDER BY communityId, label, node_id

	`

	QueryLabelprop = `
	WITH $id AS id, toInteger($hop_length) AS hop_length

// Step 1: Traverse the subgraph
MATCH (start:Operator {id: id})
CALL apoc.path.subgraphNodes(start, {
  maxLevel: hop_length,
  relationshipFilter: "CONNECTED_TO|PERFORMED_BY|WORKS_ON|FOR_PERSON",
  bfs: true
}) YIELD node
WITH collect(DISTINCT node) AS nodes

// Step 2: Drop existing projection if it exists
CALL gds.graph.exists('subgraph') YIELD exists
WITH nodes, exists
CALL apoc.do.when(
  exists,
  'CALL gds.graph.drop("subgraph") YIELD graphName RETURN graphName',
  'RETURN null AS graphName',
  {}
) YIELD value
WITH nodes

// Step 3: Project subgraph into GDS
CALL gds.graph.project.cypher(
  'subgraph',
  'UNWIND $nodes AS n RETURN id(n) AS id',
  '
    UNWIND $nodes AS n
    MATCH (n)-[r]-(m)
    WHERE m IN $nodes
    RETURN id(n) AS source, id(m) AS target
  ',
  { parameters: { nodes: nodes } }
)
YIELD graphName AS createdGraph
WITH createdGraph

// Step 4: Run Louvain algorithm
CALL gds.labelPropagation.stream(createdGraph)
YIELD nodeId, communityId
WITH gds.util.asNode(nodeId) AS node, communityId

// Step 5: Return results
RETURN communityId, labels(node)[0] AS label, node.id AS node_id
ORDER BY communityId, label, node_id

	`
)