package topsis

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sort"
)

//从json文件中获取矩阵
func GetMatrix(report *Report) (Y [][]float32) {
	M := report.Matrix
	Y = NewFloatMatrix(report.AlternativeNumber, report.CoefficientNumber)
	//Y_y := NewFloatMatrix(report.AlternativeNumber, report.CoefficientNumber)



	//矩阵正向化，前两列越大越好，后两列越小越好
	for i := 0; i < report.AlternativeNumber; i++{
		for j := 0; j < report.CoefficientNumber; j++{
			Y[i][j] = M[i][j]
		}
	}

	for i := 0; i < report.AlternativeNumber; i++{
		for j := 2; j < report.CoefficientNumber; j++{
			Y[i][j] = 1 / M[i][j]
		}
	}

	return Y
}

//原始矩阵归一化
func NormalizeMatrix(Y[][]float32) (X [][]float32) {
	S_s := NewFloatMatrix(len(Y), len(Y[0]))
	S := NewFloatVector(len(Y[0]))
	X = NewFloatMatrix(len(Y), len(Y[0]))

	for i := 0; i < len(Y); i++{
		for j := 0; j < len(Y[0]); j++{
			S_s[i][j] = Y[i][j] * Y[i][j]
		}
	}

	for j := 0; j < len(Y[0]); j++{
		for i := 0; i < len(Y); i++{
			S[j] += S_s[i][j]
		}
	}

	for i := 0; i < len(Y); i++{
		for j := 0; j < len(Y[0]); j++{
			X[i][j] = Y[i][j] / float32(math.Sqrt(float64(S[j])))
		}
	}
	return X
}


func PrintFloatVector(name string, vector []float32) {
	fmt.Println("=== ", name, " ===")
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

//生成基于熵的权值
func GetNormalizedWeights(X[][]float32) (W []float32) {
	//M := report.Matrix
	S := NewFloatVector(len(X[0]))
	P := NewFloatMatrix(len(X), len(X[0]))
	E := NewFloatVector(len(X[0]))
	D := NewFloatVector(len(X[0]))
	W = NewFloatVector(len(X[0]))

	for j := 0; j < len(X[0]); j++{
		for i := 0; i < len(X); i++{
			S[j] += X[i][j]
		}
	}

	//PrintFloatVector("S", S)//right

	for i := 0; i < len(X); i++{
		for j := 0; j < len(X[0]); j++{
			P[i][j] = X[i][j] / S[j]
		}
	}
	//PrintFloatMatrix("P", P)//right

	//熵值计算
	for j := 0; j < len(X[0]); j++{
		for i := 0; i < len(X); i++{
			E[j] = E[j] + float32(-1 / math.Log(float64(len(X))) * math.Log(float64(P[i][j])) * float64(P[i][j]))
		}
	}
	//PrintFloatVector("E", E)


	for j := 0; j < len(X[0]); j++{
			D[j] = 1 - E[j]
	}

	//PrintFloatVector("D", D)

	totalSum := Sum(D)
	for j := 0; j < len(X[0]); j++{
		W[j] = D[j] / totalSum
	}

	//PrintFloatVector("W", W)

	return W
}

//func GetNormalizedWeights(W []float32) (W_n []float32) {
//	W_n = NewFloatVector(len(W))
//	totalSum := Sum(W)
//
//	for i := range W {
//		W_n[i] = W[i] / totalSum
//	}
//
//	return
//}

//生成加权矩阵
func ApplyWeightedAverageMarks(Y [][]float32, W []float32) (Y_s [][]float32) {
	Y_s = NewFloatMatrix(len(Y), len(Y[0]))

	for i, w := range W {
		for a := range Y {
			Y_s[a][i] = Y[a][i] * w
		}
	}

	return
}


func GetReferencePoints(Y_s [][]float32) (Yplus, Yminus []float32) {
	Yplus = NewFloatVector(len(Y_s[0]))
	Yminus = NewFloatVector(len(Y_s[0]))

	for i := 0; i < len(Y_s[0]); i++ {
		column := GetColumn(Y_s, i)
		Yplus[i] = GetMax(column)
		Yminus[i] = GetMin(column)
	}

	return
}

func getEuclideanDistance(X []float32, Y[]float32) (distance float32) {
	for i := range X {
		distance += float32(math.Pow(float64(X[i] - Y[i]), float64(2)))
	}
	distance = float32(math.Sqrt(float64(distance)))
	return
}

//func GetDistancesToReferencePoints(Y_s [][]float32, Yplus []float32, Yminus []float32) (distances []*Distance) {
//	distances = make([]*Distance, len(Y_s))
//
//	for i, y := range Y_s {
//		dPlus := getEuclideanDistance(y, Yplus)
//		dMinus := getEuclideanDistance(y, Yminus)
//
//		alternativeName := "A" + strconv.Itoa(i+1)
//		distances[i] = NewDistance(alternativeName , dPlus, dMinus)
//	}
//	return distances
//}

func GetDistancesToReferencePoint(Y_s [][]float32, Yplus []float32, Yminus []float32) (h []float32) {
	h = NewFloatVector(len(Y_s))

	for i, y := range Y_s {
		dPlus := getEuclideanDistance(y, Yplus)
		dMinus := getEuclideanDistance(y, Yminus)
		h[i] = dMinus / (dPlus + dMinus)
	}
	return h
}

//RandomStr 随机生成字符串
func CreateRandomString(len int) string  {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

//切片转换为Map
func SliceToMap(Y_s [][]float32, Yplus []float32, Yminus []float32) (ResultMap map[string]float32) {

	h := NewFloatVector(len(Y_s))
	MapKey := len(h)
	ResultMap = map [string]float32{}
	var s []string

	for i, y := range Y_s {
		dPlus := getEuclideanDistance(y, Yplus)
		dMinus := getEuclideanDistance(y, Yminus)
		h[i] = dMinus / (dPlus + dMinus)
	}

	for i := 0; i < MapKey; i++{
		s = append(s, CreateRandomString(10))
		ResultMap[s[i]] = h[i]
		fmt.Printf("Key: %s, Value: %.16f\n", s[i], h[i])
	}

	return ResultMap
}

//对结果进行排序
func SortPoints(h []float32, ResultMap map[string]float32) (ResultKey string) {
	l := NewFloatVector64(len(h))
	for i := 0; i < len(h); i++{
		l[i] = float64(h[i])
	}
	sort.Float64s(l)
	//fmt.Println(l[len(h)-1])
	ResultKey = GetKey(float32(l[len(h)-1]), ResultMap)
	return ResultKey
}

//获取value对应的key
func GetKey(l float32, ResultMapResultMap map[string]float32) (k string) {

	for key := range ResultMapResultMap{
		if l == ResultMapResultMap[key]{
			k = key
		}
	}
	return k
}



