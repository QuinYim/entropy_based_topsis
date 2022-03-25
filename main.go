package main

import (
	"../entropy_based_topsis/topsis"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var fileName = flag.String("file", "data/ex200.json", "File with formatted data (json format)")

func init() {
	flag.Parse()
}

func readFile(fileName string) (report *topsis.Report, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &report)
	return
}

func PrintFloat(name string, vector []float64) {
	fmt.Println("===", name, "===")
	fmt.Printf("%10.10f", vector[len(vector)-1])
	fmt.Println()
}

func PrintFloatVector(name string, vector []float32) {
	fmt.Println("===", name, "===")
	for i := 0; i < len(vector); i++ {
		fmt.Printf("%10.10f", vector[i])
	}
	fmt.Println()
}

func PrintFloatMatrix(name string, matrix [][]float32) {
	fmt.Println("=== ", name, " ===")
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			fmt.Printf("%10.10f", matrix[i][j])
		}

		fmt.Println()
	}
}

func main() {
	fmt.Println("FILE:    ", *fileName)

	fmt.Println("1. File read")
	report, err := readFile(*fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := time.Now()

	fmt.Println("2. Compute average expert's opinions")
	Y := topsis.GetMatrix(report)
	PrintFloatMatrix("Y", Y)

	fmt.Println("2. Compute 归一")
	X := topsis.NormalizeMatrix(Y)
	PrintFloatMatrix("X", X)

	fmt.Println("3. Calculate normalized weights and apply weights to averaged marks")
	W := topsis.GetNormalizedWeights(X)
	PrintFloatVector("Normalized weights Wn", W)
	Y_s := topsis.ApplyWeightedAverageMarks(Y, W)
	PrintFloatMatrix("Y'", Y_s)

	fmt.Println("4. Calculate best and worst Y points (Y+ & Y-)")
	Yplus, Yminus := topsis.GetReferencePoints(Y_s)
	PrintFloatVector("Y+", Yplus)
	PrintFloatVector("Y-", Yminus)

	//fmt.Println("5. Calculate referenced distances")
	//distances := topsis.GetDistancesToReferencePoints(Y_s, Yplus, Yminus)
	//for _, d := range distances {
	//	fmt.Println(d)
	//}

	fmt.Println("5. Sorted distances")
	ResultMap := topsis.SliceToMap(Y_s, Yplus, Yminus)

	fmt.Println("6. Result Shows")
	l := topsis.GetDistancesToReferencePoint(Y_s, Yplus, Yminus)
	result := topsis.SortPoints(l, ResultMap)
	fmt.Println(result)

	elapsed := time.Since(t)

	fmt.Println("app elapsed:", elapsed)

}
