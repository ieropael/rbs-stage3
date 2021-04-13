package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// чтение аргументов
	urls := flag.String("urls", "urls.txt", "path to file with urls")
	resultDir := flag.String("resultdir", "./", "path to directory with html files")
	flag.Parse()

	// создание папки с результатами запросов
	os.Mkdir(*resultDir, 0777)

	// открытие файла с url-адресами
	urlsFile, err := os.Open(*urls)
	if err != nil {
		fmt.Println(err)
		return
	}

	// считывание url-адресов
	urlsScanner := bufio.NewScanner(urlsFile)
	logFunc(fmt.Sprintf("Считывание url-адресов из файла %s", *urls))

	for urlsScanner.Scan() {

		url := urlsScanner.Text()

		// создание нового файла
		fileName := fmt.Sprintf("%s/%s.html",
			*resultDir, url[strings.LastIndex(url, "/")+1:strings.LastIndex(url, ".")])
		resultFile, err := os.Create(fileName)
		if err != nil {
			log.Println(err)
			return
		}
		logFunc(fmt.Sprintf("Создание файла %s", fileName))

		// отправка GET-запроса
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		logFunc(fmt.Sprintf("Отправка GET-запроса к %s", url))

		// запись результата запроса в файл
		for {
			urlsResp := make([]byte, 1000)
			n, err := resp.Body.Read(urlsResp)
			resultFile.WriteString(string(urlsResp[:n]))
			if n == 0 || err != nil {
				break
			}
		}
		logFunc(fmt.Sprintf("Запись результата в файл %s", fileName))

		// закрытие запроса
		defer resp.Body.Close()
	}

	// закрытие файла
	defer urlsFile.Close()
}

func logFunc(logMessage string) {
	log.Println(logMessage)
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	logger := log.New(logFile, "INFO: ", log.LstdFlags)
	logger.Println(logMessage)
	defer logFile.Close()
}
