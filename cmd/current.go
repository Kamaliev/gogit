/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func readGitConfigFile(filePath string) (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var block string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			block = line[1 : len(line)-1]
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[fmt.Sprintf("%s.%s", block, key)] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

func getCurrentGitProfile() (GitProfile, string, error) {
	var configPaths = []string{
		".git/config",    // Локальный конфиг
		"~/.gitconfig",   // Глобальный конфиг
		"/etc/gitconfig", // Системный конфиг
	}

	var name, email, sshKey string
	var usePath string
	for _, path := range configPaths {
		expandedPath := path
		if path == "~/.gitconfig" {
			expandedPath = os.Getenv("HOME") + "/.gitconfig"
		} else if path == "/etc/gitconfig" {
			expandedPath = "/etc/gitconfig"
		}

		config, err := readGitConfigFile(expandedPath)
		if err != nil {
			continue
		}

		if val, ok := config["user.Name"]; ok {
			name = val
		}
		if val, ok := config["user.Email"]; ok {
			email = val
		}

		if val, ok := config["core.sshCommand"]; ok && strings.Contains(val, "-i") {
			parts := strings.Split(val, " ")
			for i, part := range parts {
				if part == "-i" && i+1 < len(parts) {
					sshKey = parts[i+1]
					break
				}
			}
		}
		usePath = path
		break
	}

	return GitProfile{Name: name, Email: email, SshKey: sshKey}, usePath, nil
}

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Текущий конфиг",
	Run: func(cmd *cobra.Command, args []string) {
		profile, path, err := getCurrentGitProfile()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("🔹 Текущий Git-профиль: %s\n", path)
		fmt.Println("👤  Имя:  ", profile.Name)
		fmt.Println("✉️  Email:", profile.Email)
		if profile.SshKey != "" {
			fmt.Println("🔑  SSH-ключ:", profile.SshKey)
		} else {
			fmt.Println("🔑  SSH-ключ: не установлен")
		}
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
