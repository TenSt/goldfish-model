package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type data struct {
	// Time     time.Time
	MeteoST1 string
	Boiler   string
	Kitchen  string
	MeteoST2 string
}

func main() {
	getJSONData()
}

func binData() {
	lines, err := readCsv("grafana_data_export_edited.csv")
	if err != nil {
		log.Println(err)
	}

	f, err := os.Create("data.csv")
	checkError("Cannot create file", err)
	defer f.Close()

	var categories []string
	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
				v = lines[i-1][j]
			}
			// binary
			// if j == 2 {
			// 	// fmt.Println(l)
			// 	// fmt.Println(l[j])
			// 	if v == "1" || v == "0" {
			// 	} else {
			// 		num, err := strconv.ParseFloat(v, 64)
			// 		if num >= 38 {
			// 			l[j] = "1"
			// 		} else {
			// 			l[j] = "0"
			// 		}
			// 		checkError("error", err)
			// 	}
			// }

			// categories
			// if j == 2 {
			// 	if len(v) == 4 {
			// 		l[j] = v[:len(v)-2]
			// 		v = v[:len(v)-2]
			// 	}
			// 	num, err := strconv.Atoi(v)
			// 	checkError("", err)
			// 	if num <= 35 {
			// 		l[j] = "0"
			// 		v = "0"
			// 	} else if num > 35 && num < 40 {
			// 		l[j] = "1"
			// 		v = "1"
			// 	} else if num >= 40 {
			// 		l[j] = "2"
			// 		v = "2"
			// 	}
			// 	if a := stringInSlice(v, categories); a == false {
			// 		categories = append(categories, v)
			// 	}
			// }

			// categories 2
			if j == 2 {
				if len(v) == 4 {
					l[j] = v[:len(v)-2]
					v = v[:len(v)-2]
				}
				num, err := strconv.Atoi(v)
				checkError("", err)
				if num <= 35 {
					l[j] = "35"
					v = "35"
				} else if num > 35 && num < 40 {
					l[j] = "40"
					v = "40"
				} else if num >= 40 {
					l[j] = "45"
					v = "45"
				}
				if a := stringInSlice(v, categories); a == false {
					categories = append(categories, v)
				}
			}
		}
		if i == 0 {
			l[0] = l[0][3:]
		}
	}

	writer := csv.NewWriter(f)
	head := []string{"Time", "C.MeteoST_bme280", "C.Gas_boiler_supply", "C.Kitchen_temp", "C.MeteoST_ds18b20"}
	err = writer.Write(head)
	err = writer.WriteAll(lines)
	checkError("Cannot write to file", err)
	sort.Strings(categories)
	fmt.Println(categories)
	fmt.Println(len(categories))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getJSONData() {
	lines, err := readCsv("grafana_data_export_edited.csv")
	if err != nil {
		log.Println(err)
	}

	var rawData []data

	var categories []string
	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
				v = lines[i-1][j]
			}

			// categories 2
			if j == 2 {
				if len(v) == 4 {
					l[j] = v[:len(v)-2]
					v = v[:len(v)-2]
				}
				num, err := strconv.Atoi(v)
				checkError("error atoi:\n", err)
				if num <= 35 {
					l[j] = "35"
					v = "35"
				} else if num > 35 && num < 40 {
					l[j] = "40"
					v = "40"
				} else if num >= 40 {
					l[j] = "45"
					v = "45"
				}
				if a := stringInSlice(v, categories); a == false {
					categories = append(categories, v)
				}
			}
		}
		if i == 0 {
			l[0] = l[0][3:]
		}
		// time1, err := time.Parse(time.RFC3339, l[0])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		d := data{
			// Time:     time1,
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

	sort.Strings(categories)
	fmt.Println(categories)
	fmt.Println(len(categories))
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
