package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

func jqueryScript(driver selenium.WebDriver, id string, value string) {
	var arguments []interface{}
	arguments = append(arguments, id, value)
	_, err := driver.ExecuteScript("$(`#${arguments[0]} option[value='${arguments[1]}']`).prop('selected', true)", arguments)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}

func waitForJsAndJquery(driver selenium.WebDriver) (error, error) {
	var err error
	var err1 error
	time.Sleep(time.Second)
	duration := time.Second * 50
	totalLoops := duration * 4
	for i := 0; i < int(totalLoops); i++ {
		js, err := jsLoadedFunc(driver)
		jquery, err1 := jQueryFunc(driver)

		if err == nil && err1 == nil && js && jquery {
			return err, err1
		}
		time.Sleep(time.Second / 4)
	}
	return err, err1
}

func waitForElementById(driver selenium.WebDriver, id string, duration time.Duration) (selenium.WebElement, error) {
	var err error
	var we selenium.WebElement
	totalLoops := duration * 4
	for i := 0; i < int(totalLoops); i++ {
		we, err := waitForElementByIdFunc(driver, id)
		if err == nil && we != nil {
			return we, err
		}
		time.Sleep(time.Second / 4)
	}
	return we, err
}
func waitForElementByXpath(driver selenium.WebDriver, id string) (selenium.WebElement, error) {
	var err error
	var we selenium.WebElement
	duration := time.Second * 50
	totalLoops := duration * 4
	for i := 0; i < int(totalLoops); i++ {
		we, err := waitForElementByXpathFunc(driver, id)
		if err == nil && we != nil {
			return we, err
		}
		time.Sleep(time.Second / 4)
	}
	return we, err
}

func scriptToClick(driver selenium.WebDriver, id string) error {
	var args []interface{}
	args = append(args, id)
	_, err := driver.ExecuteScript("document.getElementById(arguments[0]).click()", args)
	return err
}

func jQueryFunc(driver selenium.WebDriver) (bool, error) {
	r, err := driver.ExecuteScript("return jQuery.active", nil)
	v, ok := r.(float64)
	if !ok {
		return false, errors.New("não conseguiu converter para float. Tente novamente")
	}
	x := int(v)
	return x == 0, err
}

func jsLoadedFunc(driver selenium.WebDriver) (bool, error) {
	d, err := driver.ExecuteScript("return document.readyState", nil)
	return d == "complete", err
}

func waitForProcessBar1(driver selenium.WebDriver) {

	we, e := driver.FindElement(selenium.ByID, "Progress_LabelProcessando")
	if e != nil || we == nil {
		fmt.Println("primeiro erro")
		fmt.Println(e)
	}

	displayed, e := we.IsDisplayed()
	if e != nil {
		fmt.Println("entrou no segundo is displaye! ")

	}
	if !displayed {
		fmt.Println("entrou bloco ! displayed")
		time.Sleep(time.Second)

	}
	if displayed {
		driver.WaitWithTimeoutAndInterval(progressIsDisplayed, 500, time.Second*10)
	}

}

func progressIsDisplayed(driver selenium.WebDriver) (bool, error) {
	we, err := driver.FindElement(selenium.ByID, "Progress_LabelProcessando")
	displayed, _ := we.IsDisplayed()
	return !displayed, err
}

func waitForElementByIdFunc(driver selenium.WebDriver, id string) (selenium.WebElement, error) {
	d, err := driver.FindElement("id", id)
	if d != nil && err == nil {
		return d, err
	} else {
		return d, err
	}
}

func waitForElementByXpathFunc(driver selenium.WebDriver, id string) (selenium.WebElement, error) {
	d, err := driver.FindElement(selenium.ByXPATH, id)
	if d != nil && err == nil {
		time.Sleep(time.Second / 2)
		return d, err
	} else {
		return d, err
	}
}

func scriptToFillField(driver selenium.WebDriver, id, value string) {

	var arguments []interface{}
	arguments = append(arguments, id, value)
	_, err := driver.ExecuteScript("document.getElementById(arguments[0]).setAttribute('value', arguments[1])", arguments)
	errorHandle(err)
}

func checkForErrors(driver selenium.WebDriver) {
	var divError selenium.WebElement
	var err error
	for i := 0; i < 10; i++ {
		divError, err = driver.FindElement("id", "MessageBox_LabelTitulo")
		divText, err1 := driver.FindElement("id", "MessageBox_LabelMensagem")
		errorHandle(err)
		errorHandle(err1)
		displayed, _ := divError.IsDisplayed()
		if displayed {
			title, _ := divError.Text()
			message, _ := divText.Text()
			title = strings.ToLower(title)
			fmt.Println(message)
			scanner := bufio.NewScanner(os.Stdin)
			if !strings.Contains(message, "salvar") && !strings.Contains(title, "cadastro") {
				fmt.Println(title)
				var resp string
				fmt.Println("Deu Erro. Ajuste e Digite em 'Sim' para continuar")
				for scanner.Scan() {
					resp = scanner.Text()
				}
				for strings.ToLower(resp) != "sim" {
					fmt.Println("Deu Erro. Ajuste e Digite em 'Sim' para continuar")
					for scanner.Scan() {
						resp = scanner.Text()
					}
				}
			} else {
				return
			}
		}
		time.Sleep(time.Second / 2)
	}
}

func jqueryScriptWithChange(driver selenium.WebDriver, id, value string) {
	var arguments []interface{}
	arguments = append(arguments, id, value)
	driver.ExecuteScript("$(`#${arguments[0]} option[value='${arguments[1]}']`).prop('selected', true).change()", arguments)
	waitForJsAndJquery(driver)
}

func jqueryScriptWithChangeByText(driver selenium.WebDriver, id, kmInicial, kmFinal string) error {
	var arguments []interface{}
	arguments = append(arguments, id, kmInicial, kmFinal)
	_, err := driver.ExecuteScript("$(`#${arguments[0]} option:contains('${arguments[1]} - ${arguments[2]}')`).attr('selected',true)", arguments)
	waitForJsAndJquery(driver)
	return err
}

func enviaChaves(driver selenium.WebDriver, id, value string) error {

	element, e := waitForElementById(driver, id, 10)
	errorHandle(e)
	element.Clear()
	time.Sleep(time.Second / 3)
	e = element.SendKeys(value)
	waitForJsAndJquery(driver)
	return e
}

func checkForProcessando(driver selenium.WebDriver) {
	for {
		element, _ := waitForElementById(driver, idProcessando, 20)
		fmt.Println(element.IsDisplayed())

		time.Sleep(time.Second / 2)
	}
}

func waitForProcessBar(driver selenium.WebDriver, id string) (selenium.WebElement, error) {
	var err error
	var we selenium.WebElement
	var secondsToWait time.Duration = 10
	var enabled bool
	duration := time.Second * secondsToWait
	totalLoops := int(duration) * 4 / 1000000000

	for i := 0; i < totalLoops/10; i++ {

		we, e := driver.FindElement(selenium.ByID, id)
		if err != nil {
			continue
		}
		enabled, err = we.IsDisplayed()
		fmt.Printf("primeiro For. vez %d. enabled? %v \n", i, enabled)
		if we != nil && e == nil && err == nil && enabled {
			break
		}
		time.Sleep(time.Second / 10)
	}
	if enabled {
		for i := 0; i < 1000; i++ {
			we, err := driver.FindElement(selenium.ByID, id)
			enabled, e := we.IsDisplayed()
			fmt.Printf("Segundo For. vez %d. enabled? %v \n", i, enabled)
			if we != nil && err == nil && e == nil && !enabled {
				return we, err
			}
			time.Sleep(time.Second / 4)
		}
		e := errors.New("Barra Processando Nunca Desapareceu")
		return we, e
	}
	return we, err
}

func waitForElementToBeClickableAndClick(driver selenium.WebDriver, id string) error {
	var e error

	for i := 0; i < 120; i++ {
		we, err := waitForElementById(driver, id, 10)
		if we != nil {
			e = we.Click()
		}
		if e == nil && err == nil {
			return nil
		}
		time.Sleep(time.Second / 4)
	}
	return e

}

func errorHandle(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}
