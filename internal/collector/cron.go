// Copyright 2026 AdverXarial, byt3n33dl3.
//
// Licensed under the MIT License,
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

package collector

import (
	"bufio"
	"os"
	"path/filepath"
	"pypengu/internal/models"
	"strings"
	"syscall"
)

func CollectCron(verbose bool) []models.CronJob {
	jobs := []models.CronJob{}
	
	cronPaths := []string{
		"/etc/crontab",
		"/etc/cron.d",
		"/var/spool/cron",
		"/var/spool/cron/crontabs",
	}
	
	currentUID := os.Getuid()
	
	for _, cronPath := range cronPaths {
		info, err := os.Stat(cronPath)
		if err != nil {
			continue
		}
		
		if info.IsDir() {
			filepath.Walk(cronPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				
				cronJobs := parseCronFile(path, currentUID)
				jobs = append(jobs, cronJobs...)
				return nil
			})
		} else {
			cronJobs := parseCronFile(cronPath, currentUID)
			jobs = append(jobs, cronJobs...)
		}
	}
	
	if verbose {
		println("    Found", len(jobs), "cron jobs")
	}
	
	return jobs
}

func parseCronFile(path string, currentUID int) []models.CronJob {
	jobs := []models.CronJob{}
	
	fileInfo, err := os.Stat(path)
	if err != nil {
		return jobs
	}
	
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return jobs
	}
	
	writable := false
	writableBy := ""
	
	if stat.Uid == uint32(currentUID) && fileInfo.Mode().Perm()&0200 != 0 {
		writable = true
		writableBy = "owner"
	} else if fileInfo.Mode().Perm()&0020 != 0 {
		writable = true
		writableBy = "group"
	} else if fileInfo.Mode().Perm()&0002 != 0 {
		writable = true
		writableBy = "world"
	}
	
	if !writable {
		return jobs
	}
	
	file, err := os.Open(path)
	if err != nil {
		return jobs
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		
		schedule := strings.Join(parts[0:5], " ")
		command := strings.Join(parts[6:], " ")
		runAs := "root"
		
		if len(parts) >= 6 {
			runAs = parts[5]
		}
		
		job := models.CronJob{
			Path:       path,
			Schedule:   schedule,
			Command:    command,
			RunAs:      runAs,
			Writable:   true,
			WritableBy: writableBy,
		}
		
		jobs = append(jobs, job)
	}
	
	return jobs
}
