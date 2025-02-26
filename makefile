MAKEFLAGS := --silent --always-make
MAKE_CONC := $(MAKE) -j 128 -f $(lastword $(MAKEFILE_LIST)) clear=$(or $(clear),false)
VERB := $(if $(filter true,$(verb)),-v,)
FAIL := $(if $(filter false,$(fail)),,-failfast)
SHORT := $(if $(filter true,$(short)),-short,)
CLEAR := $(if $(filter false,$(clear)),,-c)
TEST_FLAGS := $(GO_FLAGS) -count=1 $(VERB) $(FAIL) $(SHORT)
TEST := test $(TEST_FLAGS) -timeout=1s -run=$(run)
BENCH := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem -benchtime=128ms
WATCH := watchexec -r $(CLEAR) -d=0 -n

default: test_w

watch:
	$(MAKE_CONC) test_w lint_w

test_w:
	gow -c -v $(TEST)

test:
	go $(TEST)

bench_w:
	gow -c -v $(BENCH)

bench:
	go $(BENCH)

lint_w:
	$(WATCH) -- $(MAKE) lint

lint:
	golangci-lint run
	echo [lint] ok

prep:
	$(MAKE_CONC) test lint

# Example: `make release tag=v0.0.1`.
release: prep
ifeq ($(tag),)
	$(error missing tag)
endif
	git pull --ff-only
	git show-ref --tags --quiet "$(tag)" || git tag "$(tag)"
	git push origin $$(git symbolic-ref --short HEAD) "$(tag)"
