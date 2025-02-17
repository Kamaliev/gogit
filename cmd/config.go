package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GitProfile struct {
	Name   string `required:"true" json:"name"`
	Email  string `required:"true" json:"email"`
	SshKey string `json:"ssh_key"`
}

type Config struct {
	Profiles map[string]GitProfile `json:"profiles"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gogitconfig")
}

func (c Config) load() (Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{Profiles: make(map[string]GitProfile)}, nil
		}
		return Config{Profiles: make(map[string]GitProfile)}, err
	}
	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{Profiles: make(map[string]GitProfile)}, err
	}

	return config, nil
}

func (c Config) save() error {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func (profile *GitProfile) GenerateSSH() {
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".gossh", profile.Name)

	err := os.MkdirAll(sshDir, os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка при создании директории:", err)
		return
	}
	sshPath := filepath.Join(sshDir, "id_rsa")
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", sshPath, "-N", "")

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	profile.SshKey = sshPath

}

func (profile *GitProfile) Activate() error {
	_, err := isGitRepository()
	if err != nil {
		return err
	}

	err = exec.Command("git", "config", "user.Name", profile.Name).Run()
	if err != nil {
		return fmt.Errorf("не удалось установить имя пользователя в локальном конфиге: %w", err)
	}
	err = exec.Command("git", "config", "user.Email", profile.Email).Run()
	if err != nil {
		return fmt.Errorf("не удалось установить Email: %w", err)
	}

	err = exec.Command("git", "config", "core.sshCommand", fmt.Sprintf("ssh -i %s", profile.SshKey)).Run()
	if err != nil {
		return fmt.Errorf("не удалось установить SSH: %w", err)
	}

	fmt.Printf("🔹 Активирован Git-профиль:\n")
	fmt.Println("👤  Имя:  ", profile.Name)
	fmt.Println("✉️  Email:", profile.Email)
	if profile.SshKey != "" {
		fmt.Println("🔑  SSH-ключ:", profile.SshKey)
	} else {
		fmt.Println("🔑  SSH-ключ: не установлен")
	}
	return nil
}

func isGitRepository() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("директория не является репозиторием Git: %w", err)
	}
	return true, nil
}
