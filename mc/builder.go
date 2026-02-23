// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package graph

import (
	"fmt"
	"pypengu/internal/models"
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

	for _, user := range users {
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
				Risk:   "medium",
				Properties: map[string]interface{}{},
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
