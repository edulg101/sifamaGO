package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	art            string = ""
	tipoOcorrencia string = ""  // cod 774 - buracos
	tipoTempoHora  string = "1" // corresponde a horas
	prazo          string = ""
	observacao     string = ""
	data           string = ""
	hora           string = ""
	uf             string = "MT"
	rodovia        string = ""
	pista          string = ""
	sentido        string = ""
	kmInicial      string = ""
	kmFinal        string = ""
)

func inicioDigitacao() {

	go keepMouseMoving()

	ops := []selenium.ServiceOption{}
	_, err := selenium.NewChromeDriverService(SELENIUMPATH, PORT, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	driver, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		panic(err)
	}

	// defer driver.Quit()

	if err := driver.Get("https://appweb1.antt.gov.br/fisn/Site/TRO/Cadastrar.aspx"); err != nil {
		panic(err)
	}
	fmt.Println("abrindo pagina do Sifama")

	usuario, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxUsuario")
	senha, e := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxSenha")
	entrar, er := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ButtonOk")

	errorHandle(err)
	errorHandle(e)
	errorHandle(er)

	usuario.SendKeys(USER)
	senha.SendKeys(PWD)

	fmt.Println("entrando com senha")

	entrar.Click()

	waitForJsAndJquery(driver)

	inicioTro(driver)

	// Alert After job is done.
	driver.ExecuteScript("alert('Terminou')", nil)

}

func inicioTro(driver selenium.WebDriver) {

	// var troList []Tro
	var t Tro

	troList := t.findAll()
	// db.Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)

	totalTro := len(troList)

	primeiro := true
	for i, tro := range troList {
		time.Sleep(time.Second)

		if !primeiro {
			waitForJsAndJquery(driver)
			driver.ExecuteScript("document.getElementById('MessageBox_ButtonOk').click()", nil)
		}
		primeiro = false
		actualTro := i + 1

		l := tro.Locais
		for _, v := range l {
			fmt.Println(v.Fotos)
		}

		waitForProcessBar(driver, idProcessando)

		registroTro(tro, driver, actualTro, totalTro)

	}
}

func registroTro(tro Tro, driver selenium.WebDriver, actualTro, totalTro int) {

	waitForJsAndJquery(driver)

	fmt.Println("Selecionando CRO na lista")

	jqueryScriptWithChange(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlConcessionaria", "19521322000104")

	waitForProcessBar(driver, idProcessando)

	fmt.Println("Seleciona Resolução")

	jqueryScriptWithChange(driver, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlLegislacao", "4071")

	waitForProcessBar(driver, idProcessando)

	waitForJsAndJquery(driver)

	locais := tro.Locais

	palavraChave := tro.PalavraChave

	artigoList := getDisposicaoLegal(palavraChave)
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

	jqueryScriptWithChange(driver, idTipoPrazo, tipoTempoHora)

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

	// consulta.checkForErrors();

	// consulta.waitForProcessBar();

	fmt.Println("Insere UF")

	jqueryScript(driver, idUf, uf)

	for _, local := range locais {

		rodovia = local.Rodovia
		pista = local.Pista
		sentido = local.Sentido
		kmInicial = local.KmInicial
		kmFinal = local.KmFinal

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

			imgpath := filepath.Join(OUTPUTIMAGEFOLDER, foto.Nome)

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

}
func getDisposicaoLegal(palavraChave string) []string {

	palavraChave = strings.ToLower(palavraChave)

	if strings.Contains(palavraChave, "buraco") {
		return []string{"6", "774"}
	} else if strings.Contains(palavraChave, "afundamen") || strings.Contains(palavraChave, "escorregamento") || strings.Contains(palavraChave, "remendo") {
		return []string{"6", "773"}
	} else if strings.Contains(palavraChave, "drenagem") {
		return []string{"6", "782"}
	} else if strings.Contains(palavraChave, "meio fio") || strings.Contains(palavraChave, "meio-fio") {
		return []string{"4", "720"}
	} else if strings.Contains(palavraChave, "desplacam") || strings.Contains(palavraChave, "deformaç") {
		return []string{"6", "773"}
	} else if strings.Contains(palavraChave, "vertical") || strings.Contains(palavraChave, "horizontal") {
		return []string{"7", "807"}
	} else if strings.Contains(palavraChave, "terrapleno") || strings.Contains(palavraChave, "talude") {
		return []string{"6", "783"}
	} else if strings.Contains(palavraChave, "defensa") {
		return []string{"7", "808"}
	} else if strings.Contains(palavraChave, "instalac") || strings.Contains(palavraChave, "instalaç") || strings.Contains(palavraChave, "edifica") {
		return []string{"5", "742"}
	} else if strings.Contains(palavraChave, "pmv") {
		return []string{"5", "752"}
	} else if strings.Contains(palavraChave, "guarda corpo") || strings.Contains(palavraChave, "guarda-corpo") {
		return []string{"7", "811"}
	} else if strings.Contains(palavraChave, "sujeira") {
		return []string{"7", "806"}
	}

	disposicaoLegal := make(map[string][]string)

	disposicaoLegal["4-v"] = []string{"4", "718"}
	disposicaoLegal["4-vi"] = []string{"4", "719"}
	disposicaoLegal["4-vii"] = []string{"4", "720"}
	disposicaoLegal["4-xii"] = []string{"4", "725"}
	disposicaoLegal["4-xiii"] = []string{"4", "726"}

	disposicaoLegal["5-iii"] = []string{"5", "742"}
	disposicaoLegal["5-v"] = []string{"5", "744"}
	disposicaoLegal["5-ix"] = []string{"5", "748"}
	disposicaoLegal["5-xii"] = []string{"5", "751"}
	disposicaoLegal["5-xiii"] = []string{"5", "752"}
	disposicaoLegal["5-xiv"] = []string{"5", "753"}
	disposicaoLegal["5-xv"] = []string{"5", "754"}
	disposicaoLegal["5-xxviii"] = []string{"5", "767"}

	disposicaoLegal["6-iii"] = []string{"6", "773"}
	disposicaoLegal["6-iv"] = []string{"6", "774"}
	disposicaoLegal["6-v"] = []string{"6", "775"}
	disposicaoLegal["6-vii"] = []string{"6", "777"}
	disposicaoLegal["6-viii"] = []string{"6", "778"}
	disposicaoLegal["6-x"] = []string{"6", "780"}
	disposicaoLegal["6-xi"] = []string{"6", "781"}
	disposicaoLegal["6-xii"] = []string{"6", "782"}
	disposicaoLegal["6-xiii"] = []string{"6", "783"}
	disposicaoLegal["6-xiv"] = []string{"6", "784"}
	disposicaoLegal["6-xvi"] = []string{"6", "786"}
	disposicaoLegal["6-xvii"] = []string{"6", "787"}
	disposicaoLegal["6-xxviii"] = []string{"6", "798"}

	disposicaoLegal["7-viii"] = []string{"7", "806"}
	disposicaoLegal["7-ix"] = []string{"7", "807"}
	disposicaoLegal["7-x"] = []string{"7", "808"}
	disposicaoLegal["7-xii"] = []string{"7", "810"}
	disposicaoLegal["7-xiii"] = []string{"7", "811"}

	disposicaoLegal["8-vii"] = []string{"8", "864"}

	disposicaoLegal["9-vii"] = []string{"9", "863"}

	return disposicaoLegal[palavraChave]
}

func getDescricaoDisposicaoLegal(codstr string) string {

	cod, err := strconv.Atoi(codstr)
	errorHandle(err)

	descricao := make(map[int]string)

	descricao[718] = "Art. 4º, V - deixar selagem em juntas de pavimento rígido ou trincas em desconformidade com o PER, por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[719] = "Art. 4º, VI - deixar de manter marcos quilométricos ou mantê-los em más condições de visibilidade, por prazo superior a 7 (sete) dias, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[720] = "Art. 4º, VII - deixar meios-fios danificados, deteriorados ou ausentes por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[725] = "Art. 4º, XII - deixar barreira de concreto de Obra-de-Arte Especial - OAE sem pintura por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[726] = "Art. 4º, XIII - deixar armaduras de OAE sem recobrimento por prazo superior a 48 (quarenta e oito horas)"

	// Art 5

	descricao[742] = "Art 5º, III - deixar de executar os serviços de conservação das instalações, áreas operacionais e bens vinculados à concessão por prazo superior a 72 horas após a ocorrência de evento que comprometa suas condições normais de uso e a integridade do bem"
	descricao[744] = "Art 5º, V - deixar de remover, da faixa de domínio, material resultante de poda, capina ou obras no prazo de 48 (quarenta e oito) horas, salvo no caso de materiais reaproveitáveis ou de bota-foras autorizados pela ANTT"
	descricao[748] = "Art 5º, IX - deixar de repor ou manter tachas, tachões e balizadores refletivos danificados ou ausentes no prazo de 72 (setenta e duas) horas"
	descricao[751] = "Art 5º, XII - deixar de adotar medidas, ainda que provisórias, para reparação de cercamento nas áreas operacionais por prazo superior a 24 (vinte e quatro) horas"
	descricao[752] = "Art 5º, XIII - deixar de adotar medidas, ainda que provisórias, para reparar painel de mensagem variável inoperante ou em condições que não permitam a transmissão de informações aos usuários, por prazo superior a 72 (setenta e duas) horas"
	descricao[752] = "Art 5º, XIV - deixar de adotar medidas, ainda que provisórias para reparação das cercas limítrofes da faixa de proteção e de seus aceiros por prazo superior a 72 (setenta e duas) horas"
	descricao[754] = "Art 5º, XV - deixar de adotar medidas, ainda que provisórias, para corrigir falha em sistema ou equipamento dos postos de pesagem no prazo de 24 (vinte e quatro) horas ou de acordo com o especificado no Contrato e/ou PER, se este fizer referência diversa"
	descricao[767] = "Art 5º, XXVIII - deixar de adotar providências para corrigir desnível entre faixas contíguas, ainda que em caráter provisório, no prazo de 24 (vinte e quatro) horas, ou, deixar de implementar a solução definitiva para correção no prazo estabelecido pela ANTT"

	// Art 6:

	descricao[773] = "Art 6º, III - deixar de corrigir depressões, abaulamentos (escorregamentos de massa asfáltica) ou áreas exsudadas na pista ou no acostamento, no prazo de 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[774] = "Art 6º, IV - deixar de corrigir/tapar buracos, panelas na pista ou no acostamento, no prazo de 24 (vinte e quatro) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[775] = "Art 6º, V - deixar de corrigir, no pavimento rígido, defeitos com grau de severidade alto, no prazo de 7 (sete) dias, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[777] = "Art 6º, VII - deixar de corrigir, no pavimento rígido, defeitos de alçamento de placa, fissura de canto, placa dividida (rompida), escalonamento ou degrau, placa bailarina, quebras localizadas e buracos no prazo de 48 (quarenta e oito) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[778] = "Art 6º, VIII - deixar de manter ou manter de forma não visível pelos usuários sinalização (vertical ou aérea) de indicação, de serviços auxiliares ou educativas, por prazo superior a 7 (sete) dias"
	descricao[780] = "Art 6º, X - deixar de manter ou manter de forma não funcional dispositivo anti-ofuscante por prazo superior a 7 (sete) dias, ou conforme previsto no Contrato de Concessão ou no PER"
	descricao[781] = "Art 6º, XI - deixar com problemas de conservação elemento de OAE, exceto guarda-corpo, por prazo superior a 30 (trinta) dias ou conforme Contrato de Concessão e/ou PER"
	descricao[782] = "Art 6º, XII - deixar de reparar, limpar ou desobstruir sistema de drenagem e Obra-de-Arte Corrente-OAC por prazo superior a 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão ou no PER"
	descricao[783] = "Art 6º, XIII - deixar de adotar providências para solucionar, ainda que de modo provisório, processo erosivo ou condição de instabilidade em talude, por prazo superior a 72 (setenta e duas) horas, ou deixar de implementar solução definitiva no prazo estabelecido pela ANTT"
	descricao[784] = "Art 6º, XIV - deixar de manter ou manter de forma não funcional o sistema de iluminação da rodovia, por prazo superior a 48 (quarenta e oito) horas"
	descricao[786] = "Art 6º, XVI - deixar de corrigir falha em equipamento de praça de pedágio no prazo de 6 (seis) horas, sem prejuízo ao atendimento dos parâmetros de desempenho estabelecidos no PER"
	descricao[787] = `Art 6º, XVII - deixar "Call Box" inoperante por prazo superior a 24 (vinte e quatro) horas, ou de acordo com o especificado no PER, se este fizer referência diversa`
	descricao[798] = "Art 6º, XXVIII - deixar de intervir, mesmo que provisoriamente, em recalque em pavimento na cabeceira de OAE e/ou OAC por prazo superior a 72 (setenta e duas) horas, desde que essa obrigação tenha sido prevista no Contrato de Concessão ou PER"

	// Art 7

	descricao[806] = "Art 7º, VIII - deixar de remover material da(s) faixa(s) de rolamento( s) ou acostamento(s) que obstrua ou comprometa a correta fluidez do tráfego no prazo de 6 (seis) horas a partir do evento que lhe deu origem"
	descricao[807] = "Art 7º, IX - deixar de manter ou manter a sinalização horizontal, vertical ou aérea, em desconformidade com as normas técnicas vigentes, por prazo superior ao estabelecido pela ANTT, excluídas as ocorrências previstas nos artigos 5°, 6° e 9°"
	descricao[808] = "Art 7º, X - deixar de recompor barreira rígida ou defensa metálica danificada no prazo de 48 horas"
	descricao[810] = "Art 7º, 	XII - deixar de intervir para restaurar a funcionalidade de elemento da rodovia quando da ocorrência de fatos oriundos da ação de terceiros ou de eventos da natureza que possam colocar em risco a segurança do usuário, no prazo de 48 (quarenta e oito) horas ou conforme estabelecido pela ANTT"
	descricao[811] = "Art 7º, XIII - deixar de recuperar, ainda que provisoriamente, guarda- corpo de OAE, inclusive passarela, por prazo superior a 24 (vinte e quatro) horas, ou, deixar de efetuar sua reposição definitiva, por prazo superior a 72 (setenta e duas) horas, ou conforme Contrato e/ou PER"

	// Art 8

	descricao[864] = "Art 8º, VII - deixar de adotar as providências cabíveis, inclusive por vias judiciais, para garantia do patrimônio da rodovia, da faixa de domínio, das edificações e dos bens da concessão, inclusive quanto à implantação de acessos irregulares e ocupações ilegais; Nos casos de constatação destas irregularidades para as concessões da 2ª etapa, há previsão contratual de prazo de 24 (vinte e quatro) horas para a correção. Deste modo, deverá ser expedido TRO enquadrado neste mesmo Art. 8º, inciso VII, da Resolução"

	// Art 9

	descricao[863] = "Art 9º, VII - deixar de manter ou manter sinalização vertical de regulamentação em desconformidade com as normas técnicas vigentes, por prazo superior ao previsto no Contrato de Concessão ou no PER"

	return descricao[cod]
}

//art 4
//        718	V - deixar selagem em juntas de pavimento rígido ou trincas em desconformidade com o PER, por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        719	VI - deixar de manter marcos quilométricos ou mantê-los em más condições de visibilidade, por prazo superior a 7 (sete) dias, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        720	VII - deixar meios-fios danificados, deteriorados ou ausentes por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        725	XII - deixar barreira de concreto de Obra-de-Arte Especial - OAE sem pintura por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        726	XIII - deixar armaduras de OAE sem recobrimento por prazo superior a 48 (quarenta e oito horas)

// Art 5

//      742	III - deixar de executar os serviços de conservação das instalações, áreas operacionais e bens vinculados à concessão por prazo superior a 72 horas após a ocorrência de evento que comprometa suas condições normais de uso e a integridade do bem
//      744	V - deixar de remover, da faixa de domínio, material resultante de poda, capina ou obras no prazo de 48 (quarenta e oito) horas, salvo no caso de materiais reaproveitáveis ou de bota-foras autorizados pela ANTT
//      748	IX - deixar de repor ou manter tachas, tachões e balizadores refletivos danificados ou ausentes no prazo de 72 (setenta e duas) horas
//      751	XII - deixar de adotar medidas, ainda que provisórias, para reparação de cercamento nas áreas operacionais por prazo superior a 24 (vinte e quatro) horas
//      752	XIII - deixar de adotar medidas, ainda que provisórias, para reparar painel de mensagem variável inoperante ou em condições que não permitam a transmissão de informações aos usuários, por prazo superior a 72 (setenta e duas) horas
//      753	XIV - deixar de adotar medidas, ainda que provisórias para reparação das cercas limítrofes da faixa de proteção e de seus aceiros por prazo superior a 72 (setenta e duas) horas
//      754	XV - deixar de adotar medidas, ainda que provisórias, para corrigir falha em sistema ou equipamento dos postos de pesagem no prazo de 24 (vinte e quatro) horas ou de acordo com o especificado no Contrato e/ou PER, se este fizer referência diversa
//      767	XXVIII - deixar de adotar providências para corrigir desnível entre faixas contíguas, ainda que em caráter provisório, no prazo de 24 (vinte e quatro) horas, ou, deixar de implementar a solução definitiva para correção no prazo estabelecido pela ANTT

//Art 6:

//       773	III - deixar de corrigir depressões, abaulamentos (escorregamentos de massa asfáltica) ou áreas exsudadas na pista ou no acostamento, no prazo de 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      774	IV - deixar de corrigir/tapar buracos, panelas na pista ou no acostamento, no prazo de 24 (vinte e quatro) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      775	V - deixar de corrigir, no pavimento rígido, defeitos com grau de severidade alto, no prazo de 7 (sete) dias, ou conforme previsto no Contrato de Concessão e/ou PER
//      777	VII - deixar de corrigir, no pavimento rígido, defeitos de alçamento de placa, fissura de canto, placa dividida (rompida), escalonamento ou degrau, placa bailarina, quebras localizadas e buracos no prazo de 48 (quarenta e oito) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      778	VIII - deixar de manter ou manter de forma não visível pelos usuários sinalização (vertical ou aérea) de indicação, de serviços auxiliares ou educativas, por prazo superior a 7 (sete) dias
//      780	X - deixar de manter ou manter de forma não funcional dispositivo anti-ofuscante por prazo superior a 7 (sete) dias, ou conforme previsto no Contrato de Concessão ou no PER
//      781	XI - deixar com problemas de conservação elemento de OAE, exceto guarda-corpo, por prazo superior a 30 (trinta) dias ou conforme Contrato de Concessão e/ou PER
//      782	XII - deixar de reparar, limpar ou desobstruir sistema de drenagem e Obra-de-Arte Corrente-OAC por prazo superior a 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão ou no PER
//      783	XIII - deixar de adotar providências para solucionar, ainda que de modo provisório, processo erosivo ou condição de instabilidade em talude, por prazo superior a 72 (setenta e duas) horas, ou deixar de implementar solução definitiva no prazo estabelecido pela ANTT
//      784	XIV - deixar de manter ou manter de forma não funcional o sistema de iluminação da rodovia, por prazo superior a 48 (quarenta e oito) horas
//      786	XVI - deixar de corrigir falha em equipamento de praça de pedágio no prazo de 6 (seis) horas, sem prejuízo ao atendimento dos parâmetros de desempenho estabelecidos no PER
//      787	XVII - deixar "Call Box" inoperante por prazo superior a 24 (vinte e quatro) horas, ou de acordo com o especificado no PER, se este fizer referência diversa
//      798	XXVIII - deixar de intervir, mesmo que provisoriamente, em recalque em pavimento na cabeceira de OAE e/ou OAC por prazo superior a 72 (setenta e duas) horas, desde que essa obrigação tenha sido prevista no Contrato de Concessão ou PER

//Art 7

//        806	VIII - deixar de remover material da(s) faixa(s) de rolamento( s) ou acostamento(s) que obstrua ou comprometa a correta fluidez do tráfego no prazo de 6 (seis) horas a partir do evento que lhe deu origem
//          807	IX - deixar de manter ou manter a sinalização horizontal, vertical ou aérea, em desconformidade com as normas técnicas vigentes, por prazo superior ao estabelecido pela ANTT, excluídas as ocorrências previstas nos artigos 5°, 6° e 9°
//        808	X - deixar de recompor barreira rígida ou defensa metálica danificada no prazo de 48 horas
//          810	XII - deixar de intervir para restaurar a funcionalidade de elemento da rodovia quando da ocorrência de fatos oriundos da ação de terceiros ou de eventos da natureza que possam colocar em risco a segurança do usuário, no prazo de 48 (quarenta e oito) horas ou conforme estabelecido pela ANTT
//          811	XIII - deixar de recuperar, ainda que provisoriamente, guarda- corpo de OAE, inclusive passarela, por prazo superior a 24 (vinte e quatro) horas, ou, deixar de efetuar sua reposição definitiva, por prazo superior a 72 (setenta e duas) horas, ou conforme Contrato e/ou PER

// Art 8

//        864	VII - deixar de adotar as providências cabíveis, inclusive por vias judiciais, para garantia do patrimônio da rodovia, da faixa de domínio, das edificações e dos bens da concessão, inclusive quanto à implantação de acessos irregulares e ocupações ilegais; Nos casos de constatação destas irregularidades para as concessões da 2ª etapa, há previsão contratual de prazo de 24 (vinte e quatro) horas para a correção. Deste modo, deverá ser expedido TRO enquadrado neste mesmo Art. 8º, inciso VII, da Re</option>

// Art 9

//        863	VII - deixar de manter ou manter sinalização vertical de regulamentação em desconformidade com as normas técnicas vigentes, por prazo superior ao previsto no Contrato de Concessão ou no PER
