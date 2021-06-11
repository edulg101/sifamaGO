package main

import (
	"bufio"
	"fmt"
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

func scriptToClick(driver selenium.WebDriver, id string) {
	var args []interface{}
	args = append(args, id)
	driver.ExecuteScript("document.getElementById(arguments[0]).click()", args)
}

func jQueryFunc(driver selenium.WebDriver) (bool, error) {
	r, err := driver.ExecuteScript("return jQuery.active", nil)
	v := r.(float64)
	x := int(v)
	return x == 0, err
}

func jsLoadedFunc(driver selenium.WebDriver) (bool, error) {
	d, err := driver.ExecuteScript("return document.readyState", nil)
	return d == "complete", err
}

func waitForProcessBar(driver selenium.WebDriver) {

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
		time.Sleep(time.Second / 2)
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
	errorHandling(err)
}

func checkForErrors(driver selenium.WebDriver) {
	var divError selenium.WebElement
	var err error
	for i := 0; i < 10; i++ {
		divError, err = driver.FindElement("id", "MessageBox_LabelTitulo")
		errorHandling(err)
		displayed, _ := divError.IsDisplayed()
		if displayed {
			message, _ := divError.Text()
			message = strings.ToLower(message)
			scanner := bufio.NewScanner(os.Stdin)
			if !strings.Contains(message, "sucesso") {
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
			}
		}
		time.Sleep(time.Second / 2)
	}
}

func errorHandling(err error) {
	fmt.Println(err)
}
