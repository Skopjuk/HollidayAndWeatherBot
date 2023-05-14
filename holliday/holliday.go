package holliday

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
	currantYear := strconv.Itoa(time.Now().Year())
	currantDay := strconv.Itoa(time.Now().Day())
	currantMonth := strconv.Itoa(int(time.Now().Month()))
	holidayCatalogue := []HolidayData{}
	var holidayList []string

	url := fmt.Sprintf("https://holidays.abstractapi.com/v1/?api_key=%s&country=%s&year=%s&month=%s&day=%s", h.token, country, currantYear, currantMonth, currantDay)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println(string(body))

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

func (h *HolidayAPI) HolidayListInString(countryCode string) (string, error) {
	var holidayListInString string

	holidayArray, err := h.MakeRequest(countryCode)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if holidayArray != nil && len(*holidayArray) > 0 {
		for i := 0; i < len(*holidayArray); i++ {
			holidayListInString += (*holidayArray)[0]
		}
		return holidayListInString, nil
	} else {
		holidayListInString = ""
	}
	return "today is no holidays in this country", nil
}
