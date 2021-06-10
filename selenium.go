package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"sifamaGO/src/util"

	"github.com/tebeka/selenium"
)

const (
	idArtigo              = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlArtigo"
	idTipoOcorrencia      = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlTipoInfracao"
	idElemento            = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlElementoOcorrencia"
	idPrazo               = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtPrazoExecucaoOcorrencia"
	idTipoPrazo           = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlExecucaoOcorrencia"
	idData                = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtDataOcorrencia"
	idHora                = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtHoraOcorrencia"
	idUf                  = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlUf"
	idRodovia             = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlRodovia"
	idPista               = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlPista"
	idSentido             = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlSentido"
	idkmInicial           = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtKmInicial"
	idKmFinal             = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtKmFinal"
	idDescricaoOcorrencia = "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_txtDescricaoOcorrencia"
	idProcessando         = "Progress_LabelProcessando"
)

var (
	art                 string
	tipoOcorrencia      string // cod 774 - buracos
	tipoTempoHora       string // 1 corresponde a horas
	prazo               string
	observacao          string
	data                string
	hora                string
	uf                  string
	rodovia             string
	pista               string
	sentido             string
	kmInicial           string
	kmFinal             string
	concessionariaValue string
	err                 error
)

func getConcessionariaValue() (string, error) {
	if util.CONCESSIONARIA == "MSVIA" {
		return "19642306000170", nil
	} else if util.CONCESSIONARIA == "CRO" {
		return "19521322000104", nil

	} else if util.CONCESSIONARIA == "ECO050" {
		return "19208022000170", nil
	}
	return "", fmt.Errorf("não foi possível determinar a concessionária.")
}

func InicioDigitacao() error {

	go KeepMouseMoving()

	ops := []selenium.ServiceOption{}
	_, err := selenium.NewChromeDriverService(util.SELENIUMPATH, util.SeleniumPORT, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	driver, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	// defer driver.Quit()

	if err := driver.Get("https://appweb1.antt.gov.br/fisn/Site/TRO/Cadastrar.aspx"); err != nil {
		panic(err)
	}
	fmt.Println("abrindo pagina do Sifama")

	usuario, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxUsuario")
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}
	senha, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxSenha")
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}
	entrar, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ButtonOk")
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err))
	}

	usuario.SendKeys(util.USER)
	senha.SendKeys(util.PWD)

	fmt.Println("entrando com senha")

	entrar.Click()

	waitForJsAndJquery(driver)

	err = inicioTro(driver)
	if err != nil {
		return err
	}

	// Alert After job is done.
	driver.ExecuteScript("alert('Terminou')", nil)

	return nil

}

func inicioTro(driver selenium.WebDriver) error {

	var t Tro

	troList := t.FindAll()

	totalTro := len(troList)

	primeiro := true
	for i, tro := range troList {
		time.Sleep(time.Second / 2)

		if !primeiro {
			waitForJsAndJquery(driver)
			driver.ExecuteScript("document.getElementById('MessageBox_ButtonOk').click()", nil)
		}
		primeiro = false
		actualTro := i + 1

		waitForProcessBar(driver, idProcessando)

		err = registroTro(tro, driver, actualTro, totalTro)
		if err != nil {
			return err
		}

	}
	return nil
}

func registroTro(tro Tro, driver selenium.WebDriver, actualTro, totalTro int) error {

	waitForJsAndJquery(driver)

	fmt.Println("Selecionando Concessionária na lista")

	concessionariaValue, err = getConcessionariaValue()
	if err != nil {
		return fmt.Errorf("não foi possível identificar a concessionária")
	}

	// time.Sleep(time.Minute / 3)

	jqueryScriptWithChange(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlConcessionaria", concessionariaValue)

	waitForProcessBar(driver, idProcessando)

	fmt.Println("Seleciona Resolução")

	jqueryScriptWithChange(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlLegislacao", "4071")

	waitForProcessBar(driver, idProcessando)

	waitForJsAndJquery(driver)

	locais := tro.Locais

	go resizeImageAndCopyMetadataFromOriginal3(&locais)

	palavraChave := tro.PalavraChave

	artigoList := GetDisposicaoLegal(palavraChave)
	art = artigoList[0]
	tipoOcorrencia = artigoList[1]

	observacao = tro.Observacao
	observacao = strings.Title(observacao)

	// verificar se data e hora não é para cada local
	data = tro.Locais[0].Data
	hora = tro.Locais[0].Hora

	fmt.Println("Seleciona Artigo da Resolução")

	//Pega da planilha
	jqueryScriptWithChange(driver, idArtigo, art)

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	fmt.Println("Seleciona TipoOcorrencia")

	jqueryScriptWithChange(driver, idTipoOcorrencia, tipoOcorrencia)

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	prazo = tro.Prazo

	enviaChaves(driver, idPrazo, prazo)

	fmt.Println("Seleciona Entre horas / dias")

	jqueryScriptWithChange(driver, idTipoPrazo, tro.TipoPrazo)

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	fmt.Println("informa data : ", data)

	we, err := waitForElementById(driver, idData, 10) // wait for clickable
	errorHandle(err)
	we.Click()
	we.Clear()
	// waitForProcessBar(driver)

	enviaChaves(driver, idData, data)

	waitForJsAndJquery(driver)

	// consulta.waitForProcessBar();

	we, err = waitForElementById(driver, idHora, 10)
	errorHandle(err)
	we.Click()

	waitForJsAndJquery(driver)
	// consulta.waitForProcessBar();

	fmt.Println("insere descrição ocorrencia")

	flag := false

	for !flag {

		we, err := waitForElementById(driver, idDescricaoOcorrencia, 10)
		errorHandle(err)
		we.Click()
		if e := enviaChaves(driver, idDescricaoOcorrencia, observacao); e == nil {
			flag = true
		}
	}

	waitForJsAndJquery(driver)
	fmt.Println("informa Hora: ", hora)

	enviaChaves(driver, idHora, hora)

	waitForJsAndJquery(driver)

	we, err = waitForElementById(driver, idDescricaoOcorrencia, 20)
	errorHandle(err)
	we.Click()

	waitForProcessBar(driver, idProcessando)

	waitForJsAndJquery(driver)

	for _, local := range locais {

		rodovia = local.Rodovia
		uf = local.Estado
		pista = local.Pista
		sentido = local.Sentido
		kmInicial = local.KmInicial
		kmFinal = local.KmFinal

		fmt.Println("Insere UF")

		jqueryScript(driver, idUf, uf)

		fmt.Println("insere rodovia")

		jqueryScript(driver, idRodovia, rodovia)
		// consulta.checkForErrors();

		fmt.Println("insere pista")

		jqueryScript(driver, idPista, pista)

		fmt.Println("insere sentido")

		jqueryScript(driver, idSentido, sentido)

		waitForProcessBar(driver, idProcessando)
		waitForJsAndJquery(driver)

		fmt.Println("insere Km Inicial e Final")

		we, err = waitForElementById(driver, idkmInicial, 30)
		errorHandle(err)

		enviaChaves(driver, idkmInicial, kmInicial)

		// consulta.checkForErrors();

		we, err = waitForElementById(driver, idKmFinal, 30)
		errorHandle(err)

		enviaChaves(driver, idKmFinal, kmFinal)

		fmt.Println("Incluindo Local .....")

		for {
			we, e := waitForElementById(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnIncluirLocal", 30)
			err = we.Click()
			if e == nil && err == nil {
				break
			} else {
				fmt.Println("não clicou no incluir local, tentando novamente")
				checkForErrors(driver)
			}

			time.Sleep(time.Second / 2)

		}
		waitForProcessBar(driver, idProcessando)
		// consulta.checkForErrors();

	}
	countImages := 0
	for _, local := range locais {

		kmInicial = local.KmInicial
		kmFinal = local.KmFinal

		fmt.Println("kmInicial: ", kmInicial)
		fmt.Println("kmFinal: ", kmFinal)

		for _, foto := range local.Fotos {
			fmt.Println(foto.Nome)

			err := jqueryScriptWithChangeByText(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlFotoLocal", kmInicial, kmFinal)
			if err != nil {
				fmt.Println(err)
			}

			imgpath := filepath.Join(util.OUTPUTIMAGEFOLDER, foto.Nome)

			fmt.Println(imgpath)

			err = enviaChaves(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_uplFotoLocal", imgpath)
			errorHandle(err)

			countImages++

			fmt.Printf("Enviando foto nº %d ............\n", countImages)

			time.Sleep(time.Second / 2)

			waitForProcessBar(driver, idProcessando)
			waitForJsAndJquery(driver)

			waitForElementById(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnIncluirFoto", 20)

			scriptToClick(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnIncluirFoto")

			waitForProcessBar(driver, idProcessando)
			waitForJsAndJquery(driver)

			time.Sleep(time.Second / 2)
			fmt.Println("OK !")
		}
	}

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	fmt.Printf("Salva o TRO %d/%d .......... ", actualTro, totalTro)

	checkForErrors(driver)

	waitForElementById(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnSalvar", 20)

	err = scriptToClick(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnSalvar")
	errorHandle(err)
	// err = waitForElementToBeClickableAndClick(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_btnSalvar")

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	checkForErrors(driver)

	err = waitForElementToBeClickableAndClick(driver, "MessageBox_ButtonOk")
	errorHandle(err)
	time.Sleep(time.Second * 2)
	checkForErrors(driver)

	waitForProcessBar(driver, idProcessando)
	waitForJsAndJquery(driver)

	fmt.Println("OK ")

	return nil

}

func resizeImageAndCopyMetadataFromOriginal3(locaisP *[]Local) {
	locais := *locaisP

	for _, local := range locais {
		fotos := local.Fotos
		for _, foto := range fotos {
			resizeImageAndCopyMetadataFromOriginal1(&foto)
		}
	}
}
func resizeImageAndCopyMetadataFromOriginal1(foto *Foto) {
	originPath := foto.OriginPath
	destPath := filepath.Join(util.OUTPUTIMAGEFOLDER, foto.Nome)
	// **
	resizeImageAndCopyMetadataFromOriginal(destPath, originPath, util.MAXIMAGEWIDTH)
}
