JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 2 --string-style s --comment-style s

all: fmt generated lint

fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

generated: mixin.libsonnet config.libsonnet $(wildcard dashboards/*.jsonnet) lib/*.libsonnet
	@mkdir -p generated
	jsonnet -J . -m generated dashboards.jsonnet
	@jq '{"__inputs": .__inputs, "__requires": .__requires} + del(.__inputs, .__requires)' \
		generated/bigip.json > generated/bigip.json.tmp && mv generated/bigip.json.tmp generated/bigip.json

lint: generated
	@find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" - || exit 1; \
		done
	@echo "Lint OK"

clean:
	rm -rf generated

.PHONY: all fmt generated lint clean
