package types

import (
	"ballerina-lang-go/compiler/bir/model"
)

type NamedNode interface {
	GetName() model.Name
}
