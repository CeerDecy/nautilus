package markdown

import "strings"

func ToMarkdownTable(data []map[string]string) string {
	if len(data) == 0 {
		return ""
	}
	builder := strings.Builder{}
	title := make([]string, 0, len(data[0]))
	split := make([]string, 0, len(data[0]))
	for index, row := range data {
		item := make([]string, 0, len(row))
		if index == 0 {
			for k := range row {
				title = append(title, k)
				split = append(split, strings.Repeat("-", len(k)))
			}
			builder.WriteString(buildTableItem(title))
			builder.WriteString(buildTableItem(split))
		}
		for _, v := range title {
			item = append(item, row[v])
		}
		builder.WriteString(buildTableItem(item))
	}
	return builder.String()
}

func buildTableItem(data []string) string {
	return "| " + strings.Join(data, " | ") + " |\n"
}
