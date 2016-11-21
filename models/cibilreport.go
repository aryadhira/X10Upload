package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type Profile struct {
	CompanyName  string
	CustomerId   int
	DealNo       string
	DunsNumber   string
	Pan          string
	Address      string
	CityTown     string
	Telephone    string
	StateUnion   string
	PinCode      string
	Country      string
	FileOpenDate string
}

type ReportSummary struct {
	Grantors                     string
	Facilities                   string
	CreditFacilities             string
	FacilitiesGuaranteed         string
	LatestCreditFacilityOpenDate string
	FirstCreditFacilityOpenDate  string
}

type ReportSummaryDetail struct {
	CreditFacilities                string
	NoOfStandard                    string
	CurrentBalanceStandard          string
	NoOfOtherThanStandard           string
	CurrentBalanceOtherThanStandard string
	NoOfLawSuits                    string
	NoOfWilfulDefaults              string
}

type CreditFacilityBorrower struct {
	NoStandard                      string
	CurrentBalanceStandard          string
	NoOtherThanStandard             string
	CurrentBalanceOtherThanStandard string
	NoLawSuits                      string
	NoWilfulDefaults                string
}

type CreditFacilityGuarantor struct {
	NoStandard                      string
	CurrentBalanceStandard          string
	NoOtherThanStandard             string
	CurrentBalanceOtherThanStandard string
	NoLawSuits                      string
	NoWilfulDefaults                string
}

type CreditTypeSummary struct {
	NoCreditFacilitiesBorrower string
	CreditType                 string
	CurrencyCode               string
	Standard                   string
	Substandard                string
	Doubtful                   string
	Loss                       string
	SpecialMention             string
	TotalCurrentBalance        string
}

type EnquirySummary struct {
	Enquiries3Month      string
	Enquiries6Month      string
	Enquiries9Month      string
	Enquiries12Month     string
	Enquiries24Month     string
	Enquiriesthan24Month string
	TotalEnquiries       string
	MostRecentDate       string
}

type CibilReportModel struct {
	orm.ModelBase     `bson:"-",json:"-"`
	Id                bson.ObjectId         `bson:"_id" , json:"_id"`
	FilePath          string                `bson:"FilePath"`
	FileName          string                `bson:"FileName"`
	ReportType        string                `bson:"ReportType"`
	IsMatch           bool                  `bson:"IsMatch"`
	Profile           Profile               `bson:"Profile"`
	ReportSummary     ReportSummary         `bson:"ReportSummary"`
	Detail            []ReportSummaryDetail `bson:"DetailReportSummary"`
	CreditTypeSummary []CreditTypeSummary   `bson:"CreditTypeSummary"`
	EnquirySummary    EnquirySummary        `bson:"EnquirySummary"`
}
