package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Profile представляет структуру профиля
type Profile struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

// ProfileManager управляет операциями с профилями
type ProfileManager struct {
	dir string
}

// NewProfileManager создает новый менеджер профилей
func NewProfileManager(dir string) *ProfileManager {
	return &ProfileManager{dir: dir}
}

// Create создает новый профиль
func (pm *ProfileManager) Create(name, user, project string) error {
	if name == "" || user == "" || project == "" {
		return fmt.Errorf("все поля (name, user, project) обязательны для заполнения")
	}

	filename := filepath.Join(pm.dir, name+".yaml")
	
	// Проверяем, существует ли уже профиль
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("профиль '%s' уже существует", name)
	}

	profile := Profile{
		User:    user,
		Project: project,
	}

	data, err := yaml.Marshal(&profile)
	if err != nil {
		return fmt.Errorf("ошибка при создании YAML: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении профиля: %w", err)
	}

	fmt.Printf("Профиль '%s' успешно создан\n", name)
	return nil
}

// Get получает профиль по имени
func (pm *ProfileManager) Get(name string) error {
	if name == "" {
		return fmt.Errorf("имя профиля обязательно")
	}

	filename := filepath.Join(pm.dir, name+".yaml")
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("профиль '%s' не найден", name)
		}
		return fmt.Errorf("ошибка при чтении профиля: %w", err)
	}

	var profile Profile
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("ошибка при разборе YAML: %w", err)
	}

	fmt.Printf("Профиль: %s\n", name)
	fmt.Printf("  User:    %s\n", profile.User)
	fmt.Printf("  Project: %s\n", profile.Project)
	return nil
}

// List выводит список всех профилей
func (pm *ProfileManager) List() error {
	files, err := ioutil.ReadDir(pm.dir)
	if err != nil {
		return fmt.Errorf("ошибка при чтении директории: %w", err)
	}

	var profiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			name := strings.TrimSuffix(file.Name(), ".yaml")
			profiles = append(profiles, name)
		}
	}

	if len(profiles) == 0 {
		fmt.Println("Профили не найдены")
		return nil
	}

	fmt.Println("Доступные профили:")
	for _, profile := range profiles {
		fmt.Printf("  - %s\n", profile)
	}
	return nil
}

// Delete удаляет профиль по имени
func (pm *ProfileManager) Delete(name string) error {
	if name == "" {
		return fmt.Errorf("имя профиля обязательно")
	}

	filename := filepath.Join(pm.dir, name+".yaml")
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("профиль '%s' не найден", name)
	}

	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("ошибка при удалении профиля: %w", err)
	}

	fmt.Printf("Профиль '%s' успешно удален\n", name)
	return nil
}

// PrintHelp выводит справку по командам
func PrintHelp() {
	fmt.Println("CLI для работы с профилями")
	fmt.Println("\nИспользование:")
	fmt.Println("  ./mws <команда> [флаги]")
	fmt.Println("\nДоступные команды:")
	fmt.Println("  profile create  Создать новый профиль")
	fmt.Println("    --name        Имя профиля (обязательно)")
	fmt.Println("    --user        Имя пользователя (обязательно)")
	fmt.Println("    --project     Название проекта (обязательно)")
	fmt.Println()
	fmt.Println("  profile get     Получить информацию о профиле")
	fmt.Println("    --name        Имя профиля (обязательно)")
	fmt.Println()
	fmt.Println("  profile list    Вывести список всех профилей")
	fmt.Println()
	fmt.Println("  profile delete  Удалить профиль")
	fmt.Println("    --name        Имя профиля (обязательно)")
	fmt.Println()
	fmt.Println("  help            Вывести эту справку")
	fmt.Println()
	fmt.Println("Примеры:")
	fmt.Println("  ./mws profile create --name=test --user=example --project=new-project")
	fmt.Println("  ./mws profile get --name=test")
	fmt.Println("  ./mws profile list")
	fmt.Println("  ./mws profile delete --name=test")
	fmt.Println("  ./mws help")
}

// parseFlags парсит флаги из аргументов командной строки
func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
			} else {
				flags[parts[0]] = ""
			}
		}
	}
	return flags
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: не указана команда")
		fmt.Println("Используйте './mws help' для получения справки")
		os.Exit(1)
	}

	// Получаем текущую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Ошибка при получении текущей директории: %v\n", err)
		os.Exit(1)
	}

	pm := NewProfileManager(currentDir)

	// Определяем команду
	command := os.Args[1]
	
	switch command {
	case "help":
		PrintHelp()

	case "profile":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: не указана подкоманда profile")
			fmt.Println("Используйте './mws help' для получения справки")
			os.Exit(1)
		}

		subcommand := os.Args[2]
		flags := parseFlags(os.Args[3:])

		switch subcommand {
		case "create":
			name := flags["name"]
			user := flags["user"]
			project := flags["project"]
			
			if err := pm.Create(name, user, project); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
				os.Exit(1)
			}

		case "get":
			name := flags["name"]
			
			if err := pm.Get(name); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
				os.Exit(1)
			}

		case "list":
			if err := pm.List(); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
				os.Exit(1)
			}

		case "delete":
			name := flags["name"]
			
			if err := pm.Delete(name); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
				os.Exit(1)
			}

		default:
			fmt.Printf("Ошибка: неизвестная подкоманда profile '%s'\n", subcommand)
			fmt.Println("Используйте './mws help' для получения справки")
			os.Exit(1)
		}

	default:
		fmt.Printf("Ошибка: неизвестная команда '%s'\n", command)
		fmt.Println("Используйте './mws help' для получения справки")
		os.Exit(1)
	}
}