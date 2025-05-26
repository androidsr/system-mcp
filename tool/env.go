package tool

import _ "embed"

var (
	WorkDir string

	//go:embed Readability.js
	Readability string
)

func Init(workDir string) {
	WorkDir = workDir
}
