#!/usr/bin/env python3
"""Check the local, immutable W1 P1 wire-decision evidence bundle."""

import argparse
import hashlib
import json
from pathlib import Path


EXPECTED_SOURCE_FILES = {
    "source/README.md", "source/proposal.md", "source/proposal.json", "source/test-vectors.json",
    "source/author-evidence.json", "source/root-acceptance-r1.json",
    "source/review-r1.md", "source/review-r1.json",
    "source/review-verification.json", "source/review-probe.py",
    "source/review-probe-results.json", "semantic-addendum-decision.json", "README.md", "verify.py",
}
ORIGINAL_SHA = "62fa8fa939835e3c745b09d6c097a86780c0187a1f6dc57bd66087275a8dbf51"
PROPOSAL_SHA = "6b1096ad1f9162be49cbf4f6be92275d827f65a9d272c8b25b25c72d867a5c1b"
REVIEW_SHA = "f7c1e2942e2ca0b5e1a40811f89a1fbff2e800c61c8897f99fcd339b2d5bb416"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def check(packet, original):
    manifest = json.loads((packet / "evidence.json").read_text())
    decision = json.loads((packet / "semantic-addendum-decision.json").read_text())
    source_proposal = json.loads((packet / "source/proposal.json").read_text())
    source_review = json.loads((packet / "source/review-r1.json").read_text())
    acceptance = json.loads((packet / "source/root-acceptance-r1.json").read_text())
    if set(manifest["files"]) != EXPECTED_SOURCE_FILES:
        raise ValueError("manifest file set changed")
    for name, expected in manifest["files"].items():
        path = packet / name
        if not path.is_file() or path.is_symlink() or sha(path) != expected:
            raise ValueError(f"missing or changed evidence file: {name}")
    if sha(original) != ORIGINAL_SHA or manifest["original_decision_sha256"] != ORIGINAL_SHA:
        raise ValueError("original accepted decision changed")
    if decision["original_decision"]["sha256"] != ORIGINAL_SHA:
        raise ValueError("new decision lost original decision binding")
    if decision["proposal"]["sha256"] != PROPOSAL_SHA or sha(packet / "source/proposal.md") != PROPOSAL_SHA:
        raise ValueError("proposal binding changed")
    if decision["independent_review"]["sha256"] != REVIEW_SHA or sha(packet / "source/review-r1.md") != REVIEW_SHA:
        raise ValueError("independent review binding changed")
    if source_proposal["proposal"]["sha256"] != PROPOSAL_SHA or source_review["review"]["sha256"] != REVIEW_SHA:
        raise ValueError("source packet self-binding changed")
    if acceptance["proposal"]["sha256"] != PROPOSAL_SHA or acceptance["review"]["sha256"] != REVIEW_SHA:
        raise ValueError("root acceptance binding changed")
    if decision["root_acceptance"]["sha256"] != sha(packet / "source/root-acceptance-r1.json"):
        raise ValueError("root acceptance byte binding changed")
    if source_review["verdict"] != "accept-bounded-semantic-addendum" or source_review["counts"] != {"blockers": 0, "majors": 0, "minors": 0}:
        raise ValueError("independent verdict changed")
    if not all("refuse" in decision["failure_scope_reconciliation"][key] for key in ["architecture_general_type_contradiction"]):
        raise ValueError("architecture failure boundary missing")
    if "retains" not in decision["failure_scope_reconciliation"]["local_p1_reference_diagnostic"]:
        raise ValueError("local P1 diagnostic boundary missing")
    if any(decision[key] for key in ["execution_admitted", "publication_admitted", "w1_complete"]):
        raise ValueError("decision grants unreviewed admission or W1 completion")
    if any(manifest[key] for key in ["execution_admitted", "publication_admitted"]):
        raise ValueError("evidence manifest grants admission")
    if len(json.loads((packet / "source/test-vectors.json").read_text()).get("vectors", [])) == 0:
        raise ValueError("semantic vector file empty")
    return len(manifest["files"])


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--packet", type=Path, default=Path(__file__).resolve().parent)
    parser.add_argument("--original-decision", type=Path)
    args = parser.parse_args()
    original = args.original_decision or args.packet.parent / "w1-followup-decisions/semantic-decision.json"
    count = check(args.packet, original)
    print(f"W1 P1 wire addendum evidence verified: {count} retained files and original decision")


if __name__ == "__main__":
    main()
