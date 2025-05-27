package tool

import _ "embed"

var (
	WorkDir     string
	browserPath string
)

func Init(workDir, browser string) {
	WorkDir = workDir
	if browserPath == "" {
		browserPath = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	} else {
		browserPath = browser
	}
}
