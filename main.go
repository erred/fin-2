package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func main() {
	var fp string
	flag.StringVar(&fp, "f", "fin.cue", "data file")
	flag.Parse()

	b, err := os.ReadFile(fp)
	if err != nil {
		log.Fatalln("read", fp, err)
	}

	cuectx := cuecontext.New()
	val := cuectx.CompileBytes(b)
	val = val.LookupPath(cue.ParsePath("out"))
	var out Output
	err = val.Decode(&out)
	if err != nil {
		log.Fatalln("decode", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', tabwriter.AlignRight)

	for _, currency := range out {
		fmt.Fprintln(tw, "currency:", currency.Currency)
		bag := make(map[string]int)
		fmt.Fprint(tw, "---\t")
		for _, holding := range currency.Holding {
			fmt.Fprint(tw, holding, "\t")
		}
		fmt.Fprintln(tw)

		for _, month := range currency.Months {
			fmt.Fprint(tw, month.Year, "-", month.Month, "\t")
			for _, transaction := range month.Transactions {
				bag[transaction.Source] -= transaction.Amount
				bag[transaction.Destination] += transaction.Amount
			}
			for _, holding := range currency.Holding {
				fmt.Fprintf(tw, "%.2f\t", float64(bag[holding])/100)
			}
			fmt.Fprintln(tw)
		}
	}
	tw.Flush()
}

type (
	Output []Currency

	Currency struct {
		Currency string
		Holding  []string
		Income   []string
		Expense  []string
		Months   []Month
	}

	Month struct {
		Year         int
		Month        int
		Transactions []Transaction
	}

	Transaction struct {
		Source      string
		Destination string
		Note        string
		Amount      int
	}
)
