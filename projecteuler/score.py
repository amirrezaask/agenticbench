#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from pathlib import Path

WORKSPACE = Path(__file__).parent.resolve()
SOLUTIONS_FILE = WORKSPACE / "solutions.txt"


def load_expected_solutions() -> dict[int, str]:
    solutions = {}
    with open(SOLUTIONS_FILE, "r") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            parts = line.split(". ", 1)
            if len(parts) == 2:
                problem_num = int(parts[0])
                answer = parts[1].strip()
                solutions[problem_num] = answer
    return solutions


def find_implementations() -> list[Path]:
    implementations = []
    for item in WORKSPACE.iterdir():
        if item.is_dir() and "-" in item.name and item.name != "problems":
            if item.name.startswith("."):
                continue
            implementations.append(item)
    return sorted(implementations)


def detect_language(impl_dir: Path) -> str | None:
    if (impl_dir / "go.mod").exists():
        return "go"
    if (impl_dir / "main.py").exists() or (impl_dir / "solutions.py").exists():
        return "python"
    if (impl_dir / "Cargo.toml").exists():
        return "rust"
    if (impl_dir / "package.json").exists():
        return "javascript"
    return None


def build_and_run_go(impl_dir: Path) -> tuple[bool, str, str]:
    build_result = subprocess.run(
        ["go", "build", "-o", "runner"],
        cwd=impl_dir,
        capture_output=True,
        text=True,
    )
    if build_result.returncode != 0:
        return False, "", f"Build failed: {build_result.stderr}"
    
    run_result = subprocess.run(
        ["./runner"],
        cwd=impl_dir,
        capture_output=True,
        text=True,
        timeout=3600,
    )
    return True, run_result.stdout, run_result.stderr


def run_python(impl_dir: Path) -> tuple[bool, str, str]:
    main_file = impl_dir / "main.py"
    if not main_file.exists():
        main_file = impl_dir / "solutions.py"
    
    if not main_file.exists():
        return False, "", "No main.py or solutions.py found"
    
    run_result = subprocess.run(
        [sys.executable, str(main_file)],
        cwd=impl_dir,
        capture_output=True,
        text=True,
        timeout=3600,
    )
    return True, run_result.stdout, run_result.stderr


def run_rust(impl_dir: Path) -> tuple[bool, str, str]:
    build_result = subprocess.run(
        ["cargo", "build", "--release"],
        cwd=impl_dir,
        capture_output=True,
        text=True,
    )
    if build_result.returncode != 0:
        return False, "", f"Build failed: {build_result.stderr}"
    
    run_result = subprocess.run(
        ["cargo", "run", "--release"],
        cwd=impl_dir,
        capture_output=True,
        text=True,
        timeout=3600,
    )
    return True, run_result.stdout, run_result.stderr


def run_javascript(impl_dir: Path) -> tuple[bool, str, str]:
    main_file = impl_dir / "main.js"
    if not main_file.exists():
        return False, "", "No main.js found"
    
    run_result = subprocess.run(
        ["node", str(main_file)],
        cwd=impl_dir,
        capture_output=True,
        text=True,
        timeout=3600,
    )
    return True, run_result.stdout, run_result.stderr


def run_implementation(impl_dir: Path) -> tuple[bool, str, str]:
    lang = detect_language(impl_dir)
    if lang is None:
        return False, "", "Unknown language/implementation type"
    
    runners = {
        "go": build_and_run_go,
        "python": run_python,
        "rust": run_rust,
        "javascript": run_javascript,
    }
    
    runner = runners.get(lang)
    if runner is None:
        return False, "", f"No runner for language: {lang}"
    
    try:
        return runner(impl_dir)
    except subprocess.TimeoutExpired:
        return False, "", "Execution timed out (1 hour limit)"
    except Exception as e:
        return False, "", f"Execution error: {str(e)}"


def parse_results(output: str) -> list[dict] | None:
    results = []
    for line in output.strip().split("\n"):
        line = line.strip()
        if not line:
            continue
        try:
            data = json.loads(line)
            if "solution" in data and "result" in data:
                results.append({
                    "problem": data["solution"],
                    "answer": str(data["result"]),
                    "time_ms": data.get("time_ms", 0),
                })
        except json.JSONDecodeError:
            continue
    return results if results else None


def normalize_answer(answer: str) -> str:
    return answer.strip()


def score_implementation(results: list[dict], expected: dict[int, str]) -> dict:
    correct = 0
    incorrect = 0
    details = []
    
    for result in results:
        problem = result["problem"]
        answer = normalize_answer(str(result["answer"]))
        time_ms = result.get("time_ms", 0)
        
        if problem not in expected:
            details.append({
                "problem": problem,
                "status": "unknown",
                "got": answer,
                "expected": None,
                "time_ms": time_ms,
            })
            continue
        
        expected_answer = normalize_answer(expected[problem])
        
        if answer == expected_answer:
            correct += 1
            details.append({
                "problem": problem,
                "status": "correct",
                "got": answer,
                "expected": expected_answer,
                "time_ms": time_ms,
            })
        else:
            incorrect += 1
            details.append({
                "problem": problem,
                "status": "incorrect",
                "got": answer,
                "expected": expected_answer,
                "time_ms": time_ms,
            })
    
    total_time = sum(d.get("time_ms", 0) for d in details)
    
    return {
        "correct": correct,
        "incorrect": incorrect,
        "total_attempted": len(results),
        "accuracy": correct / len(results) * 100 if results else 0,
        "total_time_ms": total_time,
        "details": details,
    }


def output_table(all_scores: dict[str, dict]):
    print(f"{'Implementation':<30} {'Correct':<10} {'Total':<10} {'Accuracy':<12} {'Time (s)':<10}")
    print("-" * 72)
    for impl_name, score in sorted(all_scores.items()):
        if "error" in score:
            print(f"{impl_name:<30} {'ERROR':<10} {'-':<10} {'-':<12} {'-':<10}")
        else:
            time_s = score.get('total_time_ms', 0) / 1000
            print(f"{impl_name:<30} {score['correct']:<10} {score['total_attempted']:<10} {score['accuracy']:.2f}%{'':<6} {time_s:.2f}")


def output_csv(all_scores: dict[str, dict]):
    print("implementation,correct,total,accuracy,time_s")
    for impl_name, score in sorted(all_scores.items()):
        if "error" in score:
            print(f"{impl_name},ERROR,-,-,-")
        else:
            time_s = score.get('total_time_ms', 0) / 1000
            print(f"{impl_name},{score['correct']},{score['total_attempted']},{score['accuracy']:.2f},{time_s:.2f}")


def output_json(all_scores: dict[str, dict]):
    summary = []
    for impl_name, score in sorted(all_scores.items()):
        if "error" in score:
            summary.append({"implementation": impl_name, "error": score["error"]})
        else:
            summary.append({
                "implementation": impl_name,
                "correct": score["correct"],
                "total": score["total_attempted"],
                "accuracy": round(score["accuracy"], 2),
                "time_s": round(score.get("total_time_ms", 0) / 1000, 2),
            })
    print(json.dumps(summary, indent=2))


def main():
    parser = argparse.ArgumentParser(description="Project Euler Implementation Scorer")
    parser.add_argument(
        "-f", "--format",
        choices=["table", "csv", "json"],
        default="table",
        help="Output format (default: table)"
    )
    parser.add_argument(
        "-q", "--quiet",
        action="store_true",
        help="Suppress progress messages"
    )
    args = parser.parse_args()

    expected = load_expected_solutions()
    implementations = find_implementations()
    
    if not implementations:
        if not args.quiet:
            print("No implementations found!", file=sys.stderr)
        return
    
    all_scores = {}
    
    for impl_dir in implementations:
        impl_name = impl_dir.name
        if not args.quiet:
            print(f"Running {impl_name}...", file=sys.stderr)
        
        success, stdout, stderr = run_implementation(impl_dir)
        
        if not success:
            if not args.quiet:
                print(f"  ERROR: {stderr}", file=sys.stderr)
            all_scores[impl_name] = {"error": stderr}
            continue
        
        results = parse_results(stdout)
        if results is None:
            if not args.quiet:
                print(f"  ERROR: Could not parse output", file=sys.stderr)
            all_scores[impl_name] = {"error": "Could not parse output"}
            continue
        
        score = score_implementation(results, expected)
        all_scores[impl_name] = score
    
    if args.format == "table":
        output_table(all_scores)
    elif args.format == "csv":
        output_csv(all_scores)
    elif args.format == "json":
        output_json(all_scores)
    
    output_file = WORKSPACE / "scores.json"
    with open(output_file, "w") as f:
        json.dump(all_scores, f, indent=2)


if __name__ == "__main__":
    main()
