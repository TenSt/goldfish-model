package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
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
	// new comment
}

func getJSONData() {
	lines, err := readCsv("grafana_data_export_edited.csv")
	if err != nil {
		log.Println(err)
	}

	var rawData []data3

	catBoiler := make(map[int]int)
	catMST1 := make(map[float64]int)
	catMST2 := make(map[float64]int)
	catKitchen := make(map[float64]int)
	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
				v = lines[i-1][j]
			}
		}
		if i == 0 {
			l[0] = l[0][3:]
		}
		// meteoST1
		num1, err := strconv.ParseFloat(l[1], 64)
		checkError("num1 error parse:\n", err)
		// num1 = num1 * 10 / 1000

		// boiler !!!
		numFloat2, err := strconv.ParseFloat(l[2], 64)
		checkError("error atoi:\n", err)
		margin := 4.8
		numFloat2 = numFloat2 + margin
		num2 := int(math.Round(numFloat2)) - 22
		if num2 == 34 {
			num2 = 33
		}

		// kitchen
		num3, err := strconv.ParseFloat(l[3], 64)
		checkError("num3 error parse:\n", err)
		// num3 = num3 * 10 / 1000

		// meteoST2
		num4, err := strconv.ParseFloat(l[4], 64)
		checkError("num4 error parse:\n", err)
		// num4 = num4 * 10 / 1000

		// cats
		if v, k := catMST1[num1]; k == false {
			catMST1[num1] = 1
		} else {
			catMST1[num1] = v + 1
		}

		if v, k := catBoiler[num2]; k == false {
			catBoiler[num2] = 1
		} else {
			catBoiler[num2] = v + 1
		}

		if v, k := catKitchen[num3]; k == false {
			catKitchen[num3] = 1
		} else {
			catKitchen[num3] = v + 1
		}

		if v, k := catMST2[num4]; k == false {
			catMST2[num4] = 1
		} else {
			catMST2[num4] = v + 1
		}

		d := data3{
			MeteoST1: num1,
			Boiler:   num2,
			Kitchen:  num3,
			MeteoST2: num4,
		}
		rawData = append(rawData, d)
	}

	// sort.Ints(categories)
	fmt.Println("catBoiler: \n", catBoiler)
	fmt.Println(len(catBoiler))

	fmt.Println("catMST1: \n", catMST1)
	fmt.Println(len(catMST1))

	fmt.Println("catMST2: \n", catMST2)
	fmt.Println(len(catMST2))

	fmt.Println("catKitchen: \n", catKitchen)
	fmt.Println(len(catKitchen))

	prepareData(rawData)
}

func prepareData(rawData []data3) {
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

	writeToFile(rawData, "full")
	writeToFile(trainData, "train")
	writeToFile(testData, "test")
	writeToFile(validData, "valid")
}

func intInSlice(a int, m map[int]int) bool {
	for _, b := range m {
		if b == a {
			return true
		}
	}
	return false
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
