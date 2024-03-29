package yaac_backend_database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"image"
)

type Rectangle image.Rectangle

func (rect *Rectangle) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Could not unmarshal rectangle: ", value))
	}

	result := image.Rectangle{}
	return json.Unmarshal(bytes, &result)
}

func (rect Rectangle) Value() (driver.Value, error) {
	return json.Marshal(rect)
}
