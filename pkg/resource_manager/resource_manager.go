package resource_manager

import (
	"bufio"
	"os"
	"path/filepath"
)

type Resource interface {
	Name() string
	Content() []byte
}
type StaticResource struct {
	StaticName    string
	StaticContent []byte
}

func (r *StaticResource) Name() string {
	return r.StaticName
}
func (r *StaticResource) Content() []byte {
	return r.StaticContent
}

func ReadBytes(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

func LoadResourceFromPath(path string) (Resource, error) {
	bytes, err := ReadBytes(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
	return NewStaticResource(name, bytes), nil
}

func NewStaticResource(name string, content []byte) *StaticResource {
	return &StaticResource{
		StaticName:    name,
		StaticContent: content,
	}
}
