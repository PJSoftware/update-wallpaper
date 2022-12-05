package ini_test

import (
	"testing"

	"github.com/pjsoftware/update-wallpaper/pkg/ini"
)

func TestParse(t *testing.T) {
	const testINI string = "./ini_test.ini"
	var fi ini.File
	err := fi.Parse(testINI)
	if err != nil {
		t.Errorf("Error parsing '%s'", testINI)
	}
}
