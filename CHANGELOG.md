# Changelog

All notable changes to this project are documented in this file.

## [0.2.0] - 2026-02-27

### Added

- Added comprehensive `Check` validation failures across resources, with targeted property paths and actionable error messages:
  - `Network`: non-empty `name`.
  - `Peer`: non-empty `name`.
  - `Group`: non-empty `name`, non-empty `peers[*]`, non-empty `resources[*].id`.
  - `DNS`: non-empty `name`, valid `primary`/`domains` relationship, `searchDomainsEnabled` consistency, required nameservers, nameserver `ip` and `port` validation.
  - `NetworkResource`: non-empty `name`, `networkID`, `address`, and non-empty `groupIDs[*]`.
  - `NetworkRouter`: non-empty `networkID`, non-empty `peer` when provided, non-empty `peerGroups[*]`, non-negative `metric`, and requirement that either `peer` or `peerGroups` is set.
  - `SetupKey`: non-empty `name`, valid `type` (`reusable` or `one-off`), non-negative `expiresIn` and `usageLimit`, non-empty `autoGroups[*]`.
  - `Policy`: non-empty `name`, at least one rule, non-empty rule names, source/destination presence checks, non-empty ports, valid port ranges, and non-empty IDs/types for `sourceResource` and `destinationResource`.
  - `User`: non-empty `role`, required `email` for non-service users, and non-empty `name` when provided.
- Added shared helper utilities:
  - `isBlank(string)` for consistent whitespace-aware string validation.
  - `equalResourcesPtr(*[]Resource, *[]Resource)` for stable resource list comparison in diffs.

### Changed

- Bumped provider and schema version from `0.1.4` to `0.2.0`.
- Normalized several Pulumi schema property names to camelCase for consistency:
  - `search_domains_enabled` -> `searchDomainsEnabled`
  - `network_id` -> `networkID`
  - `group_ids` -> `groupIDs`
  - `peer_groups` -> `peerGroups`
  - `posture_checks` -> `postureChecks`
  - `is_service_user` -> `isServiceUser`
  - `auto_groups` -> `autoGroups`
- Updated `User` schema contract:
  - `email`, `name`, `autoGroups`, and `blocked` are optional inputs.
  - `requiredInputs` now requires `role` (and `isServiceUser` remains required in state schema).
- Regenerated Go SDK to align with schema updates:
  - Renamed generated fields/accessors to camelCase equivalents.
  - Updated plugin/default SDK version to `0.2.0`.
  - Updated `User` generated types to pointer-capable variants where values are optional (e.g., `StringPtr`, `BoolPtr`).

### Fixed

- Corrected `DetailedDiff` semantics across resources:
  - Removed incorrect `InputDiff: true` usage where state-backed updates are expected.
  - Corrected property keys in diffs to match schema names after renames.
- Corrected immutability and replacement behavior:
  - `NetworkResource.networkID` now diffed as `UpdateReplace`.
  - `NetworkRouter.networkID` now diffed as `UpdateReplace`.
  - `SetupKey` immutable fields (`name`, `type`, `expiresIn`, `usageLimit`, `ephemeral`, `allowExtraDnsLabels`) now diff as `UpdateReplace`.
  - `User` immutable creation-only fields (`name`, `email`, `isServiceUser`) now diff as `UpdateReplace`.
- Added full `SetupKey` `Diff` implementation to avoid missing/incorrect replacement signaling.
- Fixed `User` read-path nil handling for API `IsServiceUser` to avoid nil pointer dereference.
- Fixed `User` diff comparisons:
  - Pointer-aware comparison for `name`/`email`.
  - Corrected `autoGroups` change detection logic.
  - Bool-pointer-safe comparison for `blocked`.
- Added guard in `Policy.Create` to fail fast if API response is missing policy ID.

### Breaking Changes

- Pulumi property names changed from snake_case to camelCase for several resource fields. Existing programs using old property names must be updated.
- Go SDK accessor and argument names changed to match renamed properties (for example, `Network_id()` -> `NetworkID()`, `Peer_groups` -> `PeerGroups`).
- `User` output/input type signatures changed for optional values (pointer-based inputs/outputs in generated SDK).
- Some updates that were previously treated as in-place now trigger replacement where fields are immutable by API behavior (`SetupKey`, `User` creation-only fields, router/resource `networkID`).

### Migration Notes

- Update Pulumi programs to use renamed camelCase fields listed above.
- For Go consumers, update renamed methods/fields and pointer-based optional arguments in `User`.
- Review workflows that update immutable fields:
  - Expect replacement for `SetupKey` immutable attributes.
  - Expect replacement for `User` `name`/`email`/`isServiceUser`.
  - Expect replacement when changing `networkID` on `NetworkResource` or `NetworkRouter`.
- If stack state still contains old snake_case keys, refresh/import or re-apply with updated program definitions to reconcile schema changes.
