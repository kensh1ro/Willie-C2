package clr

import (
	"fmt"

	dotnet "github.com/Ne0nd0g/go-clr"
)

func ExecuteAssembly(version string, assembly []byte, params []string) string {
	ret, _ := dotnet.ExecuteByteArray(version, assembly, params)
	return fmt.Sprintf("[+] Return exit code: %d", ret)
}
