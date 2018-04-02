package browser

import (
	"fmt"

	"github.com/tebeka/selenium"
)

type browser struct {
	HubURL  string
	HubPort int

	wd selenium.WebDriver
}

// close remote webDrive session
func (b browser) Close() error {

	return b.wd.Quit()
}

// take web page screenshot, return []byte
func (b browser) TakeScreenshot(url string) (imgData []byte, err error) {

	if err = b.wd.Get(url); err != nil {
		return nil, err
	}

	imgData, err = b.wd.Screenshot()
	if err != nil {
		return nil, err
	}

	return imgData, nil
}

// configures a new webBrowser
func NewBrowser(hubUrl string, hubPort int, client string) (*browser, error) {

	capabilities := selenium.Capabilities{"browserName": client}
	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://%s:%d/wd/hub", hubUrl, hubPort))
	if err != nil {
		return nil, err
	}

	webBrowser := &browser{
		HubURL:  hubUrl,
		HubPort: hubPort,
		wd:      wd,
	}

	return webBrowser, nil
}
