.PHONY: deploy kill clean purge

deploy:
	@./setup/deploy.sh

kill:
	@pkill spot; true

clean:
	rm -rf ~/.spot/
	rm -f $(GOPATH)/bin/spot

purge: kill clean
