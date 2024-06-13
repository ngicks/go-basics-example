package pkg2

import (
	"fmt"

	"github.com/ngicks/go-basics-example/mod-organization/pkg1"
)

func SayDouble() string {
	return fmt.Sprintf("%q%q", pkg1.Foo, pkg1.Foo)
}
