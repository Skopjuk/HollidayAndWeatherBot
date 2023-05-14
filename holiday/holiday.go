package holiday

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HolidayAPI struct {
	token string
}

type HolidayCatalogue struct {
	HolidayCatalogue []HolidayData `json:"holiday_catalogue"`
}

type HolidayData struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func NewHolidayAPI(token string) *HolidayAPI {
	return &HolidayAPI{token: token}
}

func (h *HolidayAPI) MakeRequest(country string) (*[]string, error) {
	now := time.Now()

	currentYear := strconv.Itoa(now.Year())
	currentDay := strconv.Itoa(now.Day())
	currentMonth := strconv.Itoa(int(now.Month()))
	holidayCatalogue := []HolidayData{}
	var holidayList []string

	url := fmt.Sprintf("https://holidays.abstractapi.com/v1/?api_key=%s&country=%s&year=%s&month=%s&day=%s", h.token, country, currentYear, currentMonth, currentDay)

	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	err = json.Unmarshal(body, &holidayCatalogue)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	for i := 0; i < len(holidayCatalogue); i++ {
		holidayList = append(holidayList, holidayCatalogue[i].Name)
	}

	return &holidayList, nil
}
