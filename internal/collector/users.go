// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"bufio"
	"os"
	"os/user"
	"pypengu/internal/models"
	"strings"
)

func CollectUsers(verbose bool) []models.UserInfo {
	users := []models.UserInfo{}
	
	currentUser, _ := user.Current()
	currentUsername := currentUser.Username
	
	file, err := os.Open("/etc/passwd")
	if err != nil {
		if verbose {
			println("[-] Error reading /etc/passwd:", err.Error())
		}
		return users
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}
		
		username := parts[0]
		uid := parts[2]
		gid := parts[3]
		home := parts[5]
		shell := parts[6]
		
		isCurrent := username == currentUsername
		isRoot := uid == "0"
		
		userInfo := models.UserInfo{
			Username:  username,
			UID:       uid,
			GID:       gid,
			Shell:     shell,
			Home:      home,
			Groups:    []string{},
			IsCurrent: isCurrent,
			IsRoot:    isRoot,
		}
		
		users = append(users, userInfo)
	}
	
	if verbose {
		println("    Found", len(users), "users")
	}
	
	return users
}
