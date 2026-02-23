// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"bufio"
	"os"
	"pypengu/internal/models"
	"strings"
)

func CollectGroups(verbose bool) []models.GroupInfo {
	groups := []models.GroupInfo{}
	
	privilegedGroups := map[string]bool{
		"root":   true,
		"sudo":   true,
		"wheel":  true,
		"admin":  true,
		"docker": true,
		"disk":   true,
		"lxd":    true,
		"shadow": true,
	}
	
	file, err := os.Open("/etc/group")
	if err != nil {
		if verbose {
			println("[-] Error reading /etc/group:", err.Error())
		}
		return groups
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		
		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}
		
		groupName := parts[0]
		gid := parts[2]
		membersStr := parts[3]
		
		var members []string
		if membersStr != "" {
			members = strings.Split(membersStr, ",")
		}
		
		isPrivileged := privilegedGroups[groupName]
		
		groupInfo := models.GroupInfo{
			Name:       groupName,
			GID:        gid,
			Members:    members,
			Privileged: isPrivileged,
		}
		
		groups = append(groups, groupInfo)
	}
	
	if verbose {
		println("    Found", len(groups), "groups")
	}
	
	return groups
}
