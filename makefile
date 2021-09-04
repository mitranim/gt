MAKEFLAGS := --silent --always-make
TEST      := test $(if $(filter $(verb), true), -v,) -count=1 -short -run=$(run)
BENCH     := test -count=1 -short -bench=$(or $(run),.) -benchmem

test-w:
	gow -c -v $(TEST)

test:
	go $(TEST)

bench-w:
	gow -c -v $(BENCH)

bench:
	go $(BENCH)

lint-w:
	watchexec -r -c -d=0 -n $(MAKE) lint

lint:
	golangci-lint run
	echo [lint] ok
