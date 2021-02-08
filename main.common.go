package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
)

func processFlags() {
	//-- Grab Flags
	flag.StringVar(&configInstance, "instance", "", "Hornbill Instance")
	flag.StringVar(&configAPIKey, "key", "", "Hornbill API Key")
	flag.StringVar(&configSource, "src", "en-GB", "Source Language")
	flag.StringVar(&configDestination, "dst", "", "Destination Language")
	flag.BoolVar(&configDebug, "debug", false, "Debug Mode")
	flag.BoolVar(&configLangs, "getlangs", false, "Return a list of supported languages and their codes in a file called languages.txt, then ends")
	flag.BoolVar(&configVersion, "version", false, "Return version and end")
	flag.Parse()
	//-- If configVersion just output version number and die
	if configVersion {
		fmt.Printf("%v \n", version)
		os.Exit(0)
	}
	if configInstance == "" || configAPIKey == "" {
		logr.Fatal("instance and key arguments are mandatory")
	}
	if configSource == "" || configDestination == "" {
		logr.Fatal("src and dst arguments are mandatory")
	}
	logInfo("---- Hornbill Service Manager Bulk Service Translation Tool v"+fmt.Sprintf("%v", version)+" ----", true)
	logInfo("Flag - instance "+configInstance, true)
	logInfo("Flag - src "+configSource, true)
	logInfo("Flag - dst "+configDestination, true)
	logInfo("Flag - debug "+fmt.Sprintf("%v", configDebug), true)
	logInfo("Flag - getlangs "+fmt.Sprintf("%v", configLangs), true)
}

func debugLog(s string) {
	if configDebug {
		logr.Debug(s)
	}
}

func getLangs() {
	debugLog("system::getLanguageList")
	xmlmcResponse, err := espXmlmc.Invoke("system", "getLanguageList")
	if err != nil {
		logr.Fatal("API Call failed when returning list of languages from Hornbill: " + err.Error())
	}

	var xmlRespon espGetTranslationStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		logr.Fatal("Unmarshal failed when returning list of languages from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		logr.Fatal("Error from Hornbill when returning list of languages data: " + xmlRespon.State.ErrorRet)
	}

	if configLangs {
		cwd, _ := os.Getwd()
		langFile := cwd + "/languages.txt"

		//Delete file
		os.Remove(langFile)

		//-- Open file to write langs to
		f, err := os.OpenFile(langFile, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			logr.Fatal(err)
		}
		defer f.Close()
		for _, v := range xmlRespon.Params.LanguageInfo {
			if v.Supported {
				f.WriteString(v.Language + " : " + v.Name + "\n")
			}
		}
	} else {
		for _, v := range xmlRespon.Params.LanguageInfo {
			if v.Supported {
				languages[v.Language] = v
			}
		}
	}
}

func getSMAppVer() {
	debugLog("session::getApplicationList")
	xmlmcResponse, err := espXmlmc.Invoke("session", "getApplicationList")
	if err != nil {
		logr.Fatal("API Call failed when returning list of applications from Hornbill: " + err.Error())
	}

	var xmlRespon espGetAppVerStruct
	err = xml.Unmarshal([]byte(xmlmcResponse), &xmlRespon)
	if err != nil {
		logr.Fatal("Unmarshal failed when returning list of applications from Hornbill: " + err.Error())
	}
	if xmlRespon.MethodResult != "ok" {
		logr.Fatal("Error from Hornbill when returning list of application data: " + xmlRespon.State.ErrorRet)
	}

	for _, v := range xmlRespon.Params.Application {
		if v.Name == "com.hornbill.servicemanager" && v.Build >= minSMBuild {
			smStatusSupport = true
		}
	}
}

func logInfo(s string, outputToCLI bool) {
	logFile.Info(s)
	if outputToCLI {
		logr.Info(s)
	}
}

func logWarn(s string, outputToCLI bool) {
	logFile.Warn(s)
	if outputToCLI {
		logr.Warn(s)
	}
}

func logError(e error, outputToCLI bool) {
	logFile.Error(e)
	if outputToCLI {
		logr.Error(e)
	}
}

func logDebug(s string, outputToCLI bool) {
	if configDebug {
		logFile.Debug(s)
		if outputToCLI {
			logr.Debug(s)
		}
	}
}
