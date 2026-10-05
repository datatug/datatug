#!/usr/bin/env python3
"""Verify curated small evidence bytes and preservation of original research."""
import hashlib
import json
from pathlib import Path

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def main():
    packet = Path(__file__).resolve().parent
    hub = packet.parent.parent
    manifest = json.loads((packet / "source-manifest.json").read_text())
    for item in manifest["files"]:
        path = hub / item["path"]
        assert digest(path) == item["sha256"], item["path"]
    for table, expected in manifest["preserved_baseline"].items():
        excluded = set(manifest["added_candidate_ids"]) if table == "candidates" else set()
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
    print("BigQuery packet hashes, verbatim amendments, dual-runtime requirement, original691/1177 bytes and8unreconciled candidates verified")

if __name__ == "__main__":
    main()
