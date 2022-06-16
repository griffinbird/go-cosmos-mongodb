package main

type CustomerV1 struct {
	/*
	{
		"_id": "3200B375-08DC-4A36-9F6F-E325A0B46550",
		"addresses": [
			{
				"addressLine1": "3709 Leonard Ct.",
				"addressLine2": "",
				"city": "Lake Oswego",
				"country": "US",
				"state": "OR ",
				"zipCode": "97034"
			}
		],
		"creationDate": "2013-08-18T00:00:00",
		"customerId": "3200B375-08DC-4A36-9F6F-E325A0B46550",
		"emailAddress": "emma18@adventure-works.com",
		"firstName": "Emma",
		"lastName": "Clark",
		"password": {
			"hash": "2Ox+0p9j/+nud/hqMUOJsXmxN8J8uBxzQypRifldTvg=",
			"salt": "20C44843"
		},
		"phoneNumber": "155-555-0135",
		"salesOrderCount": 1,
		"title": "",
		"type": "customer"
	}
	*/
	ID string `bson:"_id" json:"_id"`
	Addresses []struct{
		AddressLine1 string
		AddressLine2 string
		City string
		Country string
		State string
		ZipCode string
	}
	CreationDate string
	CustomerID string 
	EmailAddress string
	FirstName string
	LastName string
}

type RequestStatistics struct {
	ActivityID                                  string
	CommandName                                 string
	EstimateDelayFromRateLimitingInMilliseconds int
	RequestCharge                               float64
	RequestDurationInMilliseconds               int
	RetiredDueToRateLimiting                    bool
	OK                                          int
}