# bubbleMonitor: Senior Engineer Grade Refactoring Plan

## TL;DR

> **Quick Summary**: Comprehensive refactoring of bubbleMonitor from 4.5/10 to 10/10 through architectural decomposition, interface abstraction, test-driven development, and feature extensions while maintaining 100% backward compatibility and identical UI.

> **Deliverables**:
> - Split monolithic update.go (1,337 lines → ~300 lines with handlers in separate files)
> - Provider interfaces with adapter pattern for testability
> - Logical grouping extraction from AppState (160 fields → structured sub-objects)
> - Full TDD test coverage for core app logic
> - Container metrics (Docker/K8s) provider
> - VM/hypervisor metrics provider
> - Full remote dashboard (all tabs available remotely)
> - golangci-lint configuration
> - GitHub Actions CI/CD pipeline
> - Fixed Makefile

> **Estimated Effort**: XL (3-4 months full-time equivalent)
> **Parallel Execution**: YES - multiple waves
> **Critical Path**: Phase 1 (tests) → Phase 2 (update.go) → Phase 3 (interfaces) → Phase 4 (features) → Phase 5 (tooling)

---

## Context

### Original Request
Take bubbleMonitor from 4.5/10 to 10/10 - Senior Engineer Grade. User explicitly requested:
- No rush - thorough refactoring (quality over speed)
- Full backward compatibility (config, CLI, behavior)
- UI MUST stay 100% identical (manual verification acceptable)
- TDD approach (tests first)
- Add container/VM metrics and full remote dashboard

### Interview Summary

**Key Discussions**:
- Timeline: No rush, thorough refactor
- Backward compatibility: Required, config file format unchanged
- UI constraint: 100% identical, manual verification OK
- Testing: TDD - tests first before implementation
- Plugin system: DEFERRED to future (out of scope)
- Features confirmed: Container metrics, VM metrics, Full remote dashboard

**Research Findings**:
- Provider pattern uses tea.Command returning typed messages
- No Go interfaces exist - implicit contracts only
- 22 distinct tea.Cmd functions across providers
- ~25 message types in internal/msg/messages.go
- AppState has ~160 fields (data/state.go lines 20-162)
- update.go is 1,337 lines (single switch statement)

### Metis Review

**Identified Gaps** (addressed in plan):
- Scope inflation risk: Plugin system deferred to future
- UI verification: Manual acceptable (user confirmed)
- TDD approach: Integration tests before refactoring
- Completion criteria: Defined for each phase

---

## Work Objectives

### Core Objective
Transform bubbleMonitor from a functional but architecturally flawed project into a senior-engineer-grade codebase through:
1. Architectural decomposition (update.go, AppState)
2. Interface abstraction (providers)
3. Test-driven development (full coverage)
4. Feature extension (containers, VMs, remote)
5. Tooling maturity (linting, CI/CD)

### Concrete Deliverables
- `internal/app/handlers/` - New directory with separated message handlers
- `internal/provider/interfaces.go` - Provider interface definitions
- `internal/data/state_groups.go` - Logical AppState groupings
- `internal/provider/system/container.go` - Docker/K8s metrics
- `internal/provider/system/vm.go` - VM/hypervisor metrics
- `internal/provider/remote/` - Extended remote with full dashboard
- `.golangci.toml` - Linting configuration
- `.github/workflows/` - CI/CD pipeline

### Definition of Done

**Architecture**:
- [ ] update.go reduced to <400 lines
- [ ] Each message handler in separate file
- [ ] All providers implement interfaces
- [ ] AppState fields grouped into logical structs

**Testing**:
- [ ] Integration tests for Update() covering all message types
- [ ] Unit tests for each message handler
- [ ] Provider mocks implement interfaces
- [ ] All existing tests pass

**Features**:
- [ ] Docker container metrics displayed
- [ ] Kubernetes pod metrics displayed (if K8s detected)
- [ ] VM hypervisor metrics (if hypervisor detected)
- [ ] Remote tab shows full dashboard (not just uptime)

**Tooling**:
- [ ] golangci-lint passes with zero warnings
- [ ] GitHub Actions runs tests on PR
- [ ] Makefile fixed (no $(eval))

**Backward Compatibility**:
- [ ] Existing config files load without modification
- [ ] All keyboard shortcuts work identically
- [ ] All 30+ themes render identically

### Must Have
- 100% backward compatibility with existing configs
- Identical UI rendering (manual verification)
- All existing features functional
- TDD approach: tests before refactoring

### Must NOT Have (Guardrails)
- NO changes to UI rendering code (internal/ui/*)
- NO breaking changes to message type signatures
- NO changes to config file location or format
- NO removal of existing features or tabs
- NO plugin system (deferred to future)

---

## Verification Strategy (MANDATORY)

> **UNIVERSAL RULE: ZERO HUMAN INTERVENTION**
>
> ALL tasks in this plan MUST be verifiable WITHOUT any human action.

### Test Decision
- **Infrastructure exists**: YES (go test framework available)
- **Automated tests**: TDD - tests first
- **Framework**: Go native testing (bun test / go test)

### TDD Workflow

Each refactoring task follows RED-GREEN-REFACTOR:

**Task Structure**:
1. **RED**: Write failing test first
   - Test file: `[path]_test.go`
   - Test command: `go test -v ./...`
   - Expected: FAIL (test exists, implementation doesn't)
2. **GREEN**: Implement minimum code to pass
   - Command: `go test -v ./...`
   - Expected: PASS
3. **REFACTOR**: Clean up while keeping green
   - Command: `go test -v ./...`
   - Expected: PASS (still)

### Agent-Executed QA Scenarios

**Verification Tool**: Go native testing + bash

| Type | Tool | Verification |
|------|------|--------------|
| **Unit Tests** | Bash (go test) | go test -v ./internal/app/... |
| **Integration** | Bash (go test) | go test -v ./internal/app/... -run Integration |
| **Build** | Bash (go build) | go build -o bub ./cmd/bub |
| **Lint** | Bash (golangci-lint) | golangci-lint run ./... |
| **Format** | Bash (gofmt) | gofmt -d ./internal/... |

**Evidence Requirements**:
- Test output: go test results captured
- Coverage: go tool cover output
- Lint: golangci-lint output saved
- Build: Binary created and runs

---

## Execution Strategy

### Parallel Execution Waves

```
Phase 1: Test Infrastructure & Baseline (WEEK 1-2)
├── Task 1: Add integration tests for Update()
├── Task 2: Add test infrastructure
└── Task 3: Verify baseline tests pass

Phase 2: update.go Decomposition (WEEK 3-6)
├── Task 4: Extract key handler functions
├── Task 5: Create handler files
├── Task 6: Update main Update() to delegate
└── Task 7: Tests pass with new structure

Phase 3: Provider Interfaces (WEEK 7-10)
├── Task 8: Define provider interfaces
├── Task 9: Create adapter wrappers
├── Task 10: Add provider mocks
├── Task 11: Refactor app to use interfaces
└── Task 12: Tests pass with interfaces

Phase 4: AppState Logical Groups (WEEK 11-12)
├── Task 13: Define state group structs
├── Task 14: Migrate fields to groups
├── Task 15: Update all references
└── Task 16: Tests pass

Phase 5: New Features (WEEK 13-18)
├── Task 17: Container metrics provider
├── Task 18: K8s metrics provider
├── Task 19: VM metrics provider
├── Task 20: Remote full dashboard
└── Task 21: Feature integration tests

Phase 6: Tooling & Polish (WEEK 19-20)
├── Task 22: golangci-lint config
├── Task 23: GitHub Actions workflow
├── Task 24: Fix Makefile
└── Task 25: Final integration tests
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 4,5,6 | 2,3 |
| 2 | None | 4 | 1,3 |
| 3 | 1,2 | 7 | 5,6 |
| 4 | 1 | 7 | 5,6 |
| 5 | 1 | 7 | 4,6 |
| 6 | 1 | 7 | 4,5 |
| 7 | 4,5,6 | 8 | None |
| 8 | 7 | 9,10,11 | None |
| 9 | 8 | 11 | 10 |
| 10 | 8 | 11 | 9 |
| 11 | 9,10 | 12 | None |
| 12 | 11 | 13,14,15 | None |
| 13 | 12 | 16 | 14,15 |
| 14 | 12 | 16 | 13,15 |
| 15 | 12 | 16 | 13,14 |
| 16 | 13,14,15 | 17,18,19,20 | None |
| 17 | 16 | 21 | 18,19,20 |
| 18 | 16 | 21 | 17,19,20 |
| 19 | 16 | 21 | 17,18,20 |
| 20 | 16 | 21 | 17,18,19 |
| 21 | 17,18,19,20 | 22,23,24 | None |
| 22 | 21 | 25 | 23,24 |
| 23 | 21 | 25 | 22,24 |
| 24 | 21 | 25 | 22,23 |
| 25 | 22,23,24 | None | None |

---

## TODOs

### Phase 1: Test Infrastructure & Baseline

- [ ] 1. Add Integration Tests for Update() Message Handling

  **What to do**:
  - Create `internal/app/update_integration_test.go`
  - Test all message types: TickMsg, CpuMemMsg, ProcessesMsg, etc.
  - Test key sequences: startup, tab switch, process kill
  - Test error paths: provider failures, config errors

  **Must NOT do**:
  - Don't test UI rendering (out of scope)
  - Don't modify Update() logic yet

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Complex integration testing requiring understanding of message flows
  - **Skills**: []
    - go-testing: Integration test patterns
  - **Skills Evaluated but Omitted**:
    - frontend-ui-ux: Not applicable

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3)
  - **Blocks**: Tasks 4, 5, 6 (need baseline tests)
  - **Blocked By**: None (can start immediately)

  **References**:
  - `internal/app/update.go:1-100` - Update() function signature
  - `internal/msg/messages.go` - All message type definitions
  - `internal/data/state.go:164-276` - AppState methods for testing

  **Acceptance Criteria**:
  - [ ] Integration test file created: internal/app/update_integration_test.go
  - [ ] Tests cover: Startup flow, Tab switching, Process operations, Settings
  - [ ] go test -v ./internal/app/... -run Integration → PASS

- [ ] 2. Setup Extended Test Infrastructure

  **What to do**:
  - Add test helper functions in `internal/app/testutil/`
  - Create mock providers for system, process, remote
  - Add test fixtures for AppState
  - Configure test coverage reporting

  **Must NOT do**:
  - Don't modify production code

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Infrastructure setup, straightforward
  - **Skills**: []
  - **Skills Evaluated but Omitted**:

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3)
  - **Blocks**: Task 4 (needs test helpers)
  - **Blocked By**: None

  **Acceptance Criteria**:
  - [ ] Test utility package created: internal/app/testutil/
  - [ ] Mock providers implemented
  - [ ] go test -v ./... -cover → PASS with coverage report

- [ ] 3. Verify Baseline Tests Pass

  **What to do**:
  - Run all existing tests: go test ./...
  - Ensure no regressions
  - Capture baseline coverage metrics

  **Must NOT do**:
  - Don't modify test files yet

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Running tests, simple verification
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2)
  - **Blocks**: Task 7 (needs passing baseline)
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:
  - [ ] go test ./... → all PASS
  - [ ] Coverage baseline documented

---

### Phase 2: update.go Decomposition

- [ ] 4. Extract Key Handler Functions from update.go

  **What to do**:
  - Analyze update.go to identify handler clusters
  - Extract: keyHandler, mouseHandler, toastHandler, configHandler
  - Each handler becomes a function in internal/app/handlers/
  - Keep Update() as dispatcher only

  **Must NOT do**:
  - Don't change behavior - just refactor structure
  - Don't touch UI code

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Complex refactoring requiring deep understanding
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on Task 1)
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 7
  - **Blocked By**: Task 1

  **References**:
  - `internal/app/update.go:185-909` - Key handling (600+ lines)
  - `internal/app/update.go:52-58` - Mouse handling
  - `internal/app/update.go:85-104` - Toast handling
  - `internal/app/update.go:152-183` - Config handling

  **Acceptance Criteria**:
  - [ ] Key handler extracted: internal/app/handlers/keyboard.go
  - [ ] Mouse handler extracted: internal/app/handlers/mouse.go
  - [ ] Toast handler extracted: internal/app/handlers/toast.go
  - [ ] Config handler extracted: internal/app/handlers/config.go
  - [ ] go test ./... → PASS

- [ ] 5. Create Handler Files with Interfaces

  **What to do**:
  - Create handler interface definitions
  - Each handler: func(*Model, tea.Msg) (tea.Model, tea.Cmd)
  - Register handlers in init()
  - Update() becomes simple dispatcher

  **Must NOT do**:
  - Don't change handler logic

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Interface design required
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 7
  - **Blocked By**: Task 1

  **Acceptance Criteria**:
  - [ ] Handler interface defined
  - [ ] Update() reduced to <400 lines
  - [ ] All tests pass

- [ ] 6. Simplify Main Update() Function

  **What to do**:
  - Reduce Update() from 1,337 lines to ~200 lines
  - Make it purely a dispatcher to handlers
  - Remove duplicated code (tab switching repeated twice)
  - Extract magic numbers to constants

  **Must NOT do**:
  - Don't change behavior
  - Don't touch UI rendering

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Complex logic extraction
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 7
  - **Blocked By**: Task 1

  **References**:
  - `internal/app/update.go:399-459` - Tab handling duplicated
  - `internal/app/constants.go` - Check for existing constants

  **Acceptance Criteria**:
  - [ ] Update() < 400 lines
  - [ ] No duplicate code blocks
  - [ ] Magic numbers extracted to constants

- [ ] 7. Verify Decomposition Tests Pass

  **What to do**:
  - Run full test suite
  - Verify behavior unchanged
  - Measure code reduction

  **Must NOT do**:
  - Don't add new features yet

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Verification
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (final)
  - **Blocks**: Task 8
  - **Blocked By**: Tasks 4, 5, 6

  **Acceptance Criteria**:
  - [ ] go test ./... → all PASS
  - [ ] Update() line count < 400
  - [ ] Integration tests still pass

---

### Phase 3: Provider Interfaces

- [ ] 8. Define Provider Interfaces

  **What to do**:
  - Create `internal/provider/interfaces.go`
  - Define: SystemProvider, ProcessProvider, RemoteProvider interfaces
  - Each interface declares tea.Cmd-returning methods
  - Include capability detection in SystemProvider

  **Must NOT do**:
  - Don't break existing function signatures

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: Interface design with backward compat
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: Tasks 9, 10
  - **Blocked By**: Task 7

  **References**:
  - `internal/provider/system/metrics.go` - Existing commands
  - `internal/provider/process/list.go` - Process commands
  - `internal/provider/remote/ssh.go` - Remote commands

  **Acceptance Criteria**:
  - [ ] interfaces.go created with SystemProvider interface
  - [ ] interfaces.go created with ProcessProvider interface
  - [ ] interfaces.go created with RemoteProvider interface
  - [ ] All existing functions still work

- [ ] 9. Create Adapter Wrappers

  **What to do**:
  - Create adapter structs that implement interfaces
  - Wrap existing functions as methods on adapters
  - Maintain backward compatibility (existing functions call adapters)

  **Must NOT do**:
  - Don't change function signatures

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: Adapter pattern implementation
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 11
  - **Blocked By**: Task 8

  **Acceptance Criteria**:
  - [ ] SystemAdapter implements SystemProvider
  - [ ] ProcessAdapter implements ProcessProvider
  - [ ] RemoteAdapter implements RemoteProvider
  - [ ] go build ./... → PASS

- [ ] 10. Add Provider Mocks for Testing

  **What to do**:
  - Create mock implementations for each interface
  - Mocks return controlled test data
  - Use mock in integration tests

  **Must NOT do**:
  - Don't modify production code

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Straightforward mock implementation
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 11
  - **Blocked By**: Task 8

  **Acceptance Criteria**:
  - [ ] MockSystemProvider implemented
  - [ ] MockProcessProvider implemented
  - [ ] MockRemoteProvider implemented
  - [ ] Integration tests use mocks

- [ ] 11. Refactor App to Use Interfaces

  **What to do**:
  - Add Provider field to Model
  - Initialize with real adapters in Init()
  - Use interface methods in Update()

  **Must NOT do**:
  - Don't change behavior

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Cross-cutting refactoring
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: Task 12
  - **Blocked By**: Tasks 9, 10

  **References**:
  - `internal/app/model.go:22-25` - Model struct

  **Acceptance Criteria**:
  - [ ] Model has Provider field
  - [ ] Providers initialized in InitialModel()
  - [ ] go test ./... → PASS

- [ ] 12. Verify Interface Tests Pass

  **What to do**:
  - Run full test suite with mocks
  - Verify integration tests pass with mock providers

  **Must NOT do**:
  - Don't add features

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Verification
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (final)
  - **Blocks**: Task 13
  - **Blocked By**: Task 11

  **Acceptance Criteria**:
  - [ ] go test ./... → PASS
  - [ ] Integration tests work with mocks

---

### Phase 4: AppState Logical Groups

- [ ] 13. Define State Group Structs

  **What to do**:
  - Create state groups in `internal/data/groups.go`:
    - MetricsState: CpuHistory, MemHistory, NetHistory, etc.
    - ProcessState: Processes, ProcessFilter, SortBy, etc.
    - UIState: SelectedTab, ShowHelp, ShowSettings, etc.
    - ConfigState: Theme, BorderType, RefreshRate, etc.
    - RemoteState: RemoteUptimes, RemoteHosts, etc.

  **Must NOT do**:
  - Don't change field names
  - Don't change behavior

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: Structural design of state groups
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 16
  - **Blocked By**: Task 12

  **References**:
  - `internal/data/state.go:20-162` - All AppState fields

  **Acceptance Criteria**:
  - [ ] MetricsState struct defined
  - [ ] ProcessState struct defined
  - [ ] UIState struct defined
  - [ ] ConfigState struct defined
  - [ ] RemoteState struct defined

- [ ] 14. Migrate Fields to Groups

  **What to do**:
  - Embed group structs in AppState
  - Move fields from AppState to appropriate groups
  - Update all field references throughout codebase

  **Must NOT do**:
  - Don't change field data types
  - Don't change behavior

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Large-scale field migration
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 16
  - **Blocked By**: Task 12

  **References**:
  - `internal/data/state.go` - State definition
  - `internal/app/update.go` - Field references

  **Acceptance Criteria**:
  - [ ] AppState embeds all group structs
  - [ ] Fields migrated to appropriate groups
  - [ ] No orphan fields

- [ ] 15. Update All References

  **What to do**:
  - Update all code that references old AppState fields
  - Use group accessors: m.UIState.SelectedTab
  - Run tests to catch any missed references

  **Must NOT do**:
  - Don't change field names

  **Recommended Agent Profile**:
  - **Category**: ultrabrain
    - Reason: Large-scale code update
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4
  - **Blocks**: Task 16
  - **Blocked By**: Task 12

  **Acceptance Criteria**:
  - [ ] All references updated to use groups
  - [ ] go test ./... → PASS

- [ ] 16. Verify State Group Tests Pass

  **What to do**:
  - Run full test suite
  - Verify all state operations work

  **Must NOT do**:
  - Don't add features

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Verification
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (final)
  - **Blocks**: Tasks 17, 18, 19, 20
  - **Blocked By**: Tasks 13, 14, 15

  **Acceptance Criteria**:
  - [ ] go test ./... → PASS
  - [ ] All existing functionality works

---

### Phase 5: New Features

- [ ] 17. Container Metrics Provider (Docker/K8s)

  **What to do**:
  - Create `internal/provider/system/container.go`
  - Docker metrics: container list, CPU, memory, network, disk
  - K8s metrics: pod list, status, resource usage (if in K8s)
  - Add ContainerCmd, DockerInfoCmd, K8sInfoCmd
  - Add container tab or integrate into System tab

  **Must NOT do**:
  - Don't change existing tabs or UI
  - Don't break existing providers

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: New provider with external API integration
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Tasks 18, 19, 20)
  - **Blocks**: Task 21
  - **Blocked By**: Task 16

  **References**:
  - `internal/provider/system/hardware.go` - Hardware provider pattern
  - Docker Engine API: https://docs.docker.com/engine/api/

  **Acceptance Criteria**:
  - [ ] Container provider created
  - [ ] Docker containers listed when Docker available
  - [ ] K8s pods listed when in K8s cluster
  - [ ] Graceful degradation when not available

- [ ] 18. Kubernetes Metrics Provider

  **What to do**:
  - Extend container.go for K8s
  - Use kubeclient or kubectl for metrics
  - Show: pods, nodes, services, deployments
  - Integrate into container view

  **Must NOT do**:
  - Don't break Docker functionality

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: K8s API integration
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Tasks 17, 19, 20)
  - **Blocks**: Task 21
  - **Blocked By**: Task 16

  **Acceptance Criteria**:
  - [ ] K8s pod metrics displayed
  - [ ] Node metrics displayed
  - [ ] Works alongside Docker

- [ ] 19. VM/Hypervisor Metrics Provider

  **What to do**:
  - Create `internal/provider/system/vm.go`
  - Detect hypervisors: VMware, VirtualBox, KVM, Hyper-V, etc.
  - VM metrics: CPU, memory, disk allocated
  - Use libvirt API where available

  **Must NOT do**:
  - Don't break existing functionality

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: Platform-specific VM detection
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Tasks 17, 18, 20)
  - **Blocks**: Task 21
  - **Blocked By**: Task 16

  **Acceptance Criteria**:
  - [ ] VM provider created
  - [ ] Detects common hypervisors
  - [ ] Shows VM metrics when hypervisor detected
  - [ ] Graceful degradation when not a VM

- [ ] 20. Full Remote Dashboard

  **What to do**:
  - Extend remote/ssh.go to run all metric commands
  - Support: CPU, memory, disk, network, processes remotely
  - Add SSH key authentication support
  - Multiple hosts with aggregated/split view

  **Must NOT do**:
  - Don't change local-only tabs
  - Keep backward compat with existing remote

  **Recommended Agent Profile**:
  - **Category**: deep
    - Reason: Complex remote architecture
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Tasks 17, 18, 19)
  - **Blocks**: Task 21
  - **Blocked By**: Task 16

  **References**:
  - `internal/provider/remote/ssh.go` - Current implementation

  **Acceptance Criteria**:
  - [ ] Remote tab shows full metrics (not just uptime)
  - [ ] SSH key auth supported
  - [ ] Multiple hosts supported
  - [ ] Works alongside local-only tabs

- [ ] 21. Feature Integration Tests

  **What to do**:
  - Add tests for all new features
  - Test graceful degradation
  - Test error handling

  **Must NOT do**:
  - Don't modify existing tests

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Testing existing features
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (final)
  - **Blocks**: Tasks 22, 23, 24
  - **Blocked By**: Tasks 17, 18, 19, 20

  **Acceptance Criteria**:
  - [ ] Container tests added
  - [ ] VM tests added
  - [ ] Remote tests added
  - [ ] go test ./... → PASS

---

### Phase 6: Tooling & Polish

- [ ] 22. golangci-lint Configuration

  **What to do**:
  - Create `.golangci.toml`
  - Enable recommended linters: errcheck, gosimple, govet, ineffassign, staticcheck, unused
  - Configure for Go 1.25
  - Fix any lint errors

  **Must NOT do**:
  - Don't disable linters to make them pass

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Configuration file creation
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Tasks 23, 24)
  - **Blocks**: Task 25
  - **Blocked By**: Task 21

  **References**:
  - golangci-lint documentation: https://golangci-lint.run/

  **Acceptance Criteria**:
  - [ ] .golangci.toml created
  - [ ] golangci-lint run ./... → zero warnings
  - [ ] CI runs lint check

- [ ] 23. GitHub Actions CI/CD Pipeline

  **What to do**:
  - Create `.github/workflows/ci.yml`
  - Run: go test, golangci-lint, go build
  - Test on multiple Go versions
  - Test on multiple platforms (Linux, macOS, Windows)

  **Must NOT do**:
  - Don't break local workflow

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: CI configuration
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Tasks 22, 24)
  - **Blocks**: Task 25
  - **Blocked By**: Task 21

  **Acceptance Criteria**:
  - [ ] CI workflow created
  - [ ] Tests run on PR
  - [ ] Lint runs on PR
  - [ ] Build verified on PR

- [ ] 24. Fix Makefile

  **What to do**:
  - Replace $(eval) with proper variable assignment
  - Add .PHONY targets properly
  - Add proper dependency tracking

  **Must NOT do**:
  - Don't change build outputs

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Makefile fix
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Tasks 22, 23)
  - **Blocks**: Task 25
  - **Blocked By**: Task 21

  **Acceptance Criteria**:
  - [ ] No $(eval) in Makefile
  - [ ] All targets work
  - [ ] make help shows targets

- [ ] 25. Final Integration Tests

  **What to do**:
  - Run complete test suite
  - Verify backward compatibility
  - Document any issues

  **Must NOT do**:
  - Don't introduce new issues

  **Recommended Agent Profile**:
  - **Category**: quick
    - Reason: Final verification
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 6 (final)
  - **Blocks**: None
  - **Blocked By**: Tasks 22, 23, 24

  **Acceptance Criteria**:
  - [ ] go test ./... → PASS
  - [ ] golangci-lint → zero warnings
  - [ ] go build → success
  - [ ] Backward compat verified

---

## Future Enhancements (Out of Scope)

These are documented here for future reference but NOT included in this plan:

1. **Plugin System**: External script execution for custom metrics
   - Implement as separate phase after this refactoring
   - Requires careful scope definition to avoid inflation

2. **WebSocket Remote**: Real-time remote updates instead of SSH polling
   - Would require significant architecture change
   - Consider after remote dashboard is stable

3. **Additional Themes**: Theme marketplace or user-submitted themes
   - Consider data-driven theme loading

4. **Metrics Export**: Prometheus/InfluxDB export
   - Could leverage new provider interfaces

---

## Success Criteria

### Verification Commands
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Lint check
golangci-lint run ./...

# Build
go build -o bub ./cmd/bub

# Verify backward compatibility
# (Load existing config file - config format unchanged)
```

### Final Checklist
- [ ] All "Must Have" present
- [ ] All "Must NOT Have" absent
- [ ] All tests pass
- [ ] golangci-lint zero warnings
- [ ] UI identical (manual verification)
- [ ] Config backward compatible
- [ ] All new features functional
- [ ] Phase 1-6 complete (25 tasks done)
