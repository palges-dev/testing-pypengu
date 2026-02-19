package collector

import (
	"main/internal/models"
)

func CollectUsers(verbose bool) []models.UserInfo {
	return []models.UserInfo{}
}

func CollectGroups(verbose bool) []models.GroupInfo {
	return []models.GroupInfo{}
}

func CollectSudo(verbose bool) []models.SudoEntry {
	return []models.SudoEntry{}
}

func CollectSUID(verbose bool) []models.SuidBinary {
	return []models.SuidBinary{}
}

func CollectDocker(verbose bool) models.DockerInfo {
	return models.DockerInfo{}
}

func CollectServices(verbose bool) []models.ServiceInfo {
	return []models.ServiceInfo{}
}

func CollectCron(verbose bool) []models.CronJob {
	return []models.CronJob{}
}

func CollectKernel(verbose bool) models.KernelInfo {
	return models.KernelInfo{}
}
