package main

import (
	"github.com/BenediktBertsch/cf_ddns/httpclient"
	"strconv"
	"os"
	"strings"
	"time"
	"flag"
	"fmt"
)

///Global variables (Environment)

//TOKENS = Bearer Token for authorization
var TOKENS []string
//ZONES = Predefined Cloudflare DNS Zone
var ZONES []string
//DOMAINS = Domain ex: test.com
var DOMAINS []string
//PROXIES = if the traffic should get proxied by Cloudflare
var PROXIES []bool
//IPV6 = if you want to update IPv6 as well
var IPV6 []bool
//INTERVAL = Minuteinterval in which the DNS should get updated
var INTERVAL uint64

func main() {
	var ticker *time.Ticker
	//Only for debugging purposes or if you want to run without environment variables
	//setEnvVariables()
	if checkConfig(){
		runEnv()
		CheckDuration := flag.Duration("duration", time.Duration(INTERVAL), "update interval (ex. 15s, 1m, 6h); if not specified or set to 0s, run only once and exit")
		flag.Parse()
		ticker = time.NewTicker(*CheckDuration*time.Minute)
	
		runddns()
	
		for range ticker.C {
			runddns()
		}
	}
}

func runEnv(){
	//Split Env variables because of only string input
	splitEnvVariables()
}

func setEnvVariables(){
	os.Setenv("CF_TOKENS", "")
	os.Setenv("CF_ZONES", "")
	os.Setenv("CF_DOMAINS", "")
	os.Setenv("CF_PROXIES", "")
	os.Setenv("CF_IPV6", "true,true")
	os.Setenv("CF_INTERVAL", "")
}

func debugEnvVariables(){
	fmt.Println(os.Getenv("CF_TOKEN"))
	fmt.Println(os.Getenv("CF_ZONES"))
	fmt.Println(os.Getenv("CF_DOMAINS"))
	fmt.Println(os.Getenv("CF_PROXIES"))
	fmt.Println(os.Getenv("CF_IPV6"))
	fmt.Println(os.Getenv("CF_INTERVAL"))
}

func splitEnvVariables(){
	TOKENS = strings.Split(os.Getenv("CF_TOKENS"), ",")
	ZONES = strings.Split(os.Getenv("CF_ZONES"), ",")
	DOMAINS = strings.Split(os.Getenv("CF_DOMAINS"), ",")
	temp := strings.Split(os.Getenv("CF_PROXIES"), ",")
	for i := 0; i < len(TOKENS); i++ {
		if len(temp)-1 < i {
			PROXIES = append(PROXIES, false)
		} else {
			b, _ := strconv.ParseBool(temp[i])
			PROXIES = append(PROXIES, b)
		}
	}
	temp = strings.Split(os.Getenv("CF_IPV6"), ",")
	for i := 0; i < len(TOKENS); i++ {
		if len(temp)-1 < i {
			IPV6 = append(IPV6, false)
		} else {
			b, _ := strconv.ParseBool(temp[i])
			IPV6 = append(IPV6, b)
		}
	}
	INTERVAL, _ = strconv.ParseUint(os.Getenv("CF_INTERVAL"), 10, 64)
	_, err := strconv.ParseUint(os.Getenv("CF_INTERVAL"), 10, 64)
	if err != nil {
		INTERVAL = 1
	}
}

func runddns() {
	//GetIPv4 and if IPv6 enabled this as well
	ipv4, _ := httpclient.GetAddressIpv4()
	ipv6, _ := httpclient.GetAddressIpv6()
	//Loop over all Cloudflare data
	//First check if ENV data are set
	fmt.Println("Checking for updates:", time.Now().Format("15.01.2006 15:04:05"))
	for i := 0; i < len(TOKENS); i++ {
		var IDa = httpclient.CheckUpdate("A", ipv4, DOMAINS[i], ZONES[i], TOKENS[i])
		if IDa != "" {
			httpclient.Update(ZONES[i], IDa, TOKENS[i], ipv4, PROXIES[i], DOMAINS[i], "A", httpclient.PREVIOUSIP4)
		}else {
			fmt.Println("IPv4 of " + DOMAINS[i] + " is still the same.")
		}
		if IPV6[i] {
			var IDaaaa = httpclient.CheckUpdate("AAAA", ipv6, DOMAINS[i], ZONES[i], TOKENS[i])
			if IDaaaa != "" {
				httpclient.Update(ZONES[i], IDaaaa, TOKENS[i], ipv6, PROXIES[i], DOMAINS[i], "AAAA", httpclient.PREVIOUSIP6)
			} else {
				fmt.Println("IPv6 of " + DOMAINS[i] + " is still the same.")
			}
		}
	}
}

func checkConfig() bool {
	if os.Getenv("CF_TOKENS") == ""{
		fmt.Println("No CF_TOKENS set. This parameter is needed.")
		return false
	}
	if os.Getenv("CF_ZONES") == ""{
		fmt.Println("No CF_ZONES set. This parameter is needed.")
		return false
	}
	if os.Getenv("CF_DOMAINS") == ""{
		fmt.Println("No CF_DOMAINS set. This parameter is needed.")
		return false
	}
	if os.Getenv("CF_PROXIES") == ""{
		fmt.Println("No CF_PROXIES set. Will use default: false")
	}
	if os.Getenv("CF_IPV6") == ""{
		fmt.Println("No CF_IPV6 set. Will use default: false")
	}
	if os.Getenv("CF_INTERVAL") == ""{
		fmt.Println("No CF_INTERVAL set. Will use default: 1")
	}
	return true
}


