package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// listProfiles выводит список доступных профилей
func listProfiles() {
	config, err := Config{}.load()
	if err != nil {
		fmt.Println("Ошибка чтения конфига:", err)
		return
	}
	if len(config.Profiles) == 0 {
		fmt.Println("\033[1;31mGoGit - Пусто:(\033[0m")
		return
	}

	// Заголовок
	fmt.Println("\033[1;32mGoGit - Доступные профили:\033[0m")
	fmt.Println("--------------------------------------------------")
	// Вывод каждого профиля с его деталями
	for _, profile := range config.Profiles {
		fmt.Printf("Профиль: %-20s | Email: %-30s | SSH путь: %s\n", profile.Name, profile.Email, profile.SshKey)
	}
	fmt.Println("--------------------------------------------------")
}

// useCmd представляет команду use
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Активировать профиль Git",
	Run: func(cmd *cobra.Command, args []string) {
		// Проверяем флаг --list
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			// Если флаг --list указан, выводим список профилей
			listProfiles()
			return
		}

		// Если не указан флаг, активируем профиль
		if len(args) == 0 {
			fmt.Println("Нужно указать профиль")
			return
		}

		config, err := Config{}.load()
		if err != nil {
			fmt.Println(err)
			return
		}

		profileName := args[0]

		if findProfile, ok := config.Profiles[profileName]; ok {
			err = findProfile.Activate()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else {
			fmt.Printf("Профиль '%s' не найден\n", profileName)
		}
	},
}

func init() {
	// Добавляем команду use
	rootCmd.AddCommand(useCmd)

	useCmd.Flags().BoolP("list", "l", false, "Список Git профилей")
}
