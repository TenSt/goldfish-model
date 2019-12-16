package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

type data3 struct {
	MeteoST1 float64
	Boiler   int
	Kitchen  float64
	MeteoST2 float64
}

type data2 struct {
	X string `json:"x"`
	Y string `json:"y"`
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

	var rawData []data3

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
				switch num {
				case 40:
					l[j] = "0"
					v = "0"
				case 41:
					l[j] = "1"
					v = "1"
				case 42:
					l[j] = "2"
					v = "2"
				case 43:
					l[j] = "3"
					v = "3"
				case 44:
					l[j] = "4"
					v = "4"
				case 45:
					l[j] = "5"
					v = "5"
				case 46:
					l[j] = "6"
					v = "6"
				case 47:
					l[j] = "7"
					v = "7"
				case 48:
					l[j] = "8"
					v = "8"
				case 49:
					l[j] = "9"
					v = "9"
				}
				if num <= 40 {
					l[j] = "0"
					v = "0"
				} else if num >= 50 {
					l[j] = "9"
					v = "9"
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
		num1, err := strconv.ParseFloat(l[1], 64)
		checkError("num1 error parse:\n", err)
		num1 = num1 / 100
		num1 = toFixed(num1, 3)
		num2, err := strconv.Atoi(l[2])
		checkError("error atoi:\n", err)
		num3, err := strconv.ParseFloat(l[3], 64)
		checkError("num3 error parse:\n", err)
		num3 = num3 / 100
		num3 = toFixed(num3, 3)
		num4, err := strconv.ParseFloat(l[4], 64)
		checkError("num4 error parse:\n", err)
		num4 = num4 / 100
		num4 = toFixed(num4, 3)
		d := data3{
			// Time:     time1,
			MeteoST1: num1,
			Boiler:   num2,
			Kitchen:  num3,
			MeteoST2: num4,
		}
		rawData = append(rawData, d)
	}

	// var rawData2 []data2
	// for _, v := range rawData {
	// 	d := data2{}
	// 	d.X = v.MeteoST1 + " " + v.MeteoST2 + " " + v.Kitchen
	// 	d.Y = v.Boiler
	// 	rawData2 = append(rawData2, d)
	// }

	// data
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rawData), func(i, j int) { rawData[i], rawData[j] = rawData[j], rawData[i] })

	splitTrain := float64(len(rawData)) * 0.7
	splitVal := float64(len(rawData)) * 0.9
	fmt.Println(int64(splitTrain))
	fmt.Println(int64(splitVal) - int64(splitTrain))
	fmt.Println(int64(len(rawData)) - int64(splitVal))

	trainData := rawData[:int64(splitTrain)]
	validData := rawData[int64(splitTrain):int64(splitVal)]
	testData := rawData[int64(splitVal):]

	writeToFile(trainData, "train")
	writeToFile(testData, "test")
	writeToFile(validData, "valid")

	// data2
	// rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(len(rawData2), func(i, j int) { rawData2[i], rawData2[j] = rawData2[j], rawData2[i] })

	// fmt.Println(int64(splitTrain))
	// fmt.Println(int64(splitVal) - int64(splitTrain))
	// fmt.Println(int64(len(rawData2)) - int64(splitVal))

	// trainData2 := rawData2[:int64(splitTrain)]
	// validData2 := rawData2[int64(splitTrain):int64(splitVal)]
	// testData2 := rawData2[int64(splitVal):]

	// writeToFile2(trainData2, "train2")
	// writeToFile2(testData2, "test2")
	// writeToFile2(validData2, "valid2")

	sort.Strings(categories)
	fmt.Println(categories)
	fmt.Println(len(categories))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func writeToFile(rawData []data3, name string) {
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

func writeToFile2(rawData2 []data2, name string) {
	fmt.Println(rawData2[0])
	jsonData, err := json.Marshal(rawData2)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonData[:10]))
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
