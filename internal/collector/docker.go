package collector

import (
	"os"
	"os/user"
	"pypengu/internal/models"
)

func CollectDocker(verbose bool) models.DockerInfo {
	info := models.DockerInfo{
		SocketPath: "/var/run/docker.sock",
	}
	
	currentUser, err := user.Current()
	if err != nil {
		return info
	}
	
	groups, err := currentUser.GroupIds()
	if err != nil {
		return info
	}
	
	dockerGroup, err := user.LookupGroup("docker")
	if err == nil {
		for _, gid := range groups {
			if gid == dockerGroup.Gid {
				info.HasDockerGroup = true
				break
			}
		}
	}
	
	sockInfo, err := os.Stat("/var/run/docker.sock")
	if err == nil {
		mode := sockInfo.Mode()
		if mode.Perm()&0002 != 0 {
			info.SocketWritable = true
		} else if mode.Perm()&0020 != 0 {
			info.SocketWritable = true
		}
	}
	
	if verbose {
		if info.HasDockerGroup {
			println("    User is in docker group")
		}
		if info.SocketWritable {
			println("    Docker socket is writable")
		}
	}
	
	return info
}
