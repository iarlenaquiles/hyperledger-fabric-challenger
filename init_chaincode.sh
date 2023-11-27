# Package o chaincode
peer lifecycle chaincode package leilao.tar.gz --path localhost:7050 --lang golang --label leilao_1.0

# Instale o chaincode em todos os peers
peer lifecycle chaincode install leilao.tar.gz

# Aprovar o chaincode
peer lifecycle chaincode queryinstalled

PACKAGE_ID=<package_id>  # Substitua <package_id> pelo ID real do pacote retornado no comando anterior
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.test.com --channelID channel1 --name leilao --version 1.0 --package-id $PACKAGE_ID --sequence 1

# Verificar a aprovação
peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name leilao --version 1.0 --sequence 1 --output json

# Confirmar o chaincode
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.test.com --channelID channel1 --name leilao --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA
