package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// LeilaoChaincode implements Hyperledger Fabric's Chaincode interface
type LeilaoChaincode struct {
}

// Bond represents a bond with a unique identifier (CUSIP) and an associated amount
type Bond struct {
	CUSIP  string `json:"cusip"`
	Quantia int    `json:"quantia"`
}

// Offer represents an offer from a bidder with a price
type Oferta struct {
	Licitante string `json:"licitante"`
	Preco     int    `json:"preco"`
}

// Auction status maintains the title/bond and current bids
type EstadoLeilao struct {
	Titulo   Bond    `json:"titulo"`
	Ofertas  []Oferta `json:"ofertas"`
	Vendido   bool    `json:"vendido"`
	Vencedor  string  `json:"vencedor"`
}

// Main function that initializes the chaincode
func main() {
	if err := shim.Start(new(LeilaoChaincode)); err != nil {
		fmt.Printf("Erro ao iniciar o chaincode: %s", err)
	}
}

// Init is called when the chaincode is instantiated
func (t *LeilaoChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called when a transaction is proposed
func (t *LeilaoChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "iniciarLeilao":
		return t.iniciarLeilao(stub, args)
	case "fazerOferta":
		return t.fazerOferta(stub, args)
	case "encerrarLeilao":
		return t.encerrarLeilao(stub, args)
	case "consultarLeilao":
		return t.consultarLeilao(stub, args)
	default:
		return shim.Error("Função inválida")
	}
}

// Start the auction with a specific security/bond
func (t *LeilaoChaincode) iniciarLeilao(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Número incorreto de argumentos. Esperado: 2")
	}

	cusip := args[0]
	quantia, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Erro ao converter quantia para inteiro")
	}

	titulo := Bond{
		CUSIP:  cusip,
		Quantia: quantia,
	}

	estadoLeilao := EstadoLeilao{
		Titulo:   titulo,
		Ofertas:  make([]Oferta, 0),
		Vendido:   false,
		Vencedor:  "",
	}

	estadoJSON, err := json.Marshal(estadoLeilao)
	if err != nil {
		return shim.Error("Erro ao converter estado do leilão para JSON")
	}

	err = stub.PutState("leilao", estadoJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao salvar estado do leilão: %s", err))
	}

	return shim.Success(nil)
}

// Make a bid on an ongoing auction
func (t *LeilaoChaincode) fazerOferta(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Número incorreto de argumentos. Esperado: 2")
	}

	licitante := args[0]
	preco, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Erro ao converter preço para inteiro")
	}

	estadoJSON, err := stub.GetState("leilao")
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao obter estado do leilão: %s", err))
	}

	var estadoLeilao EstadoLeilao
	err = json.Unmarshal(estadoJSON, &estadoLeilao)
	if err != nil {
		return shim.Error("Erro ao converter estado do leilão para struct")
	}

	if estadoLeilao.Vendido {
		return shim.Error("O leilão já foi encerrado")
	}

	oferta := Oferta{
		Licitante: licitante,
		Preco:     preco,
	}

	estadoLeilao.Ofertas = append(estadoLeilao.Ofertas, oferta)

	estadoJSON, err = json.Marshal(estadoLeilao)
	if err != nil {
		return shim.Error("Erro ao converter estado do leilão para JSON")
	}

	err = stub.PutState("leilao", estadoJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao salvar estado do leilão: %s", err))
	}

	return shim.Success(nil)
}

// Closes the auction and awards the title/bond to the bidder with the highest price
func (t *LeilaoChaincode) encerrarLeilao(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	estadoJSON, err := stub.GetState("leilao")
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao obter estado do leilão: %s", err))
	}

	var estadoLeilao EstadoLeilao
	err = json.Unmarshal(estadoJSON, &estadoLeilao)
	if err != nil {
		return shim.Error("Erro ao converter estado do leilão para struct")
	}

	if estadoLeilao.Vendido {
		return shim.Error("O leilão já foi encerrado")
	}

	if len(estadoLeilao.Ofertas) == 0 {
		return shim.Error("Não há ofertas no leilão")
	}

	// Find the offer with the highest price
	precoMaximo := estadoLeilao.Ofertas[0].Preco
	vencedor := estadoLeilao.Ofertas[0].Licitante

	for _, oferta := range estadoLeilao.Ofertas {
		if oferta.Preco > precoMaximo {
			precoMaximo = oferta.Preco
			vencedor = oferta.Licitante
		}
	}

	// Assign the title/bond to the winner
	estadoLeilao.Vendido = true
	estadoLeilao.Vencedor = vencedor

	estadoJSON, err = json.Marshal(estadoLeilao)
	if err != nil {
		return shim.Error("Erro ao converter estado do leilão para JSON")
	}

	err = stub.PutState("leilao", estadoJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao salvar estado do leilão: %s", err))
	}

	return shim.Success(nil)
}

// Check the current status of the auction
func (t *LeilaoChaincode) consultarLeilao(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	estadoJSON, err := stub.GetState("leilao")
	if err != nil {
		return shim.Error(fmt.Sprintf("Erro ao obter estado do leilão: %s", err))
	}

	if estadoJSON == nil {
		return shim.Error("Leilão não encontrado")
	}

	return shim.Success(estadoJSON)
}
