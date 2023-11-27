# Geração de certificados
cryptogen generate --config=./crypto-config.yaml

# Criação do bloco de gênese
configtxgen -profile TwoOrgsOrdererGenesis -channelID sys-channel -outputBlock ./channel-artifacts/genesis.block

# Criação do canal
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID channel1

# Atualizar âncoras dos peers
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID channel1 -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID channel1 -asOrg Org2MSP