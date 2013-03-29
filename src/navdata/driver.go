package navdata

import (
	"os"
)

const DefaultTTYPath = "/dev/ttyO1"

type Driver struct{
	*Decoder
	file *os.File
}

func NewDriver(ttyPath string) (*Driver, error) {
	file, err := os.OpenFile(ttyPath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		file: file,
		Decoder: NewDecoder(file),
	}
	
	if _, err := file.Write([]byte{3}); err != nil {
		return nil, err
	}

	return driver, nil
}
