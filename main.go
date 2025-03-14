package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja Bem-vindo(a) ao Quiz-Game!")
	fmt.Println("Escreva o seu nome abaixo:")

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler a string")
	}

	validate_name := regexp.MustCompile(`^[a-zA-Z\s]`)
	if len(name) == 0 || !validate_name.MatchString(name) {
		panic("Por favor, insira um nome válido (apenas letras e espaços).")
	}

	g.Name = name
	fmt.Printf("\033[34m Olá %s! Jogo começando em 5 segundos. Boa sorte!\n", strings.TrimSpace(g.Name))
	time.Sleep(5 * time.Second)
}

func (g *GameState) ProcessCSV() {
	f, err := os.Open("Arquives/quiz-go.csv")
	if err != nil {
		panic("Erro ao ler o arquivo!")
	}

	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Erro ao ler o arquivo CSV")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	for index, question := range g.Questions {
		fmt.Printf("\n\033[33m %d. %s \033[0m\n", index+1, question.Text)

		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}

		fmt.Printf("Digite uma alternativa:\n")
		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = toInt(read[:len(read)-2])

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			break
		}

		if answer != question.Answer {
			fmt.Printf("\033[31m Você errou a resposta da questão [%d].\033[0m\n", index+1)
			break
		}

		g.Points += 10
		fmt.Printf("\033[32m Parabéns %s! \033[32mVocê acertou a resposta da questão %d. e recebeu\033[34m %d pontos!\033[0m\n", strings.TrimSpace(g.Name), index+1, g.Points)
	}
}

func main() {
	game := &GameState{}
	go game.ProcessCSV()
	game.Init()
	game.Run()

	if game.Points == 10 {
		fmt.Printf("\033[0m Olá, %s! O jogo chegou ao fim e você acertou apenas uma questão. Quase lá! Você conquistou \033[34m%d pontos\033[0m. Continue se esforçando, você vai conseguir mais na próxima!\n", strings.TrimSpace(game.Name), game.Points)
	} else if game.Points == 20 {
		fmt.Printf("\033[32m Olá, %s! Parabéns! Você venceu o jogo acertando todas as questões e conquistou \033[34m%d pontos\033[0m. Excelente trabalho!\033[0m\n", strings.TrimSpace(game.Name), game.Points)
	} else if game.Points == 30 {
		fmt.Printf("\033[32m Olá, %s! Parabéns! Você quase acertou tudo e obteve \033[34m%d pontos\033[0m. Que desempenho incrível! Só mais um pouquinho para alcançar a perfeição.\033[0m\n", strings.TrimSpace(game.Name), game.Points)
	} else if game.Points == 40 {
		fmt.Printf("\033[32m Olá, %s! Uau! Você acertou todas as questões e venceu o jogo com \033[34m%d pontos\033[0m.\033[0m\n", strings.TrimSpace(game.Name), game.Points)
	} else if game.Points == 0 {
		fmt.Printf("\033[31m Olá, %s! Infelizmente, você não marcou pontos desta vez. Mas não desanime! Cada erro é uma chance de aprender. Seu total foi \033[34m%d pontos\033[0m. Vamos tentar de novo e fazer melhor!\033[0m\n", strings.TrimSpace(game.Name), game.Points)
	}

	fmt.Println("Pressione a tecla [ENTER] para fechar o jogo!")
	fmt.Scanln()
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("não é permitido caracteres diferentes de números, por favor insira um número")
	}

	return i, nil
}
