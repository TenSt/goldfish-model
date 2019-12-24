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

type data3 struct {
	MeteoST1 float64
	Boiler   int
	Kitchen  float64
	MeteoST2 float64
}

type data4 struct {
	Meteo  int `json:"meteo"`
	Boiler int `json:"boiler"`
	House  int `json:"house"`
}

func main() {
	// getJSONData()
	dataPrep()
}

func dataPrep() {
	lines, err := readCsv("grafana_data_export_edited.csv")
	if err != nil {
		log.Println(err)
	}

	for i, l := range lines {
		for j, v := range l {
			if v == "null" {
				l[j] = lines[i-1][j]
				v = lines[i-1][j]
			}
		}
	}
	for i, l := range lines {
		if i == 0 {
			l[0] = l[0][3:]
		}
		// lines[i] = lines[i][1:]
	}
	k := 1
	temps := make(map[string][][]string)
	var rawData []data4
	for i, l := range lines {

		if i == 0 {
			l[0] = l[0][3:]
		}

		var s string
		// meteoST1
		num1, err := strconv.ParseFloat(l[k], 64)
		checkError("num1 error parse:\n", err)
		// meteoST2
		num4, err := strconv.ParseFloat(l[k+3], 64)
		checkError("num4 error parse:\n", err)
		// meteo average
		num5 := (num1 + num4) / 2.0
		num5 = math.Round(num5)
		s = strconv.Itoa(int(num5))
		lines[i][k+3] = s

		// boiler
		numFloat2, err := strconv.ParseFloat(l[k+1], 64)
		checkError("error atoi:\n", err)
		margin := 4.8
		numFloat2 = numFloat2 + margin
		num2 := math.Round(numFloat2)
		s = strconv.Itoa(int(num2))
		lines[i][k+1] = s

		// kitchen
		num3, err := strconv.ParseFloat(l[k+2], 64)
		checkError("num3 error parse:\n", err)
		num3 = math.Round(num3)
		s = strconv.Itoa(int(num3))
		s1 := s
		lines[i][k+2] = s

		lines[i] = lines[i][2:]
		temps[s1] = append(temps[s1], lines[i])

		d := data4{
			Meteo:  int(num5),
			Boiler: int(num2),
			House:  int(num3),
		}
		rawData = append(rawData, d)
	}

	sort.Slice(rawData, func(i, j int) bool {
		if rawData[i].House < rawData[j].House {
			return true
		}
		if rawData[i].House > rawData[j].House {
			return false
		}
		return rawData[i].Meteo < rawData[j].Meteo
	})

	var pickedData []data4
	// var m, b, h int
	var temp data4
	for i, d := range rawData {
		if i == 0 {
			temp = d
			continue
		}
		if i == len(rawData)-1 {
			if temp.House == d.House && temp.Meteo == d.Meteo {
				if temp.Boiler < d.House {
					temp = d
				}
			}
			pickedData = append(pickedData, temp)
		}
		if temp.House == d.House && temp.Meteo == d.Meteo {
			if temp.Boiler < d.House {
				temp = d
				continue
			}
			continue
		}
		pickedData = append(pickedData, temp)
		temp = d
	}

	// for i, l := range temps {
	// 	writeCSV(l, "./newData/"+i)
	// }

	writeCSV(lines, "./newData/fullData")
	// writeFile(rawData, "./newData/fullDataSorted")
	writeFile(pickedData, "./newData/fullDataPicked")
}

func writeCSV(lines [][]string, name string) {
	f, err := os.Create(name + ".csv")
	checkError("Cannot create file", err)
	writer := csv.NewWriter(f)

	defer func() {
		writer.Flush()
		f.Close()
	}()

	err = writer.WriteAll(lines)
	checkError("Cannot write to file", err)
}

func writeFile(data []data4, name string) {
	f, err := os.Create(name + ".csv")
	checkError("Cannot create file", err)
	writer := csv.NewWriter(f)

	defer func() {
		writer.Flush()
		f.Close()
	}()

	for _, d := range data {
		var line []string
		line = append(line, strconv.Itoa(d.House), strconv.Itoa(d.Meteo), strconv.Itoa(d.Boiler))
		err = writer.Write(line)
		checkError("Cannot write to file", err)
	}

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
	}

	for i, l := range lines {
		if i == 0 {
			l[0] = l[0][3:]
		}
		// meteoST1
		num1, err := strconv.ParseFloat(l[1], 64)
		checkError("num1 error parse:\n", err)
		// num1 = num1 * 10 / 1000
		num1 = math.Round(num1)
		// if num1 > 15 {
		// 	continue
		// }

		// boiler !!!
		numFloat2, err := strconv.ParseFloat(l[2], 64)
		checkError("error atoi:\n", err)
		margin := 4.8
		numFloat2 = numFloat2 + margin
		num2 := int(math.Round(numFloat2))

		if num2 < 40 {
			num2 = 0
		} else if num2 >= 40 && num2 < 42 {
			num2 = 1
		} else if num2 >= 42 && num2 < 44 {
			num2 = 2
		} else if num2 >= 44 && num2 < 46 {
			num2 = 3
		} else if num2 >= 46 {
			num2 = 4
		}

		// br := 39
		// if num2 < br-1 {
		// 	continue
		// } else if num2 == 38 || num2 == 39 {
		// 	num2 = 39 - br
		// } else if num2 >= 40 && num2 < 44 {
		// 	num2 = 40 - br
		// } else if num2 >= 44 {
		// 	num2 = 44 - br - 3
		// }

		// num2 = num2 - br // the lowes cat
		// if num2 < 33 {
		// 	num2 = 33
		// } else if num2 > 46 {
		// 	num2 = 46
		// }
		// num2 = num2 - 33 // the lowes cat

		// kitchen
		num3, err := strconv.ParseFloat(l[3], 64)
		checkError("num3 error parse:\n", err)
		// num3 = num3 * 10 / 1000
		// num3 = math.Round(num3)

		// meteoST2
		num4, err := strconv.ParseFloat(l[4], 64)
		checkError("num4 error parse:\n", err)
		// num4 = num4 * 10 / 1000
		// num4 = math.Round(num4)
		// if num4 > 15 {
		// 	continue
		// }

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

	// fmt.Println("catMST1: \n", catMST1)
	// fmt.Println(len(catMST1))
	// fmt.Println("catMST2: \n", catMST2)
	// fmt.Println(len(catMST2))
	// fmt.Println("catKitchen: \n", catKitchen)
	// fmt.Println(len(catKitchen))

	prepareData(rawData)
}

func prepareData(rawData []data3) {
	// data
	writeToFile(rawData, "full")
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
}

func intInSlice(a int, m map[int]int) bool {
	for _, b := range m {
		if b == a {
			return true
		}
	}
	return false
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
