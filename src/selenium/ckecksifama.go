package selenium

import (
	"fmt"
	"sifamaGO/src/config"
	"sifamaGO/src/util"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const (
	port = 9515
)

var seleniumPath string

func CheckSifama(user, pass string) (string, error) {

	err := config.GetEnv()
	if err != nil {
		panic(err)
	}

	seleniumPath = util.SELENIUMPATH

	ops := []selenium.ServiceOption{}
	_, err = selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	driver, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		if err != nil {
			return "", err
		}
	}

	defer driver.Quit()

	if err := driver.Get("https://appweb1.antt.gov.br/fisn/Site/TRO/Cadastrar.aspx"); err != nil {
		panic(err)
	}
	fmt.Println("abrindo pagina do Sifama")

	usuario, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxUsuario")
	if err != nil {
		return "", err
	}
	senha, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_TextBoxSenha")
	if err != nil {
		return "", err
	}
	entrar, err := driver.FindElement(selenium.ByID, "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ButtonOk")
	if err != nil {
		return "", err
	}

	usuario.SendKeys(user)
	senha.SendKeys(pass)
	entrar.Click()

	waitForJsAndJquery(driver)

	err = checkForErrors(driver)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	tros := getInfoFromExcel()
	if len(tros) <= 1 {
		return "", fmt.Errorf("não foi possivel importar TROs da planilha - verifique a planilha.")

	}

	waitForJsAndJquery(driver)

	quit := make(chan string)
	go KeepMouseMoving(quit)

	i := 1
	mainWindow, _ := driver.CurrentWindowHandle()
	var returnMessage string

	for ; i < len(tros); i++ {
		tro := tros[i]
		driver.SwitchWindow(mainWindow)
		if i == (len(tros) - 1) {
			go func() {
				quit <- "quit"
			}()
		}
		returnMessage, err = inicioVerificacao(driver, tro, pass, tros, i)
		if err != nil {
			return returnMessage, err
		}

	}

	driver.SwitchWindow(mainWindow)
	fmt.Println("Done")

	driver.ExecuteScript("alert('Terminou')", nil)

	return returnMessage, err

}

func inicioVerificacao(driver selenium.WebDriver, tro []string, pass string, tros [][]string, count int) (string, error) {
	mainWindow, _ := driver.CurrentWindowHandle()
	troNumber := tro[0]
	troHora := tro[1]
	troText := tro[2]
	var codAtendimento string = ""

	if strings.ToLower(tro[3]) == "s" || strings.ToLower(tro[3]) == "sim" {
		codAtendimento = "2"
	} else if strings.ToLower(tro[3]) == "n" || strings.ToLower(tro[3]) == "não" || strings.ToLower(tro[3]) == "nao" {
		codAtendimento = "3"
	} else {
		return "", fmt.Errorf("nao foi possivel identificar o codigo de atendimento ('sim ou nao'")
	}

	waitForJsAndJquery(driver)

	we, err := waitForElementByXpath(driver, "/html/body/div[1]/div[1]/div[1]/div[1]")
	if err != nil {
		return "", err
	}
	we.Click()
	fmt.Println("abrindo a lista de TROs ........")

	time.Sleep(time.Second * 2)

	waitForJsAndJquery(driver)

	listaTros, err := driver.FindElements(selenium.ByXPATH, "//div[@class='wingsDivNomeTarefa']")
	if err != nil {
		return "", err
	}

	totalTros := len(tros)
	sucessMessage := fmt.Sprintf("Foram Registratos %d TROs com Sucesso !", totalTros-1)

	var listaTrosEmOrdem []int
	for _, we := range listaTros {
		text, err := we.Text()
		if err != nil {
			return "", err
		}
		troStr := reg(text)
		troInt, _ := strconv.Atoi(troStr)
		listaTrosEmOrdem = append(listaTrosEmOrdem, troInt)
	}

	sort.Ints(listaTrosEmOrdem)

	fmt.Println("Tros disponiveis para análise:")

	for _, v := range listaTrosEmOrdem {
		fmt.Println(v)
	}

	if count == 1 {
		err := checkForMissingTros(tros, listaTrosEmOrdem)
		if err != nil {

			return "", fmt.Errorf("%s\n%s", err.Error(), "Remova esses TROs da Planilha e Tente Novamente...")
		}
	}

	for _, x := range listaTros {
		text, _ := x.Text()
		if strings.Contains(text, troNumber+"202") {
			getTBody, err := x.FindElement(selenium.ByXPATH, "../../..")
			if err != nil {
				return "", err
			}
			divToClick, _ := getTBody.FindElement(selenium.ByXPATH, "./tr[5]/td/div[2]")
			err = divToClick.Click()
			if err != nil {
				return "", err
			}
			bytes := []byte(text)
			bytes = bytes[25:]
			fmt.Println("Entrando no TRO n. " + string(bytes))
		}
	}

	time.Sleep(time.Second * 10)

	waitForJsAndJquery(driver)

	flag := true
	for i := 0; i < 120 && flag; i++ {
		handles, err := driver.WindowHandles()
		if err != nil {
			return "", err
		}
		for _, wh := range handles {
			if wh != mainWindow && flag {
				err = driver.SwitchWindow(wh)
				if err != nil {
					return "", err
				}
				we, err := driver.FindElement("id", "ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_div1")
				if we != nil && err == nil {
					flag = false
					break
				}

			}
			time.Sleep(time.Second / 2)
		}
		time.Sleep(time.Second / 2)
		fmt.Println(i)
	}

	todayDate := time.Now()
	today := todayDate.Format("02/01/2006")

	waitForJsAndJquery(driver)

	dataCampo, err := waitForElementByXpath(driver, `//input[@name="ctl00$ctl00$ctl00$ContentPlaceHolderCorpo$ContentPlaceHolderCorpo$ContentPlaceHolderCorpo$txtDataVerificacao"]`)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Second * 3)
	fmt.Println("Insere data")
	dataCampo.Clear()

	time.Sleep(time.Second * 3)
	dataCampo.Click()

	time.Sleep(time.Second * 3)
	dataCampo.SendKeys(today)

	passwordElement, err := driver.FindElement("id", PASSWORDCAMPO)
	if err != nil {
		return "", err
	}
	passwordElement.Click()

	waitForJsAndJquery(driver)

	_, err = waitForElementById(driver, HORACAMPO, time.Second*30)

	if err != nil {
		return "", err
	}

	fmt.Println("insere hora")

	we, err = driver.FindElement("id", HORACAMPO)
	if err != nil {
		return "", err
	}
	we.Clear()
	we.SendKeys(troHora)

	we.Clear()
	we.SendKeys(troHora)

	scriptToFillField(driver, HORACAMPO, troHora)

	waitForJsAndJquery(driver)

	fmt.Println("marca como atendido")

	jqueryScript(driver, ATENDIDOCAMPOSELECT, codAtendimento)
	waitForJsAndJquery(driver)
	time.Sleep(time.Second)
	fmt.Println("codigo atendimento:", codAtendimento)

	response, _ := getResponseFromScript(driver, "$('#ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlResultadoAnaliseExecucao').val()")
	waitForJsAndJquery(driver)

	fmt.Print("response:")
	fmt.Println(response)
	fmt.Println(response == codAtendimento)

	// doublecheck if codAtendimento has been correctly changed.
	actualLoop := 0
	for response != codAtendimento {
		jqueryScript(driver, ATENDIDOCAMPOSELECT, codAtendimento)
		time.Sleep(time.Second + time.Second*time.Duration(actualLoop))
		waitForJsAndJquery(driver)

		response, err = getResponseFromScript(driver, "$('#ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ContentPlaceHolderCorpo_ddlResultadoAnaliseExecucao').val()")
		actualLoop++
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("response:", response)

		fmt.Printf("tentativa %d para inserir cod atendimento\n", actualLoop)

		if actualLoop > 5 {
			return "", fmt.Errorf("não foi possivel alterar o campo Atendimento. abortando. internet ruim?")
		}
	}

	scriptToClick(driver, PASSWORDCAMPO)

	jqueryScript(driver, ATENDIDOCAMPOSELECT, codAtendimento)

	waitForJsAndJquery(driver)

	fmt.Println("Insere senha")

	for i := 0; i < 3; i++ {
		we, err = waitForElementById(driver, PASSWORDCAMPO, time.Second*30)
		if err != nil {
			return "", err
		}
		we.Clear()
		we.SendKeys(pass)
		time.Sleep(7000)

	}

	waitForElementById(driver, IFRAMEOBS, time.Second*30)

	driver.SwitchFrame(IFRAMEOBS)

	fmt.Println("insere texto na Observação")

	waitForElementById(driver, OBSCAMPO, time.Second*30)

	we, err = driver.FindElement("id", OBSCAMPO)
	if err != nil {
		return "", err
	}

	we.SendKeys(troText)

	driver.SwitchFrame(nil)
	time.Sleep(time.Second / 2)
	waitForJsAndJquery(driver)

	fmt.Println("Envia formulario")

	scriptToClick(driver, SALVARBUTTON)

	time.Sleep(time.Second)

	checkForErrorsCkeckSifama(driver)

	scriptToClick(driver, "MessageBox_ButtonOk")

	time.Sleep(time.Second)

	fmt.Printf("Salva Tro n. %v \n", troNumber)

	return sucessMessage, err

}
