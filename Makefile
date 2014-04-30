default: assets

assets: 
	@go-bindata -pkg=sieve -prefix=assets/ assets

.PHONY: assets
