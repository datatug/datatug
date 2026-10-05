#!/usr/bin/env python3
"""Offline consistency checks for the public W1 follow-up receipt."""

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

HERE = Path(__file__).resolve().parent
HUB = HERE.parents[4]


def sha256(data):
    return hashlib.sha256(data).hexdigest()


def read_json(name):
    return json.loads((HERE / name).read_text())


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--repo",
        action="append",
        default=[],
        metavar="GITHUB_URL=LOCAL_CHECKOUT",
        help="optionally rehash immutable Git objects; repeat for the three source repositories",
    )
    args = parser.parse_args()
    repo_dirs = dict(value.split("=", 1) for value in args.repo)

    proof = read_json("evidence.json")
    sources = read_json("source-manifest.json")
    receipts = read_json("root-publication-receipt-index.json")
    delta = read_json("resource-source-deltas-r3.json")

    assert proof["kind"] == "w1-reviewed-followup-evidence"
    assert len(sources["original_reviewed_documents"]) == 38
    for name, entry in sources["original_reviewed_documents"].items():
        assert "/" not in name and re.fullmatch(r"[0-9a-f]{64}", entry["sha256"]) and entry["bytes"] > 0, name
    assert len(receipts["receipts"]) == proof["programme"]["root_publication_receipts_before_pr193"] == 38
    assert sha256((HERE / "root-publication-receipt-index.json").read_bytes()) == proof["programme"]["root_publication_receipt_index_sha256"]
    assert sha256((HERE / "resource-source-deltas-r3.json").read_bytes()) == proof["official_resource"]["r3"]["source_deltas_sha256"]
    assert len(delta["geonames_places"]["added"]) == 4
    assert len(delta["geonames_places"]["modified"]) == 13
    assert len(delta["geonames_alternate_names"]["added"]) == 49
    assert len(delta["geonames_alternate_names"]["modified"]) == 2

    decision = HUB / "spec/research/public-data-fabric/evidence/w1-ror-user-decision.json"
    assert sha256(decision.read_bytes()) == proof["ror_input_scope"]["original_decision_sha256"]
    assert proof["ror_input_scope"]["changed_data_accepted"] is False
    assert proof["ror_input_scope"]["new_semantic_acceptance"] is False
    assert proof["official_resource"]["r2"]["cap_passed"] is False
    assert proof["official_resource"]["r3"]["cap_passed"] is True
    assert proof["official_resource"]["r3"]["candidate_promoted"] is False
    assert proof["official_resource"]["accepted_pins_changed"] is False
    assert sum(item["bytes"] for item in proof["official_resource"]["r3"]["source_bodies"]) == 670251902
    assert proof["programme"]["w1_complete"] is False
    assert proof["programme"]["b1_open"] is True

    implementation = HUB / "spec/research/public-data-fabric/implementation/$records"
    records = list(implementation.glob("w1-*.json"))
    assert len(records) == 5
    for path in records:
        record = json.loads(path.read_text())
        assert record["status"] == "implementing"
        assert record["publication_evidence_ids"] == []
        assert record["query_evidence_ids"] == []

    checked_objects = 0
    for obj in sources["immutable_git_objects"]:
        repo = repo_dirs.get(obj["repository"])
        if not repo:
            continue
        result = subprocess.run(
            ["git", "-C", repo, "show", f"{obj['revision']}:{obj['path']}"],
            capture_output=True,
            check=True,
        )
        assert sha256(result.stdout) == obj["sha256"], obj
        checked_objects += 1

    for path in HERE.iterdir():
        if path.suffix in {".json", ".md"}:
            content = path.read_bytes()
            assert not re.search(rb"/(?:private/tmp|Users/[^/]+)/", content), path.name
    print(
        "checked 38 original-document hash declarations (originals not rehashed), "
        "rehash-verified 2 packet files and the hub decision, "
        f"checked 38 receipt summaries and 5 implementation states, "
        f"rehash-verified {checked_objects} optional Git objects"
    )


if __name__ == "__main__":
    main()
