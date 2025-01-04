package onlinesim

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	getTariffsEndpoint = "getTariffs.php"
	getNumEndpoint 	   = "getNum.php"
	getStateEndpoint   = "getState.php"
)

type GetTariffsOptions struct {
	LocalePrice    string
	Country        string
	FilterCountry  string
	FilterService  string
	Count          string
	Page           string
	Lang           string
}


type GetTariffsResp struct {
	Response            string                      `json:"response"`
	Countries           map[string]Country         `json:"countries"`
	Services            map[string]Service         `json:"services"`
	FavoriteCountries   map[string]interface{}     `json:"favorite_countries"`
	FavoriteServices    []interface{}              `json:"favorite_services"`
	Page                int                        `json:"page"`
	Country             int                        `json:"country"`
	Filter              string                     `json:"filter"`
	SubscriptionTariffs []SubscriptionTariff       `json:"subscription_tariffs"`
	End                 bool                       `json:"end"`
	Favorites           map[string]interface{}     `json:"favorites"`
}

type Country struct {
	Name     string `json:"name"`
	Original string `json:"original"`
	Code     int    `json:"code"`
	Pos      int    `json:"pos"`
	Other    bool   `json:"other"`
	New      bool   `json:"new"`
	Enable   bool   `json:"enable"`
}

type Service struct {
	ID      int    `json:"id"`
	Count   int    `json:"count"`
	Price   string `json:"price"`
	Service string `json:"service"`
	Slug    string `json:"slug"`
}

type SubscriptionTariff struct {
	ID             int     `json:"id"`
	CountOps       int     `json:"count_operations"`
	Price          string  `json:"price"`
	LifeDays       int     `json:"life_days"`
	IsBest         bool    `json:"is_best"`
	IsCustom       bool    `json:"is_custom"`
	Enabled        bool    `json:"enabled"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	Sum            string  `json:"sum"`
	CurrencyRatio  float64 `json:"currency_ratio"`
	Symbol         string  `json:"symbol"`
}

type GetNumResp struct {
	Response     int    `json:"response"`
	Tzid         int    `json:"tzid"`
	Number       string `json:"number"`
	Country      int    `json:"country"`
	Time         int    `json:"time"`
	Service      string `json:"service"`
	Title        string `json:"title"`
	ResponseText string `json:"response_text"`
}

type NumState struct {
	Country  int     `json:"country"`
	Sum      float32 `json:"sum"`
	Service  string  `json:"service"`
	Number   string  `json:"number"`
	Response string  `json:"response"`
	Tzid     int     `json:"tzid"`
	Time     int     `json:"time"`
	Form     string  `json:"form"`
	Msg      string  `json:"msg"`
}


func (c *Client) GetTariffs(options GetTariffsOptions) (GetTariffsResp, error) {
	URL := baseURL + getTariffsEndpoint
	
	query := url.Values{}

	if options.LocalePrice != "" {
		query.Add("locale_price", options.LocalePrice)
	}
	if options.Country != "" {
		query.Add("country", options.Country)
	}
	if options.FilterCountry != "" {
		query.Add("filter_country", options.FilterCountry)
	}
	if options.FilterService != "" {
		query.Add("filter_service", options.FilterService)
	}
	if options.Count != "" {
		query.Add("count", options.Count)
	}
	if options.Page != "" {
		query.Add("page", options.Page)
	}
	if options.Lang != "" {
		query.Add("lang", options.Lang)
	}


	req, _ := http.NewRequest(
		http.MethodGet,
		URL + "?" + query.Encode(),
		nil,
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetTariffsResp{}, err
	}
	
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetTariffsResp{}, err
	}

	getTariffsResp := GetTariffsResp{}

	err = json.Unmarshal(respBody, &getTariffsResp)
	if err != nil {
		return GetTariffsResp{}, err
	}

	return getTariffsResp, nil
}

func (c *Client) GetNum(service string, country int) (GetNumResp, error) {
	URL := baseURL + getNumEndpoint

	query := url.Values{}

	query.Add("service", service)
	query.Add("country", strconv.Itoa(country))
	query.Add("number", "true")
	query.Add("apikey", c.ApiKey)

	req, _ := http.NewRequest(
		http.MethodGet,
		URL + "?" + query.Encode(),
		nil,
	)
	req.Header.Set("Authorization", "Bearer " + c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetNumResp{}, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetNumResp{}, err
	}
	
	getNumResp := GetNumResp{}
	
	err = json.Unmarshal(respBody, &getNumResp)
	if err != nil {
		return GetNumResp{}, err
	}
	return getNumResp, err
}

func (c *Client) GetState(tzid int, messageToCode int) ([]NumState, error) {
	URL := baseURL + getStateEndpoint

	query := url.Values{}
	query.Add("tzid", strconv.Itoa(tzid))
	query.Add("message_to_code", strconv.Itoa(messageToCode))
	query.Add("apikey", c.ApiKey)

	req, _ := http.NewRequest(
		http.MethodGet,
		URL + "?" + query.Encode(),
		nil,
	)
	req.Header.Add("Authorization", "Bearer " + c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	states := []NumState{}
	
	err = json.Unmarshal(respBody, &states)
	if err != nil {
		return nil, err
	}

	return states, nil
}
