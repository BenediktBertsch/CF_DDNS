package main

import (
	"strconv"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
	"encoding/json"
	"time"
	"flag"
	"fmt"
)

//Global variables (Environment)
var TOKENS []string
var MAILS []string
var ZONES []string
var DOMAINS []string
var PROXIES []bool
var IPV6 []bool
var INTERVAL uint64

func main() {
	var runonce bool
	var ticker *time.Ticker
	runEnv()
	CheckDuration := flag.Duration("duration", time.Duration(INTERVAL), "update interval (ex. 15s, 1m, 6h); if not specified or set to 0s, run only once and exit")
	flag.Parse()
	if *CheckDuration == time.Duration(0) {
		runonce = true
	} else {
		ticker = time.NewTicker(*CheckDuration*time.Minute)
	}
	runddns()

	if runonce {
		os.Exit(0)
	}

	for range ticker.C {
		runddns()
	}
}

func runEnv(){
	//Split Env variables because of only string input
	splitEnvVariables()
}

func setEnvVariables(){
	os.Setenv("CF_TOKENS", "")
	os.Setenv("CF_MAILS", "")
	os.Setenv("CF_ZONES", "")
	os.Setenv("CF_DOMAINS", "")
	os.Setenv("CF_PROXIES", "")
	os.Setenv("CF_IPV6", "")
	os.Setenv("CF_INTERVAL", "")
}

func debugEnvVariables(){
	println(os.Getenv("CF_TOKENS"))
	println(os.Getenv("CF_MAILS"))
	println(os.Getenv("CF_ZONES"))
	println(os.Getenv("CF_DOMAINS"))
	println(os.Getenv("CF_PROXIES"))
	println(os.Getenv("CF_IPV6"))
	println(os.Getenv("CF_INTERVAL"))
}

func splitEnvVariables(){
	TOKENS = strings.Split(os.Getenv("CF_TOKENS"), ",")
	MAILS = strings.Split(os.Getenv("CF_MAILS"), ",")
	ZONES = strings.Split(os.Getenv("CF_ZONES"), ",")
	DOMAINS = strings.Split(os.Getenv("CF_DOMAINS"), ",")
	var temp = strings.Split(os.Getenv("CF_PROXIES"), ",")
	for i := 0; i < len(temp); i++ {
		b, err := strconv.ParseBool(temp[i])
		if err != nil{
			println(err)
		}
		PROXIES = append(PROXIES,b)
	}
	temp = strings.Split(os.Getenv("CF_IPV6"), ",")
	for i := 0; i < len(temp); i++ {
		b, err := strconv.ParseBool(temp[i])
		if err != nil{
			println(err)
		}
		IPV6 = append(IPV6,b)
	}
	INTERVAL, _ = strconv.ParseUint(os.Getenv("CF_INTERVAL"), 10, 64)
}

func runddns() {
	//GetIPv4 and if IPv6 enabled this as well
	var ipv4 = getAddressIpv4()
	var ipv6 = getAddressIpv6()
	//Loop over all Cloudflare data
    fmt.Println("Checking for updates:", time.Now().Format("15.01.2006 15:04:05"))
	for i := 0; i < len(TOKENS); i++ {
		var IDa = checkUpdate("A", ipv4, DOMAINS[i], ZONES[i], MAILS[i], TOKENS[i])
		if IDa != "" {
			update(ZONES[i], IDa, MAILS[i], TOKENS[i], ipv4, PROXIES[i], DOMAINS[i], "A")
		}else {
			print("IPv4 of " + DOMAINS[i] + " is still the same.")
		}
		if IPV6[i] {
			var IDaaaa = checkUpdate("AAAA", ipv6, DOMAINS[i], ZONES[i], MAILS[i], TOKENS[i])
			if IDaaaa != "" {
				update(ZONES[i], IDaaaa, MAILS[i], TOKENS[i], ipv6, PROXIES[i], DOMAINS[i], "AAAA")
			} else {
				println("IPv6 of " + DOMAINS[i] + " is still the same.")
			}
		}
	}
}

func update(zone string, id string, mail string, token string, ip string, proxy bool, domain string, recordtype string){
	client := resty.New()
	var status = recordUpdate{}
	resp, _ := client.R().SetHeaders(map[string]string{
        "Content-Type": "application/json",
		"X-Auth-Email": mail,
		"X-Auth-Key": token,
	  }).SetBody(`{
		"type": "` + recordtype + `",
		"name": "` + domain + `",
		"content": "` + ip + `",
		"ttl": 120,
		"proxied": ` + strconv.FormatBool(proxy) + `
	}`).Put("https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records/" + id)
	err := json.Unmarshal(resp.Body(), &status)
	if err != nil {
        println(err)
	}
	if status.Success {
		println("Domain: " + domain + " got updated with IP: " + ip)
	} else {
		println("Error on updating the IP: " + status.Errors[0].Message)
	}
}

func checkUpdate(recordtype string, currentIP string, domain string, zone string, mail string, token string) string {
	client := resty.New()
	var record = records{}
	resp, _  := client.R().SetHeaders(map[string]string{
        "Content-Type": "application/json",
		"X-Auth-Email": mail,
		"X-Auth-Key": token,
	  }).Get("https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records?type=" + recordtype)
	err := json.Unmarshal(resp.Body(), &record)
	if err != nil {
        println(err)
	}
	for i := 0; i < len(record.Result); i++ {
		if record.Result[i].Content != currentIP && record.Result[i].Name == domain{
			return record.Result[i].ID
		}
	}
	return ""
}

func getAddressIpv4() string{
	client := resty.New()
	address := address{}
	resp, _  := client.R().Get("https://api.ipify.org?format=json")
	err := json.Unmarshal(resp.Body(), &address)
	if err != nil {
        println(err)
	}
	return address.IP
}

func getAddressIpv6() string{
	client := resty.New()
	address := address{}
	resp, _  := client.R().Get("https://api6.ipify.org?format=json")
	err := json.Unmarshal(resp.Body(), &address)
	if err != nil {
        println(err)
	}
	return address.IP
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