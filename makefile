MAKEFLAGS  := --silent --always-make
PAR        := $(MAKE) -j 128
VERB       := $(if $(filter $(verb), true), -v,)
SHORT      := $(if $(filter $(short), true), -short,)
TEST_FLAGS := -count=1 $(VERB) $(SHORT)
TEST       := test $(TEST_FLAGS) -timeout=8s -run=$(run)
BENCH      := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem
WATCH      := watchexec -r -c -d=0 -n

default: test_w

watch:
	$(PAR) test_w lint_w

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
