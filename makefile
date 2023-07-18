MAKEFLAGS  := --silent --always-make
MAKE_PAR   := $(MAKE) -j 128
VERB       := $(if $(filter $(verb), true), -v,)
SHORT      := $(if $(filter $(short), true), -short,)
TEST_FLAGS := -count=1 $(VERB) $(SHORT)
TEST       := test $(TEST_FLAGS) -timeout=8s -run=$(run)
BENCH      := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem -benchtime=128ms
WATCH      := watchexec -r -c -d=0 -n

default: test_w

watch:
	$(MAKE_PAR) test_w lint_w

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
	$(MAKE_PAR) test lint

# Example: `make release tag=v0.0.1`.
release: prep
ifeq ($(tag),)
	$(error missing tag)
endif
	git pull --ff-only
	git show-ref --tags --quiet "$(tag)" || git tag "$(tag)"
	git push origin $$(git symbolic-ref --short HEAD) "$(tag)"
