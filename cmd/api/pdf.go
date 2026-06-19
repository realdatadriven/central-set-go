package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// USING CHROMIUN
// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
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
	err = os.WriteFile(output_path, pdf, 0644)
	if err != nil {
		return err
	}
	// capture pdf
	/*var buf []byte
	err := chromedp.Run(ctx, printToPDF("data:text/html,"+html, &buf))
	if err != nil {
		return err
	}
	err = os.WriteFile(output_path, buf, 0o644)
	if err != nil {
		return err
	}*/
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
	// latex = LatexEscape(latex)
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	dir := filepath.Dir(output_path)
	base := filepath.Base(output_path)
	ext := filepath.Ext(base)
	base_no_ext := strings.Replace(base, ext, "", 1)
	temptex, err := os.CreateTemp("", fmt.Sprintf("%s-*.tex", base_no_ext))
	if err != nil {
		return err
	}
	fmt.Println(temptex.Name(), output_path)
	//defer os.Remove(temptex.Name())
	defer temptex.Close()
	_, err = temptex.WriteString(latex)
	if err != nil {
		return err
	}
	temptex.Close()
	jobname := fmt.Sprintf("-jobname=%s", base_no_ext)
	output_directory := fmt.Sprintf("-output-directory=%s", dir)
	cmd := exec.Command("pdflatex", jobname, output_directory, temptex.Name())
	err = cmd.Run()
	// output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	/*/ Move generated PDF to output path
	tempPDF := strings.TrimSuffix(temptex.Name(), ".tex") + ".pdf"
	err = os.Rename(tempPDF, output_path)
	if err != nil {
		return err
	}*/
	return nil
}
