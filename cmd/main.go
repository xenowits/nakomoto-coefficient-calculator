package main

import (
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/binance"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/cosmos"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/mina"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/osmosis"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/polygon"
)

func main() {
	binance.Binance()
	cosmos.Cosmos()
	osmosis.Osmosis()
	polygon.Polygon()
	mina.Mina()
}
