package model

import (
	"github.com/MarcGrol/astTools/model/annotation/eventAnno"
	"github.com/MarcGrol/astTools/model/annotation/restAnno"
)

func init() {
	eventAnno.Register()
	restAnno.Register()
}
