# fin

[![Go Reference](https://pkg.go.dev/badge/go.seankhliao.com/fin.svg)](https://pkg.go.dev/go.seankhliao.com/fin)
[![License](https://img.shields.io/github/license/seankhliao/fin.svg?style=flat-square)](LICENSE)

Calculates summaries for my personal expenditures.

Example nput looks like:

```cue
out: [ for _cur, _recs in _raw {
	currency: _cur
	holding:  _recs.holdings
	income:   _recs.income
	expense:  _recs.expenses
	months: [ for _rec in _recs.records {
		year:  _rec[0]
		month: _rec[1]
		transactions: [ for _tr in _rec[2] {
			source:      _tr[0]
			destination: _tr[1]
			amount:      _tr[2]
			if len(_tr) > 3 {
				note: _tr[3]
			}
		}]
	}]
}]

let FOO = "foo_x"
let BAR = "bars"
let FIZ = "fizz"

_raw: eur: {
	holdings: [FOO]
	income: [FIZ]
	expenses: [BAR]

	records: [
		[2022, 6, [
		  [FIZ, FOO, 1234567, "income"],
		  [FOO, BAR, 200],
		  [FOO, BAR, 2000, "bars..."],
		]],
	]
}
```
