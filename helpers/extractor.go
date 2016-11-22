package x10upload

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "eaciit/x10upload/models"

	"github.com/eaciit/dbox"
	. "github.com/eaciit/textsearch"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type Text struct {
	Top     int    `xml:"top,attr"`
	Left    int    `xml:"left,attr"`
	Content string `xml:",chardata"`
	Width   int    `xml:"width,attr"`
	Inline  string `xml:",innerxml"`
}

type Page struct {
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Texts  []Text `xml:"text"`
}

type Pdf2xml struct {
	Pages []Page `xml:"page"`
}

func ExtractCompanyCibilReport(PathFrom string, Filename string) CibilReportModel {
	XmlFile := PathFrom + "/" + Filename
	v := &Pdf2xml{}
	rawdata, err := ioutil.ReadFile(XmlFile)
	if err != nil {
		tk.Println(err.Error())
	}
	xml.Unmarshal(rawdata, &v)

	topcreditsummarylayout := 0
	botcreditsummarylayout := 0
	topageindexcreditsummary := 0
	botpageindexcreditsummary := 0
	Profiles := Profile{}
	ReportSummarys := ReportSummary{}
	CreditTypeSummarys := []CreditTypeSummary{}
	CreditTypeSummaryData := CreditTypeSummary{}
	ReportSummaryDetails := ReportSummaryDetail{}
	DetailReportSummary := []ReportSummaryDetail{}
	EnquirySummarys := EnquirySummary{}
	CreditType := ""
	Address := ""
	nametop := 0
	pantop := 0
	citytop := 0
	statetop := 0
	countrytop := 0
	dunsnumbertop := 0
	addresstop := 0
	telephonetop := 0
	pincodetop := 0
	fileopentop := 0
	creditgranttop := 0
	creditfacilityguaranttop := 0
	standardtop := 0
	topenquirysummarylayout := 0
	topenquirysummaryindex := 0

	//Create Layout
	for i, page := range v.Pages {
		for _, text := range page.Texts {
			if text.Content == "Mention A/C" {
				topcreditsummarylayout = text.Top
				topageindexcreditsummary = i
			}
			if text.Content == "Enquiry Summary" {
				botcreditsummarylayout = text.Top
				topenquirysummaryindex = i
				botpageindexcreditsummary = i
			}
			if text.Content == "No. of Enquiries" {
				topenquirysummarylayout = text.Top
			}
		}
	}

	for _, page := range v.Pages {
		for _, text := range page.Texts {
			//Extract Profile
			if text.Content == "Name" {
				nametop = text.Top
			}
			if text.Content == "PAN" && pantop == 0 {
				pantop = text.Top
			}
			if text.Content == "City / Town" && citytop == 0 {
				citytop = text.Top
			}
			if text.Content == "State / Union Territory" && statetop == 0 {
				statetop = text.Top
			}
			if text.Content == "Country" && countrytop == 0 {
				countrytop = text.Top
			}
			if text.Content == "Short Name" && dunsnumbertop == 0 {
				dunsnumbertop = text.Top
			}
			if text.Content == "Address" && addresstop == 0 {
				addresstop = text.Top
			}
			if text.Content == "Telephone Number" && telephonetop == 0 {
				telephonetop = text.Top
			}
			if text.Content == "PIN Code" && pincodetop == 0 {
				pincodetop = text.Top
			}
			if text.Content == "File Open Date" && fileopentop == 0 {
				fileopentop = text.Top
			}
			if text.Content == "No. of Credit Grantors" {
				creditgranttop = text.Top
			}
			if text.Content == "No. of Credit Facilities" {
				creditfacilityguaranttop = text.Top
			}
		}
	}
	//End Of Create Layout

	for i, page := range v.Pages {
		for _, text := range page.Texts {
			if i == 0 {
				if text.Top == nametop && text.Left == 275 {
					Profiles.CompanyName = text.Content
				}
				if text.Top == pantop && text.Left == 275 {
					Profiles.Pan = text.Content
				}
				if (text.Top == citytop || text.Top == citytop-1) && text.Left == 275 {
					Profiles.CityTown = text.Content
				}
				if (text.Top == statetop || text.Top == statetop-1) && text.Left == 275 {
					Profiles.StateUnion = text.Content
				}
				if (text.Top == countrytop || text.Top == countrytop-1) && text.Left == 275 {
					Profiles.Country = text.Content
				}
				if text.Top == dunsnumbertop && text.Left == 626 {
					Profiles.DunsNumber = text.Content
				}
				if text.Top >= addresstop && text.Top < telephonetop {
					if text.Left == 626 {
						Address = Address + " " + text.Content
						Profiles.Address = Address
					}
				}
				if (text.Top == telephonetop || text.Top == telephonetop-1) && text.Left == 626 {
					Profiles.Telephone = text.Content
				}
				if (text.Top == pincodetop || text.Top == pincodetop-1) && text.Left == 626 {
					Profiles.PinCode = text.Content
				}
				if (text.Top == fileopentop || text.Top == fileopentop-1) && text.Left == 626 {
					Profiles.FileOpenDate = text.Content
				}
				if text.Content == "Standard" {
					standardtop = text.Top
				}
			}
			//End Of Extract Profile

			//Extract Report Summary
			if (text.Top == creditgranttop-1 || text.Top == creditgranttop) && text.Left == 248 {
				ReportSummarys.Grantors = text.Content
			}
			if (text.Top == creditgranttop-1 || text.Top == creditgranttop) && text.Left == 469 {
				ReportSummarys.Facilities = text.Content
			}
			if (text.Top == creditgranttop-1 || text.Top == creditgranttop) && text.Left == 691 {
				ReportSummarys.CreditFacilities = text.Content
			}
			if (text.Top == creditfacilityguaranttop-1 || text.Top == creditfacilityguaranttop) && text.Left == 248 {
				ReportSummarys.FacilitiesGuaranteed = text.Content
			}
			if (text.Top == creditfacilityguaranttop-1 || text.Top == creditfacilityguaranttop) && text.Left == 469 {
				ReportSummarys.LatestCreditFacilityOpenDate = text.Content
			}
			if (text.Top == creditfacilityguaranttop-1 || text.Top == creditfacilityguaranttop) && text.Left == 691 {
				ReportSummarys.FirstCreditFacilityOpenDate = text.Content
			}
			if i == 0 {
				if text.Top > standardtop && text.Top < topcreditsummarylayout-21 {
					if text.Left >= 100 && text.Left <= 724 {
						if text.Left == 100 {
							if text.Content != "Credit Type Summary" {
								ReportSummaryDetails.CreditFacilities = text.Content
							}
						}
						if text.Left == 200 {
							ReportSummaryDetails.NoOfStandard = text.Content
						}
						if text.Left == 300 {
							ReportSummaryDetails.CurrentBalanceStandard = text.Content
						}
						if text.Left == 401 {
							ReportSummaryDetails.NoOfOtherThanStandard = text.Content
						}
						if text.Left == 501 {
							ReportSummaryDetails.CurrentBalanceOtherThanStandard = text.Content
						}
						if text.Left == 601 {
							ReportSummaryDetails.NoOfLawSuits = text.Content
						}
						if text.Left == 701 {
							ReportSummaryDetails.NoOfWilfulDefaults = text.Content
							DetailReportSummary = append(DetailReportSummary, ReportSummaryDetails)
						}
					}
				}
			}
			//End Of Extract Report Summary

			//Extract Credit Type Summary
			if topageindexcreditsummary != botpageindexcreditsummary {
				if i == topageindexcreditsummary {
					if text.Top > topcreditsummarylayout {
						if text.Left >= 100 && text.Left <= 724 {
							if text.Left == 100 {
								CreditTypeSummaryData.NoCreditFacilitiesBorrower = text.Content
							}
							if text.Left == 178 {
								CreditType = CreditType + " " + text.Content
								CreditTypeSummaryData.CreditType = CreditType
							}
							if text.Left == 256 {
								CreditTypeSummaryData.CurrencyCode = text.Content
							}
							if text.Left == 334 {
								CreditTypeSummaryData.Standard = text.Content
							}
							if text.Left == 724 {
								CreditTypeSummaryData.TotalCurrentBalance = text.Content
								if CreditTypeSummaryData.CurrencyCode != "Total" {
									CreditTypeSummarys = append(CreditTypeSummarys, CreditTypeSummaryData)
								}
								CreditType = ""
							}
						}
					}
				}

				if i == botpageindexcreditsummary {
					if text.Top < botcreditsummarylayout {
						if text.Left >= 100 && text.Left <= 724 {
							if text.Left == 100 {
								CreditTypeSummaryData.NoCreditFacilitiesBorrower = text.Content
							}
							if text.Left == 178 {
								CreditType = CreditType + " " + text.Content
								CreditTypeSummaryData.CreditType = CreditType
							}
							if text.Left == 256 {
								CreditTypeSummaryData.CurrencyCode = text.Content
							}
							if text.Left == 334 {
								CreditTypeSummaryData.Standard = text.Content
							}
							if text.Left == 724 {
								CreditTypeSummaryData.TotalCurrentBalance = text.Content
								if CreditTypeSummaryData.CurrencyCode != "Total" {
									CreditTypeSummarys = append(CreditTypeSummarys, CreditTypeSummaryData)
								}
								CreditType = ""
							}
						}
					}
				}
			} else {
				if i == topageindexcreditsummary {
					if text.Top > topcreditsummarylayout {
						if text.Left >= 100 && text.Left <= 724 {
							if text.Left == 100 {
								CreditTypeSummaryData.NoCreditFacilitiesBorrower = text.Content
							}
							if text.Left == 178 {
								CreditType = CreditType + " " + text.Content
								CreditTypeSummaryData.CreditType = CreditType
							}
							if text.Left == 256 {
								CreditTypeSummaryData.CurrencyCode = text.Content
							}
							if text.Left == 334 {
								CreditTypeSummaryData.Standard = text.Content
							}
							if text.Left == 724 {
								CreditTypeSummaryData.TotalCurrentBalance = text.Content
								if CreditTypeSummaryData.CurrencyCode != "Total" {
									CreditTypeSummarys = append(CreditTypeSummarys, CreditTypeSummaryData)
								}
								CreditType = ""
							}
						}
					}
				}
			}
			//End Of Credit Type Summary
			//Extract Enquiry Summary
			if i == topenquirysummaryindex {
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 205 {
					EnquirySummarys.Enquiries3Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 275 {
					EnquirySummarys.Enquiries6Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 345 {
					EnquirySummarys.Enquiries9Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 416 {
					EnquirySummarys.Enquiries12Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 486 {
					EnquirySummarys.Enquiries24Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 556 {
					EnquirySummarys.Enquiriesthan24Month = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 626 {
					EnquirySummarys.TotalEnquiries = text.Content
				}
				if (text.Top == topenquirysummarylayout || text.Top == topenquirysummarylayout-1) && text.Left == 696 {
					EnquirySummarys.MostRecentDate = text.Content
				}
			}
			//End Of Extract Enquiry Summary
		}
	}

	CibilReport := CibilReportModel{}
	CibilReport.ReportType = "Company"
	CibilReport.Profile = Profiles
	CibilReport.Detail = DetailReportSummary
	CibilReport.ReportSummary = ReportSummarys
	CibilReport.CreditTypeSummary = CreditTypeSummarys
	CibilReport.EnquirySummary = EnquirySummarys
	return CibilReport

}

func ExtractIndividualCibilReport(PathFrom string, Filename string) ReportData {
	XmlFile := PathFrom + "/" + Filename
	v := &Pdf2xml{}
	rawdata, err := ioutil.ReadFile(XmlFile)
	if err != nil {
		tk.Println(err.Error())
	}
	xml.Unmarshal(rawdata, &v)

	nametop := 0
	nameleft := 96
	bodtop := 0
	bodleft := 128
	gendertop := 0
	genderleft := 508
	datereporttop := 0
	datereportleft := 712
	timereporttop := 0
	timereportleft := 711
	cibilscoreversiontop := 450
	cibilscoreversionleft := 42
	cibilscoretop := 432
	cibilscoreleft := 248
	scoringfactortop := 440
	scoringfactorleft := 331
	scoringfactorbot := 0
	ishavescoringfactor := false
	incometaxtop := 0
	incometaxleft := 291
	passportnumbertop := 0
	passportnumberleft := 291
	telephonetop := 0
	telephonebot := 0
	emailtop := 0
	emailbot := 0
	emailtopindex := 0
	emailbotindex := 0
	addresstop := 0
	addressbot := 0
	addresstopindex := 0
	addressbotindex := 0
	accounttop := 0
	accountbot := 0
	totalacctop := 0
	totalaccleft := 289
	overduetop := 0
	overdueleft := 289
	zerobalancetop := 0
	zerobalanceleft := 289
	highcreditsanctiontop := 0
	highcreditsancright := 0
	//highcreditsanctionleft := 465
	currentbalancetop := 0
	currentbalanceright := 0
	currentbalanceleft := 631
	overduebalancetop := 0
	overduebalanceright := 0
	//overduebalanceleft := 648
	dateopenrecenttop := 0
	dateopenrecentleft := 800
	dateopenoldesttop := 0
	dateopenoldestleft := 800
	enquirytop := 0
	enquiryright := 0
	//enquirytotalleft := 322
	//enquirypast30left := 455
	enquirypast30right := 0
	enquiryrecentdateleft := 800
	enquirybot := 0
	addressdetailtop := 0
	addressdetailpermanentleft := 124
	addressdetailleft := 109
	addresscategorytop := 0
	addresscategoryleft := 119
	addressdatereportedtop := 0
	addressdatereportedleft := 706
	scoringfactors := []string{}
	telephonedata := ReportTelephone{}
	telephones := []ReportTelephone{}
	emails := []string{}
	addressdetail := ReportAddress{}
	addressdetails := []ReportAddress{}
	consumerinfo := ConsumerInfo{}
	reportdata := ReportData{}
	layout := "02-01-2006"
	layoutdatetime := "15:04:05"

	//Create Layout
	for i, page := range v.Pages {
		for _, text := range page.Texts {
			if text.Inline == "<b> CONSUMER: </b>" && nametop == 0 {
				nametop = text.Top - 1
			}
			if text.Inline == "<i>DATE OF BIRTH: </i>" && bodtop == 0 {
				bodtop = text.Top - 1
			}
			if text.Inline == "<i>GENDER: </i>" && gendertop == 0 {
				gendertop = text.Top - 1
			}
			if text.Inline == "<b>DATE:</b>" && datereporttop == 0 {
				datereporttop = text.Top - 1
			}
			if text.Inline == "<b>TIME: </b>" && timereporttop == 0 {
				timereporttop = text.Top - 1
			}
			if (text.Top == scoringfactortop || text.Top == 458) && text.Left == scoringfactorleft && text.Content != "" {
				ishavescoringfactor = true
			}
			if text.Content == "POSSIBLE RANGE FOR CIBIL TRANSUNION SCORE VERSION 1.0" && scoringfactorbot == 0 {
				scoringfactorbot = text.Top
			}
			if text.Content == "INCOME TAX ID NUMBER (PAN) " && incometaxtop == 0 {
				incometaxtop = text.Top
			}
			if text.Content == "PASSPORT NUMBER " && passportnumbertop == 0 {
				passportnumbertop = text.Top
			}
			if text.Inline == "<b>TELEPHONE TYPE</b>" && telephonetop == 0 {
				telephonetop = text.Top
			}
			if text.Inline == "<b>EMAIL CONTACT(S):</b>" && telephonebot == 0 {
				telephonebot = text.Top
			}
			if text.Inline == "<b>EMAIL CONTACT(S):</b>" && emailtop == 0 {
				emailtop = text.Top
				emailtopindex = i
			}
			if text.Inline == "<b>ADDRESS(ES): </b>" && emailbot == 0 {
				emailbot = text.Top
				emailbotindex = i
			}
			if text.Inline == "<b>ADDRESS(ES): </b>" && addresstop == 0 {
				addresstop = text.Top
				addresstopindex = i
			}
			if text.Inline == "<b>ACCOUNT TYPE</b>" && addressbot == 0 {
				addressbot = text.Top
				addressbotindex = i
			}
			if text.Inline == "<b>ACCOUNT(S) </b>" && accounttop == 0 {
				accounttop = text.Top
			}
			if text.Inline == "<i>TOTAL: </i>" && totalacctop == 0 {
				totalacctop = text.Top - 1
			}
			if text.Inline == "<i>OVERDUE: </i>" && overduetop == 0 {
				overduetop = text.Top - 1
			}
			if text.Inline == "<i>ZERO-BALANCE: </i>" && zerobalancetop == 0 {
				zerobalancetop = text.Top - 1
			}
			if text.Inline == "<i>HIGH CR/SANC. AMT: </i>" && highcreditsanctiontop == 0 {
				highcreditsanctiontop = text.Top - 1
				highcreditsancright = text.Left + text.Width
			}
			if text.Inline == "<i>CURRENT: </i>" && currentbalancetop == 0 {
				currentbalancetop = text.Top - 1
				currentbalanceright = text.Left + text.Width
			}
			if text.Inline == "<i>OVERDUE: </i>" && overduebalancetop == 0 {
				overduebalancetop = text.Top - 1
				overduebalanceright = text.Left + text.Width
			}
			if text.Inline == "<i>RECENT: </i>" && dateopenrecenttop == 0 {
				dateopenrecenttop = text.Top - 1
			}
			if text.Inline == "<i>OLDEST: </i>" && dateopenoldesttop == 0 {
				dateopenoldesttop = text.Top - 1
			}
			if text.Inline == "<b>ENQUIRIES </b>" && accountbot == 0 {
				accountbot = text.Top
			}
			// if text.Inline == "<b>ENQUIRIES </b>" && enquirytop == 0 {
			// 	enquirytop = text.Top
			// }
			if text.Inline == "<b>All Enquiries</b>" && enquirybot == 0 {
				enquirybot = text.Top
			}
			if text.Inline == "<b>TOTAL </b>" && enquirytop == 0 {
				enquirytop = text.Top
				enquiryright = text.Left + text.Width
			}
			if text.Inline == "<b>PAST 30 DAYS </b>" {
				enquirypast30right = text.Left + text.Width
			}

		}
	}

	for i, page := range v.Pages {
		for _, text := range page.Texts {
			// if (text.Top == nametop || text.Top == 235) && text.Left == nameleft {
			// 	tk.Println(text.Top)
			// 	consumerinfo.ConsumerName = text.Content
			// }
			if text.Top == bodtop && text.Left == bodleft {
				bodval, _ := time.Parse(layout, text.Content)
				consumerinfo.DateOfBirth = bodval
			}
			if text.Top == gendertop && text.Left == genderleft {
				consumerinfo.Gender = text.Content
			}
			if i == 0 {
				if (text.Top == nametop || text.Top == 235) && text.Left == nameleft {
					consumerinfo.ConsumerName = text.Content
				}
				if text.Top == datereporttop && text.Left == datereportleft {
					//dateval = dateval + text.Content
					dates, _ := time.Parse(layout, text.Content)
					reportdata.DateOfReport = dates
				}
				if text.Top == timereporttop && text.Left == timereportleft {
					times, _ := time.Parse(layoutdatetime, text.Content)
					reportdata.TimeOfReport = times
				}
				if text.Top == passportnumbertop && text.Left == passportnumberleft {
					reportdata.PassportNumber = text.Content
				}
			}
			if text.Top == cibilscoreversiontop && text.Left == cibilscoreversionleft {
				reportdata.CibilScoreVersion = text.Content
			}
			if (text.Top == cibilscoretop || text.Top == 450) && text.Left == cibilscoreleft {
				score, _ := strconv.Atoi(text.Content)
				reportdata.CibilScore = score
			}
			if ishavescoringfactor == true {
				if text.Top >= scoringfactortop-2 && text.Top < scoringfactorbot && text.Left == 345 {
					scoringfactors = append(scoringfactors, text.Content)
				}
			}
			if text.Top == incometaxtop && text.Left == incometaxleft {
				reportdata.IncomeTaxIdNumber = text.Content
			}
			// if text.Top == passportnumbertop && text.Left == passportnumberleft {
			// 	tk.Println(passportnumbertop)
			// 	reportdata.PassportNumber = text.Content
			// }
			if text.Top > telephonetop && text.Top < telephonebot && text.Left == 42 {
				telephonedata.Type = text.Content
			}
			if text.Top > telephonetop && text.Top < telephonebot && text.Left == 333 {
				telephonedata.Number = text.Content
				telephones = append(telephones, telephonedata)
			}
			if emailtopindex != emailbotindex {
				if i == emailtopindex {
					if text.Top > emailtop && text.Left == 42 {
						emails = append(emails, text.Content)
					}
				}
				if i == emailbotindex {
					if text.Top < emailbot && text.Left == 42 {
						emails = append(emails, text.Content)
					}
				}
			} else {
				if text.Top > emailtop && text.Top < emailbot && text.Left == 42 {
					emails = append(emails, text.Content)
				}
			}
			if text.Top == totalacctop && (text.Left == totalaccleft || text.Left == 283) {
				strs := strings.Split(text.Content, " ")
				if len(strs) > 0 {
					totalacc, _ := strconv.Atoi(strs[0])
					reportdata.TotalAccount = totalacc
				} else {
					tk.Println("else")
					totalacc, _ := strconv.Atoi(text.Content)
					reportdata.TotalAccount = totalacc
				}
			}
			if text.Top == overduetop && (text.Left == overdueleft || text.Left == 283) {
				strs := strings.Split(text.Content, " ")
				if len(strs) > 0 {
					overdue, _ := strconv.Atoi(strs[0])
					reportdata.TotalOverdues = overdue
				} else {
					overdue, _ := strconv.Atoi(text.Content)
					reportdata.TotalOverdues = overdue
				}
			}
			if text.Top == zerobalancetop && (text.Left == zerobalanceleft || text.Left == 283) {
				strs := strings.Split(text.Content, " ")
				if len(strs) > 0 {
					zerobalance, _ := strconv.Atoi(strs[0])
					reportdata.TotalZeroBalanceAcc = zerobalance
				} else {
					zerobalance, _ := strconv.Atoi(text.Content)
					reportdata.TotalZeroBalanceAcc = zerobalance
				}

			}
			if text.Top == highcreditsanctiontop && text.Left >= highcreditsancright && text.Left < currentbalanceleft {
				val := ReplaceString(text.Content)
				highcredit, _ := strconv.ParseFloat(val, 64)
				reportdata.HighCreditSanctionAmount = highcredit
			}
			if text.Top == currentbalancetop && text.Left >= currentbalanceright && text.Left < dateopenrecentleft {
				val := ReplaceString(text.Content)
				currentbalance, _ := strconv.ParseFloat(val, 64)
				reportdata.CurrentBalance = currentbalance
			}
			if text.Top == overduebalancetop && text.Left >= overduebalanceright && text.Left < dateopenoldestleft {
				val := ReplaceString(text.Content)
				overduebalance, _ := strconv.ParseFloat(val, 64)
				reportdata.OverdueBalance = overduebalance
			}
			if text.Top == dateopenrecenttop && text.Left == dateopenrecentleft {
				dateopenrecent, _ := time.Parse(layout, text.Content)
				reportdata.DateOpenedRecent = dateopenrecent
			}
			if text.Top == dateopenoldesttop && text.Left == dateopenoldestleft {
				dateopenoldest, _ := time.Parse(layout, text.Content)
				reportdata.DateOpenedOldest = dateopenoldest
			}
			if text.Top == enquirybot && text.Left <= enquiryright {
				enquirytotal, _ := strconv.Atoi(text.Content)
				reportdata.TotalEnquiries = enquirytotal
			}
			if text.Top == enquirybot && text.Left <= enquirypast30right {
				enquiry30, _ := strconv.Atoi(text.Content)
				reportdata.TotalEnquiries30Days = enquiry30
			}
			if text.Top == enquirybot && text.Left == enquiryrecentdateleft {
				enquiryrecentdate, _ := time.Parse(layout, text.Content)
				reportdata.RecentEnquiriesDates = enquiryrecentdate
			}

			if addresstopindex != addressbotindex {
				if i == addresstopindex {
					if text.Content == "ADDRESS:" || text.Content == "ADDRESS(e):" {
						addressdetailtop = text.Top
					}
					if text.Content == "CATEGORY:" {
						addresscategorytop = text.Top
					}
					if text.Content == "DATE REPORTED:" {
						addressdatereportedtop = text.Top
					}
					if text.Top == addressdetailtop && text.Left == addressdetailpermanentleft {
						addressdetail.AddressPinCode = text.Content
					}
					if text.Top == addressdetailtop && text.Left == addressdetailleft {
						addressdetail.AddressPinCode = text.Content
					}
					if text.Top == addressdetailtop+14 && text.Left == 48 {
						addressdetail.AddressPinCode = addressdetail.AddressPinCode + " " + text.Content
					}
					if text.Top == addresscategorytop && text.Left == addresscategoryleft {
						addressdetail.Category = text.Content
					}
					if text.Top == addressdatereportedtop && text.Left == addressdatereportedleft {
						addressreport, _ := time.Parse(layout, text.Content)
						addressdetail.DateReported = addressreport
						addressdetails = append(addressdetails, addressdetail)
					}
				}
				if i == addressbotindex {
					if text.Content == "ADDRESS:" || text.Content == "ADDRESS(e):" {
						addressdetailtop = text.Top
					}
					if text.Content == "CATEGORY:" {
						addresscategorytop = text.Top
					}
					if text.Content == "DATE REPORTED:" {
						addressdatereportedtop = text.Top
					}
					if text.Top == addressdetailtop && text.Left == addressdetailpermanentleft {
						addressdetail.AddressPinCode = text.Content
					}
					if text.Top == addressdetailtop && text.Left == addressdetailleft {
						addressdetail.AddressPinCode = text.Content
					}
					if text.Top == addressdetailtop+14 && text.Left == 48 {
						addressdetail.AddressPinCode = addressdetail.AddressPinCode + " " + text.Content
					}
					if text.Top == addresscategorytop && text.Left == addresscategoryleft {
						addressdetail.Category = text.Content
					}
					if text.Top == addressdatereportedtop && text.Left == addressdatereportedleft {
						addressreport, _ := time.Parse(layout, text.Content)
						addressdetail.DateReported = addressreport
						addressdetails = append(addressdetails, addressdetail)
					}
				}
			} else {
				if text.Content == "ADDRESS:" || text.Content == "ADDRESS(e):" {
					addressdetailtop = text.Top
				}
				if text.Content == "CATEGORY:" {
					addresscategorytop = text.Top
				}
				if text.Content == "DATE REPORTED:" {
					addressdatereportedtop = text.Top
				}
				if text.Top == addressdetailtop && text.Left == addressdetailpermanentleft {
					addressdetail.AddressPinCode = text.Content
				}
				if text.Top == addressdetailtop && text.Left == addressdetailleft {
					addressdetail.AddressPinCode = text.Content
				}
				if text.Top == addresscategorytop && text.Left == addresscategoryleft {
					addressdetail.Category = text.Content
				}
				if text.Top == addressdatereportedtop && text.Left == addressdatereportedleft {
					addressreport, _ := time.Parse(layout, text.Content)
					addressdetail.DateReported = addressreport
					addressdetails = append(addressdetails, addressdetail)
				}
			}
		}
	}
	reportdata.ConsumersInfos = consumerinfo
	reportdata.ScoringFactor = scoringfactors
	reportdata.Telephones = telephones
	reportdata.EmailAddress = emails
	reportdata.AddressData = addressdetails
	reportdata.ReportType = "Individual"
	return reportdata
}

func ExtractPdfDataCibilReport(PathFrom string, PathTo string, FName string, ReportType string, XmlName string, inbox string, success string, failed string, webapps string) {
	tk.Println("Extracting", FName)
	//Name := strings.TrimRight(FName, ".pdf")

	conn, err := PrepareConnection()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	if ReportType == "Company" {
		reportobj := ExtractCompanyCibilReport(PathTo, XmlName)

		filename := strings.TrimRight(FName, ".pdf")
		timestamp := time.Now().UTC()
		datestr := timestamp.String()
		dates := strings.Split(datestr, " ")
		times := strings.Split(dates[1], ".")
		newfilename := filename + "_" + dates[0] + "_" + times[0] + ".pdf"
		os.Rename(inbox+"/"+FName, inbox+"/"+newfilename)
		formattedName := strings.Replace(newfilename, " ", "\\ ", -1)

		if reportobj.Profile.CompanyName == "" {
			tk.Println("Undefined Company Name")
			MoveFile(inbox+"/"+formattedName, failed)
			os.RemoveAll(PathFrom + "/" + XmlName)
		} else {
			customer := strings.Split(reportobj.Profile.CompanyName, " ")
			res := []tk.M{}
			filter := []*dbox.Filter{}
			isMatch := false
			customerid := 0
			dealno := ""

			for _, splited := range customer {
				if len(splited) > 3 && splited != "PVT" && splited != "LTD" && splited != "PRIVATE" && splited != "LIMITED" {
					tk.Println(splited)
					filter = append(filter, dbox.Contains("applicantdetail.CustomerName", splited))
				}
			}

			cursor, err := conn.NewQuery().Select().From("CustomerProfile").Where(filter...).Cursor(nil)
			if err != nil {
				tk.Println(err.Error())
			}
			err = cursor.Fetch(&res, 0, false)
			defer cursor.Close()

			if len(res) > 0 {
				for _, val := range res {
					customername := val["applicantdetail"].(tk.M)["CustomerName"].(string)
					app := val.Get("applicantdetail").(tk.M)
					customerid = app.GetInt("CustomerID")
					dealno = val["applicantdetail"].(tk.M)["DealNo"].(string)
					custpan := val["applicantdetail"].(tk.M)["CustomerPan"].(string)

					setting := NewSimilaritySetting()
					setting.SplitDelimeters = []rune{' ', '.', '-'}
					splitedreportname := strings.Split(reportobj.Profile.CompanyName, " ")
					splitedcpname := strings.Split(customername, " ")
					simreportname := ""
					simcpname := ""

					for _, reportname := range splitedreportname {
						if reportname != "PVT" || reportname != "LTD" || reportname != "PRIVATE" || reportname != "LIMITED" {
							simreportname = simreportname + " " + reportname
						}
					}

					for _, cpname := range splitedcpname {
						if cpname != "PVT" || cpname != "LTD" || cpname != "PRIVATE" || cpname != "LIMITED" {
							simcpname = simcpname + " " + cpname
						}
					}

					tk.Println("test", simreportname, simcpname)
					similar := Similarity(simreportname, simcpname, setting)

					if reportobj.Profile.Pan != "" {
						if similar >= 50 && custpan == reportobj.Profile.Pan {
							isMatch = true
						}
					} else {
						if similar >= 70 {
							isMatch = true
						}
					}
				}
			} else {
				tk.Println("else")
				reportobj.Id = bson.NewObjectId()
				reportobj.FilePath = PathFrom + "/" + FName
				reportobj.FileName = newfilename
				reportobj.IsMatch = isMatch
				query := conn.NewQuery().From("CibilReport").Save()
				err = query.Exec(tk.M{
					"data": reportobj,
				})
				if err != nil {
					tk.Println(err.Error())
				}
				query.Close()

				os.RemoveAll(PathFrom + "/" + XmlName)
				CopyFile(inbox+"/"+formattedName, webapps)
				MoveFile(inbox+"/"+formattedName, success)
			}

			if isMatch {
				reportobj.Id = bson.NewObjectId()
				reportobj.Profile.CustomerId = customerid
				reportobj.Profile.DealNo = dealno
				reportobj.FilePath = PathFrom + "/" + FName
				reportobj.FileName = newfilename
				reportobj.IsMatch = isMatch
				query := conn.NewQuery().From("CibilReport").Save()
				err = query.Exec(tk.M{
					"data": reportobj,
				})
				if err != nil {
					tk.Println(err.Error())
				}
				query.Close()

				os.RemoveAll(PathFrom + "/" + XmlName)
				CopyFile(inbox+"/"+formattedName, webapps)
				MoveFile(inbox+"/"+formattedName, success)
			} else {
				os.RemoveAll(PathFrom + "/" + XmlName)
				MoveFile(inbox+"/"+formattedName, success)
			}
		}
	}

	if ReportType == "Individual" {
		reportobj := ExtractIndividualCibilReport(PathTo, XmlName)

		filename := strings.TrimRight(FName, ".pdf")
		timestamp := time.Now().UTC()
		datestr := timestamp.String()
		dates := strings.Split(datestr, " ")
		times := strings.Split(dates[1], ".")
		newfilename := filename + "_" + dates[0] + "_" + times[0] + ".pdf"
		os.Rename(inbox+"/"+FName, inbox+"/"+newfilename)
		formattedName := strings.Replace(newfilename, " ", "\\ ", -1)

		if reportobj.CibilScore == 0 {
			MoveFile(inbox+"/"+formattedName, failed)
			os.RemoveAll(PathFrom + "/" + XmlName)
		} else {
			customer := strings.Split(reportobj.ConsumersInfos.ConsumerName, " ")
			res := []tk.M{}
			filter := []*dbox.Filter{}
			isMatch := false
			customerid := 0
			dealno := ""

			for _, splited := range customer {
				if len(splited) > 2 && splited != "JAIN" && splited != "PATEL" && splited != "SHAH" {
					tk.Println(splited)
					filter = append(filter, dbox.Contains("detailofpromoters.biodata.Name", splited))
				}
			}

			cursor, err := conn.NewQuery().Select().From("CustomerProfile").Where(dbox.Or(filter...)).Cursor(nil)
			if err != nil {
				tk.Println(err.Error())
			}
			err = cursor.Fetch(&res, 0, false)
			defer cursor.Close()

			if len(res) > 0 {
				for _, val := range res {
					isMatch = false
					customername := val.Get("detailofpromoters").(tk.M)["biodata"]
					bio := customername.([]interface{})
					app := val.Get("applicantdetail").(tk.M)
					customerid = app.GetInt("CustomerID")
					dealno = val["applicantdetail"].(tk.M)["DealNo"].(string)

					for _, vals := range bio {
						data := vals.(tk.M)
						setting := NewSimilaritySetting()
						setting.SplitDelimeters = []rune{' ', '.', '-'}
						similar := Similarity(reportobj.ConsumersInfos.ConsumerName, data.GetString("Name"), setting)
						dob, isdate := data.Get("DateOfBirth").(time.Time)

						if isdate {
							if similar >= 50 && reportobj.ConsumersInfos.DateOfBirth == dob.UTC() {
								isMatch = true
								break
							} else if data.GetString("PAN") != "" {
								if reportobj.IncomeTaxIdNumber == data.GetString("PAN") {
									isMatch = true
									break
								}
							}
						} else {
							datestring := data.GetString("DateOfBirth")
							datesplitted := strings.Split(datestring, "T")
							layout := "2006-01-02"
							strdate := datesplitted[0]
							t, err := time.Parse(layout, strdate)

							if err != nil {
								tk.Println(err)
							} else {
								if similar >= 50 && reportobj.ConsumersInfos.DateOfBirth == t {
									isMatch = true
									break
								} else if data.GetString("PAN") != "" {
									if reportobj.IncomeTaxIdNumber == data.GetString("PAN") {
										isMatch = true
										break
									}
								}
							}
						}
					}

					if isMatch {
						tk.Println("PDF Match")
						tk.Println("Where", customerid, dealno, reportobj.ConsumersInfos.ConsumerName)
						filter := []*dbox.Filter{}
						filter = append(filter, dbox.Eq("ConsumerInfo.ConsumerName", reportobj.ConsumersInfos.ConsumerName))
						filter = append(filter, dbox.Eq("ConsumerInfo.CustomerId", customerid))
						filter = append(filter, dbox.Eq("ConsumerInfo.DealNo", dealno))
						cursor, err = conn.NewQuery().Select().From("CibilReportPromotorFinal").Where(filter...).Cursor(nil)
						if err != nil {
							tk.Println(err.Error())
						}
						result := []tk.M{}

						err = cursor.Fetch(&result, 0, false)

						if len(result) == 0 {
							reportobj.Id = bson.NewObjectId()
							reportobj.ConsumersInfos.CustomerId = customerid
							reportobj.ConsumersInfos.DealNo = dealno
							reportobj.FilePath = PathFrom + "/" + FName
							reportobj.FileName = newfilename
							reportobj.StatusCibil = 0
							reportobj.IsMatch = isMatch
							query := conn.NewQuery().From("CibilReportPromotorFinal").Save()
							err = query.Exec(tk.M{
								"data": reportobj,
							})
							if err != nil {
								tk.Println(err.Error())
							}
							query.Close()

						} else {
							for _, existdata := range result {
								if existdata.GetInt("StatusCibil") != 1 {
									datereport := existdata.Get("DateOfReport").(time.Time)
									timereport := existdata.Get("TimeOfReport").(time.Time)
									if datereport.Before(reportobj.DateOfReport) || datereport == reportobj.DateOfReport && timereport.Before(reportobj.TimeOfReport) {
										wh := []*dbox.Filter{}
										ids := existdata.Get("_id").(bson.ObjectId)
										tk.Println("ID", ids)
										wh = append(wh, dbox.Eq("_id", ids))
										err = conn.NewQuery().From("CibilReportPromotorFinal").Delete().Where(wh...).Exec(nil)
										if err != nil {
											tk.Println(err.Error())
										}

										reportobj.Id = bson.NewObjectId()
										reportobj.ConsumersInfos.CustomerId = customerid
										reportobj.ConsumersInfos.DealNo = dealno
										reportobj.FilePath = PathFrom + "/" + FName
										reportobj.FileName = newfilename
										reportobj.StatusCibil = 0
										reportobj.IsMatch = isMatch
										query := conn.NewQuery().From("CibilReportPromotorFinal").Save()
										err = query.Exec(tk.M{
											"data": reportobj,
										})
										if err != nil {
											tk.Println(err.Error())
										}
										query.Close()

									} else {
										reportobj.Id = bson.NewObjectId()
										reportobj.ConsumersInfos.CustomerId = customerid
										reportobj.ConsumersInfos.DealNo = dealno
										reportobj.FilePath = PathFrom + "/" + FName
										reportobj.FileName = newfilename
										reportobj.StatusCibil = 0
										reportobj.IsMatch = isMatch
										query := conn.NewQuery().From("CibilReportPromotorFinal").Save()
										err = query.Exec(tk.M{
											"data": reportobj,
										})
										if err != nil {
											tk.Println(err.Error())
										}
										query.Close()
									}
								} else {
									isMatch = false
								}
							}
						}
					}
				}
			}

			if isMatch == false {
				tk.Println("PDF Unmatch")
				reportobj.Id = bson.NewObjectId()
				reportobj.FilePath = PathFrom + "/" + FName
				reportobj.FileName = newfilename
				reportobj.StatusCibil = 0
				reportobj.IsMatch = isMatch
				query := conn.NewQuery().From("CibilReportPromotorFinal").Save()
				err = query.Exec(tk.M{
					"data": reportobj,
				})
				if err != nil {
					tk.Println(err.Error())
				}
				query.Close()

			}
			os.RemoveAll(PathFrom + "/" + XmlName)
			CopyFile(inbox+"/"+formattedName, webapps)
			MoveFile(inbox+"/"+formattedName, success)
		}
	}

	tk.Println("Extracting Finish")
}

func ReplaceString(number string) string {
	rex := regexp.MustCompile("[^0-9]")
	valStr := rex.ReplaceAllString(number, "")
	return valStr
}
