package httpclient

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
	"bytes"
)

//PREVIOUSIP4 = The previous used IPv4
var PREVIOUSIP4 string
//PREVIOUSIP6 = The previous used IPv4
var PREVIOUSIP6 string

//GetAddressIpv4 = Returns the current IPv4 of the Client
func GetAddressIpv4() (string, error) {
	address := address{}
	resp, err  := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&address)
	if err != nil {
		return "", err
	}
	return address.IP, err
}

//GetAddressIpv6 = Returns the current IPv4 of the Client
func GetAddressIpv6() (string, error) {
	address := address{}
	resp, err  := http.Get("https://api6.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&address)
	if err != nil {
		return "", err
	}
	return address.IP, err
}

type records struct {
	Result []record `json:"result"`
}

type record struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
}

type recordUpdate struct {
	Success bool `json:"success"`
	Errors []errors `json:"errors"`
}

type errors struct {
	Message string `json:"message"`
} 

type address struct {
	IP string `json:"ip"`
}

//Update = Updates the IP
func Update(zone string, id string, token string, ip string, proxy bool, domain string, recordtype string, previousip string) (string, error) {
	status := recordUpdate{}
	body := []byte(`{
		"type": "` + recordtype + `",
		"name": "` + domain + `",
		"content": "` + ip + `",
		"ttl": 120,
		"proxied": ` + strconv.FormatBool(proxy) + `
	}`)
	req, err := http.NewRequest("PUT", "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records/" + id, bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer " + token)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&status)
	if err != nil {
		return "", err
	}
	if status.Success {
		return "Domain: " + domain + " got updated with IP: " + ip + " Prvious IP: " + previousip, err
	} else {
		return "Error on updating the IP: " + status.Errors[0].Message, err
	}
}

//CheckUpdate = Checks if the update is needed, returns string if empty not needed else returns the zone which needs to be updated
func CheckUpdate(recordtype string, currentIP string, domain string, zone string, token string) (string, error) {
	record := records{}
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records?type=" + recordtype, nil)
	req.Header.Add("Authorization", "Bearer " + token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&record)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	for i := 0; i < len(record.Result); i++ {
		if record.Result[i].Content != currentIP && record.Result[i].Name == domain{
			switch recordtype {
			case "A": PREVIOUSIP4 = record.Result[i].Content
				break;
			case "AAAA": PREVIOUSIP6 = record.Result[i].Content
				break;
			}
			return record.Result[i].ID, err
		}
	}
	return "", err
}