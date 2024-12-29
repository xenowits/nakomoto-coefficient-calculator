package chains

import (
	"fmt"
	"log"
)

// Chain contains details of a particular Chain.
type Chain struct {
	PrevNCVal int
	CurrNCVal int
}

// Token represents the name of token for a blockchain.
// For example, ATOM for cosmos.
// It is used to identify a particular Chain.
type Token string

// ChainState contains complete NC information for all supported chains.
type ChainState map[Token]Chain

// Append new chains in alphabetical order only.
const (
	ADA   Token = "ADA"
	ALGO  Token = "ALGO"
	APT   Token = "APT"
	ATOM  Token = "ATOM"
	AVAX  Token = "AVAX"
	BLD   Token = "BLD"
	BNB   Token = "BNB"
	DOT   Token = "DOT"
	EGLD  Token = "EGLD"
	GRT   Token = "GRT"
	HBAR  Token = "HBAR"
	JUNO  Token = "JUNO"
	MATIC Token = "MATIC"
	MINA  Token = "MINA"
	NEAR  Token = "NEAR"
	OSMO  Token = "OSMO"
	PLS   Token = "PLS"
	REGEN Token = "REGEN"
	RONIN Token = "RONIN"
	RUNE  Token = "RUNE"
	SEI   Token = "SEI"
	SOL   Token = "SOL"
	STARS Token = "STARS"
	SUI   Token = "SUI"
	TIA   Token = "TIA"
)

// ChainName returns the name of the chain given the token name.
func (t Token) ChainName() string {
	switch t {
	case ADA:
		return "Cardano"
	case ALGO:
		return "Algo"
	case APT:
		return "Aptos"
	case ATOM:
		return "Cosmos"
	case AVAX:
		return "Avalanche"
	case BLD:
		return "Agoric"
	case BNB:
		return "BNB Smart Chain"
	case DOT:
		return "Polkadot"
	case EGLD:
		return "MultiversX"
	case GRT:
		return "Graph Protocol"
	case HBAR:
		return "Hedera"
	case JUNO:
		return "Juno"
	case MATIC:
		return "Polygon"
	case MINA:
		return "Mina Protocol"
	case NEAR:
		return "Near Protocol"
	case OSMO:
		return "Osmosis"
	case PLS:
		return "Pulsechain"
	case REGEN:
		return "Regen Network"
	case RONIN:
		return "Ronin Chain"
	case RUNE:
		return "Thorchain"
	case SEI:
		return "Sei"
	case SOL:
		return "Solana"
	case STARS:
		return "Stargaze"
	case SUI:
		return "Sui Protocol"
	case TIA:
		return "Celestia"
	default:
		return "Unknown"
	}
}


var Tokens = []Token{ADA, ALGO, APT, ATOM, AVAX, BLD, BNB, DOT, EGLD, GRT, HBAR, JUNO, MATIC, MINA, NEAR, OSMO, PLS, REGEN, RONIN, RUNE, SEI, SOL, STARS, SUI, TIA}


// NewState returns a new fresh state.
func NewState() ChainState {
	state := make(ChainState)

	return RefreshChainState(state)
}

func RefreshChainState(prevState ChainState) ChainState {
	newState := make(ChainState)
	for _, token := range Tokens {
		currVal, err := newValues(token)
		if err != nil {
			log.Println("Failed to update chain info:", token, err)
			continue
		}

		newState[token] = Chain{
			PrevNCVal: prevState[token].CurrNCVal,
			CurrNCVal: currVal,
		}
	}

	return newState
}

func newValues(token Token) (int, error) {
	var (
		currVal int
		err     error
	)

	log.Printf("Calculating Nakamoto coefficient for %s", token.ChainName())

	switch token {
	case ADA:
		currVal, err = Cardano()
	case ALGO:
		currVal, err = Algorand()
	case APT:
		currVal, err = Aptos()
	case ATOM:
		currVal, err = Cosmos()
	case AVAX:
		currVal, err = Avalanche()
	case BLD:
		currVal, err = Agoric()
	case BNB:
		currVal, err = BSC()
	case DOT:
		currVal, err = Polkadot()
	case EGLD:
		currVal, err = MultiversX()
	case GRT:
		currVal, err = Graph()
	case HBAR:
		currVal, err = Hedera()
	case JUNO:
		currVal, err = Juno()
	case MATIC:
		currVal, err = Polygon()
	case MINA:
		currVal, err = Mina()
	case NEAR:
		currVal, err = Near()
	case OSMO:
		currVal, err = Osmosis()
	case PLS:
		currVal, err = Pulsechain()
	case REGEN:
		currVal, err = Regen()
	case RONIN:
		currVal, err = Ronin()
	case RUNE:
		currVal, err = Thorchain()
	case SEI:
		log.Println("Attempting to calculate Sei Nakamoto coefficient...")
		currVal, err = Sei()
		if err != nil {
			log.Printf("Error calculating Sei Nakamoto coefficient: %v", err)
		}
	case SOL:
		currVal, err = Solana()
	case STARS:
		log.Println("Attempting to calculate Stargaze Nakamoto coefficient...")
		currVal, err = Stargaze()
		if err != nil {
			log.Printf("Error calculating Stargaze Nakamoto coefficient: %v", err)
		}
	case SUI:
		currVal, err = Sui()
	case TIA:
		currVal, err = Celestia()
	default:
		return 0, fmt.Errorf("chain not found: %s", token)
	}

	if err != nil {
		log.Printf("Error in chain %s: %v", token.ChainName(), err)
	} else {
		log.Printf("Successfully calculated Nakamoto coefficient for %s: %d", token.ChainName(), currVal)
	}

	return currVal, err
}
