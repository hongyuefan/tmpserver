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

func GetDescribe(key, scr, url string) error {

	var mapResult map[int64]int

	points, err := features.GetFeatures(key, scr, url)

	if err != nil {
		return err
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
	mapResult[types.MOUTH_LIPS_RATIO] = Computer(features.RateMoseLip(points.Landmark150), mapStand[types.MOUTH_LIPS_RATIO].Upper, mapStand[types.MOUTH_LIPS_RATIO].Lower)
	mapResult[types.CHIN_WIDTH] = Computer(features.RateChinWidth(points.Landmark150), mapStand[types.CHIN_WIDTH].Upper, mapStand[types.CHIN_WIDTH].Lower)

	return nil
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
