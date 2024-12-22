package application

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lazerphp/calculator/pkg/calculator"
)

type Config struct {
	Port string
}

func CreateConfig() *Config {
	config := Config{}
	if len(os.Args) >= 2 {
		if _, err := strconv.Atoi(os.Args[1]); err == nil {
			config.Port = os.Args[1]
		} else {
			panic("передаваемыйпорт не число")
		}
	} else {
		config.Port = os.Getenv("PORT")
	}
	if config.Port == "" {
		config.Port = "8080"
	}

	return &config
}

type Application struct {
	config *Config
}

func (a *Application) Run() error {
	for {
		log.Println("input")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')

		if err != nil {
			log.Println("ошибка ввода")
			return nil
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			return nil
		}

		result, err := calculator.Calc(text)
		if err != nil {
			log.Println(text, "ошибка калькулятора: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

func New() *Application {
	return &Application{
		config: CreateConfig(),
	}
}

func (a *Application) RunServer() error {
	file, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	log.SetOutput(file)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", Outer(Validation(CalcHandler)))
	return http.ListenAndServe(":"+a.config.Port, mux)
}
