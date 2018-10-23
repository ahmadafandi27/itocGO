package main

import (
	"fmt"
	"strconv"
	"splunkadapter"
	"os"
	"time"
)


func main() {
	start := time.Now()
	fmt.Printf("start splunk parser\n")

	appname := "reflex premium"
	query := `sourcetype=tselreflexpremiumtdr
	| eval timex=strftime(_time,"%Y-%m-%d %H:%M")
	| eval timexepoch=round(strptime(timex,"%Y-%m-%d %H:%M"),0)
	| fillnull value=null host, adn, appsID, taskID, description_trx, clientID, statusCode, detailedCode, productCategoryLookup, productCategoryLookupBasis, productLookup
	| stats count as trx, sum(elapTime) as totRT, max(elapTime) as maxRT, min(elapTime) as minRT by timexepoch, host, adn, appsID, taskID, description_trx, clientID, statusCode, detailedCode, productCategoryLookup, productCategoryLookupBasis, productLookup`
	earliestTime := "-5m"
	latestTime := "-4m"

	responseSearch, errorSearch := splunkadapter.SplunkSearch(query,earliestTime,latestTime)
	if errorSearch != nil {
		fmt.Print(errorSearch.Error())
	}else{
		fmt.Printf(strconv.Itoa(responseSearch.StatusCode))
	}

	elapsed := time.Since(start)

	logString := appname+"|"+strconv.Itoa(responseSearch.StatusCode)+"|"+responseSearch.Status+"|"+fmt.Sprintf("%f", elapsed.Seconds())+"\n"
	fmt.Print(logString)
	ServiceLogging("../../log",logString)
}

func ServiceLogging(path, text string) error {
	logfile := path+"/splunkcollector.log"
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

