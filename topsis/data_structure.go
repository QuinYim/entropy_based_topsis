package topsis

import (
	"fmt"
)

type Expert [][]float32

type Report struct {
	AlternativeNumber int        `json:"alternative_number"`   //行号
	CoefficientNumber int        `json:"coefficient_number"`   //列号
    Matrix [][]float32          `json:"matrix"`
	//Experts           []Expert        `json:"experts"`
	//Weights           []float32       `json:"weights"`
}

type Distance struct {
	AName  string
	DPlus  float32
	DMinus float32
	H      float64
}

func NewDistance(name string, dPlus float32, dMinus float32) *Distance {
	h := float64(dMinus / (dPlus + dMinus))

	return &Distance{
		AName:name,
		DPlus: dPlus,
		DMinus: dMinus,
		H: h,
	}
}


func (d *Distance) String() string {

	return fmt.Sprintf("H(%s) = %3.10f",d.AName, d.H)
}