package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

const (
	chromeDriverPath = `/Users/vadimgeshiktor/libraries/chromewebdriver/chrome-mac-arm64/Google Chrome for Testing.app/Contents/MacOS/chromewebdrv`
)

func main() {
	// Set up ChromeDriver
	opts := []selenium.ServiceOption{}
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444, opts...)
	if err != nil {
		fmt.Println("Error starting the ChromeDriver service:", err)
		return
	}
	defer service.Stop()

	// configure the browser options
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		//"--headless-new", // comment out this line for testing
	}})

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {	
		fmt.Println("Error connecting to WebDriver:", err)
		return
	}
	defer wd.Quit()

	// Open WhatsApp Web
	wd.Get("https://web.whatsapp.com")

	// Wait for manual QR code scan
	fmt.Println("Please scan the QR code within 15 seconds...")
	time.Sleep(15 * time.Second)

	// Locate the search bar and input contact name or number
	searchBox, err := wd.FindElement(selenium.ByCSSSelector, "div[contenteditable='true']")
	if err != nil {
		fmt.Println("Error finding the search box:", err)
		return
	}
	searchBox.SendKeys("Your Contact Name")

	time.Sleep(2 * time.Second) // wait for search to show contact

	// Select contact from search results
	contact, err := wd.FindElement(selenium.ByCSSSelector, "span[title='Your Contact Name']")
	if err != nil {
		fmt.Println("Error finding contact:", err)
		return
	}
	contact.Click()

	// Type message
	messageBox, err := wd.FindElement(selenium.ByCSSSelector, "div[contenteditable='true']")
	if err != nil {
		fmt.Println("Error finding message box:", err)
		return
	}
	messageBox.SendKeys("Hello from Golang!" + selenium.EnterKey)

	fmt.Println("Message sent!")
}
