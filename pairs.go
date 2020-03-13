package main

//	"8873":   "Dow Jones 30",
//	"8839":   "S&P 500",
//	"8874":   "Nasdaq 100",
//	"8864":   "Russell 2000",
//	"8884":   "VIX",
//	"8826":   "DAX",
//	"15288":  "TecDAX",
//	"8867":   "Euro Stoxx 50",
//	"8838":   "FTSE 100",
//	"8828":   "IBEX 35",
//	"8853":   "CAC 40",
//	"8858":   "FTSE MIB",
//	"8878":   "AEX",
//	"8837":   "SMI",
//	"8863":   "WIG20",
//	"941612": "Ibovespa",
//	"8840":   "South Africa 40",
//	"104396": "RTS",
//	"8824":   "S&P/ASX 200",
//	"8987":   "KOSPI",
//	"8859":   "Nikkei",
//	"8897":   "Singapore MSCI",
//	"8985":   "Nifty 50",
//	"8984":   "Hang Seng",
//	"8982":   "China H-Shares",
//	"44486":  "China A50",

// "8827": "Dollar Index",

//	"8827": "Dollar Index",
//	"8830": "Gold",
//	"8833": "Crude Oil",
//
//	"3": "USD/JPY",
//
//	"8907": "US 30Y Futures",
//	"8880": "US 10Y Futures",
//	"8905": "US 5Y Futures",
//	"8906": "US 2Y Futures",

var defaultPairsSlice = []string{
	"8827",  // DXY
	"3",     // USD/JPY
	"0000",  // New Line '\n'
	"8907",  // US 30Y Futures
	"8880",  // US 10Y Futures
	"8905",  // US 5Y Futures
	"8906",  // US 2Y Futures
	"0000",  // New Line
	"8884",  // VIX
	"8830",  // Gold
	"8833",  // Crude Oil
	"0000",  // New Line
	"8873",  // Dow Jones 30
	"8839",  // S&P 500
	"8874",  // Nasdaq 100
	"8864",  // Russell 2000
	"0000",  // New Line
	"8838",  // FTSE 100
	"8826",  // DAX
	"8867",  // Euro Stoxx 50
	"0000",  // New Line
	"8859",  // Nikkei
	"8987",  // KOSPI
	"8984",  // Hang Seng
	"8982",  // China H-Shares
	"44486", // China A50
	"8824",  // S&P/ASX 200
	"8897",  // Singapore MSCI
	"8985",  // Nifty 50
}
