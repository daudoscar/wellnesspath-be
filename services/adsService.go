package services

import (
	"fmt"
	"time"
	"wellnesspath/dto"
	"wellnesspath/helpers"
)

type AdsService struct{}

func (s *AdsService) GetAds(UserID uint64, AdsID uint64) (dto.AdsResponseDTO, error) {
	var ext string
	if AdsID == 3 {
		ext = ".mp4"
	} else {
		ext = ".png"
	}
	blobName := "ads/ads" + fmt.Sprint(AdsID) + ext
	AdsURL, err := helpers.GenerateSASURLAds(blobName, time.Hour)
	if err != nil {
		AdsURL = ""
	}

	return dto.AdsResponseDTO{
		AdsURL: AdsURL,
	}, nil
}
