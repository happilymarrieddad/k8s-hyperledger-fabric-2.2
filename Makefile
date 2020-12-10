
clean:
	sudo rm -rf ./state

full-clean:
	sudo rm -rf \
		./crypto-config \
		./state \
		./orderer \
		./channels
	mkdir -p ./orderer ./channels
