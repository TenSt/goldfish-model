package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type data struct {
	Time     time.Time
	MeteoST1 string
	Boiler   string
	Kitchen  string
	MeteoST2 string
}

func main() {
	lines, err := readCsv("grafana_data_export_edited.csv")
	if err != nil {
		log.Println(err)
	}

	f, err := os.Create("data.csv")
	checkError("Cannot create file", err)
	defer f.Close()

	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
			}
		}
		if i == 0 {
			l[0] = l[0][3:]
		}
	}

	writer := csv.NewWriter(f)
	err = writer.WriteAll(lines)
	checkError("Cannot write to file", err)
}

func getJSONData() {
	lines, err := readCsv("grafana_data_export.csv")
	if err != nil {
		log.Println(err)
	}

	var rawData []data

	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
			}
		}
		if i == 0 {
			l[0] = l[0][3:]
		}
		time1, err := time.Parse(time.RFC3339, l[0])
		if err != nil {
			log.Fatal(err)
		}
		d := data{
			Time:     time1,
			MeteoST1: l[1],
			Boiler:   l[2],
			Kitchen:  l[3],
			MeteoST2: l[4],
		}
		rawData = append(rawData, d)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rawData), func(i, j int) { rawData[i], rawData[j] = rawData[j], rawData[i] })

	splitTrain := float64(len(rawData)) * 0.7
	splitTest := float64(len(rawData)) * 0.9
	fmt.Println(int64(splitTrain))
	fmt.Println(int64(splitTest) - int64(splitTrain))
	fmt.Println(int64(len(rawData)) - int64(splitTest))

	trainData := rawData[:int64(splitTrain)]
	testData := rawData[int64(splitTrain):int64(splitTest)]
	validData := rawData[int64(splitTest):]

	writeToFile(trainData, "train")
	writeToFile(testData, "test")
	writeToFile(validData, "valid")
}

func writeToFile(rawData []data, name string) {
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(name + ".json")
	checkError("Cannot create file", err)
	defer f.Close()

	_, err = f.Write(jsonData)
	checkError("Cannot write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func readCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	reader := csv.NewReader(f)
	reader.Comma = rune(';')
	lines, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
