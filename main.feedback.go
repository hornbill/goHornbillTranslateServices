package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

func getFeedbacks(serviceID int, language string) ([]feedbackDetailsStruct, error) {
	var returnParams []feedbackDetailsStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "ServiceFeedbackQuestions")
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
	debugLog("data::entityBrowseRecords2::ServiceFeedbackQuestions")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning ServiceFeedbackQuestions from Hornbill: " + err.Error())
	}

	var xmlRespon espGetFeedbackStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning ServiceFeedbackQuestions from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning ServiceFeedbackQuestions data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.FeedbackDetails, err
}

func makeFeedbackExtraParams(feedbackDetails feedbackDetailsStruct) string {
	var extraParams feedbackExtraParamsStruct
	extraParams.HServiceID = feedbackDetails.HServiceID
	extraParams.HRequestType = feedbackDetails.HRequestType
	extraParams.HFieldRequired = feedbackDetails.HFieldRequired
	extraParams.HFieldType = feedbackDetails.HFieldType
	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeFeedbackEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "ServiceFeedbackQuestions"
	entityObj.Title = "h_question"
	entityObj.EntityColumn = "h_id"
	entityObj.LinkedColumn = "h_question_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
