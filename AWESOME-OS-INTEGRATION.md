# Awesome-OS Integration: `awesome-os/maclab`

> Upstream source: [`chriskearney/maclab`](https://github.com/chriskearney/maclab) (MIT License, preserved intact in [`LICENSE`](LICENSE)).

`awesome-os/maclab` is the Apple Silicon hardware-in-the-loop (HIL) kernel test lab for the Awesome-OS Apple hardware ecosystem. It provides a controller daemon (`labd`), on-device agent (`lab-agent`), out-of-band DFU/serial recovery daemon (`oobd`), CLI (`lab`), Web UI, and an **MCP server (`lab mcp`)** for autonomous AI agents.

## Role in Awesome-OS

`maclab` automates one-shot kernel booting, crash/hang recovery, log collection, and GUI/suspend verification on physical Apple Silicon Macs for:

- **[`awesome-os/t2linux`](https://github.com/awesome-os/t2linux)**: Stages and boots candidate Linux kernels (GRUB or Limine UKI on ESP) once via one-shot boot flags (`maclab_next` / `LoaderEntryOneShot`), runs `boot-health`, Wayland/Hyprland GUI tests (`gui-smoke`, `gui-terminal`), and `suspend-devices` cycles, collects `dmesg`/`journalctl`/serial console logs, and automatically restores the known-good kernel on panic or watchdog timeout.
- **[`awesome-os/t2xnu`](https://github.com/awesome-os/t2xnu)**: Coordinates out-of-band DFU/serial console monitoring (`oobd`), device leasing (`lab lease`), and automated kernel bring-up regression loops on Apple Silicon hardware.

## Agent & MCP Workflow

Autonomous Awesome-OS agents interact with the lab through the built-in Model Context Protocol (MCP) server:

```sh
# Start the MCP server for agent tool calls
lab mcp

# Or invoke directly from CLI
lab run m1air --kernel kernel-t2linux.tar.zst --gui-test gui-smoke
```

Key capabilities exposed to agents over MCP:
- **Device discovery & leasing**: Reserve physical Macs (`lab lease <device> --ttl 3h`) so concurrent kernel CI runs never collide.
- **Safe one-shot kernel testing**: Boot experimental `t2linux` / `t2xnu` builds without ever mutating the known-good fallback kernel.
- **Automated failure classification**: Classifies `pass`, `tests_failed`, `booted_unhealthy`, `panicked`, `hung`, and `crashed_or_hung` outcomes with `new_error_lines` diffs against baseline.

## Upstream Attribution

All original Go sources, CLI/Web UI/MCP implementations, TLA+ specs, and MIT licensing are preserved from [`chriskearney/maclab`](https://github.com/chriskearney/maclab).
