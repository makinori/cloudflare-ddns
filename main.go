package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

var (
	currentIPV4 string
	currentIPV6 string
)

func simpleRecordTypeToString(recordType dns.RecordListParamsType) string {
	switch recordType {
	case dns.RecordListParamsTypeA:
		return "A"
	case dns.RecordListParamsTypeAAAA:
		return "AAAA"
	}
	return "?"
}

func updateRecord(
	client *cloudflare.Client, zone *zones.Zone,
	recordToUpdate string, ipv6 bool,
) {
	recordType := dns.RecordListParamsTypeA
	if ipv6 {
		recordType = dns.RecordListParamsTypeAAAA
	}

	recordQuery, err := client.DNS.Records.List(context.Background(),
		dns.RecordListParams{
			ZoneID:  cloudflare.F(zone.ID),
			Page:    cloudflare.Float(1),
			PerPage: cloudflare.Float(1),
			Type:    cloudflare.F(recordType),
			Name: cloudflare.F(dns.RecordListParamsName{
				Exact: cloudflare.F(recordToUpdate),
			}),
		},
	)

	if err != nil {
		log.Printf(
			"failed to query %s record: %s: %s\n",
			simpleRecordTypeToString(recordType), recordToUpdate,
			err.Error(),
		)
		return
	}

	recordExists := len(recordQuery.Result) > 0

	var verb string

	if recordExists {
		verb = "update"

		record := recordQuery.Result[0]

		var recordUpdateParams dns.RecordUpdateParamsBodyUnion

		// TODO: the whole point of params is that you dont have to copy data

		if !ipv6 {
			recordUpdateParams = dns.ARecordParam{
				Name:    cloudflare.F(recordToUpdate),
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(currentIPV4),
				// keep
				TTL:     cloudflare.F(record.TTL),
				Comment: cloudflare.F(record.Comment),
				Proxied: cloudflare.F(record.Proxied),
			}
		} else {
			recordUpdateParams = dns.AAAARecordParam{
				Name:    cloudflare.F(recordToUpdate),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
				Content: cloudflare.F(currentIPV6),
				// keep
				TTL:     cloudflare.F(record.TTL),
				Comment: cloudflare.F(record.Comment),
				Proxied: cloudflare.F(record.Proxied),
			}
		}

		_, err = client.DNS.Records.Update(
			context.Background(), record.ID, dns.RecordUpdateParams{
				ZoneID: cloudflare.F(zone.ID),
				Body:   recordUpdateParams,
			},
		)
	} else {
		verb = "create"

		var recordNewParams dns.RecordNewParamsBodyUnion

		if !ipv6 {
			recordNewParams = dns.ARecordParam{
				Name:    cloudflare.F(recordToUpdate),
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(currentIPV4),
			}
		} else {
			recordNewParams = dns.AAAARecordParam{
				Name:    cloudflare.F(recordToUpdate),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
				Content: cloudflare.F(currentIPV6),
			}
		}

		_, err = client.DNS.Records.New(context.Background(),
			dns.RecordNewParams{
				ZoneID: cloudflare.F(zone.ID),
				Body:   recordNewParams,
			},
		)
	}

	if err != nil {
		log.Printf(
			"failed to %s %s record: %s: %s\n",
			verb, simpleRecordTypeToString(recordType), recordToUpdate,
			err.Error(),
		)
	} else {
		log.Printf(
			"%sd %s record: %s\n",
			verb, simpleRecordTypeToString(recordType), recordToUpdate,
		)
	}
}

func updateZone(
	client *cloudflare.Client, zoneName string, recordsToUpdate []string,
) {
	zoneQuery, err := client.Zones.List(context.Background(), zones.ZoneListParams{
		Page:    cloudflare.Float(1),
		PerPage: cloudflare.Float(1),
		Status:  cloudflare.F(zones.ZoneListParamsStatusActive),
		// TODO: defaults to equal but how to explicitly state?
		Name: cloudflare.F(zoneName),
	})

	if err != nil {
		log.Printf(
			"failed to query zone: %s: %s\n",
			zoneName, err.Error(),
		)
		return
	}

	if len(zoneQuery.Result) == 0 {
		log.Printf("failed to find zone: %s\n", zoneName)
		return
	}

	zone := zoneQuery.Result[0]

	for _, recordToUpdate := range recordsToUpdate {
		go updateRecord(client, &zone, recordToUpdate, false)
		if settings.IPV6 {
			go updateRecord(client, &zone, recordToUpdate, true)
		}
	}
}

func updateAccount(account *Account) {
	client := cloudflare.NewClient(
		option.WithAPIEmail(account.Email),
		option.WithAPIKey(account.Key),
	)

	for zone, recordsToUpdate := range account.Zones {
		go updateZone(client, zone, recordsToUpdate)
	}
}

func onInterval() {
	ipv4, err := getMyIP(settings.IPV6)
	if err != nil {
		log.Println("failed to get ipv4: " + err.Error())
		return
	}

	var ipv6 string
	if settings.IPV6 {
		ipv6, err = getMyIP(settings.IPV6)
		if err != nil {
			log.Println("failed to get ipv6: " + err.Error())
			return
		}
	}

	if currentIPV4 != ipv4 || currentIPV6 != ipv6 {
		currentIPV4 = ipv4
		currentIPV6 = ipv6

		log.Println("new ipv4: " + ipv4)
		if settings.IPV6 {
			log.Println("new ipv6: " + ipv4)
		}

		for i := range settings.Accounts {
			go updateAccount(&settings.Accounts[i])
		}
	}
}

func main() {
	loadSettings()

	log.Printf("interval set to %d minutes\n", settings.Interval)

	ticker := time.NewTicker(
		time.Minute * time.Duration(settings.Interval),
	)

	onInterval()
	for {
		<-ticker.C
		onInterval()
	}
}
