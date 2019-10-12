package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"strings"
)

func updateDNS(hostname string, ipv4 string, ipv6 string) error {
	hostparts := strings.Split(hostname, ".")
	length := len(hostparts)

	if length < 3 {
		return errors.New(fmt.Sprintf("hostname \"%s\" does not have enought parts", hostname))
	}

	zone := fmt.Sprintf("%s.%s.", hostparts[length-2], hostparts[length-1])

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(conf.AWS.UserID, conf.AWS.Key, ""),
	})

	if err != nil {
		return err
	}

	svc := route53.New(sess)

	hostedZones, err := svc.ListHostedZonesByName(nil)
	if err != nil {
		return err
	}

	matchingZone := &route53.HostedZone{}

	for _, hostedZone := range hostedZones.HostedZones {
		if zone == *hostedZone.Name {
			matchingZone = hostedZone
			break
		}
	}

	if matchingZone == nil {
		return errors.New(fmt.Sprintf("Could not find zone: %s", zone))
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(hostname),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ipv4),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("A"),
					},
				},
			},
		},
		HostedZoneId: matchingZone.Id,
	}

	_, err = svc.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}

	return nil
}