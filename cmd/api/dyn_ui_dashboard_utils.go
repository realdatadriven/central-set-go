package main

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type SQLBlock struct {
	Language string
	Code     string
}

func ExtractSQLBlocks(markdown string) ([]SQLBlock, string) {
	source := text.NewReader([]byte(markdown))
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	doc := md.Parser().Parse(source)
	var blocks []SQLBlock
	// Collect the actual fenced blocks from the source.
	type removal struct {
		start int
		end   int
	}
	var removals []removal
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		block, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}
		lang := string(block.Language([]byte(markdown)))
		// Only SQL blocks.
		if !strings.EqualFold(lang, "sql") {
			return ast.WalkContinue, nil
		}
		var code strings.Builder
		lines := block.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write([]byte(markdown)[line.Start:line.Stop])
		}
		blocks = append(blocks, SQLBlock{
			Language: lang,
			Code:     code.String(),
		})
		// Find the fenced block in the original source.
		//
		// We search backwards from the code content for the opening
		// fence and forwards for the closing fence.
		if lines.Len() == 0 {
			return ast.WalkContinue, nil
		}
		first := lines.At(0)
		last := lines.At(lines.Len() - 1)
		start := first.Start
		for start > 0 && []byte(markdown)[start-1] != '\n' {
			start--
		}
		// Include the opening fence line.
		if start > 0 {
			start--
			for start > 0 && []byte(markdown)[start-1] != '\n' {
				start--
			}
		}
		end := last.Stop
		// Include closing fence.
		for end < len([]byte(markdown)) && []byte(markdown)[end] != '\n' {
			end++
		}
		if end < len([]byte(markdown)) {
			end++
		}
		removals = append(removals, removal{
			start: start,
			end:   end,
		})
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, markdown
	}
	// Remove from the end to the beginning.
	result := []byte(markdown)
	for i := len(removals) - 1; i >= 0; i-- {
		r := removals[i]
		result = append(result[:r.start], result[r.end:]...)
	}
	return blocks, string(result)
}

func extractQueryNameFromSQL(sqlContent string) (string, string) {
	// Define a regex to match the first line with a comment containing the query name
	re := regexp.MustCompile(`(?m)^--\s*(\w+)\s*$`)
	// Find the query name using the regex
	matches := re.FindStringSubmatch(sqlContent)
	queryName := ""
	if len(matches) > 1 {
		queryName = matches[1] // Capture the query name
	}
	// Remove the matched line from the content
	updatedContent := re.ReplaceAllString(sqlContent, "")
	// Remove any leading newline or whitespace from the updated content
	updatedContent = strings.TrimLeft(updatedContent, "\n")
	return queryName, updatedContent
}

func ExtractSQLBlocksV2(markdown string) ([]Dict, string) {
	source := text.NewReader([]byte(markdown))
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	doc := md.Parser().Parse(source)
	var blocks []Dict
	// Collect the actual fenced blocks from the source.
	type removal struct {
		start int
		end   int
	}
	var removals []removal
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		block, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}
		lang := string(block.Language([]byte(markdown)))
		info := string(block.Info.Segment.Value([]byte(markdown)))
		key := strings.TrimSpace(strings.TrimPrefix(info, "sql"))
		if key == "" {
			// If no key in the info, try to extract from the first comment line
			key, _ = extractQueryNameFromSQL(markdown)
			// fmt.Println("NOT A NAMED QUERY, FIND -- name instead", key, contentFinal)
		}
		// fmt.Println(info, key)
		// Only SQL blocks.
		if !strings.EqualFold(lang, "sql") {
			return ast.WalkContinue, nil
		}
		var code strings.Builder
		lines := block.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write([]byte(markdown)[line.Start:line.Stop])
		}
		/*blocks = append(blocks, SQLBlock{
			Language: lang,
			Code:     code.String(),
		})*/
		blocks = append(blocks, Dict{
			"key":  key,
			"lang": lang,
			"code": code.String(),
		})
		// Find the fenced block in the original source.
		//
		// We search backwards from the code content for the opening
		// fence and forwards for the closing fence.
		if lines.Len() == 0 {
			return ast.WalkContinue, nil
		}
		first := lines.At(0)
		last := lines.At(lines.Len() - 1)
		start := first.Start
		for start > 0 && []byte(markdown)[start-1] != '\n' {
			start--
		}
		// Include the opening fence line.
		if start > 0 {
			start--
			for start > 0 && []byte(markdown)[start-1] != '\n' {
				start--
			}
		}
		end := last.Stop
		// Include closing fence.
		for end < len([]byte(markdown)) && []byte(markdown)[end] != '\n' {
			end++
		}
		if end < len([]byte(markdown)) {
			end++
		}
		removals = append(removals, removal{
			start: start,
			end:   end,
		})
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, markdown
	}
	// Remove from the end to the beginning.
	result := []byte(markdown)
	for i := len(removals) - 1; i >= 0; i-- {
		r := removals[i]
		result = append(result[:r.start], result[r.end:]...)
	}
	return blocks, string(result)
}

// get the

/*func main() {
	md := `# Test

Some text.

` + "```sql\n" + `CREATE SECRET api_auth (
  TYPE http_bearer,
  TOKEN '<JWT_TOKEN>',
  SCOPE 'http://localhost:4444/'
);
` + "```\n\n" + `More text.

` + "```sql\n" + `ATTACH IF NOT EXISTS 'http://localhost:4444/odata/ADMIN' AS admin (TYPE odata);
` + "```\n"

	blocks, remaining := ExtractSQLBlocks(md)

	for _, block := range blocks {
		fmt.Printf("LANGUAGE: %s\n", block.Language)
		fmt.Printf("CODE:\n%s\n", block.Code)
	}

	fmt.Println("REMAINING:")
	fmt.Println(remaining)
}*/
