package main

import (
	"encoding/xml"
	"errors"
	"strconv"
)

func translateData(title, description string) (translationParamsStruct, error) {
	var returnParams translationParamsStruct
	espXmlmc.SetParam("title", title)
	if description != "" {
		espXmlmc.SetParam("description", description)
	}
	espXmlmc.SetParam("sourceLanguage", configSource)
	espXmlmc.SetParam("targetLanguage", configDestination)

	var XMLSTRING = espXmlmc.GetParam()
	debugLog("apps/com.hornbill.servicemanager::translateData")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("apps/com.hornbill.servicemanager", "translateData")
	if err != nil {
		return returnParams, errors.New("API Call failed when returning Translations from Hornbill: " + err.Error())
	}

	var xmlRespon espGetTranslationStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return returnParams, errors.New("Unmarshal failed when returning Translations from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return returnParams, errors.New("Error from Hornbill when returning Translations data: " + xmlRespon.State.ErrorRet)
	}
	return xmlRespon.Params, nil
}

func translateAddLanguage(entityObj, inputTitle, inputDescription, language, extraParams string, linkedID int) error {
	espXmlmc.SetParam("entityObj", entityObj)
	espXmlmc.SetParam("linkedId", strconv.Itoa(linkedID))
	espXmlmc.SetParam("inputTitle", inputTitle)
	espXmlmc.SetParam("inputDescription", inputDescription)
	espXmlmc.SetParam("language", language)
	espXmlmc.SetParam("extraParams", extraParams)

	var XMLSTRING = espXmlmc.GetParam()
	debugLog("apps/com.hornbill.servicemanager::translateAddLanguage")
	debugLog(XMLSTRING)

	xmlmcResponse, err := espXmlmc.Invoke("apps/com.hornbill.servicemanager", "translateAddLanguage")
	if err != nil {
		return errors.New("API Call failed when adding Translation into Hornbill: " + err.Error())
	}

	var xmlRespon espGetTranslationStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		return errors.New("Unmarshal failed when addint Translation into Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		return errors.New("Error from Hornbill when adding Translation data: " + xmlRespon.State.ErrorRet)
	}
	return nil
}
