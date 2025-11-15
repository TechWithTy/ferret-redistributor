ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
PY   := python

.PHONY: all validate snapshot

all: validate

validate:
	@echo "[make] validating TOON schemas"
	@$(PY) $(ROOT)/shared/tools/run_codegen.py --lang go

snapshot:
	@echo "[make] snapshotting schema metadata"
	@$(PY) $(ROOT)/shared/tools/run_codegen.py --lang python



