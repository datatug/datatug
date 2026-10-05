# Focused SDK and policy compile proof

Exact released dependencies and in-memory transports are preserved with original test bytes. Run `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -v ./...` with a private writable cache. The author proof has five tests; independent review adds two tests. Original independent source bytes are retained in `contract_test.go.txt`; `contract_test.go` is a gofmt-only runnable copy, bound by `format-receipt.json`, to satisfy normal WB hooks. No real credentials or cloud jobs are used. This proves SDK/policy API compatibility, not adapter acceptance.

## Open Questions

Production raw-wire corpus and both actual Go/JavaScript live gates remain required.
