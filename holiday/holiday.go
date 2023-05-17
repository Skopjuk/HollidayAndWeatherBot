package holiday

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
)

type HolidayAPI struct {
	token      string
	apiAddress string
}

type HolidayCatalogue struct {
	HolidayCatalogue []HolidayData `json:"holiday_catalogue"`
}

type HolidayData struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func NewHolidayAPI(token string, apiAddress string) *HolidayAPI {
	return &HolidayAPI{
		token:      token,
		apiAddress: apiAddress,
	}
}

func (h *HolidayAPI) MakeRequest(country string) ([]string, error) {
	now := time.Now()

	currentYear := strconv.Itoa(now.Year())
	currentDay := strconv.Itoa(now.Day())
	currentMonth := strconv.Itoa(int(now.Month()))
	holidayCatalogue := []HolidayData{}
	var holidayList []string

	url := fmt.Sprintf("%sv1/?api_key=%s&country=%s&year=%s&month=%s&day=%s", h.apiAddress, h.token, country, currentYear, currentMonth, currentDay)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &holidayCatalogue)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(holidayCatalogue); i++ {
		holidayList = append(holidayList, holidayCatalogue[i].Name)
	}

	return holidayList, nil
}

func (h *HolidayAPI) TransformListOfHolidaysToStr(pressedButton string) (string, error) {
	var holidayListInString string

	holidayArray, err := h.MakeRequest(pressedButton)
	if err != nil {
		logrus.Error(err)
		return "", fmt.Errorf("can't get a list of holidays for %s: %w", pressedButton, err)
	}

	for i := 0; i < len(holidayArray); i++ {
		holidayListInString += (holidayArray)[0]
	}

	if len(holidayListInString) > 0 {
		return holidayListInString, nil
	}

	return "Today is no holidays in this country.", nil
}
