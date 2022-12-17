//go:build nobuild
// +build nobuild

package third_party

import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/matryer/moq"
)
