package types

import "ballerina-lang-go/compiler/common"

type NamedNode interface {
	GetName() common.Name
}
