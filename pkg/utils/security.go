package utils

import (
	"fmt"
	"strings"
)

func EscapeIdentifier(name string) string {
	// Escapa backticks internos reemplazándolos por doble backtick
	name = strings.ReplaceAll(name, "`", "``")
	return fmt.Sprintf("`%s`", name)
}
