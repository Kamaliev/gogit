package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func ReadCmdLine(consoleMsg string, validate *func(string) (bool, error)) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(consoleMsg)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if validate != nil {
			isLine, _ := (*validate)(line)
			if !isLine {
				continue
			}
		}

		return line, err
	}

}

func validateEmail(s string) (bool, error) {
	if !strings.Contains(s, "@") {
		fmt.Println("Некорректный Email")
		return false, nil
	}
	return true, nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить гит профиль",
	Run: func(cmd *cobra.Command, args []string) {
		emailValidator := validateEmail
		email, _ := ReadCmdLine("Email: ", &emailValidator)
		username, _ := ReadCmdLine("Username: ", nil)

		config, _ := Config{}.load()

		gitProfile := GitProfile{Name: username, Email: email}
		gitProfile.GenerateSSH()
		if config.Profiles == nil {
			config.Profiles = make(map[string]GitProfile)
		}
		config.Profiles[username] = gitProfile
		err := config.save()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("\n🔹 Создан Git-профиль:")
		fmt.Println("👤  Имя:  ", gitProfile.Name)
		fmt.Println("✉️  Email:", gitProfile.Email)
		if gitProfile.SshKey != "" {
			fmt.Println("🔑  SSH-ключ:", gitProfile.SshKey)
		} else {
			fmt.Println("🔑  SSH-ключ: не установлен")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
