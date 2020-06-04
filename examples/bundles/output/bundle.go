package testing

import (
	"context"
	"fmt"
	"strconv"

	"github.com/iancoleman/strcase"
)

// we can even import remote code

// make conditional functions. Conversion helpers is a great use case here
func IDToGraphQL(column, tableName string) string {
	return strcase.ToLowerCamel(tableName) + "-" + column
}

// render dynamic funcs
func IDToString(id int) string {
	return strconv.Itoa(id)
}

func partThree_MultipleImports(ctx context.Context) {
	fmt.Println(ctx.Value("some-random-key"))
}

func partTwo() {
	fmt.Println("what's 1 + 1?")
	answer := 2
	fmt.Println("Answer:", answer)
}
