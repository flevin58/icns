// Defines the structure of an 'icns' file.
// This package is based on the specs at: https://en.wikipedia.org/wiki/Apple_Icon_Image_format
// Also, here below is an article to work with PNG files (load, save, resize, etc.)
// https://gist.github.com/logrusorgru/570d64fd6a051e0441014387b89286ca

package icns

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type String4 [4]byte

func (s *String4) String() string {
	result := strings.Builder{}
	for _, ch := range s {
		result.WriteByte(ch)
	}
	return result.String()
}

type IconFile struct {
	Header IcnsHeader
	Data   []IconData
}

type IcnsHeader struct {
	Magic    String4
	Filesize uint32
}

type IconData struct {
	IconType String4
	IconLen  uint32
	IconData []byte
}

type IconType struct {
	ostype    string
	length    int
	width     int
	height    int
	osversion float32
	desc      string
}

var IconTypes = []IconType{
	{"ic14", 0, 512, 512, 10.8, "JPEG 2000 or PNG format (256x256@2x retina)"},
	{"ic13", 0, 256, 256, 10.8, "JPEG 2000 or PNG format (128x128@2x retina)"},
	{"ic12", 0, 64, 64, 10.8, "JPEG 2000 or PNG format (32x32@2x retina)"},
	{"ic11", 0, 32, 32, 10.8, "JPEG 2000 or PNG format (16x16@2x retina)"},
	{"ic10", 0, 1024, 1024, 10.7, "JPEG 2000 or PNG format (512x512@2x retina)"},
}

func NormalizeIconName(path string) string {
	var base, name string
	parts := strings.Split(strings.ToLower(path), "/")
	for _, part := range parts {
		if strings.HasSuffix(part, ".app") {
			base = strings.Replace(part, ".app", "", 1)
		}
		if strings.HasSuffix(part, ".icns") {
			name = strings.Replace(part, ".icns", "", 1)
		}
	}
	if base == "" {
		base = parts[0]
	}
	return fmt.Sprintf("%s_%s", base, name)
}

func GetIconSize(itype string) int {
	for _, it := range IconTypes {
		if itype == it.ostype {
			return it.width
		}
	}
	return 0
}

func invalidHeaderType(headertype String4) bool {
	if headertype.String() == "icns" {
		return false
	}
	if headertype[0] == 0x89 && string(headertype[1:4]) == "PNG" {
		return false
	}
	return true
}

func (ih *IcnsHeader) LoadFromReader(r io.Reader) error {
	binary.Read(r, binary.BigEndian, ih)
	if invalidHeaderType(ih.Magic) {
		return fmt.Errorf("bad magic number %s", ih.Magic.String())
	}
	return nil
}

func (ih *IcnsHeader) String() string {
	return fmt.Sprintf("magic: '%s', size:%d", ih.Magic.String(), ih.Filesize)
}

func invalidIconType(icontype String4) bool {
	for _, itype := range IconTypes {
		if icontype.String() == itype.ostype {
			return false
		}
	}
	return true
}

func (id *IconData) LoadFromReader(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &id.IconType); err != nil {
		return err
	}
	if invalidIconType(id.IconType) {
		return fmt.Errorf("invalid icon type: %s", id.IconType.String())
	}

	if err := binary.Read(r, binary.BigEndian, &id.IconLen); err != nil {
		return err
	}

	id.IconData = make([]byte, id.IconLen)
	if err := binary.Read(r, binary.BigEndian, id.IconData); err != nil {
		return err
	}

	return nil
}

func (id *IconData) String() string {
	return fmt.Sprintf("type: '%s', length: %d", id.IconType.String(), id.IconLen)
}

func NewIconFromFile(filename string) (*IconFile, error) {
	fh, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	// Load header
	icf := new(IconFile)
	if err := icf.Header.LoadFromReader(fh); err != nil {
		return nil, err
	}

	// Load all icons
	for i := 0; ; i++ {
		curpos, _ := fh.Seek(0, io.SeekCurrent)
		if curpos >= int64(icf.Header.Filesize) {
			break
		}
		data := IconData{}
		if err := data.LoadFromReader(fh); err == nil {
			icf.Data = append(icf.Data, data)
		}
	}

	return icf, nil
}
