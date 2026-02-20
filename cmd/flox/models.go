package ein

import "github.com/pengu-apm/BloodPengu/"

func NewIngestibleRelationship(source IngestibleEndpoint, target IngestibleEndpoint, rel IngestibleRel) IngestibleRelationship {
	if rel.RelProps == nil {
		rel.RelProps = make(map[string]any)
	}

	return IngestibleRelationship{
		Source:   source,
		Target:   target,
		RelProps: rel.RelProps,
		RelType:  rel.RelType,
	}
}

// IngestMatchStrategy defines how a node should be matched during ingestion—
// either by its object ID (default) or by its name.
type IngestMatchStrategy string

const (
	MatchByID   IngestMatchStrategy = "id"
	MatchByName IngestMatchStrategy = "name"
)

// IngestibleEndpoint represents a node reference in a relationship to be ingested.
type IngestibleEndpoint struct {
	Value   string              // The actual lookup value (either objectid or name)
	MatchBy IngestMatchStrategy // Strategy used to resolve the node
	Kind    graph.Kind          // Optional kind filter to help disambiguate nodes
}

type IngestibleRel struct {
	RelProps map[string]any
	RelType  graph.Kind
}

// IngestibleRelationship represents a directional relationship between two nodes
// intended for ingestion into the graph database. Both endpoints include resolution
// strategies and optional kind filters.
type IngestibleRelationship struct {
	Source IngestibleEndpoint
	Target IngestibleEndpoint

	RelProps map[string]any
	RelType  graph.Kind
}

func (s IngestibleRelationship) IsValid() bool {
	return s.Target.Value != "" && s.Source.Value != "" && s.RelProps != nil
}

type IngestibleSession struct {
	Source    string
	Target    string
	LogonType int
}

type IngestibleNode struct {
	ObjectID    string
	PropertyMap map[string]any
	Labels      []graph.Kind
}

func (s IngestibleNode) IsValid() bool {
	return s.ObjectID != ""
}

type ParsedLocalGroupData struct {
	Relationships []IngestibleRelationship
	Nodes         []IngestibleNode
}

type ParsedDomainTrustData struct {
	TrustRelationships []IngestibleRelationship
	ExtraNodeProps     []IngestibleNode
}

type ParsedGroupMembershipData struct {
	RegularMembers           []IngestibleRelationship
	DistinguishedNameMembers []IngestibleRelationship
}
