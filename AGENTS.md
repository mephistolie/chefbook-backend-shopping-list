# Shopping List Service Agents Guide

This service owns shopping list behavior.

## Scope

- shopping list entities and item lifecycle
- shopping list APIs, persistence, and domain rules

## Working Rules

- Keep shopping list-specific behavior isolated from recipe or profile services unless integration is required by contract.
- Be careful with synchronization semantics and item state transitions.
- Validate consumer impact if list item schemas or events change.
