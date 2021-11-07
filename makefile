MAKEFLAGS  := --silent --always-make
PAR        := $(MAKE) -j 128
TEST_FLAGS := $(if $(filter $(verb), true), -v,) -count=1
TEST       := test $(TEST_FLAGS) -timeout=1s -run=$(run)
BENCH      := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem
WATCH      := watchexec -r -c -d=0 -n

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
