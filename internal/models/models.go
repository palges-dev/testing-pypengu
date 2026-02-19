package models

type GraphData struct {
	Metadata Metadata `json:"metadata"`
	Nodes    []Node   `json:"nodes"`
	Edges    []Edge   `json:"edges"`
	Stats    Stats    `json:"stats"`
}

type Metadata struct {
	Hostname      string `json:"hostname"`
	IP            string `json:"ip,omitempty"`
	OS            string `json:"os"`
	Kernel        string `json:"kernel"`
	Arch          string `json:"arch"`
	CollectedAt   string `json:"collected_at"`
	Collector     string `json:"collector"`
	Version       string `json:"collector_version"`
	CollectedAs   string `json:"collected_as"`
	UID           string `json:"uid"`
}

type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

type Edge struct {
	ID         string                 `json:"id"`
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Risk       string                 `json:"risk"`
	Properties map[string]interface{} `json:"properties"`
}

type Stats struct {
	TotalNodes   int `json:"total_nodes"`
	TotalEdges   int `json:"total_edges"`
	PathsToRoot  int `json:"paths_to_root"`
}

type SudoEntry struct {
	User     string
	Command  string
	NoPasswd bool
	RunAs    string
	Entry    string
}

type SuidBinary struct {
	Path     string
	Owner    string
	GTFOBin  bool
	GTFOUrl  string
}

type ServiceInfo struct {
	Name      string
	Path      string
	RunAs     string
	State     string
	Writable  bool
	WritableBy string
}

type CronJob struct {
	Path      string
	Schedule  string
	Command   string
	RunAs     string
	Writable  bool
	WritableBy string
}

type UserInfo struct {
	Username  string
	UID       string
	GID       string
	Shell     string
	Home      string
	Groups    []string
	IsCurrent bool
	IsRoot    bool
}

type GroupInfo struct {
	Name       string
	GID        string
	Members    []string
	Privileged bool
}

type DockerInfo struct {
	HasDockerGroup  bool
	SocketWritable  bool
	SocketPath      string
}

type KernelInfo struct {
	Version     string
	CVEs        []CVEInfo
}

type CVEInfo struct {
	ID          string
	Description string
	Reference   string
}
