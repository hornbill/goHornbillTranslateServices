package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"strconv"
)

func getServices(language string) []smServiceStruct {
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "Services")
	espXmlmc.OpenElement("searchFilter")
	espXmlmc.SetParam("column", "h_language")
	espXmlmc.SetParam("value", language)
	espXmlmc.SetParam("matchType", "exact")
	espXmlmc.CloseElement("searchFilter")

	var XMLSTRING = espXmlmc.GetParam()
	debugLog("data::entityBrowseRecords2::Services")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("data", "entityBrowseRecords2")
	if err != nil {
		log.Fatal(errors.New("API Call failed when returning Services from Hornbill: " + err.Error()))
	}

	var xmlRespon espGetServicesStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		log.Fatal("Unmarshal failed when returning Services from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		logr.Fatal("Error from Hornbill when returning Services data: ", xmlRespon.State.ErrorRet)
	}
	return xmlRespon.Rows
}

func getServiceDetails(serviceID int) (smServiceDetailsStruct, error) {
	var serviceDetails smServiceDetailsStruct
	espXmlmc.SetParam("serviceId", strconv.Itoa(serviceID))
	var XMLSTRING = espXmlmc.GetParam()
	debugLog("apps/com.hornbill.servicemanager/Services::smGetServiceDetails")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("apps/com.hornbill.servicemanager/Services", "smGetServiceDetails")
	if err != nil {
		return serviceDetails, errors.New("API Call failed when returning default Service record from Hornbill: " + err.Error())
	}
	var xmlRespon espGetTranslationStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return serviceDetails, errors.New("Unmarshal failed when returning default Service record from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return serviceDetails, errors.New("Error from Hornbill when returning default Service record: " + xmlRespon.State.ErrorRet)
	}
	err = json.Unmarshal([]byte(xmlRespon.Params.ServiceDetails), &serviceDetails)
	if err != nil {
		return serviceDetails, errors.New("JSON Unmarshal failed when returning default Service record from Hornbill: " + err.Error())
	}
	return serviceDetails, nil
}

func makeServiceExtraParams(serviceDetails smServiceDetailsStruct) string {
	var extraParams serviceExtraParamsStruct
	extraParams.HServiceCatalogCategory, _ = strconv.Atoi(serviceDetails.HServiceCatalogCategory)
	extraParams.HServiceCatalogDomain = serviceDetails.HServiceCatalogDomain
	extraParams.HFkServicecategory = serviceDetails.HFkServicecategory
	extraParams.HPortfolioStatus = serviceDetails.HPortfolioStatus
	extraParams.HIcon = serviceDetails.HIcon
	extraParams.HAccess = serviceDetails.HAccess
	extraParams.HAvailable = serviceDetails.HAvailable
	extraParams.HStatus = serviceDetails.HStatus
	extraParams.HMbid = serviceDetails.HMbid

	returnParams, _ := json.Marshal(extraParams)
	return string(returnParams)
}

func makeServiceEntityObj() string {
	var entityObj entityObjStruct
	entityObj.Name = "Services"
	entityObj.Title = "h_servicename"
	entityObj.Description = "h_servicedescription"
	entityObj.EntityColumn = "h_pk_serviceid"
	entityObj.LinkedColumn = "h_linked_service_id"
	returnParams, _ := json.Marshal(entityObj)
	return string(returnParams)
}
