# coname

A Go library to normalize company names by stripping legal entity designators.

Useful for company name matching, deduplication, search indexing, and data cleaning. Handles international company suffixes like Ltd, Inc, GmbH, S.A., B.V., Pty Ltd, and 500+ more from 20+ countries.

## The Problem

Company names come in many forms:

- "Apple Inc." vs "Apple" vs "Apple, Inc"
- "Volkswagen AG" vs "Volkswagen"
- "The Coca-Cola Company" vs "Coca-Cola"

This makes matching and deduplication difficult. `coname` strips legal suffixes and prefixes to extract the core company name.

## Usage

```go
import "github.com/muzzii255/coname"

coname.NormalizeWords("Acme Ltd.")            // "Acme"
coname.NormalizeWords("Müller GmbH & Co. KG") // "Müller"
coname.NormalizeWords("THE COCA-COLA CO")     // "COCA-COLA"
coname.NormalizeWords("Toyota 株式会社")        // "Toyota"
coname.NormalizeWords("SAP SE")               // "SAP"
```

## Features

- **500+ legal entity designators** from 20+ countries
- **Multi-language support**: English, German, French, Spanish, Portuguese, Italian, Dutch, Nordic (Swedish, Norwegian, Danish, Finnish), Polish, Hungarian, Czech, Russian, Japanese, Korean, Chinese, Malaysian, Arabic, and more
- **Case and period insensitive**: `L.L.C.` = `LLC` = `llc`
- **Stacked suffix handling**: `GmbH & Co. KG` fully stripped
- **Prefix stripping**: removes "The" from company names
- **Word boundary protection**: `Tabasco` keeps its `co`, `Costco` stays intact
- **Full-name collision protection**: short names like `SAP` won't be stripped even if they match abbreviations
- **Zero runtime dependencies**: all designators baked into the binary
- **Fast**: ~1.7μs per normalization, 1 allocation

## Supported Designators

Includes but not limited to:

| Region             | Examples                                                 |
| ------------------ | -------------------------------------------------------- |
| English            | Ltd, Inc, Corp, LLC, LLP, PLC, Co, Limited, Incorporated |
| German             | GmbH, AG, KG, OHG, GbR, e.V., KGaA                       |
| French             | SA, SARL, SAS, SNC, EURL                                 |
| Spanish/Portuguese | S.A., S.L., S.R.L., Ltda, S.A. de C.V.                   |
| Italian            | S.p.A., S.r.l.                                           |
| Dutch/Belgian      | N.V., B.V., V.O.F.                                       |
| Nordic             | AB, AS, A/S, Oy, Oyj, ApS                                |
| Eastern European   | Sp. z o.o., s.r.o., Kft., OOO, ЗАО                       |
| Asian              | 株式会社, 有限公司, 주식회사, Sdn. Bhd., Pte. Ltd.       |

## Install

```bash
go get github.com/muzzii255/coname
```

## Benchmark

```text
goos: linux
goarch: amd64
pkg: github.com/muzzii255/coname
cpu: AMD Ryzen 7 9700X 8-Core Processor
BenchmarkNormalizeWords-16    	  686209	      1741 ns/op	     208 B/op	       1 allocs/op
BenchmarkNormalizeWords-16    	  657469	      1712 ns/op	     208 B/op	       1 allocs/op
BenchmarkNormalizeWords-16    	  707410	      1760 ns/op	     208 B/op	       1 allocs/op
BenchmarkNormalizeWords-16    	  700401	      1753 ns/op	     208 B/op	       1 allocs/op
```

## Updating Designators

Designators are generated from `company_designator.yml` at build time:

```bash
go generate ./...
```

## Credits

Designator list based on [ProfoundNetworks/company_designator](https://github.com/ProfoundNetworks/company_designator).

## License

MIT
