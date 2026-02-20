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
