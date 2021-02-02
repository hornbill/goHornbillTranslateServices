package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	apiLib "github.com/hornbill/goApiLib"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {
	//Setup logging
	cwd, _ := os.Getwd()
	logPath := cwd + "/log"
	//-- If Folder Does Not Exist then create it
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.Mkdir(logPath, 0777)
		if err != nil {
			fmt.Println("Error Creating Log Folder ", logPath, ": ", err)
			os.Exit(101)
		}
	}
	logFileName := logPath + "/Service_Translation_" + time.Now().Format("20060102150405") + ".log"
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		os.Exit(1)
	}

	//Setup logr to FILE ONLY
	logFile = &logrus.Logger{
		Out:   f,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "%time% [%lvl%] %msg%\n",
		},
	}

	logr = &logrus.Logger{
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%] %msg%\n",
		},
	}
	logr.SetOutput(os.Stdout)
}

func main() {
	processFlags()
	espXmlmc = apiLib.NewXmlmcInstance(configInstance)
	espXmlmc.SetAPIKey(configAPIKey)

	if configLangs {
		getLangs()
		os.Exit(0)
	}

	srcServices = getServices(configSource)
	dstServices = getServices(configDestination)

	for _, srcService := range srcServices {
		translateServiceRecord := true

		//-- HLinkedServiceID contains the ID for the DEFAULT service language.
		for _, dstService := range dstServices {
			if dstService.HLinkedServiceID == srcService.HLinkedServiceID {
				translateServiceRecord = false
			}
		}

		if translateServiceRecord {
			defaultServiceDetails, err := getServiceDetails(srcService.HLinkedServiceID)
			if err != nil {
				logError(err, true)
				counters.ServicesError++
				continue
			} else {
				logInfo("[SERVICE] ["+srcService.HServicename+"] translating from "+configSource+" to "+configDestination, true)
				translated, err := translateData(srcService.HServicename, srcService.HServicedescription)
				if err != nil {
					logError(err, true)
					counters.ServicesError++
					continue
				}
				extraParams := makeServiceExtraParams(defaultServiceDetails)
				entityObj := makeServiceEntityObj()
				err = translateAddLanguage(entityObj, translated.TranslatedTitle, translated.TranslatedDescription, configDestination, extraParams, srcService.HLinkedServiceID)
				if err != nil {
					logError(err, true)
					counters.ServicesError++
				} else {
					logInfo("[SERVICE] ["+srcService.HServicename+"] successfully translated from "+configSource+" to "+configDestination, true)
					counters.ServicesTranslated++
				}
			}
		} else {
			counters.ServicesSkipped++
			logInfo("[SERVICE] ["+srcService.HServicename+"] has already been translated to "+configDestination, true)
		}

		//Get Catalog records for current Service
		sourceLangCatalogs, err := getCatalogs(srcService.HLinkedServiceID, configSource)
		if err != nil {
			logError(err, true)
			continue
		}
		destLangCatalogs, err := getCatalogs(srcService.HLinkedServiceID, configDestination)
		if err != nil {
			logError(err, true)
			continue
		}

		//Process Catalog records for current Service
		for _, sc := range sourceLangCatalogs {
			transCat := true
			for _, dc := range destLangCatalogs {
				if dc.HRequestCatalogID == sc.HRequestCatalogID && dc.HLanguage == configDestination {
					transCat = false
				}
			}
			if transCat {
				logInfo("[CATALOGITEM] ["+sc.HCatalogTitle+"] translating from "+configSource+" to "+configDestination, true)
				translated, err := translateData(sc.HCatalogTitle, sc.HCatalogDescription)
				if err != nil {
					logError(err, true)
					counters.CatalogsError++
					continue
				}
				extraParams := makeCatalogExtraParams(sc)
				entityObj := makeCatalogEntityObj()
				err = translateAddLanguage(entityObj, translated.TranslatedTitle, translated.TranslatedDescription, configDestination, extraParams, sc.HRequestCatalogID)
				if err != nil {
					logError(err, true)
					counters.CatalogsError++
				} else {
					logInfo("[CATALOGITEM] ["+sc.HCatalogTitle+"] successfully translated from "+configSource+" to "+configDestination, true)
					counters.CatalogsCreated++
				}
			} else {
				logInfo("[CATALOGITEM] ["+sc.HCatalogTitle+"] has already been translated to "+configDestination, true)
				counters.CatalogsSkipped++
			}
		}

		//Get Bulletin records for current Service
		sourceLangBulletins, err := getBulletins(srcService.HLinkedServiceID, configSource)
		if err != nil {
			logError(err, true)
			continue
		}
		destLangBulletins, err := getBulletins(srcService.HLinkedServiceID, configDestination)
		if err != nil {
			logError(err, true)
			continue
		}

		//Process Bulletin records for current Service
		for _, sb := range sourceLangBulletins {
			transBul := true
			for _, db := range destLangBulletins {
				if db.HServiceBulletinID == sb.HServiceBulletinID && db.HLanguage == configDestination {
					transBul = false
				}
			}
			if transBul {
				logInfo("[BULLETIN] Translating ["+sb.HBulletinTitle+"] from "+configSource+" to "+configDestination, true)
				translated, err := translateData(sb.HBulletinTitle, sb.HBulletinDescription)
				if err != nil {
					logError(err, true)
					counters.BulletinsError++
					continue
				}
				extraParams := makeBulletinExtraParams(sb)
				entityObj := makeBulletinEntityObj()
				err = translateAddLanguage(entityObj, translated.TranslatedTitle, translated.TranslatedDescription, configDestination, extraParams, sb.HServiceBulletinID)
				if err != nil {
					logError(err, true)
					counters.BulletinsError++
					continue
				} else {
					logInfo("[BULLETIN] ["+sb.HBulletinTitle+"] successfully translated from "+configSource+" to "+configDestination, true)
					counters.BulletinsCreated++
				}
			} else {
				logInfo("[BULLETIN] ["+sb.HBulletinTitle+"] has already been translated to "+configDestination, true)
				counters.BulletinsSkipped++
			}
		}

		//Get FAQ records for current Service
		sourceLangFAQs, err := getFAQs(srcService.HLinkedServiceID, configSource)
		if err != nil {
			logError(err, true)
			continue
		}
		destLangFAQs, err := getFAQs(srcService.HLinkedServiceID, configDestination)
		if err != nil {
			logError(err, true)
			continue
		}

		//Process FAQ records for current Service
		for _, sf := range sourceLangFAQs {
			transFAQ := true
			for _, df := range destLangFAQs {
				if df.HFAQID == sf.HFAQID && df.HLanguage == configDestination {
					transFAQ = false
				}
			}
			if transFAQ {
				logInfo("[FAQ] Translating ["+sf.HFAQQuestion+"] from "+configSource+" to "+configDestination, true)
				translated, err := translateData(sf.HFAQQuestion, sf.HFAQAnswer)
				if err != nil {
					logError(err, true)
					counters.FAQsError++
					continue
				}
				extraParams := makeFAQExtraParams(sf)
				entityObj := makeFAQEntityObj()
				err = translateAddLanguage(entityObj, translated.TranslatedTitle, translated.TranslatedDescription, configDestination, extraParams, sf.HFAQID)
				if err != nil {
					logError(err, true)
					counters.FAQsError++
					continue
				} else {
					logInfo("[FAQ] ["+sf.HFAQQuestion+"] successfully translated from "+configSource+" to "+configDestination, true)
					counters.FAQsCreated++
				}
			} else {
				logInfo("[FAQ] ["+sf.HFAQQuestion+"] has already been translated to "+configDestination, true)
				counters.FAQsSkipped++
			}
		}

		//Get Feedback records for current Service
		sourceLangFeedbacks, err := getFeedbacks(srcService.HLinkedServiceID, configSource)
		if err != nil {
			logError(err, true)
			continue
		}
		destLangFeedbacks, err := getFeedbacks(srcService.HLinkedServiceID, configDestination)
		if err != nil {
			logError(err, true)
			continue
		}

		//Process Feedback records for current Service
		for _, sf := range sourceLangFeedbacks {
			transFeedback := true
			for _, df := range destLangFeedbacks {
				if df.HQuestionID == sf.HQuestionID && df.HLanguage == configDestination {
					transFeedback = false
				}
			}
			if transFeedback {
				logInfo("[FEEDBACK] Translating ["+sf.HQuestion+"] from "+configSource+" to "+configDestination, true)
				translated, err := translateData(sf.HQuestion, "")
				if err != nil {
					logError(err, true)
					counters.FeedbackError++
					continue
				}
				extraParams := makeFeedbackExtraParams(sf)
				entityObj := makeFeedbackEntityObj()
				err = translateAddLanguage(entityObj, translated.TranslatedTitle, translated.TranslatedDescription, configDestination, extraParams, sf.HQuestionID)
				if err != nil {
					logError(err, true)
					counters.FeedbackError++
					continue
				} else {
					logInfo("[FEEDBACK] ["+sf.HQuestion+"] successfully translated from "+configSource+" to "+configDestination, true)
					counters.FeedbackCreated++
				}
			} else {
				logInfo("[FEEDBACK] ["+sf.HQuestion+"] has already been translated to "+configDestination, true)
				counters.FeedbackSkipped++
			}
		}

	}
	logInfo("---- Service Bulk Translation Complete ----", true)
	logInfo("Service Translations Created: "+strconv.Itoa(counters.ServicesTranslated), true)
	logInfo("Service Translations Skipped: "+strconv.Itoa(counters.ServicesSkipped), true)
	logInfo("Service Translations Errors "+strconv.Itoa(counters.ServicesError), true)
	logInfo("Catalog Item Translations Created: "+strconv.Itoa(counters.CatalogsCreated), true)
	logInfo("Catalog Item Translations Skipped: "+strconv.Itoa(counters.CatalogsSkipped), true)
	logInfo("Catalog Item Translations Errors "+strconv.Itoa(counters.CatalogsError), true)
	logInfo("Bulletin Translations Created: "+strconv.Itoa(counters.BulletinsCreated), true)
	logInfo("Bulletin Translations Skipped: "+strconv.Itoa(counters.BulletinsSkipped), true)
	logInfo("Bulletin Translations Errors "+strconv.Itoa(counters.BulletinsError), true)
	logInfo("FAQ Translations Created: "+strconv.Itoa(counters.FAQsCreated), true)
	logInfo("FAQ Translations Skipped: "+strconv.Itoa(counters.FAQsSkipped), true)
	logInfo("FAQ Translations Errors "+strconv.Itoa(counters.FAQsError), true)
	logInfo("Feedback Question Translations Created: "+strconv.Itoa(counters.FeedbackCreated), true)
	logInfo("Feedback Question Translations Skipped: "+strconv.Itoa(counters.FeedbackSkipped), true)
	logInfo("Feedback Question Translations Errors "+strconv.Itoa(counters.FeedbackError), true)

}
