package util

import (
	"runtime"
	"strings"
)

const (
	PORT = 9515
	USER = ""
	PWD  = ""
)

var (
	OUTPUTIMAGEFOLDER = getOutputFolder()
	ROOTPATH          = getRootPath()
	SELENIUMPATH      = getSeleniumPath()
	SPREADSHEETPATH   = getSpreadSheetPath()
	FONTPATH          = getFontPath()

	ORIGINIMAGEPATH string = ""
	TITLE           string = ""
)

func getOutputFolder() string {
	if strings.Contains(runtime.GOOS, "window") {
		return "D:\\sifamadocs\\imagens"
	}
	return "/home/eduardo/Documentos/projetos/sifamaSources/imagens"
}

func getRootPath() string {
	if strings.Contains(runtime.GOOS, "window") {
		return "D:\\Documentos\\Users\\Eduardo\\Documentos\\ANTT\\OneDrive - ANTT- Agencia Nacional de Transportes Terrestres\\CRO\\Relatórios RTA"
	}
	return "/home/eduardo/Documentos/projetos/sifamaSources/Anexos"
}

func getSeleniumPath() string {
	if strings.Contains(runtime.GOOS, "window") {
		return "D:\\chromedriver.exe"
	}
	return "/home/eduardo/automation/chromedriver"
}

func getSpreadSheetPath() string {
	if strings.Contains(runtime.GOOS, "window") {
		return "D:\\Documentos\\Users\\Eduardo\\Documentos\\ANTT\\OneDrive - ANTT- Agencia Nacional de Transportes Terrestres\\sistema\\sifamadocs\\planilha\\tros.xlsx"
	}
	return "/home/eduardo/Documentos/projetos/sifamaSources/tros.xlsx"
}

func getFontPath() string {
	if strings.Contains(runtime.GOOS, "window") {
		return "C:\\Windows\\Fonts\\arial.ttf"
	}
	return "/home/eduardo/Documentos/projetos/sifamaSources/arial.ttf"
}
