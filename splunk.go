package splunkadapter

import (
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"fmt"
	"strings"
)

type SplunkResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Body string // splunk response data in csv format
	ContentLength int64
}

func SplunkSearch(query, earliestTime, lastestTime string) (SplunkResponse, error) {

	url := "https://10.251.38.60:8088/services/search/jobs/export"
	

	payload := strings.NewReader("search=search%20sourcetype%3Dtselreflexpremiumtdr%20%0A%7C%20eval%20timex%3Dstrftime(_time%2C%22%25Y-%25m-%25d%20%25H%3A%25M%22)%20%0A%7C%20eval%20timexepoch%3Dround(strptime(timex%2C%22%25Y-%25m-%25d%20%25H%3A%25M%22)%2C0)%20%0A%7C%20fillnull%20value%3Dnull%20host%2C%20adn%2C%20appsID%2C%20taskID%2C%20description_trx%2C%20clientID%2C%20statusCode%2C%20detailedCode%2C%20productCategoryLookup%2C%20productCategoryLookupBasis%2C%20productLookup%0A%7C%20stats%20count%20as%20trx%2C%20sum(elapTime)%20as%20totRT%2C%20max(elapTime)%20as%20maxRT%2C%20min(elapTime)%20as%20minRT%20by%20timexepoch%2C%20host%2C%20adn%2C%20appsID%2C%20taskID%2C%20description_trx%2C%20clientID%2C%20statusCode%2C%20detailedCode%2C%20productCategoryLookup%2C%20productCategoryLookupBasis%2C%20productLookup%20&earliest_time=-5m&latest_time=-4m&output_mode=csv&adhoc_search_level=fast")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authorization", "Basic aXRvYzppdG9jNTg1")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "58ca893c-4bfe-00eb-93dd-4c0ad506eab5")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)


	var splunkResp SplunkResponse

	if res!=nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		splunkResp.Status = res.Status
		splunkResp.StatusCode = res.StatusCode
		splunkResp.ContentLength = res.ContentLength
		fmt.Print(splunkResp.StatusCode)
		splunkResp.Body = string(body)
	}

	return splunkResp, err
}