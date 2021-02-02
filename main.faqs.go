package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

func getFAQs(serviceID int, language string) ([]faqDetailsStruct, error) {
	var returnParams []faqDetailsStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "FAQs")
	espXmlmc.OpenElement("searchFilter")
	espXmlmc.SetParam("column", "h_language")
	espXmlmc.SetParam("value", language)
	espXmlmc.SetParam("matchType", "exact")
	espXmlmc.CloseElement("searchFilter")
	espXmlmc.OpenElement("searchFilter")
	espXmlmc.SetParam("column", "h_service_id")
	espXmlmc.SetParam("value", strconv.Itoa(serviceID))
	espXmlmc.SetParam("matchType", "exact")
	espXmlmc.CloseElement("searchFilter")

	var XMLSTRING = espXmlmc.GetParam()
	debugLog("data::entityBrowseRecords2::FAQs")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning FAQs from Hornbill: " + err.Error())
	}

	var xmlRespon espGetFAQsStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning FAQs from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning FAQs data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.FAQDetails, err
}

func makeFAQExtraParams(faqDetails faqDetailsStruct) string {
	var extraParams faqExtraParamsStruct
	extraParams.HServiceID = faqDetails.HServiceID
	extraParams.HServiceName = faqDetails.HServiceName
	extraParams.HViewCount = faqDetails.HViewCount
	extraParams.HMediaLink = faqDetails.HMediaLink
	extraParams.HVisibility = faqDetails.HVisibility
	extraParams.HCreatedByUserID = faqDetails.HCreatedByUserID
	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeFAQEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "Faqs"
	entityObj.Title = "h_question"
	entityObj.Description = "h_answer"
	entityObj.EntityColumn = "h_id"
	entityObj.LinkedColumn = "h_faq_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
