package wp_spotlight

import "testing"

func TestNewFilename(t *testing.T) {
	const nfDesc string = "New File Name"
	const nfCRn string = "Bob Smith"
	const nfCRo1 string = "**Stock  Photos  Inc** E+"
	const nfCRo2 string = "SPI"
	const nfCRoX string = "Stock Photos Inc + E+ + SPI"
	const crSym string = "©"
	nfCR := crSym + nfCRn + " | " + nfCRo1 + " | " + nfCRo2

	a := new(asset)
	a.description = nfDesc
	a.copyright = nfCR
	input := a.description + " / " + a.copyright

	nfX := nfDesc + " " + crSym + " " + nfCRn + " + " + nfCRoX
	a.newFilename()

	if a.newName != nfX {
		t.Errorf("NewName incorrect:\ninput: '%s'\n  got: '%s'\n want: '%s'", input, a.newName, nfX)
	}
}
