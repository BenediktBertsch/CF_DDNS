package main

import (
	"github.com/BenediktBertsch/cf_ddns/src/httpclient"
    "testing"
)

//APITest = Unit test if the ipv4 works.
func TestGetAddressIpv4(t *testing.T) {
    _, err := httpclient.GetAddressIpv4()
    if err != nil {
        t.Errorf("Failed: Could not fetch the IP String!")
    } else {
        t.Logf("Success: Could fetch the IP String!")
    }
}

//APITest = Unit test if the ipv6 works.
func TestGetAddressIpv6(t *testing.T) {
    _, err := httpclient.GetAddressIpv6()
    if err != nil {
        t.Errorf("Failed: Could not fetch the IP String!")
    } else {
        t.Logf("Success: Could fetch the IP String!")
    }
}