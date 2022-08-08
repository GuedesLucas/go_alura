package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 5
const delay = 5

func main() {
	showIntroduction()
	for {
		showMenu()
		command := scanCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Exibindo Logs...")
			showLogs()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	name := "Senhor"
	version := 1.1
	fmt.Println("Olá, sr.", name)
	fmt.Println("Este programa esta na versão", version)
}

func showMenu() {
	fmt.Println("1- Iniciar monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func scanCommand() int {
	var command int
	fmt.Scan(&command)

	fmt.Println("O command escolhido foi ", command)
	return command
}

func startMonitoring() {
	fmt.Println("Monitorando...")
	sites := readSitesArchive()
	for i := 0; i < monitoring; i++ {
		for _, site := range sites {
			validateSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func validateSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code:", resp.StatusCode)
		registerLog(site, false)
	}
}

func validateArchive() {
	valid, err := os.Open("sites.txt")
	if err != nil {
		createFileSites()
	}
	valid.Close()
}

func readSitesArchive() []string {
	var sites []string
	validateArchive()
	archive, _ := os.Open("sites.txt")
	reader := bufio.NewReader(archive)
	for {
		line, err := reader.ReadString('\n')
		sites = append(sites, strings.TrimSpace(line))
		if err == io.EOF {
			break
		}
	}
	archive.Close()
	return sites
}

func createFileSites() {
	newFile, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	sites := []string{"https://random-status-code.herokuapp.com/", "https://google.com/", "https://www.alura.com.br/", "https://www.caelum.com.br/"}
	for _, site := range sites {
		if site == sites[len(sites)-1] {
			newFile.WriteString(site)
		} else {
			newFile.WriteString(site + "\n")
		}
	}
	newFile.Close()
}

func registerLog(site string, status bool) {
	archive, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	archive.WriteString(time.Now().Format("02/01/2006 15:04:05") + "-" + site + "- online: " + strconv.FormatBool(status) + "\n")
	archive.Close()
}

func showLogs() {
	archive, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(archive))
}
