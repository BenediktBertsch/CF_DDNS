# How to use with Docker:
Use the environment variables to customize the Application for yourself. Here is the [Dockerimage](https://hub.docker.com/repository/docker/xaviius/ddns). The variables are: 

  `CF_TOKENS` = NEEDED, sets the Bearer Token
  
  
  `CF_ZONES` = NEEDED, sets the DNS Zone
  
  
  `CF_DOMAINS` = NEEDED, sets the Domain Name
  
  
  `CF_PROXIES` = Sets if the traffic should get proxied by Cloudflare (Default: false)
  
  
  `CF_IPV6` = Sets if AAAA Record should be updated (Default: false)
  
  
  `CF_INTERVAL` = Sets the minute interval in which the DNS gets checked (Default: 1)
  

To run it for example: docker run -e CF_TOKENS="token0,token1" -e CF_ZONES="zone0,zone1" -e CF_DOMAINS="test.com,test.net" xaviius/ddns:latest