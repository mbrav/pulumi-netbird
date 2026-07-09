# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

## [0.5.1] - 2026-07-09

### Changed

- Reverted the v0.5.0 breaking change: `Group` again accepts optional `peers` and `resources` inputs, so group membership and resource associations can be managed directly from the `Group` resource.
- Bumped provider version from `0.5.0` to `0.5.1`.

### Fixed

- `Group` resource and peer lists are now written to state in a stable sorted order (resources by type then ID, peers alphabetically) from `Create`, `Update`, and `Read` alike. Previously only `Read` sorted resources while `Create`/`Update` stored them in raw API order, and `Update` never sorted peers — the nondeterministic NetBird API ordering produced spurious `resources`/`peers` diffs on every `pulumi refresh`/`preview`.
- Restored the `toAPIResourceList` and `equalResourcesPtr` helpers in `common.go` that back the restored inputs (order-insensitive resource comparison in `Diff`).

## [0.5.0] - 2026-06-26

### Changed

- **Breaking:** `Group` no longer accepts `peers` or `resources` as inputs. Group membership remains visible as read-only outputs, but should be managed from the resources that own the assignment to avoid drift and destroy-order failures.
- Bumped provider version from `0.4.1` to `0.5.0`.

### Fixed

- Pinned `github.com/pulumi/pulumi/pkg/v3` and `github.com/pulumi/pulumi/sdk/v3` to `v3.232.0` (down from `v3.248.0`). `pulumi-go-provider v1.3.2` is built against Pulumi `v3.232.0`, and `v3.248.0` changed the schema codegen `pkg.Provider` field from `schema.ResourceSpec` to `*schema.ResourceSpec`, which broke compilation of the go-provider dependency and surfaced as a `golangci-lint` typecheck error (`could not import .../infer`). Keep these two modules in lockstep with whatever `pulumi-go-provider` requires.
- Cleaned up lint fallout from the `Group` input change: added explicit `nil` `Peers`/`Resources` fields to the `nbapi.GroupRequest` and `GroupState` struct literals in `group.go` (`exhaustruct`), and removed the now-unused `toAPIResourceList` and `equalResourcesPtr` helpers from `common.go` (`unused`).

## [0.4.1] - 2026-06-08

### Added

- **Experimental components** (proof of concept) — two composite resources that bundle multiple related NetBird resources into a single Pulumi declaration. These are an exploration of the `pulumi-go-provider` component API; the interface may change without notice and they should not be relied upon in production.
  - **`netbird:component:NetworkBundle`** — declares a `Network`, a `NetworkRouter`, and one `NetworkResource` (subnet) per entry in `subnets[]` as a single unit. The `networkID` is wired automatically between all child resources. Inputs: `name`, `description?`, `router` (`enabled`, `masquerade`, `metric`, `peerGroups?`, `peer?`), `subnets[]` (`name`, `address`, `enabled`, `groupIDs`, `description?`). Outputs: `networkId`, `routerId`, `subnetIds[]`.
  - **`netbird:component:DNSZoneBundle`** — declares a `DNSZone` and one `DNSRecord` per entry in `records[]` as a single unit. The `zoneID` is wired automatically into each record. Inputs: `name`, `domain`, `enabled`, `enableSearchDomain`, `distributionGroups[]`, `records[]` (`name`, `type`, `content`, `ttl`). Outputs: `zoneId`, `recordIds[]`.
- **`provider/component/`** package with `networkBundle.go`, `dnsZoneBundle.go`, and `all.go`. Components are registered via `infer.Component` . Child resources use `var child pulumi.CustomResourceState` directly — no wrapper struct needed.

### Changed

- Bumped provider version from `0.4.0` to `0.4.1`.
- Regenerated Go SDK (`sdk/go/netbird/component/`) to expose `NewNetworkBundle` and `NewDNSZoneBundle` with typed `Args` and output structs.

## [0.4.0] - 2026-06-08

### Added

- **Invoke functions** (data sources) — 6 read-only provider functions that query live NetBird state without managing resources. Useful for cross-stack references and referencing objects that exist outside the current Pulumi program:
  - **`netbird:function:getPeers`** — list all peers in the account, with an optional `groupId` filter to return only peers belonging to a specific group. Returns `peers[]` with `peerId`, `name`, `ip`, `dnsLabel`, `connected`, `hostname`, and `groups[]`.
  - **`netbird:function:lookupGroup`** — look up a group by `name`. Returns `groupId`, `peersCount`, `resourcesCount`, `peers[]`, and `resources[]`. Useful for referencing groups not managed by Pulumi (e.g. the built-in "All" group).
  - **`netbird:function:lookupPeer`** — look up a peer by `name`. Returns `peerId`, `ip`, `dnsLabel`, `connected`, `hostname`, `os`, and `groups[]`. Peers cannot be created via the API, so this is the standard way to resolve a peer ID from its hostname.
  - **`netbird:function:lookupRoute`** — look up the first route whose `network` CIDR matches the input. Returns `routeId`, `description`, `network`, `domains[]`, `enabled`, `masquerade`, `metric`, `peer`, `peerGroups[]`, and `groups[]`.
  - **`netbird:function:lookupSetupKey`** — look up a setup key by `name`. Returns `setupKeyId`, `type`, `state`, `revoked`, `ephemeral`, `usageLimit`, `autoGroups[]`, `expires`, and `lastUsed`.
  - **`netbird:function:lookupUser`** — look up a user by `email`. Returns `userId`, `name`, `email`, `role`, `isBlocked`, and `autoGroups[]`.
- **`provider/resource/all.go`** — `resource.All()` helper that returns all 16 registered resources as `[]infer.InferredResource`. `provider.go` now calls `WithResources(resource.All()...)`, keeping the registration list in one place.
- **`provider/function/all.go`** — `function.All()` helper that returns all 6 functions as `[]infer.InferredFunction`, used by `WithFunctions(function.All()...)` in `provider.go`.
- Function examples added to `examples/yaml/Pulumi.yaml`: three `variables` blocks using `fn::invoke` for `lookupGroup` (All group), `getPeers` (all peers), and `getPeers` filtered by `${group-devops.id}`, with results exported as stack outputs.
- Function examples added to `examples/go/main.go`: `function.LookupGroup`, `function.GetPeers` (unfiltered), and `function.GetPeers` filtered by the resolved All-group ID, with peer counts exported as stack outputs.

### Changed

- `provider.go` refactored to use `resource.All()...` and `function.All()...` spreads instead of listing every resource and function inline. Adding a new resource or function now requires editing only its registration file.
- Bumped provider, schema, and Go SDK metadata from `0.3.8` to `0.4.0`.
- `examples/go/go.mod` updated to require `github.com/mbrav/pulumi-netbird/sdk v0.4.0` with a local `replace` directive (`../../sdk`) for development builds prior to the published release.
- Regenerated Go SDK (`sdk/go/netbird/function/`) to expose all six functions with typed `Args`, `Result`, `OutputArgs`, and `ResultOutput` variants for each.

## [0.3.8] - 2026-06-06

### Changed

- Bumped `github.com/netbirdio/netbird` dependency from `v0.72.0` to `v0.72.1`.
- Bumped provider, schema, and Go SDK metadata from `0.3.7` to `0.3.8`.

## [0.3.7] - 2026-06-05

### Added

- Added `private` and `accessGroups` fields to the `ReverseProxyService` resource, tracking new fields introduced in NetBird v0.72.0's `ServiceRequest` / `Service` API types.
  - `private` (`*bool`, optional) — when `true`, the service is NetBird-only: inbound peers authenticate via WireGuard tunnel identity and an ACL policy is auto-generated from `accessGroups`. Requires `mode=http`. Mutually exclusive with SSO/bearer auth.
  - `accessGroups` (`*[]string`, optional) — NetBird group IDs whose peers may reach this private service over the tunnel. Required when `private=true`; ignored otherwise.
- Added `authorizedGroups` (`*map[string][]string`, optional) to `PolicyRule` inputs and state: a map of NetBird group IDs to lists of local users for network access authorization. Sent to the API on every `Create` and `Update`, round-tripped through `Read`, and included in the per-rule `Diff` comparison via a new `equalMapStringSlice` helper.

### Changed

- Regenerated Go SDK to expose `Private`, `AccessGroups` on `ReverseProxyServiceArgs`/`ReverseProxyServiceState` and `AuthorizedGroups` on `PolicyRuleArgs`/`PolicyRuleState`.

### Fixed

- Added `GOPATH: /home/runner/go` to the lint job and split its cache step into separate module-cache (shared key `linux-go-*`) and lint-analysis-cache (`linux-golangci-lint-*`) entries. Added explicit `go mod download` step so the module cache directory is created before the post-job cache save runs.
- Switched the build job from `actions/setup-go` built-in cache to explicit `actions/cache` steps using the same `linux-go-*` key as the lint job, enabling true module-cache sharing between the two jobs.

## [0.3.6] - 2026-06-05

### Fixed

- Fixed persistent description diffs on `Policy`, `Network`, `NetworkResource`, and `PostureCheck` when the `description` field is not declared in the Pulumi program (inputs.Description is nil). The `Diff` now skips the comparison entirely when `inputs.Description` is nil, and `Read` returns `nil` for `State.Description` in that case — so a refresh always sees nil vs nil and no spurious update is planned. Resources that do declare `description` continue to track it and detect changes normally.

## [0.3.5] - 2026-06-05

### Fixed

- Applied consistent 404 / drift handling across all resources (`Group`, `Network`, `Policy`, `NetworkResource`, `NetworkRouter`, `Route`, `DNS`, `DNSZone`): `Read` now returns an empty ID when the resource is gone so Pulumi removes it from state, and `Delete` treats 404 as success (idempotent destroy).
- Fixed false-positive description diffs in `Network`, `NetworkResource`, and `Policy` caused by a nil guard that prevented detecting when the user removed an optional description field.
- Fixed false-positive `Policy` rule diffs: rule `description` is now included in the per-rule comparison inside `Diff`.
- Fixed `Route.Diff`: `networkId` is now diffed as `UpdateReplace` (changing the parent network requires resource replacement).
- Fixed in-place slice mutation in `NetworkResource` and `DNSZone`: replaced `slices.Sort(req.Inputs.GroupIDs / DistributionGroups)` with `sortedStrings()` so the original input slice is never modified.
- Fixed `NetworkResource.Diff` group comparison: replaced sorted in-place + `slices.Equal` with `equalSlice` (order-insensitive, non-mutating).
- Fixed `NetworkResource.Read` and `NetworkRouter.Read` compound import ID parsing: replaced manual `strings.SplitN` with `parseNestedID`, which validates both parts are non-empty.

### Added

- Added `isNotFoundErr(err error) bool` helper in `util.go`: matches "not found" in the lowercased error message, since the NetBird REST client returns plain `errors.New` strings with no typed 404 sentinel.
- Added `parseNestedID(kind, id string) (string, string, error)` helper in `util.go`: splits `<parentID>/<childID>` compound import IDs and returns a clear error if either part is blank.
- Added `sortedStrings(s []string) []string` helper in `util.go`: returns a sorted clone, leaving the original slice unmodified.
- Moved shared resource-list helpers (`equalResourcesPtr`, `sortedResources`, `compareResources`) from `group.go` into `common.go` so they can be reused across resources.
- Extracted `routeCheckArgs` helper in `route.go` to satisfy `gocognit`/`cyclop` limits; added complete validation: blank `networkId`, `network`/`domains` mutual exclusion, blank entries in `domains`, `groups`, `peerGroups`, `accessControlGroups`, and `peer`/`peerGroups` at-least-one constraint.
- Updated `CLAUDE.md` with an "Established patterns" section covering: the 404 read/delete pattern, optional-field diff guards, safe slice handling, compound import IDs, parent-ID `UpdateReplace`, shared helpers in `common.go`, and Check validation completeness.

## [0.3.4] - 2026-06-05

### Added

- Added the `Route` resource (`netbird:resource:Route`) for managing NetBird network routes through routing peers or peer groups.
- Added generated Go SDK support for `Route`, including `NewRoute`, `GetRoute`, route inputs/outputs, and package construction wiring.

## [0.3.3] - 2026-06-05

### Fixed

- Normalized ordering for API-returned slice fields to reduce false positive diffs after refresh/import:
  - `DNS.domains` and `DNS.groups`
  - `Policy.postureChecks`
  - policy rule source and destination groups
  - `User.autoGroups`
- Treated nil and empty peer/resource lists as equivalent in `Group` diffs.
- Avoided tracking externally-populated `Group.resources` unless resources are explicitly declared in the Pulumi inputs.

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
