
clean:
	sudo rm -rf ./state

full-clean:
	sudo rm -rf \
		./crypto-config \
		./state \
		./orderer \
		./channels
	mkdir -p ./orderer ./channels

docker-fix-osx:
	brew install docker-machine
	docker-machine create -d virtualbox default
	@echo "set the CA path in your ENV"
	@echo "https://stackoverflow.com/questions/33169122/docker-error-dial-unix-var-run-docker-sock-no-such-file-or-directory"
