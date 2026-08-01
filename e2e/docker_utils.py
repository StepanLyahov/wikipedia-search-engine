"""Thin wrappers around `docker compose` used to drive the stack for end-to-end tests."""
from __future__ import annotations

import json
import subprocess
import time
from pathlib import Path
from typing import Any

REPO_ROOT = Path(__file__).resolve().parent.parent


def run_compose(*args: str, timeout: int = 120) -> subprocess.CompletedProcess:
    return subprocess.run(
        ["docker", "compose", *args],
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        timeout=timeout,
    )


def up(build: bool = True, timeout: int = 600) -> None:
    args = ["up", "-d"]
    if build:
        args.append("--build")
    result = run_compose(*args, timeout=timeout)
    if result.returncode != 0:
        raise RuntimeError(f"docker compose up failed:\nstdout: {result.stdout}\nstderr: {result.stderr}")


def down(volumes: bool = True, timeout: int = 120) -> None:
    args = ["down"]
    if volumes:
        args.append("-v")
    run_compose(*args, timeout=timeout)


def resolved_config() -> dict[str, Any]:
    """Return the fully resolved compose config (env vars, ports interpolated)."""
    result = run_compose("config", "--format", "json")
    if result.returncode != 0:
        raise RuntimeError(f"docker compose config failed: {result.stderr}")
    return json.loads(result.stdout)


def published_port(config: dict[str, Any], service: str, target_port: int) -> int:
    for port in config["services"][service].get("ports", []):
        if int(port["target"]) == target_port:
            return int(port["published"])
    raise KeyError(f"service {service!r} does not publish port {target_port}")


def service_environment(config: dict[str, Any], service: str) -> dict[str, str]:
    return config["services"][service].get("environment", {})


def ps(all_states: bool = True) -> list[dict[str, Any]]:
    args = ["ps", "--format", "json"]
    if all_states:
        args.insert(1, "-a")
    result = run_compose(*args)
    if result.returncode != 0:
        raise RuntimeError(f"docker compose ps failed: {result.stderr}")

    containers = []
    for line in result.stdout.splitlines():
        line = line.strip()
        if line:
            containers.append(json.loads(line))
    return containers


def wait_for_one_shot_completion(service: str, timeout: float) -> None:
    """Wait for a one-shot service (crawler, indexer) to exit, and fail loudly on a nonzero
    exit code so a broken pipeline stage is reported at the fixture boundary, not as a
    confusing downstream assertion failure."""
    deadline = time.monotonic() + timeout
    last_state = None
    while time.monotonic() < deadline:
        for container in ps():
            if container.get("Service") != service:
                continue
            last_state = container.get("State")
            if last_state == "exited":
                exit_code = container.get("ExitCode")
                if exit_code != 0:
                    raise RuntimeError(
                        f"service {service!r} exited with code {exit_code}; "
                        f"check `docker compose logs {service}`"
                    )
                return
        time.sleep(1)

    raise TimeoutError(f"service {service!r} did not exit within {timeout}s (last state: {last_state})")


def wait_for_healthy(service: str, timeout: float) -> None:
    deadline = time.monotonic() + timeout
    last_state = None
    while time.monotonic() < deadline:
        for container in ps():
            if container.get("Service") != service:
                continue
            last_state = (container.get("State"), container.get("Health"))
            if container.get("State") == "running" and container.get("Health") in ("healthy", ""):
                return
        time.sleep(1)

    raise TimeoutError(f"service {service!r} did not become healthy within {timeout}s (last state: {last_state})")
