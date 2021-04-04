.PHONY: dynamicbeat dist clean

define build_dist
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build -o build/dynamicbeat-$(1)-$(2) github.com/scorestack/scorestack/dynamicbeat
	zip -j build/dynamicbeat-$(1)-$(2).zip build/dynamicbeat-$(1)-$(2)
endef

dynamicbeat:
	CGO_ENABLED=0 go build github.com/scorestack/scorestack/dynamicbeat

dist:
	$(call build_dist,linux,amd64)
	$(call build_dist,windows,amd64)
	$(call build_dist,darwin,amd64)
	$(call build_dist,darwin,arm64)

clean:
	rm -rf dynamicbeat
