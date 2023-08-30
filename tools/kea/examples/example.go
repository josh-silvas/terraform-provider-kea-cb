package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

func main() {
	c := kea.New(kea.WithAuth(os.Getenv("KEA_USERNAME"), os.Getenv("KEA_PASSWORD")))

	res, err := c.RemoteSubnet4Set("kea-primary.example.com", []kea.NewRemoteSubnet4{
		{
			ID:     1921682270,
			Subnet: "192.168.227.0/24",
			Pools: []kea.Pool{
				{
					Pool: "192.168.227.110-192.168.227.120",
				},
			},
			OptionData: []kea.OptionData{
				{
					Name: "routers",
					Data: "192.168.227.1",
				},
				{
					Name:       "domain-name-servers",
					Data:       "4.2.2.2, 8.8.8.8",
					AlwaysSend: true,
				},
				{
					Name: "domain-name",
					Data: "example.com",
				},
			},
			Relay: kea.Relay{
				IPAddresses: []string{"192.168.227.1"},
			},
			UserContext: map[string]string{
				"site":        "San Antonio",
				"description": "San Antonio test subnet",
			},
		},
	})
	if err != nil {
		spew.Dump(err)
		return
	}
	spew.Dump(res)
}
