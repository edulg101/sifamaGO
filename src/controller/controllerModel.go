package controller

type Request struct {
	StartDigitacao bool
	Restart        bool
	Compact        bool
	Folder         string
	Title          string
	User           string
	Passd          string
}
