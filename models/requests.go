package models

type BlastRad struct {
	Nodes         []Node         `json:"nodes"`         // List of nodes
	Relationships []Relationship `json:"relationships"` // List of relationships
}

type Node struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	// Properties map[string]string `json:"properties"`
}

type Relationship struct {
	Source     string            `json:"source"` // Source node ID
	Target     string            `json:"target"` // Target node ID
	Type       string            `json:"type"`   // Relationship type
	Properties map[string]string `json:"properties"`
}

type Communities struct {
	Communities []Community `json:"communities"` // List of communities
}

type Community struct {
	ID    string   `json:"id"`    // Community ID
	Nodes []string `json:"nodes"` // List of node IDs in the community
}
