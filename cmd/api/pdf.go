package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// USING CHROMIUN
func (app *application) GenPDFFromHTML(html, output_path string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var pdf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+html),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err
	}
	err = os.WriteFile(output_path, pdf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LatexEscape(v any) string {
	s := fmt.Sprint(v)

	replacer := strings.NewReplacer(
		"\\", "\\textbackslash{}",
		"&", "\\&",
		"%", "\\%",
		"$", "\\$",
		"#", "\\#",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
	)

	return replacer.Replace(s)
}

// sudo apt install texlive-latex-base
func (app *application) GenPDFFromLatex(latex, output_path string) error {
	latex = LatexEscape(latex)
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	temptex, err := os.CreateTemp("", "*.tex")
	if err != nil {
		return err
	}
	defer os.Remove(temptex.Name())
	defer temptex.Close()
	_, err = temptex.WriteString(latex)
	if err != nil {
		return err
	}
	temptex.Close()
	cmd := exec.Command("pdflatex", temptex.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}
	// Move generated PDF to output path
	tempPDF := strings.TrimSuffix(temptex.Name(), ".tex") + ".pdf"
	err = os.Rename(tempPDF, output_path)
	if err != nil {
		return err
	}
	return nil
}
