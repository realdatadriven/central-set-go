package main

import (
	"context"
	"os"

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
