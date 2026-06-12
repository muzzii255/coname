package coname

import (
	"testing"
)

var tests = []struct {
	input string
	want  string
}{
	// English basics
	{"Acme Ltd.", "Acme"},
	{"Acme Limited", "Acme"},
	{"Acme Inc.", "Acme"},
	{"Acme Inc", "Acme"},
	{"Acme Corp", "Acme"},
	{"Acme Corp.", "Acme"},
	{"Acme Corporation", "Acme"},
	{"Acme LLC", "Acme"},
	{"Acme L.L.C.", "Acme"},
	{"Acme LLP", "Acme"},
	{"Acme L.L.P.", "Acme"},
	{"Acme Co.", "Acme"},
	{"Acme Co", "Acme"},
	{"Acme Holdings", "Acme"},
	{"Acme Holdings Ltd", "Acme"},
	{"Acme Holding", "Acme"},
	{"Acme PLC", "Acme"},
	{"Acme P.L.C.", "Acme"},
	{"Acme Partners", "Acme"},
	{"Acme Foundation", "Acme"},
	{"Acme Cooperative", "Acme"},
	{"Acme Limited Partnership", "Acme"},
	{"Acme Limited Liability Company", "Acme"},
	{"Acme Private Limited Company", "Acme"},
	{"Acme Public Limited Company", "Acme"},

	// Stacked suffixes
	{"Foo GmbH & Co. KG", "Foo"},
	{"Bar GmbH & Co KG", "Bar"},
	{"Baz GmbH und Co. KG", "Baz"},
	{"Qux AG & Co. KGaA", "Qux"},

	// German
	{"Müller GmbH", "Müller"},
	{"Siemens AG", "Siemens"},
	{"Volkswagen AG", "Volkswagen"},
	{"Deutsche Bahn AG", "Deutsche Bahn"},
	{"Bosch GmbH", "Bosch"},
	{"Henkel KGaA", "Henkel"},
	{"Schmidt OHG", "Schmidt"},
	{"Weber e.K.", "Weber"},
	{"Fischer GbR", "Fischer"},
	{"Bauer e.V.", "Bauer"},

	// French
	{"Total SA", "Total"},
	{"Renault S.A.", "Renault"},
	{"BNP Paribas SARL", "BNP Paribas"},
	{"Carrefour S.A.R.L.", "Carrefour"},
	{"Michelin SAS", "Michelin"},
	{"Orange SNC", "Orange"},

	// Spanish / Portuguese
	{"Telefónica S.A.", "Telefónica"},
	{"CEMEX S.A. de C.V.", "CEMEX"},
	{"Santander S.L.", "Santander"},
	{"Petrobras Ltda.", "Petrobras"},
	{"Embraer S.R.L.", "Embraer"},

	// Italian
	{"Fiat S.p.A.", "Fiat"},
	{"Ferrari SpA", "Ferrari"},
	{"Gucci S.r.l.", "Gucci"},
	{"Prada Srl", "Prada"},

	// Dutch / Belgian
	{"Philips N.V.", "Philips"},
	{"Shell B.V.", "Shell"},
	{"Heineken NV", "Heineken"},
	{"ASML BV", "ASML"},

	// Nordic
	{"Volvo AB", "Volvo"},
	{"Ericsson AB", "Ericsson"},
	{"Statoil ASA", "Statoil"},
	{"Nordea Oyj", "Nordea"},
	{"Nokia Oy", "Nokia"},
	{"Maersk A/S", "Maersk"},
	{"Carlsberg ApS", "Carlsberg"},

	// Polish
	{"Orlen Sp. z o.o.", "Orlen"},
	{"PKO S.A.", "PKO"},

	// Hungarian
	{"MOL Kft.", "MOL"},
	{"OTP Rt.", "OTP"},
	{"Richter Nyrt.", "Richter"},

	// Czech
	{"Škoda s.r.o.", "Škoda"},

	// Russian (Cyrillic)
	{"Газпром ОАО", "Газпром"},
	{"Сбербанк ООО", "Сбербанк"},
	{"Лукойл ЗАО", "Лукойл"},

	// Japanese
	{"トヨタ 株式会社", "トヨタ"},
	{"Sony K.K.", "Sony"},

	{"Toyota Motor Corporation", "Toyota Motor"},
	{"Sony Group Corporation", "Sony Group"},

	// Korean
	{"삼성 주식회사", "삼성"},

	// Chinese
	{"阿里巴巴 有限公司", "阿里巴巴"},
	{"腾讯 股份有限公司", "腾讯"},

	// Malaysian
	{"Petronas Sdn. Bhd.", "Petronas"},
	{"Maybank Bhd.", "Maybank"},

	// Singapore / India
	{"DBS Pte. Ltd.", "DBS"},
	{"Tata Pvt. Ltd.", "Tata"},
	{"Infosys Pvt Ltd", "Infosys"},

	// Prefix stripping
	{"The Coca-Cola Co", "Coca-Cola"},
	{"THE COCA-COLA COMPANY", "COCA-COLA"},
	{"The Walt Disney Company", "Walt Disney"},
	{"The Boeing Company", "Boeing"},
	{"The Procter & Gamble Company", "Procter & Gamble"},

	// Word boundary protection - should NOT strip
	{"Tabasco", "Tabasco"},
	{"Costco", "Costco"},
	{"Nabisco", "Nabisco"},
	{"Sysco", "Sysco"},
	{"Crisco", "Crisco"},
	{"Monaco", "Monaco"},
	{"Texaco", "Texaco"},
	{"Alco", "Alco"},
	{"Bronco", "Bronco"},
	{"Franco", "Franco"},
	{"Saab", "Saab"},
	{"Ikea", "Ikea"},

	// Period insensitivity
	{"Acme L.L.C.", "Acme"},
	{"Acme LLC", "Acme"},
	{"Acme Llc", "Acme"},
	{"Test P.L.L.C.", "Test"},
	{"Test PLLC", "Test"},

	// Case insensitivity
	{"ACME LTD", "ACME"},
	{"acme ltd", "acme"},
	{"Acme LTD", "Acme"},
	{"Acme ltd", "Acme"},
	{"ACME gmbh", "ACME"},

	// Edge cases - empty / whitespace
	{"", ""},
	{"   ", ""},
	{"\t\n", ""},

	// Edge cases - only suffix (preserved, not stripped to empty)
	{"Ltd", "Ltd"},
	{"GmbH", "GmbH"},
	{"Inc.", "Inc."},
	{"SAP SE", "SAP"}, // SAP preserved, not stripped due to S.A.P. collision
	{"The Company", "The"}, // prefix only strips if followed by content

	// Edge cases - punctuation
	{"Acme, Inc.", "Acme"},
	{"Acme - Ltd.", "Acme"},
	{"Acme & Co.", "Acme"},
	{"Acme & Co., Ltd.", "Acme"},
	{"Smith & Wesson Corp", "Smith & Wesson"},

	// Whitespace handling
	{"Coca  Cola Co", "Coca Cola"},
	{"The  Walt  Disney  Company", "Walt Disney"},
	{"Acme\tWidgets Ltd", "Acme Widgets"},
	{"Acme\nWidgets\tLtd", "Acme Widgets"},
	{"  Acme Ltd  ", "Acme"},

	// Mixed scripts
	{"München GmbH", "München"},
	{"Société Générale SA", "Société Générale"},
	{"Citroën S.A.", "Citroën"},

	// Numbers in names
	{"3M Company", "3M"},
	{"7-Eleven Inc", "7-Eleven"},
	{"20th Century Fox Corporation", "20th Century Fox"},

	// Complex real-world examples
	{"Berkshire Hathaway Inc.", "Berkshire Hathaway"},
	{"Johnson & Johnson", "Johnson & Johnson"},
	{"Procter & Gamble Co.", "Procter & Gamble"},
	{"JPMorgan Chase & Co.", "JPMorgan Chase"},
	{"Goldman Sachs Group, Inc.", "Goldman Sachs Group"},
	{"BASF SE", "BASF"},
	{"Daimler AG", "Daimler"},
	// Note: "SAP SE" → "" because "SAP" collides with "S.A.P." (Sociedad de Ahorro y Prestamo)
}

func TestNormalizeWords(t *testing.T) {
	for _, tt := range tests {
		got := NormalizeWords(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeWords(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func BenchmarkNormalizeWords(b *testing.B) {
	for b.Loop() {
		NormalizeWords("The Coca-Cola Company Ltd.")
	}
}
