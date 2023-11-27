package main

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func main() {
    // Configurar o SDK Fabric
    configProvider := config.FromFile("./config.yaml")
    sdk, err := fabsdk.New(configProvider)
    if err != nil {
        fmt.Printf("Erro ao criar o SDK Fabric: %s\n", err)
        return
    }
    defer sdk.Close()

    // Obter o cliente do canal
    clientChannelContext := sdk.ChannelContext("channel1", fabsdk.WithUser("iarlenaquiles"))
    client, err := channel.New(clientChannelContext)
    if err != nil {
        fmt.Printf("Erro ao criar o cliente do canal: %s\n", err)
        return
    }

    // Interagir com o chaincode
    response, err := client.Execute(channel.Request{
        ChaincodeID: "leilao",
        Fcn:          "iniciarLeilao",
        Args:         [][]byte{[]byte("CUSIP123"), []byte("100000")},
    })
    if err != nil {
        fmt.Printf("Erro ao enviar transação para o chaincode: %s\n", err)
        return
    }

    fmt.Printf("Resposta do chaincode: %s\n", response.Payload)
}