package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)
const urlVA = "http://45.127.134.13:8082/"
const skeyheader = "skey"
const secretkey ="9casisd8emg49ops"

func GetVA(va string) string{
	var json string
	var mapper  ="{\"vaNo\":\""+va+"\"}"
	bt  := GenerateSecureIVVector()
	encbt := Encrypt([]byte(GetVectorIVAsString(bt)),[]byte(secretkey))
	mapper = AESEncryptWithIV(bt,mapper,secretkey)
	json = "{\"val\":\""+encbt+"\" , \"body\":\""+mapper+"\"}"
	fmt.Print("Request Body ")
	fmt.Println(json)

	url := urlVA+"vaapps/belanjaonline/getva"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("key", skeyheader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)
}
func CreateVA(mapper string) string{
	var json string
	bt  := GenerateSecureIVVector()
	encbt := Encrypt([]byte(GetVectorIVAsString(bt)),[]byte(secretkey))
	mapper = AESEncryptWithIV(bt,mapper,secretkey)
	json = "{\"val\":\""+encbt+"\" , \"body\":\""+mapper+"\"}"
	fmt.Print("Request Body ")
	fmt.Println(json)

	url := urlVA+"vaapps/belanjaonline/newva"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("key", skeyheader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println(resp.StatusCode)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	if(resp.StatusCode!=200){
		return fmt.Sprintf("{\"message\":\"Failed creating VA with resp: "+"aa"+resp.Status+" %v\"}",resp.StatusCode)
	}
	respafterdecrypt := DecrypRespBodyAPI(string(body))
	fmt.Println("response Body after decrypt:",respafterdecrypt)
	return string(respafterdecrypt)
}

func GeneratReqBodyForCreateVA(amount ,description,name,va string)string {
	currentTime := time.Now().Add(72*time.Hour)
	var expiredDate = currentTime.Format(time.RFC3339)
	fmt.Println("Current Time in String: ", currentTime.Format(time.RFC3339))
	var body = "{\"amount\": "+amount+", \"description\": \""+description+"\",\"expiredDate\": \""+expiredDate+"\",\"idMerchant\": {\"id\": \"10001\"},\"name\":\""+name+"\",\"vaNo\": \""+va+"\"}"
	return body
}