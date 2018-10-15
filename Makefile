
.PHONY: snapshot
snapshot:
	goreleaser --snapshot --rm-dist
