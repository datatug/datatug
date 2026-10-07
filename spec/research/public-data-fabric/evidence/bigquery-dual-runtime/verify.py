#!/usr/bin/env python3
"""Verify curated small evidence bytes and preservation of original research."""
import hashlib
import json
import sys
from pathlib import Path

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def main():
    packet = Path(__file__).resolve().parent
    hub = packet.parent.parent
    sys.path.insert(0, str(hub))
    from verify_no_copy_decision import approved_record_additions, preserved_bytes
    additions = approved_record_additions(hub.parents[2])
    manifest = json.loads((packet / "source-manifest.json").read_text())
    for item in manifest["files"]:
        path = hub / item["path"]
        retained = preserved_bytes(hub.parents[2], str(path.relative_to(hub.parents[2])), item["sha256"])
        assert hashlib.sha256(retained).hexdigest() == item["sha256"], item["path"]
    for table, expected in manifest["preserved_baseline"].items():
        excluded = set(manifest["added_candidate_ids"]) if table == "candidates" else set()
        excluded.update(additions[table])
        entries = [{"path": str(p.relative_to(hub)), "sha256": digest(p)}
                   for p in sorted((hub / table / "$records").glob("*.json"))
                   if p.stem not in excluded]
        actual = hashlib.sha256(json.dumps(entries, separators=(",", ":"), sort_keys=True).encode()).hexdigest()
        assert len(entries) == expected["count"], table + " count"
        assert actual == expected["sorted_path_sha256_index_sha256"], table + " original bytes"
    proposal = json.loads((packet / "proposal-r3.json").read_text())
    assert digest(packet / "proposal-r3.md") == proposal["proposal_sha256"]
    assert proposal["required_runtimes"] == ["Go", "JavaScript"]
    review = (packet / "history/review-r1.md").read_text()
    for amendment in proposal["mandatory_amendments"].values():
        assert amendment in review, "verbatim amendment provenance"
        assert amendment in (packet / "proposal-r3.md").read_text(), "integrated amendment"
    for ident in manifest["added_candidate_ids"]:
        row = json.loads((hub / "candidates/$records" / (ident + ".json")).read_text())
        assert row["status"] == "needs-review", ident
    a0 = packet / "a0"
    curated = json.loads((a0 / "manifest.json").read_text())
    for name, expected in curated["files"].items():
        path = a0 / name
        assert digest(path) == expected["sha256"], name
        assert path.stat().st_size == expected["bytes"], name + " bytes"
    decision = json.loads((a0 / "root-decision-r1.json").read_text())
    assert digest(a0 / "root-decision-r1.json") == "0a9516308f2a059fb9d4565e5a3530df903e357ff4dbfbe648b5ac23ec309e86"
    assert decision["decision"] == "accept-A0-with-mandatory-amendments"
    for source, expected in decision["input_sha256"].items():
        assert digest(a0 / Path(source).name) == expected, source
    review = json.loads((a0 / "review-r1.json").read_text())
    for source, expected in review["input_sha256"].items():
        assert digest(a0 / Path(source).name) == expected, source
    accepted = json.loads((a0 / "accepted-contract.json").read_text())
    assert accepted["mandatory_amendments"] == decision["mandatory_amendments"]
    assert accepted["required_runtimes"] == ["Go", "JavaScript"]
    assert not accepted["runtime_accepted"] and not accepted["source_admitted"]
    assert digest(a0 / "accepted-contract.md") == accepted["accepted_contract_sha256"]
    text = (a0 / "accepted-contract.md").read_text()
    for amendment in decision["mandatory_amendments"]:
        assert amendment["amendment"] in text, amendment["id"]
        for case in amendment["conformance"]:
            assert case in text, case
    assert curated["canonical_r3_sha256"] == digest(packet / "proposal-r3.md")
    original_index = json.loads((a0 / "hash-index.json").read_text())
    # Source-only historical blobs remain digest declarations. Rehash every
    # original indexed file that is curated here (proposal and compile evidence).
    for item in original_index["files"]:
        path = a0 / item["path"]
        if path.is_file():
            assert digest(path) == item["sha256"], item["path"]
    for table, expected in curated["preserved_record_indexes"].items():
        entries = [{"path": str(p.relative_to(hub)), "sha256": digest(p)}
                   for p in sorted((hub / table / "$records").glob("*.json"))
                   if p.stem not in additions[table]]
        assert len(entries) == expected["count"], table
        actual = hashlib.sha256(json.dumps(entries, separators=(",", ":"), sort_keys=True).encode()).hexdigest()
        assert actual == expected["sha256"], table + " full A0 baseline"
    print("BigQuery r3/A0 bytes, review/root bindings, exact minor conformance, both runtimes and full 699/1177 record preservation verified")

if __name__ == "__main__":
    main()
