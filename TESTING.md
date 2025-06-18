# TESTING.md

## Testing Philosophy

Our goal is **100% coverage** with **meaningful, robust tests** that ensure the reliability and correctness of our codebase.

We use **mutation testing** as a powerful technique to measure the effectiveness of our tests. Mutation testing helps us verify that tests not only execute code but also catch subtle logical errors and regressions.

---

## Coverage Expectations & Thresholds

- Aim for **100% line and branch coverage** in unit tests.
- Mutation testing thresholds:
  - **Test Efficacy:** 100% — all covered mutants should be caught by tests.
  - **Mutation Coverage:** as close as possible to 100%, with an allowance for justified exceptions.
- We regularly run mutation tests and update coverage targets as the project evolves.

---

## When to Write Unit Tests

- Cover all **core logic paths** and common error conditions.
- Write tests that exercise:
  - Normal operation
  - Boundary conditions
  - Expected failure modes that can be easily simulated

---

## Exceptions & Pragmatic Decisions

Some rare or complex error paths are **not worth over-engineering in unit tests**, because:

- They are difficult or impractical to simulate accurately (e.g., `Flush` failures in buffered IO).
- Attempting to mock these edge cases often results in fragile, unreadable tests.
- The risk/benefit ratio favors testing these in **functional or integration tests** with real environments or external dependencies.

Such exceptions are documented in code comments or this document and revisited periodically.

---

## Functional & Integration Testing

- Complex edge cases, real-world error conditions, and environment-specific behaviors belong in functional tests.
- We leverage real VMs, live network scenarios, or external dependencies to validate system-level robustness.

---

## Continuous Integration

- Mutation testing and coverage analysis are integrated in CI pipelines.
- Tests must pass with coverage and efficacy above configured thresholds before code merges.
- Exceptions and known gaps must be acknowledged and tracked.

---

## Summary

We prioritize **test quality over quantity**, focusing on tests that provide true value and confidence.

By combining:

- High coverage,
- Mutation testing,
- Pragmatic exceptions,
- Functional tests for hard-to-simulate cases,

We maintain a solid and maintainable testing foundation for long-term project health.
