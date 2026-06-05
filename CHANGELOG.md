# Changelog

All notable changes to this project are documented in this file.

## [0.3.4] - 2026-06-05

### Added

- Added the `Route` resource (`netbird:resource:Route`) for managing NetBird network routes through routing peers or peer groups.
- Added generated Go SDK support for `Route`, including `NewRoute`, `GetRoute`, route inputs/outputs, and package construction wiring.

### Changed

- Bumped provider, schema, and Go SDK metadata from `0.3.3` to `0.3.4`.
- Updated the Go example module to consume SDK `v0.3.3`.

## [0.3.3] - 2026-06-05

### Fixed

- Normalized ordering for API-returned slice fields to reduce false positive diffs after refresh/import:
  - `DNS.domains` and `DNS.groups`
  - `Policy.postureChecks`
  - policy rule source and destination groups
  - `User.autoGroups`
- Treated nil and empty peer/resource lists as equivalent in `Group` diffs.
- Avoided tracking externally-populated `Group.resources` unless resources are explicitly declared in the Pulumi inputs.

### Changed

- Bumped provider, schema, and Go SDK metadata from `0.3.2` to `0.3.3`.
- Updated the Go example module to consume SDK `v0.3.2`.

## [0.3.2] - 2026-06-05

### Fixed

- Fixed `NetworkRouter` and `NetworkResource` imports by supporting compound import IDs in the form `<networkID>/<routerID>` and `<networkID>/<resourceID>`.
- Fixed `Peer.approvalRequired` and peer management booleans so imported peers no longer require cloud-only or client-controlled fields in Pulumi programs.
- Fixed persistent `Policy.rules` diffs caused by nil-vs-empty optional slices and API-generated rule descriptions.
- Fixed policy imports by reconstructing rule inputs from the NetBird API when stored Pulumi inputs are empty.
- Reduced false positive description diffs for `Network`, `NetworkResource`, and `Policy` when NetBird returns API-generated descriptions.

### Changed

- `Peer.Create` now returns an explicit error because peers must be imported and cannot be created through the NetBird management API.
- Documented compound import ID formats in generated `NetworkRouter` and `NetworkResource` resource descriptions.
- Bumped provider, schema, and Go SDK metadata from `0.3.1` to `0.3.2`.
- Regenerated the Go SDK so optional peer fields and updated resource descriptions are reflected for Go consumers.

## [0.3.1] - 2026-06-05

### Fixed

- Added `replace` directives in `go.mod` to redirect `github.com/dexidp/dex` and `github.com/dexidp/dex/api/v2` to NetBird's own fork (`github.com/netbirdio/dex`). Go replace directives from dependencies do not propagate to dependent modules, causing `go mod tidy` to fail with "module does not contain package `github.com/dexidp/dex/server/signer`".
- Added missing `Ipv6` field (`nil`) to `PeerRequest` struct literal in `peer.go` to satisfy the `exhaustruct` linter.

### Changed

- Bumped provider version from `0.3.0` to `0.3.1`.
- Updated CI actions: `pulumi/actions` v6 → v7, `softprops/action-gh-release` v2 → v3.
- Updated golangci-lint container image to `v2.12.2`.
- Disabled `goconst` and `gomoddirectives` linters in `.golangci.yml` to allow the dex replace override and repeated string literals.

## [0.3.0] - 2026-04-04

> **Note:** A minor version bump (0.2.x → 0.3.0) was necessary due to the large number of new resources and the scope of internal changes introduced in this release.

### Added

- **`DNSRecord`** (`netbird_dns_record`) — manage DNS records within a DNS zone. Supports A, AAAA, and CNAME record types. Fields: `zoneID`, `name`, `content`, `ttl`, `type`.
- **`DNSSettings`** (`netbird_dns_settings`) — singleton resource managing global DNS settings. Fields: `disabledManagementGroups`. Create/Update call the settings endpoint; Delete is a no-op since global settings cannot be removed.
- **`DNSZone`** (`netbird_dns_zone`) — manage DNS zones. Fields: `name`, `domain`, `enabled`, `enableSearchDomain`, `distributionGroups`. Domain changes trigger resource replacement.
- **`PostureCheck`** (`netbird_posture_check`) — manage peer posture checks used by policies. Supports all check types:
  - `geoLocationCheck` — allow/deny by country code and city, with `PostureGeoLocationAction` enum.
  - `nbVersionCheck` — minimum NetBird client version (`minVersion`).
  - `osVersionCheck` — per-platform minimum versions: `android`, `darwin`, `ios` (minVersion); `linux`, `windows` (minKernelVersion).
  - `peerNetworkRangeCheck` — allow/deny by CIDR range list, with `PosturePeerNetworkRangeAction` enum.
  - `processCheck` — required running processes per platform (`linuxPath`, `macPath`, `windowsPath`).
- **`ReverseProxyDomain`** (`netbird_reverse_proxy_domain`) — manage custom reverse proxy domains. Fields: `domain`, `targetCluster`. Read-only outputs: `type`, `validated`, `requireSubdomain`, `supportsCustomPorts`. No Update endpoint exists — all input changes trigger replacement.
- **`ReverseProxyService`** (`netbird_reverse_proxy_service`) — manage reverse proxy services. Fields: `name`, `domain`, `enabled`, `mode`, `targets[]`, `passHostHeader`, `rewriteRedirects`, `listenPort`. Read-only outputs: `proxyCluster`, `status`. New enums: `ReverseProxyServiceMode` (http/tcp/tls/udp), `ReverseProxyTargetProtocol` (http/https/tcp/udp), `ReverseProxyTargetType` (domain/host/peer/subnet). Domain changes trigger replacement.
- New `DNSRecordType` enum: `A`, `AAAA`, `CNAME`.
- New `ReverseProxyDomainType` enum: `custom`, `free`.

### Changed

- Bumped provider version from `0.2.0` to `0.3.0`.
- Refactored `equalPtr` in `util.go` from a string-only function to a generic `equalPtr[T comparable]`, replacing the separate `equalBoolPtr`, `equalIntPtr`, `equalStringPtr`, and `equalReverseProxyServiceModePtr` helpers.
- OS version check conversion logic in `PostureCheck` extracted into dedicated `toAPIOSVersionCheck` / `fromAPIOSVersionCheck` helpers to reduce nesting complexity.

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
