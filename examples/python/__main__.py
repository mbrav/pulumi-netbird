from pulumi_netbird import resource

# ── Groups ────────────────────────────────────────────────────────────────────
# Peer groups used to scope policies, network resources, and DNS zones.

group_devops = resource.Group("group-devops", name="DevOps", peers=[])
group_dev = resource.Group("group-dev", name="Dev", peers=[])
group_backoffice = resource.Group("group-backoffice", name="Backoffice", peers=[])
group_hr = resource.Group("group-hr", name="HR", peers=[])

# ── Networks ──────────────────────────────────────────────────────────────────
# Overlay network for Region 1 that groups related subnets and routers.

net_r1 = resource.Network(
    "net-r1",
    name="R1",
    description="Network for Region 1",
)

# ── Network Resources ─────────────────────────────────────────────────────────
# Subnet resources attached to net_r1; accessible by the DevOps group.

resource.NetworkResource(
    "netres-r1-net-01",
    name="Region 1 Net 01",
    description="Network 01 in Region 1",
    network_id=net_r1.id,
    address="10.10.1.0/24",
    enabled=True,
    group_ids=[group_devops.id],
)

netres_r1_net02 = resource.NetworkResource(
    "netres-r1-net-02",
    name="Region 1 Net 02",
    description="Network 02 in Region 1",
    network_id=net_r1.id,
    address="10.10.2.0/24",
    enabled=True,
    group_ids=[group_devops.id],
)

resource.NetworkResource(
    "netres-r1-net-03",
    name="Region 1 Net 03",
    description="Network 03 in Region 1",
    network_id=net_r1.id,
    address="10.10.3.0/24",
    enabled=True,
    group_ids=[group_devops.id],
)

# ── Network Router ────────────────────────────────────────────────────────────
# Masquerading router for net_r1; uses DevOps group as peer group.

resource.NetworkRouter(
    "router-r1",
    network_id=net_r1.id,
    enabled=True,
    masquerade=True,
    metric=10,
    peer_groups=[group_devops.id],
)

# ── Posture Check ─────────────────────────────────────────────────────────────
# Validates peer properties before granting policy access.
# Combines: client version, OS version, geo location, network range, and process.

posture_devops = resource.PostureCheck(
    "posture-devops",
    name="DevOps Posture",
    description="Enforce client version, OS, geo, network range, and process checks",
    checks=resource.PostureChecksConfigArgs(
        # Require minimum NetBird client version 0.28.0.
        nb_version_check=resource.PostureMinVersionCheckArgs(
            min_version="0.28.0",
        ),
        # Require minimum OS versions per platform.
        os_version_check=resource.PostureOSVersionCheckArgs(
            darwin=resource.PostureMinVersionCheckArgs(min_version="13.0"),
            linux=resource.PostureMinKernelVersionCheckArgs(min_kernel_version="5.15"),
            windows=resource.PostureMinKernelVersionCheckArgs(min_kernel_version="10.0"),
        ),
        # Allow peers only from Germany (DE) or United States (US).
        geo_location_check=resource.PostureGeoLocationCheckArgs(
            action=resource.PostureGeoLocationAction.ALLOW,
            locations=[
                resource.PostureLocationArgs(country_code="DE"),
                resource.PostureLocationArgs(country_code="US", city_name="New York"),
            ],
        ),
        # Deny peers whose local network is a private RFC-1918 /8 range.
        peer_network_range_check=resource.PosturePeerNetworkRangeCheckArgs(
            action=resource.PosturePeerNetworkRangeAction.DENY,
            ranges=["10.0.0.0/8"],
        ),
        # Require the NetBird agent binary to be present on each platform.
        process_check=resource.PostureProcessCheckArgs(
            processes=[
                resource.PostureProcessArgs(
                    linux_path="/usr/bin/netbird",
                    mac_path="/usr/local/bin/netbird",
                    windows_path=r"C:\Program Files\NetBird\netbird.exe",
                ),
            ],
        ),
    ),
)

# ── Policies ──────────────────────────────────────────────────────────────────
# Policy: DevOps/Dev → Region 1 Net 02 (subnet destination, with posture check).

resource.Policy(
    "policy-ssh-grp-src-net-dest",
    name="SSH Policy - Group to Subnet",
    description="Allow SSH (22/TCP) from DevOps and Dev groups to Region 1 Net 02",
    enabled=True,
    posture_checks=[posture_devops.id],
    rules=[
        resource.PolicyRuleArgsArgs(
            name="SSH Access - Group → Subnet",
            description="Allow unidirectional SSH from DevOps & Dev groups to Net 02",
            bidirectional=False,
            action=resource.RuleAction.ACCEPT,
            enabled=True,
            protocol=resource.Protocol.TCP,
            ports=["22"],
            sources=[group_devops.id, group_dev.id],
            destination_resource=resource.ResourceArgs(
                type=resource.Type.SUBNET,
                id=netres_r1_net02.id,
            ),
        ),
    ],
)

# Policy: DevOps → Backoffice (group-to-group destination, no posture).

resource.Policy(
    "policy-ssh-grp-src-grp-dest",
    name="SSH Policy - Group to Group",
    description="Allow SSH (22/TCP) from DevOps to Backoffice group resources",
    enabled=True,
    posture_checks=[],
    rules=[
        resource.PolicyRuleArgsArgs(
            name="SSH Access - Group → Group",
            description="SSH from DevOps group to Backoffice group",
            bidirectional=False,
            action=resource.RuleAction.ACCEPT,
            enabled=True,
            protocol=resource.Protocol.TCP,
            ports=["22"],
            sources=[group_devops.id],
            destinations=[group_backoffice.id],
        ),
    ],
)

# ── DNS Settings ──────────────────────────────────────────────────────────────
# Singleton resource — only one exists per account.
# disabled_management_groups: groups whose peers resolve DNS outside NetBird.

resource.DNSSettings(
    "dns-settings",
    disabled_management_groups=[],
)

# ── DNS Zone ──────────────────────────────────────────────────────────────────
# Internal DNS zone for corp.example.com distributed to DevOps and Dev groups.
# enable_search_domain: allows bare hostnames to be resolved within the zone.

dns_zone = resource.DNSZone(
    "dns-zone-corp",
    name="corp-internal",
    domain="corp.example.com",
    enabled=True,
    enable_search_domain=True,
    distribution_groups=[group_devops.id, group_dev.id],
)

# ── DNS Records ───────────────────────────────────────────────────────────────
# A record pointing the gateway hostname to an internal IPv4 address.

resource.DNSRecord(
    "dns-record-gw-a",
    zone_id=dns_zone.id,
    name="gw.corp.example.com",
    type=resource.DNSRecordType.A,
    content="10.10.1.1",
    ttl=300,
)

# CNAME record aliasing api.corp.example.com to the gateway hostname.
resource.DNSRecord(
    "dns-record-api-cname",
    zone_id=dns_zone.id,
    name="api.corp.example.com",
    type=resource.DNSRecordType.CNAME,
    content="gw.corp.example.com",
    ttl=300,
)

# ── Setup Key ─────────────────────────────────────────────────────────────────
# Reusable setup key for onboarding new peers into the DevOps group.
# usage_limit=0: unlimited uses; expires_in=0: no expiry.

resource.SetupKey(
    "setup-key-devops",
    name="DevOps Onboarding",
    type=resource.SetupKeyType.REUSABLE,
    expires_in=0,
    usage_limit=0,
    ephemeral=False,
    allow_extra_dns_labels=False,
    auto_groups=[group_devops.id],
)

# ── Route ─────────────────────────────────────────────────────────────────────
# Network route advertising 192.168.10.0/24 through the DevOps peer group.
# masquerade=True hides the source IP behind the router's address.

resource.Route(
    "route-r1-mgmt",
    network_id="route-r1-mgmt",
    description="Management subnet route via Region 1",
    enabled=True,
    network="192.168.10.0/24",
    masquerade=True,
    metric=100,
    keep_route=True,
    groups=[group_devops.id],
    peer_groups=[group_devops.id],
)

# ── Service User ──────────────────────────────────────────────────────────────
# Automation service user with admin role; placed in the DevOps group.

resource.User(
    "user-ci-bot",
    role="admin",
    is_service_user=True,
    name="ci-bot",
    auto_groups=[group_devops.id],
)

# ── Reverse Proxy Domain ──────────────────────────────────────────────────────
# Custom domain validated against a specific proxy cluster.
# Domain changes trigger resource replacement (no Update endpoint exists).

rp_domain = resource.ReverseProxyDomain(
    "rp-domain-corp",
    domain="proxy.corp.example.com",
    target_cluster="eu-central-1",
)

# ── Reverse Proxy Service ─────────────────────────────────────────────────────
# HTTP (L7) reverse proxy routing traffic to an internal backend host.
# pass_host_header: preserves the original Host header forwarded to the backend.
# rewrite_redirects: rewrites Location headers in backend responses.

resource.ReverseProxyService(
    "rp-svc-api",
    name="api-service",
    domain=rp_domain.domain,
    enabled=True,
    mode=resource.ReverseProxyServiceMode.HTTP,
    pass_host_header=True,
    rewrite_redirects=False,
    targets=[
        resource.ReverseProxyTargetArgs(
            enabled=True,
            host="10.10.1.10",
            port=8080,
            protocol=resource.ReverseProxyTargetProtocol.HTTP,
            target_type=resource.ReverseProxyTargetType.HOST,
            target_id="",
        ),
    ],
)
