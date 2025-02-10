package markdown

import (
	"fmt"
	"testing"
)

func TestBuildTableItem(t *testing.T) {
	fmt.Println(buildTableItem([]string{"a", "b", "c"}))
}

func TestToMarkdownTable(t *testing.T) {
	fmt.Println(ToMarkdownTable([]map[string]string{
		{
			"name":      "nginx-saf3e",
			"namespace": "default",
			"phase":     "Pending",
		},
	}))
}
