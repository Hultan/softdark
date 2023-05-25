package softdark

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type SoftBuilder struct {
	builder *gtk.Builder
}

func newSoftBuilder(gladeFileName string) *SoftBuilder {
	builder := &SoftBuilder{}
	gladePath, err := getResourcePath(gladeFileName)
	if err != nil {
		panic(err)
	}

	b, err := gtk.BuilderNewFromFile(gladePath)
	if err != nil {
		panic(err)
	}
	builder.builder = b

	return builder
}

func (s *SoftBuilder) getObject(name string) glib.IObject {
	obj, err := s.builder.GetObject(name)
	if err != nil {
		panic(err)
	}

	return obj
}
