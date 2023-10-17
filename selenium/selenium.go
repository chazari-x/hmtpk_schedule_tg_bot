package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
)

type Selenium struct {
	wd selenium.WebDriver
}

func NewSelenium() (*Selenium, selenium.WebDriver, error) {
	caps := selenium.Capabilities{"browserName": "chrome", "chrome.switches": []string{"--headless"}}

	wd, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		return nil, nil, err
	}

	return &Selenium{wd: wd}, wd, nil
}

func (c *Selenium) OpenURL(url string) error {
	return c.wd.Get(url)
}

func (c *Selenium) GetURL() (string, error) {
	return c.wd.CurrentURL()
}

func (c *Selenium) GetTitle() (string, error) {
	return c.wd.Title()
}

func (c *Selenium) GetElementValue(by, value string) (string, error) {
	element, err := c.wd.FindElement(by, value)
	if err != nil {
		return "", fmt.Errorf("find element err: %s", err)
	}

	return element.GetAttribute("value")
}

func (c *Selenium) GetElementText(by, value string) (string, error) {
	element, err := c.wd.FindElement(by, value)
	if err != nil {
		return "", err
	}

	return element.Text()
}
