package api

import (
	"github.com/hongyuefan/facedetect/features"
	"github.com/hongyuefan/tmpserver/models"
	"github.com/hongyuefan/tmpserver/types"
)

type Stand struct {
	Id    int64
	Name  string
	Upper float64
	Lower float64
}

func NewStand(id int64, name string, upper, lower float64) *Stand {
	return &Stand{
		Id:    id,
		Name:  name,
		Upper: upper,
		Lower: lower,
	}
}

var mapStand map[int64]*Stand

func init() {

	m := make(map[string]string, 1)

	mapStand = make(map[int64]*Stand, 0)

	mls, err := models.GetFeatureStander(m, []string{}, []string{}, []string{}, -1, -1)

	if err != nil {
		panic(err)
	}

	for _, ml := range mls {

		id := ml.(models.FeatureStand).Id

		mapStand[id] = NewStand(id, ml.(models.FeatureStand).Name, ml.(models.FeatureStand).LimitUpper, ml.(models.FeatureStand).LimitLower)

	}
}

func GetResult(key, scr, url string) (string, error) {

	var (
		mapResult map[int64]int
		results   []int64
	)

	points, err := features.GetFeatures(key, scr, url)

	if err != nil {
		return "", err
	}

	mapResult[types.EYE_LENGTH] = Computer(features.RateEyeSize(points.Landmark150), mapStand[types.EYE_LENGTH].Upper, mapStand[types.EYE_LENGTH].Lower)
	mapResult[types.EYE_HEIGHT] = Computer(features.RateEyeWidth(points.Landmark150), mapStand[types.EYE_HEIGHT].Upper, mapStand[types.EYE_HEIGHT].Lower)
	mapResult[types.EYE_ANGLE] = Computer(features.AngleEye(points.Landmark150), mapStand[types.EYE_ANGLE].Upper, mapStand[types.EYE_ANGLE].Lower)
	mapResult[types.EYE_TO_EYE] = Computer(features.RateEyeDistance(points.Landmark150), mapStand[types.EYE_TO_EYE].Upper, mapStand[types.EYE_TO_EYE].Lower)
	mapResult[types.EYE_TO_EYEBROW] = Computer(features.RateEyeToBrow(points.Landmark150), mapStand[types.EYE_TO_EYEBROW].Upper, mapStand[types.EYE_TO_EYEBROW].Lower)
	mapResult[types.EYEBROW_LENGTH] = Computer(features.RateEyeBrowEye(points.Landmark150), mapStand[types.EYEBROW_LENGTH].Upper, mapStand[types.EYEBROW_LENGTH].Lower)
	mapResult[types.EYEBROW_TO_EYEBROW] = Computer(features.RateEyeBrowToEyeBrow(points.Landmark150), mapStand[types.EYEBROW_TO_EYEBROW].Upper, mapStand[types.EYEBROW_TO_EYEBROW].Lower)
	mapResult[types.EYEBROW_ANGLE] = Computer(features.AngleEyeBrow(points.Landmark150), mapStand[types.EYEBROW_ANGLE].Upper, mapStand[types.EYEBROW_ANGLE].Lower)
	mapResult[types.EYEBROW_MAX_RATIO] = Computer(features.RateEyeBrow(points.Landmark150), mapStand[types.EYEBROW_MAX_RATIO].Upper, mapStand[types.EYEBROW_MAX_RATIO].Lower)
	mapResult[types.NOSE_WIDTH] = Computer(features.RateNoseWidth(points.Landmark150), mapStand[types.NOSE_WIDTH].Upper, mapStand[types.NOSE_WIDTH].Lower)
	mapResult[types.NOSE_LENGTH] = Computer(features.RateFaceLength(points.Landmark150), mapStand[types.NOSE_LENGTH].Upper, mapStand[types.NOSE_LENGTH].Lower)
	mapResult[types.NOSE_RATIO] = Computer(features.RateNoseEagle(points.Landmark150), mapStand[types.NOSE_RATIO].Upper, mapStand[types.NOSE_RATIO].Lower)
	mapResult[types.PHILTRUM_LENGTH] = Computer(features.RateRenZLength(points.Landmark150), mapStand[types.PHILTRUM_LENGTH].Upper, mapStand[types.PHILTRUM_LENGTH].Lower)
	mapResult[types.MOUTH_WIDTH] = Computer(features.RateMouseLength(points.Landmark150), mapStand[types.MOUTH_WIDTH].Upper, mapStand[types.MOUTH_WIDTH].Lower)
	mapResult[types.MOUTH_THICKNESS] = Computer(features.RateMouthLipThickness(points.Landmark150), mapStand[types.MOUTH_THICKNESS].Upper, mapStand[types.MOUTH_THICKNESS].Lower)
	mapResult[types.MOUTH_LIPS_RATIO] = Computer(features.RateMouthLip(points.Landmark150), mapStand[types.MOUTH_LIPS_RATIO].Upper, mapStand[types.MOUTH_LIPS_RATIO].Lower)
	mapResult[types.CHIN_WIDTH] = Computer(features.RateChinWidth(points.Landmark150), mapStand[types.CHIN_WIDTH].Upper, mapStand[types.CHIN_WIDTH].Lower)
	mapResult[types.MOUTH_LIPS_EQUAL] = Computer(features.AngleMouth(points.Landmark150), mapStand[types.MOUTH_LIPS_EQUAL].Upper, mapStand[types.CHIN_WIDTH].Lower)

	if features.RateMouthLip(points.Landmark150) >= 46 && features.RateMouthLip(points.Landmark150) <= 56 {
		mapResult[types.MOUTH_LIPS_EQUAL] = 1
	} else {
		mapResult[types.MOUTH_LIPS_EQUAL] = 0
	}
	if points.Emotion.Type == "angry" {
		mapResult[types.FACE_ANGRY] = 1
	} else {
		mapResult[types.FACE_ANGRY] = 0
	}
	if points.Faceshap.Type == "square" {
		mapResult[types.FACE_SHAP] = 1
	} else if points.Faceshap.Type == "oval" {
		mapResult[types.FACE_SHAP] = -1
	} else {
		mapResult[types.FACE_SHAP] = 0
	}

	results = append(results, EyeResult(mapResult)...)
	results = append(results, EyeBrowResult(mapResult)...)
	results = append(results, MouthResult(mapResult)...)
	results = append(results, NoseResult(mapResult)...)
	results = append(results, RenZResult(mapResult)...)
	results = append(results, FaceResult(mapResult)...)
	results = append(results, ChinResult(mapResult)...)

	return GetDescribe(results), nil
}

func GetDescribe(results []int64) string {

	var s string

	for _, result := range results {

		des := &models.Describe{
			Id: result,
		}

		if err := models.GetDescribeById(des); err != nil {
			continue
		}

		s = s + des.Des
	}

	return s
}

func EyeResult(mapResult map[int64]int) []int64 {

	var result []int64

	switch mapResult[types.EYE_LENGTH] {
	case 1:
		switch mapResult[types.EYE_HEIGHT] {
		case 1:
			result = append(result, 2)
		case 0:
			result = append(result, 2)
		case -1:
			if mapResult[types.EYE_ANGLE] == 1 {
				result = append(result, 9)
			} else {
				result = append(result, 3)
			}

		}

	case 0:
		switch mapResult[types.EYE_HEIGHT] {
		case 1:
			result = append(result, 2)
		case -1:
			result = append(result, 1)
		}

	case -1:
		switch mapResult[types.EYE_HEIGHT] {
		case 1:
			result = append(result, 4)
		case 0:
			result = append(result, 1)
		case -1:
			result = append(result, 1)
		}
	}

	if mapResult[types.EYE_ANGLE] == 1 {
		result = append(result, 8)
	} else if mapResult[types.EYE_ANGLE] == -1 {
		result = append(result, 7)
	}

	if mapResult[types.EYE_TO_EYE] == 1 {
		result = append(result, 6)
	} else if mapResult[types.EYE_TO_EYE] == -1 {
		result = append(result, 5)
	}

	return result
}

func EyeBrowResult(mapResult map[int64]int) []int64 {

	var result []int64

	switch mapResult[types.EYEBROW_TO_EYEBROW] {
	case 1:
		result = append(result, 15)
	case 0:
		result = append(result, 14)
	case -1:
		result = append(result, 13)
	}

	if mapResult[types.EYEBROW_ANGLE] == -1 {
		result = append(result, 19) //八字眉
	}

	if mapResult[types.EYEBROW_HEIGHT] != -1 && mapResult[types.EYEBROW_LENGTH] != -1 && mapResult[types.EYEBROW_MAX_RATIO] == 0 {
		result = append(result, 20) //一字眉
	}

	if mapResult[types.EYEBROW_LENGTH] == 1 && mapResult[types.EYEBROW_HEIGHT] == -1 {
		result = append(result, 10) //细长
	}

	if mapResult[types.EYEBROW_LENGTH] != -1 && mapResult[types.EYEBROW_MAX_RATIO] == 1 && mapResult[types.EYEBROW_ANGLE] != -1 {
		result = append(result, 18) //三角
	}

	if mapResult[types.EYEBROW_LENGTH] == -1 {
		result = append(result, 12) //眉形短
	}

	if mapResult[types.EYE_TO_EYEBROW] == 1 {
		result = append(result, 17)
	}

	if mapResult[types.EYE_TO_EYEBROW] == -1 {
		result = append(result, 16)
	}
	return result

}

func NoseResult(mapResult map[int64]int) []int64 {

	var result []int64

	if mapResult[types.NOSE_WIDTH] == 1 {
		result = append(result, 21)
	}
	if mapResult[types.NOSE_WIDTH] == -1 {
		result = append(result, 22)
	}
	if mapResult[types.NOSE_LENGTH] == 1 {
		result = append(result, 27)
	}
	if mapResult[types.NOSE_LENGTH] == -1 {
		result = append(result, 28)
	}
	if mapResult[types.NOSE_RATIO] == 1 {
		result = append(result, 26)
	}
	return result
}

func RenZResult(mapResult map[int64]int) []int64 {

	var result []int64

	if mapResult[types.PHILTRUM_LENGTH] == -1 {
		result = append(result, 29)
	}
	return result
}

func MouthResult(mapResult map[int64]int) []int64 {

	var result []int64

	if mapResult[types.MOUTH_WIDTH] == 1 {
		result = append(result, 30)
	}
	if mapResult[types.MOUTH_WIDTH] == -1 {
		result = append(result, 31)
	}
	if mapResult[types.MOUTH_THICKNESS] == 1 {
		result = append(result, 33)
	}
	if mapResult[types.MOUTH_THICKNESS] == -1 {
		result = append(result, 32)
	}
	if mapResult[types.MOUTH_LIPS_RATIO] == -1 {
		result = append(result, 38)
	}
	if mapResult[types.MOUTH_LIPS_RATIO] == 1 {
		result = append(result, 39)
	}
	if mapResult[types.MOUTH_LIPS_EQUAL] == 1 {
		result = append(result, 34)
	}
	if mapResult[types.MOUTH_ANGLE] == 1 {
		result = append(result, 36)
	}
	if mapResult[types.MOUTH_ANGLE] == -1 {
		result = append(result, 37)
	}
	return result
}

func ChinResult(mapResult map[int64]int) []int64 {

	var result []int64

	if mapResult[types.CHIN_WIDTH] == 1 && mapResult[types.FACE_SHAP] != 1 {
		result = append(result, 40)
	}
	if mapResult[types.CHIN_WIDTH] == -1 {
		result = append(result, 41)
	}

	return result
}

func FaceResult(mapResult map[int64]int) []int64 {

	var result []int64

	if mapResult[types.FACE_SHAP] == -1 {
		result = append(result, 45)
	}

	if mapResult[types.FACE_SHAP] == 1 && mapResult[types.CHIN_WIDTH] == 1 {
		result = append(result, 42)
	}

	if mapResult[types.FACE_ANGRY] == 1 {
		result = append(result, 46)
	}

	return result
}

func Computer(o, u, l float64) int {
	if o > u {
		return 1
	} else if o < l {
		return -1
	} else {
		return 0
	}
}