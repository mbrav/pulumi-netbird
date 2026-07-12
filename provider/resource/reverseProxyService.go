package resource

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ReverseProxyService represents a NetBird reverse proxy service resource.
type ReverseProxyService struct{}

// Annotate adds a description to the ReverseProxyService resource type.
func (r *ReverseProxyService) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r, "A NetBird reverse proxy service.")
}

// ReverseProxyServiceArgs defines input fields for creating or updating a reverse proxy service.
type ReverseProxyServiceArgs struct {
	Name               string                          `pulumi:"name"`
	Domain             string                          `pulumi:"domain"`
	Enabled            bool                            `pulumi:"enabled"`
	Mode               *ReverseProxyServiceMode        `pulumi:"mode,optional"`
	Targets            []ReverseProxyTarget            `pulumi:"targets"`
	PassHostHeader     *bool                           `pulumi:"passHostHeader,optional"`
	RewriteRedirects   *bool                           `pulumi:"rewriteRedirects,optional"`
	ListenPort         *int                            `pulumi:"listenPort,optional"`
	Private            *bool                           `pulumi:"private,optional"`
	AccessGroups       *[]string                       `pulumi:"accessGroups,optional"`
	Auth               *ReverseProxyAuth               `pulumi:"auth,optional"`
	AccessRestrictions *ReverseProxyAccessRestrictions `pulumi:"accessRestrictions,optional"`
}

// Annotate provides documentation for ReverseProxyServiceArgs fields.
func (r *ReverseProxyServiceArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Name, "Service name.")
	annotator.Describe(&r.Domain, "Domain for the service.")
	annotator.Describe(&r.Enabled, "Whether the service is enabled.")
	annotator.Describe(&r.Mode, `Service mode: "http" for L7 reverse proxy, "tcp"/"udp"/"tls" for L4 passthrough.`)
	annotator.Describe(&r.Targets, "List of target backends for this service.")
	annotator.Describe(&r.PassHostHeader, "When true, the original client Host header is passed through to the backend.")
	annotator.Describe(&r.RewriteRedirects, "When true, Location headers in backend responses are rewritten to the public-facing domain.")
	annotator.Describe(&r.ListenPort, "Port the proxy listens on (L4/TLS only). Set to 0 for auto-assignment.")
	annotator.Describe(&r.Private, "When true, the service is NetBird-only: peers authenticate via WireGuard tunnel identity and an ACL policy is auto-generated from accessGroups. Requires mode=http. Mutually exclusive with SSO/bearer auth.")
	annotator.Describe(&r.AccessGroups, "NetBird group IDs whose peers may reach this private service over the tunnel. Required when private=true; ignored otherwise.")
	annotator.Describe(&r.Auth, "Authentication configuration for the service (bearer/header/link/password/pin). Mutually exclusive with private=true.")
	annotator.Describe(&r.AccessRestrictions, "Connection-level access restrictions based on IP address or geography. Applies to both HTTP and L4 services.")
}

// ReverseProxyServiceState represents the output state of a reverse proxy service resource.
type ReverseProxyServiceState struct {
	Name               string                          `pulumi:"name"`
	Domain             string                          `pulumi:"domain"`
	Enabled            bool                            `pulumi:"enabled"`
	Mode               *ReverseProxyServiceMode        `pulumi:"mode,optional"`
	Targets            []ReverseProxyTarget            `pulumi:"targets"`
	PassHostHeader     *bool                           `pulumi:"passHostHeader,optional"`
	RewriteRedirects   *bool                           `pulumi:"rewriteRedirects,optional"`
	ListenPort         *int                            `pulumi:"listenPort,optional"`
	Private            *bool                           `pulumi:"private,optional"`
	AccessGroups       *[]string                       `pulumi:"accessGroups,optional"`
	Auth               *ReverseProxyAuth               `pulumi:"auth,optional"`
	AccessRestrictions *ReverseProxyAccessRestrictions `pulumi:"accessRestrictions,optional"`
	ProxyCluster       *string                         `pulumi:"proxyCluster,optional"`
	Status             *ReverseProxyServiceStatus      `pulumi:"status,optional"`
	Terminated         *bool                           `pulumi:"terminated,optional"`
	PortAutoAssigned   *bool                           `pulumi:"portAutoAssigned,optional"`
}

// Annotate provides documentation for ReverseProxyServiceState fields.
func (r *ReverseProxyServiceState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Name, "Service name.")
	annotator.Describe(&r.Domain, "Domain for the service.")
	annotator.Describe(&r.Enabled, "Whether the service is enabled.")
	annotator.Describe(&r.Mode, `Service mode: "http" for L7 reverse proxy, "tcp"/"udp"/"tls" for L4 passthrough.`)
	annotator.Describe(&r.Targets, "List of target backends for this service.")
	annotator.Describe(&r.PassHostHeader, "When true, the original client Host header is passed through to the backend.")
	annotator.Describe(&r.RewriteRedirects, "When true, Location headers in backend responses are rewritten to the public-facing domain.")
	annotator.Describe(&r.ListenPort, "Port the proxy listens on (L4/TLS only).")
	annotator.Describe(&r.Private, "When true, the service is NetBird-only: peers authenticate via WireGuard tunnel identity and an ACL policy is auto-generated from accessGroups. Requires mode=http. Mutually exclusive with SSO/bearer auth.")
	annotator.Describe(&r.AccessGroups, "NetBird group IDs whose peers may reach this private service over the tunnel. Required when private=true; ignored otherwise.")
	annotator.Describe(&r.Auth, "Authentication configuration for the service.")
	annotator.Describe(&r.AccessRestrictions, "Connection-level access restrictions based on IP address or geography.")
	annotator.Describe(&r.ProxyCluster, "The proxy cluster handling this service (derived from domain).")
	annotator.Describe(&r.Status, "Current status of the service.")
	annotator.Describe(&r.Terminated, "Whether the service has been terminated. Terminated services cannot be updated.")
	annotator.Describe(&r.PortAutoAssigned, "Whether the listen port was auto-assigned.")
}

// ReverseProxyServiceMode defines the allowed service modes.
type ReverseProxyServiceMode string

const (
	// ReverseProxyServiceModeHTTP represents L7 HTTP reverse proxy mode.
	ReverseProxyServiceModeHTTP ReverseProxyServiceMode = ReverseProxyServiceMode(nbapi.ServiceRequestModeHttp)
	// ReverseProxyServiceModeTCP represents L4 TCP passthrough mode.
	ReverseProxyServiceModeTCP ReverseProxyServiceMode = ReverseProxyServiceMode(nbapi.ServiceRequestModeTcp)
	// ReverseProxyServiceModeTLS represents L4 TLS passthrough mode.
	ReverseProxyServiceModeTLS ReverseProxyServiceMode = ReverseProxyServiceMode(nbapi.ServiceRequestModeTls)
	// ReverseProxyServiceModeUDP represents L4 UDP passthrough mode.
	ReverseProxyServiceModeUDP ReverseProxyServiceMode = ReverseProxyServiceMode(nbapi.ServiceRequestModeUdp)
)

// Values returns the valid enum values for ReverseProxyServiceMode.
func (ReverseProxyServiceMode) Values() []infer.EnumValue[ReverseProxyServiceMode] {
	return []infer.EnumValue[ReverseProxyServiceMode]{
		{Name: "http", Value: ReverseProxyServiceModeHTTP, Description: "L7 HTTP reverse proxy mode."},
		{Name: "tcp", Value: ReverseProxyServiceModeTCP, Description: "L4 TCP passthrough mode."},
		{Name: "tls", Value: ReverseProxyServiceModeTLS, Description: "L4 TLS passthrough mode."},
		{Name: "udp", Value: ReverseProxyServiceModeUDP, Description: "L4 UDP passthrough mode."},
	}
}

// ReverseProxyTargetProtocol defines the protocol used to connect to a backend target.
type ReverseProxyTargetProtocol string

const (
	// ReverseProxyTargetProtocolHTTP represents the HTTP protocol.
	ReverseProxyTargetProtocolHTTP ReverseProxyTargetProtocol = ReverseProxyTargetProtocol(nbapi.ServiceTargetProtocolHttp)
	// ReverseProxyTargetProtocolHTTPS represents the HTTPS protocol.
	ReverseProxyTargetProtocolHTTPS ReverseProxyTargetProtocol = ReverseProxyTargetProtocol(nbapi.ServiceTargetProtocolHttps)
	// ReverseProxyTargetProtocolTCP represents the TCP protocol.
	ReverseProxyTargetProtocolTCP ReverseProxyTargetProtocol = ReverseProxyTargetProtocol(nbapi.ServiceTargetProtocolTcp)
	// ReverseProxyTargetProtocolUDP represents the UDP protocol.
	ReverseProxyTargetProtocolUDP ReverseProxyTargetProtocol = ReverseProxyTargetProtocol(nbapi.ServiceTargetProtocolUdp)
)

// Values returns the valid enum values for ReverseProxyTargetProtocol.
func (ReverseProxyTargetProtocol) Values() []infer.EnumValue[ReverseProxyTargetProtocol] {
	return []infer.EnumValue[ReverseProxyTargetProtocol]{
		{Name: "http", Value: ReverseProxyTargetProtocolHTTP, Description: "HTTP protocol."},
		{Name: "https", Value: ReverseProxyTargetProtocolHTTPS, Description: "HTTPS protocol."},
		{Name: "tcp", Value: ReverseProxyTargetProtocolTCP, Description: "TCP protocol."},
		{Name: "udp", Value: ReverseProxyTargetProtocolUDP, Description: "UDP protocol."},
	}
}

// ReverseProxyTargetType defines the type of backend target.
type ReverseProxyTargetType string

const (
	// ReverseProxyTargetTypeCluster represents a proxy-cluster target.
	ReverseProxyTargetTypeCluster ReverseProxyTargetType = ReverseProxyTargetType(nbapi.ServiceTargetTargetTypeCluster)
	// ReverseProxyTargetTypeDomain represents a domain-based target.
	ReverseProxyTargetTypeDomain ReverseProxyTargetType = ReverseProxyTargetType(nbapi.ServiceTargetTargetTypeDomain)
	// ReverseProxyTargetTypeHost represents a host-based target.
	ReverseProxyTargetTypeHost ReverseProxyTargetType = ReverseProxyTargetType(nbapi.ServiceTargetTargetTypeHost)
	// ReverseProxyTargetTypePeer represents a peer-based target.
	ReverseProxyTargetTypePeer ReverseProxyTargetType = ReverseProxyTargetType(nbapi.ServiceTargetTargetTypePeer)
	// ReverseProxyTargetTypeSubnet represents a subnet-based target.
	ReverseProxyTargetTypeSubnet ReverseProxyTargetType = ReverseProxyTargetType(nbapi.ServiceTargetTargetTypeSubnet)
)

// Values returns the valid enum values for ReverseProxyTargetType.
func (ReverseProxyTargetType) Values() []infer.EnumValue[ReverseProxyTargetType] {
	return []infer.EnumValue[ReverseProxyTargetType]{
		{Name: "cluster", Value: ReverseProxyTargetTypeCluster, Description: "Proxy-cluster target."},
		{Name: "domain", Value: ReverseProxyTargetTypeDomain, Description: "Domain-based target."},
		{Name: "host", Value: ReverseProxyTargetTypeHost, Description: "Host-based target."},
		{Name: "peer", Value: ReverseProxyTargetTypePeer, Description: "Peer-based target."},
		{Name: "subnet", Value: ReverseProxyTargetTypeSubnet, Description: "Subnet-based target."},
	}
}

// ReverseProxyPathRewrite controls how the request path is rewritten before forwarding to the backend.
type ReverseProxyPathRewrite string

// ReverseProxyPathRewritePreserve keeps the full original request path (default strips the matched prefix).
const ReverseProxyPathRewritePreserve ReverseProxyPathRewrite = ReverseProxyPathRewrite(nbapi.ServiceTargetOptionsPathRewritePreserve)

// Values returns the valid enum values for ReverseProxyPathRewrite.
func (ReverseProxyPathRewrite) Values() []infer.EnumValue[ReverseProxyPathRewrite] {
	return []infer.EnumValue[ReverseProxyPathRewrite]{
		{Name: "preserve", Value: ReverseProxyPathRewritePreserve, Description: "Keep the full original request path instead of stripping the matched prefix."},
	}
}

// ReverseProxyCrowdsecMode defines the CrowdSec IP reputation mode.
type ReverseProxyCrowdsecMode string

const (
	// ReverseProxyCrowdsecModeEnforce blocks connections flagged by CrowdSec.
	ReverseProxyCrowdsecModeEnforce ReverseProxyCrowdsecMode = ReverseProxyCrowdsecMode(nbapi.AccessRestrictionsCrowdsecModeEnforce)
	// ReverseProxyCrowdsecModeObserve only logs connections flagged by CrowdSec.
	ReverseProxyCrowdsecModeObserve ReverseProxyCrowdsecMode = ReverseProxyCrowdsecMode(nbapi.AccessRestrictionsCrowdsecModeObserve)
)

// Values returns the valid enum values for ReverseProxyCrowdsecMode.
func (ReverseProxyCrowdsecMode) Values() []infer.EnumValue[ReverseProxyCrowdsecMode] {
	return []infer.EnumValue[ReverseProxyCrowdsecMode]{
		{Name: "enforce", Value: ReverseProxyCrowdsecModeEnforce, Description: "Block connections flagged by CrowdSec."},
		{Name: "observe", Value: ReverseProxyCrowdsecModeObserve, Description: "Only log connections flagged by CrowdSec."},
	}
}

// ReverseProxyServiceStatus reflects the current lifecycle status of a service (output only).
type ReverseProxyServiceStatus string

const (
	// ReverseProxyServiceStatusActive means the service is provisioned and serving.
	ReverseProxyServiceStatusActive ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusActive)
	// ReverseProxyServiceStatusCertificateFailed means TLS certificate issuance failed.
	ReverseProxyServiceStatusCertificateFailed ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusCertificateFailed)
	// ReverseProxyServiceStatusCertificatePending means TLS certificate issuance is in progress.
	ReverseProxyServiceStatusCertificatePending ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusCertificatePending)
	// ReverseProxyServiceStatusError means the service is in an error state.
	ReverseProxyServiceStatusError ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusError)
	// ReverseProxyServiceStatusPending means the service is being provisioned.
	ReverseProxyServiceStatusPending ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusPending)
	// ReverseProxyServiceStatusTunnelNotCreated means the underlying tunnel has not been created yet.
	ReverseProxyServiceStatusTunnelNotCreated ReverseProxyServiceStatus = ReverseProxyServiceStatus(nbapi.ServiceMetaStatusTunnelNotCreated)
)

// Values returns the valid enum values for ReverseProxyServiceStatus.
func (ReverseProxyServiceStatus) Values() []infer.EnumValue[ReverseProxyServiceStatus] {
	return []infer.EnumValue[ReverseProxyServiceStatus]{
		{Name: "active", Value: ReverseProxyServiceStatusActive, Description: "Service is provisioned and serving."},
		{Name: "certificate_failed", Value: ReverseProxyServiceStatusCertificateFailed, Description: "TLS certificate issuance failed."},
		{Name: "certificate_pending", Value: ReverseProxyServiceStatusCertificatePending, Description: "TLS certificate issuance is in progress."},
		{Name: "error", Value: ReverseProxyServiceStatusError, Description: "Service is in an error state."},
		{Name: "pending", Value: ReverseProxyServiceStatusPending, Description: "Service is being provisioned."},
		{Name: "tunnel_not_created", Value: ReverseProxyServiceStatusTunnelNotCreated, Description: "Underlying tunnel has not been created yet."},
	}
}

// ReverseProxyTargetOptions defines advanced per-target proxy options.
type ReverseProxyTargetOptions struct {
	CustomHeaders      *map[string]string       `pulumi:"customHeaders,optional"`
	DirectUpstream     *bool                    `pulumi:"directUpstream,optional"`
	PathRewrite        *ReverseProxyPathRewrite `pulumi:"pathRewrite,optional"`
	ProxyProtocol      *bool                    `pulumi:"proxyProtocol,optional"`
	RequestTimeout     *string                  `pulumi:"requestTimeout,optional"`
	SessionIdleTimeout *string                  `pulumi:"sessionIdleTimeout,optional"`
	SkipTLSVerify      *bool                    `pulumi:"skipTlsVerify,optional"`
}

// Annotate provides documentation for ReverseProxyTargetOptions fields.
func (r *ReverseProxyTargetOptions) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.CustomHeaders, "Extra headers sent to the backend. Hop-by-hop and proxy-managed headers are rejected.")
	annotator.Describe(&r.DirectUpstream, "When true, the proxy dials this target via the host's network stack instead of through its embedded NetBird client.")
	annotator.Describe(&r.PathRewrite, `How the request path is rewritten before forwarding. Default strips the matched prefix; "preserve" keeps the full path.`)
	annotator.Describe(&r.ProxyProtocol, "Send PROXY Protocol v2 header to this backend (TCP/TLS only).")
	annotator.Describe(&r.RequestTimeout, `Per-target response timeout as a Go duration string (e.g. "30s", "2m").`)
	annotator.Describe(&r.SessionIdleTimeout, `Idle timeout before a UDP session is reaped, as a Go duration string (e.g. "30s", "2m").`)
	annotator.Describe(&r.SkipTLSVerify, "Skip TLS certificate verification for this backend.")
}

// ReverseProxyTarget defines a single backend target for a reverse proxy service.
type ReverseProxyTarget struct {
	TargetID   string                     `pulumi:"targetId"`
	Enabled    bool                       `pulumi:"enabled"`
	Host       *string                    `pulumi:"host,optional"`
	Port       int                        `pulumi:"port"`
	Protocol   ReverseProxyTargetProtocol `pulumi:"protocol"`
	TargetType ReverseProxyTargetType     `pulumi:"targetType"`
	Path       *string                    `pulumi:"path,optional"`
	Options    *ReverseProxyTargetOptions `pulumi:"options,optional"`
}

// Annotate provides documentation for ReverseProxyTarget fields.
func (r *ReverseProxyTarget) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.TargetID, "Target ID (assigned by the server).")
	annotator.Describe(&r.Enabled, "Whether this target is enabled.")
	annotator.Describe(&r.Host, "Backend IP or domain for this target.")
	annotator.Describe(&r.Port, "Backend port for this target.")
	annotator.Describe(&r.Protocol, "Protocol to use when connecting to the backend.")
	annotator.Describe(&r.TargetType, "Target type.")
	annotator.Describe(&r.Path, "URL path prefix for this target (HTTP only).")
	annotator.Describe(&r.Options, "Advanced per-target proxy options.")
}

// ReverseProxyBearerAuth configures bearer (SSO) authentication.
type ReverseProxyBearerAuth struct {
	Enabled            bool      `pulumi:"enabled"`
	DistributionGroups *[]string `pulumi:"distributionGroups,optional"`
}

// Annotate provides documentation for ReverseProxyBearerAuth fields.
func (r *ReverseProxyBearerAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Enabled, "Whether bearer auth is enabled.")
	annotator.Describe(&r.DistributionGroups, "List of group IDs that can use bearer auth.")
}

// ReverseProxyHeaderAuth configures header-based authentication.
type ReverseProxyHeaderAuth struct {
	Enabled bool   `pulumi:"enabled"`
	Header  string `pulumi:"header"`
	Value   string `pulumi:"value"`
}

// Annotate provides documentation for ReverseProxyHeaderAuth fields.
func (r *ReverseProxyHeaderAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Enabled, "Whether this header auth entry is enabled.")
	annotator.Describe(&r.Header, "The header name to match.")
	annotator.Describe(&r.Value, "The header value to match.")
}

// ReverseProxyLinkAuth configures link-based authentication.
type ReverseProxyLinkAuth struct {
	Enabled bool `pulumi:"enabled"`
}

// Annotate provides documentation for ReverseProxyLinkAuth fields.
func (r *ReverseProxyLinkAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Enabled, "Whether link auth is enabled.")
}

// ReverseProxyPasswordAuth configures password-based authentication.
type ReverseProxyPasswordAuth struct {
	Enabled  bool   `pulumi:"enabled"`
	Password string `provider:"secret" pulumi:"password"`
}

// Annotate provides documentation for ReverseProxyPasswordAuth fields.
func (r *ReverseProxyPasswordAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Enabled, "Whether password auth is enabled.")
	annotator.Describe(&r.Password, "The password required to access the service.")
}

// ReverseProxyPINAuth configures PIN-based authentication.
type ReverseProxyPINAuth struct {
	Enabled bool   `pulumi:"enabled"`
	Pin     string `provider:"secret" pulumi:"pin"`
}

// Annotate provides documentation for ReverseProxyPINAuth fields.
func (r *ReverseProxyPINAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.Enabled, "Whether PIN auth is enabled.")
	annotator.Describe(&r.Pin, "The PIN required to access the service.")
}

// ReverseProxyAuth defines the authentication configuration for a service.
type ReverseProxyAuth struct {
	BearerAuth   *ReverseProxyBearerAuth   `pulumi:"bearerAuth,optional"`
	HeaderAuths  *[]ReverseProxyHeaderAuth `pulumi:"headerAuths,optional"`
	LinkAuth     *ReverseProxyLinkAuth     `pulumi:"linkAuth,optional"`
	PasswordAuth *ReverseProxyPasswordAuth `pulumi:"passwordAuth,optional"`
	PinAuth      *ReverseProxyPINAuth      `pulumi:"pinAuth,optional"`
}

// Annotate provides documentation for ReverseProxyAuth fields.
func (r *ReverseProxyAuth) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.BearerAuth, "Bearer (SSO) authentication configuration.")
	annotator.Describe(&r.HeaderAuths, "Header-based authentication entries.")
	annotator.Describe(&r.LinkAuth, "Link-based authentication configuration.")
	annotator.Describe(&r.PasswordAuth, "Password-based authentication configuration.")
	annotator.Describe(&r.PinAuth, "PIN-based authentication configuration.")
}

// ReverseProxyAccessRestrictions defines connection-level access restrictions.
type ReverseProxyAccessRestrictions struct {
	AllowedCidrs     *[]string                 `pulumi:"allowedCidrs,optional"`
	AllowedCountries *[]string                 `pulumi:"allowedCountries,optional"`
	BlockedCidrs     *[]string                 `pulumi:"blockedCidrs,optional"`
	BlockedCountries *[]string                 `pulumi:"blockedCountries,optional"`
	CrowdsecMode     *ReverseProxyCrowdsecMode `pulumi:"crowdsecMode,optional"`
}

// Annotate provides documentation for ReverseProxyAccessRestrictions fields.
func (r *ReverseProxyAccessRestrictions) Annotate(annotator infer.Annotator) {
	annotator.Describe(&r.AllowedCidrs, "CIDR allowlist. If non-empty, only IPs matching these CIDRs are allowed.")
	annotator.Describe(&r.AllowedCountries, "ISO 3166-1 alpha-2 country codes to allow. If non-empty, only these countries are permitted.")
	annotator.Describe(&r.BlockedCidrs, "CIDR blocklist. Connections from these CIDRs are rejected. Evaluated after allowedCidrs.")
	annotator.Describe(&r.BlockedCountries, "ISO 3166-1 alpha-2 country codes to block.")
	annotator.Describe(&r.CrowdsecMode, "CrowdSec IP reputation mode. Only available when the proxy cluster supports CrowdSec.")
}

// toAPITargetOptions converts Pulumi target options to API target options.
func toAPITargetOptions(options *ReverseProxyTargetOptions) *nbapi.ServiceTargetOptions {
	if options == nil {
		return nil
	}

	var pathRewrite *nbapi.ServiceTargetOptionsPathRewrite

	if options.PathRewrite != nil {
		pr := nbapi.ServiceTargetOptionsPathRewrite(*options.PathRewrite)
		pathRewrite = &pr
	}

	return &nbapi.ServiceTargetOptions{
		CustomHeaders:      options.CustomHeaders,
		DirectUpstream:     options.DirectUpstream,
		PathRewrite:        pathRewrite,
		ProxyProtocol:      options.ProxyProtocol,
		RequestTimeout:     options.RequestTimeout,
		SessionIdleTimeout: options.SessionIdleTimeout,
		SkipTlsVerify:      options.SkipTLSVerify,
	}
}

// fromAPITargetOptions converts API target options to Pulumi target options.
func fromAPITargetOptions(options *nbapi.ServiceTargetOptions) *ReverseProxyTargetOptions {
	if options == nil {
		return nil
	}

	var pathRewrite *ReverseProxyPathRewrite

	if options.PathRewrite != nil {
		pr := ReverseProxyPathRewrite(*options.PathRewrite)
		pathRewrite = &pr
	}

	return &ReverseProxyTargetOptions{
		CustomHeaders:      options.CustomHeaders,
		DirectUpstream:     options.DirectUpstream,
		PathRewrite:        pathRewrite,
		ProxyProtocol:      options.ProxyProtocol,
		RequestTimeout:     options.RequestTimeout,
		SessionIdleTimeout: options.SessionIdleTimeout,
		SkipTLSVerify:      options.SkipTlsVerify,
	}
}

// toAPITargets converts Pulumi targets to API targets.
func toAPITargets(targets []ReverseProxyTarget) []nbapi.ServiceTarget {
	apiTargets := make([]nbapi.ServiceTarget, len(targets))
	for idx, target := range targets {
		apiTargets[idx] = nbapi.ServiceTarget{
			TargetId:   target.TargetID,
			Enabled:    target.Enabled,
			Host:       target.Host,
			Port:       target.Port,
			Protocol:   nbapi.ServiceTargetProtocol(target.Protocol),
			TargetType: nbapi.ServiceTargetTargetType(target.TargetType),
			Path:       target.Path,
			Options:    toAPITargetOptions(target.Options),
		}
	}

	return apiTargets
}

// fromAPITargets converts API targets to Pulumi targets.
func fromAPITargets(targets []nbapi.ServiceTarget) []ReverseProxyTarget {
	result := make([]ReverseProxyTarget, len(targets))
	for idx, apiTarget := range targets {
		result[idx] = ReverseProxyTarget{
			TargetID:   apiTarget.TargetId,
			Enabled:    apiTarget.Enabled,
			Host:       apiTarget.Host,
			Port:       apiTarget.Port,
			Protocol:   ReverseProxyTargetProtocol(apiTarget.Protocol),
			TargetType: ReverseProxyTargetType(apiTarget.TargetType),
			Path:       apiTarget.Path,
			Options:    fromAPITargetOptions(apiTarget.Options),
		}
	}

	return result
}

// toAPIAuth converts Pulumi auth config to an API auth config.
func toAPIAuth(auth *ReverseProxyAuth) *nbapi.ServiceAuthConfig {
	if auth == nil {
		return nil
	}

	var bearer *nbapi.BearerAuthConfig

	if auth.BearerAuth != nil {
		bearer = &nbapi.BearerAuthConfig{
			Enabled:            auth.BearerAuth.Enabled,
			DistributionGroups: auth.BearerAuth.DistributionGroups,
		}
	}

	var headers *[]nbapi.HeaderAuthConfig

	if auth.HeaderAuths != nil {
		list := make([]nbapi.HeaderAuthConfig, len(*auth.HeaderAuths))
		for idx, header := range *auth.HeaderAuths {
			list[idx] = nbapi.HeaderAuthConfig{
				Enabled: header.Enabled,
				Header:  header.Header,
				Value:   header.Value,
			}
		}

		headers = &list
	}

	var link *nbapi.LinkAuthConfig

	if auth.LinkAuth != nil {
		link = &nbapi.LinkAuthConfig{Enabled: auth.LinkAuth.Enabled}
	}

	var password *nbapi.PasswordAuthConfig

	if auth.PasswordAuth != nil {
		password = &nbapi.PasswordAuthConfig{
			Enabled:  auth.PasswordAuth.Enabled,
			Password: auth.PasswordAuth.Password,
		}
	}

	var pin *nbapi.PINAuthConfig

	if auth.PinAuth != nil {
		pin = &nbapi.PINAuthConfig{
			Enabled: auth.PinAuth.Enabled,
			Pin:     auth.PinAuth.Pin,
		}
	}

	return &nbapi.ServiceAuthConfig{
		BearerAuth:   bearer,
		HeaderAuths:  headers,
		LinkAuth:     link,
		PasswordAuth: password,
		PinAuth:      pin,
	}
}

// fromAPIAuth converts an API auth config to a Pulumi auth config.
// Returns nil when no sub-configuration is set, keeping state clean.
func fromAPIAuth(auth nbapi.ServiceAuthConfig) *ReverseProxyAuth {
	result := &ReverseProxyAuth{
		BearerAuth:   nil,
		HeaderAuths:  nil,
		LinkAuth:     nil,
		PasswordAuth: nil,
		PinAuth:      nil,
	}

	if auth.BearerAuth != nil {
		result.BearerAuth = &ReverseProxyBearerAuth{
			Enabled:            auth.BearerAuth.Enabled,
			DistributionGroups: auth.BearerAuth.DistributionGroups,
		}
	}

	if auth.HeaderAuths != nil {
		list := make([]ReverseProxyHeaderAuth, len(*auth.HeaderAuths))
		for idx, header := range *auth.HeaderAuths {
			list[idx] = ReverseProxyHeaderAuth{
				Enabled: header.Enabled,
				Header:  header.Header,
				Value:   header.Value,
			}
		}

		result.HeaderAuths = &list
	}

	if auth.LinkAuth != nil {
		result.LinkAuth = &ReverseProxyLinkAuth{Enabled: auth.LinkAuth.Enabled}
	}

	if auth.PasswordAuth != nil {
		result.PasswordAuth = &ReverseProxyPasswordAuth{
			Enabled:  auth.PasswordAuth.Enabled,
			Password: auth.PasswordAuth.Password,
		}
	}

	if auth.PinAuth != nil {
		result.PinAuth = &ReverseProxyPINAuth{
			Enabled: auth.PinAuth.Enabled,
			Pin:     auth.PinAuth.Pin,
		}
	}

	if result.BearerAuth == nil && result.HeaderAuths == nil && result.LinkAuth == nil &&
		result.PasswordAuth == nil && result.PinAuth == nil {
		return nil
	}

	return result
}

// toAPIAccessRestrictions converts Pulumi access restrictions to API access restrictions.
func toAPIAccessRestrictions(restrictions *ReverseProxyAccessRestrictions) *nbapi.AccessRestrictions {
	if restrictions == nil {
		return nil
	}

	var crowdsec *nbapi.AccessRestrictionsCrowdsecMode

	if restrictions.CrowdsecMode != nil {
		mode := nbapi.AccessRestrictionsCrowdsecMode(*restrictions.CrowdsecMode)
		crowdsec = &mode
	}

	return &nbapi.AccessRestrictions{
		AllowedCidrs:     restrictions.AllowedCidrs,
		AllowedCountries: restrictions.AllowedCountries,
		BlockedCidrs:     restrictions.BlockedCidrs,
		BlockedCountries: restrictions.BlockedCountries,
		CrowdsecMode:     crowdsec,
	}
}

// fromAPIAccessRestrictions converts API access restrictions to Pulumi access restrictions.
func fromAPIAccessRestrictions(restrictions *nbapi.AccessRestrictions) *ReverseProxyAccessRestrictions {
	if restrictions == nil {
		return nil
	}

	var crowdsec *ReverseProxyCrowdsecMode

	if restrictions.CrowdsecMode != nil {
		mode := ReverseProxyCrowdsecMode(*restrictions.CrowdsecMode)
		crowdsec = &mode
	}

	return &ReverseProxyAccessRestrictions{
		AllowedCidrs:     restrictions.AllowedCidrs,
		AllowedCountries: restrictions.AllowedCountries,
		BlockedCidrs:     restrictions.BlockedCidrs,
		BlockedCountries: restrictions.BlockedCountries,
		CrowdsecMode:     crowdsec,
	}
}

// serviceStateFromAPI builds a ReverseProxyServiceState from an API Service response.
func serviceStateFromAPI(svc *nbapi.Service) ReverseProxyServiceState {
	var mode *ReverseProxyServiceMode

	if svc.Mode != nil {
		m := ReverseProxyServiceMode(*svc.Mode)
		mode = &m
	}

	statusVal := ReverseProxyServiceStatus(svc.Meta.Status)

	return ReverseProxyServiceState{
		Name:               svc.Name,
		Domain:             svc.Domain,
		Enabled:            svc.Enabled,
		Mode:               mode,
		Targets:            fromAPITargets(svc.Targets),
		PassHostHeader:     svc.PassHostHeader,
		RewriteRedirects:   svc.RewriteRedirects,
		ListenPort:         svc.ListenPort,
		Private:            svc.Private,
		AccessGroups:       svc.AccessGroups,
		Auth:               fromAPIAuth(svc.Auth),
		AccessRestrictions: fromAPIAccessRestrictions(svc.AccessRestrictions),
		ProxyCluster:       svc.ProxyCluster,
		Status:             &statusVal,
		Terminated:         svc.Terminated,
		PortAutoAssigned:   svc.PortAutoAssigned,
	}
}

// buildServiceRequest constructs an API ServiceRequest from Pulumi inputs.
func buildServiceRequest(args ReverseProxyServiceArgs) nbapi.ServiceRequest {
	var mode *nbapi.ServiceRequestMode

	if args.Mode != nil {
		m := nbapi.ServiceRequestMode(*args.Mode)
		mode = &m
	}

	apiTargets := toAPITargets(args.Targets)

	return nbapi.ServiceRequest{
		Name:               args.Name,
		Domain:             args.Domain,
		Enabled:            args.Enabled,
		Mode:               mode,
		Targets:            &apiTargets,
		PassHostHeader:     args.PassHostHeader,
		RewriteRedirects:   args.RewriteRedirects,
		ListenPort:         args.ListenPort,
		Private:            args.Private,
		AccessGroups:       args.AccessGroups,
		Auth:               toAPIAuth(args.Auth),
		AccessRestrictions: toAPIAccessRestrictions(args.AccessRestrictions),
	}
}

// Create creates a new NetBird reverse proxy service.
func (*ReverseProxyService) Create(ctx context.Context, req infer.CreateRequest[ReverseProxyServiceArgs]) (infer.CreateResponse[ReverseProxyServiceState], error) {
	p.GetLogger(ctx).Debugf("Create:ReverseProxyService name=%s, domain=%s", req.Inputs.Name, req.Inputs.Domain)

	if req.DryRun {
		return infer.CreateResponse[ReverseProxyServiceState]{
			ID:     "preview",
			Output: dryRunState(req.Inputs, nil, nil),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[ReverseProxyServiceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	svc, err := client.ReverseProxyServices.Create(ctx, buildServiceRequest(req.Inputs))
	if err != nil {
		return infer.CreateResponse[ReverseProxyServiceState]{}, fmt.Errorf("creating reverse proxy service failed: %w", err)
	}

	return infer.CreateResponse[ReverseProxyServiceState]{
		ID:     svc.Id,
		Output: serviceStateFromAPI(svc),
	}, nil
}

// Read reads a reverse proxy service from NetBird.
func (*ReverseProxyService) Read(ctx context.Context, req infer.ReadRequest[ReverseProxyServiceArgs, ReverseProxyServiceState]) (infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState], error) {
	p.GetLogger(ctx).Debugf("Read:ReverseProxyService[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	svc, err := client.ReverseProxyServices.Get(ctx, req.ID)
	if err != nil {
		if isNotFoundErr(err) {
			return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{
				ID:     "",
				Inputs: ReverseProxyServiceArgs{},  //nolint:exhaustruct
				State:  ReverseProxyServiceState{}, //nolint:exhaustruct
			}, nil
		}

		return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{}, fmt.Errorf("reading reverse proxy service failed: %w", err)
	}

	state := serviceStateFromAPI(svc)

	return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{
		ID: req.ID,
		Inputs: ReverseProxyServiceArgs{
			Name:               state.Name,
			Domain:             state.Domain,
			Enabled:            state.Enabled,
			Mode:               state.Mode,
			Targets:            state.Targets,
			PassHostHeader:     state.PassHostHeader,
			RewriteRedirects:   state.RewriteRedirects,
			ListenPort:         state.ListenPort,
			Private:            state.Private,
			AccessGroups:       state.AccessGroups,
			Auth:               state.Auth,
			AccessRestrictions: state.AccessRestrictions,
		},
		State: state,
	}, nil
}

// Update updates a reverse proxy service in NetBird.
func (*ReverseProxyService) Update(ctx context.Context, req infer.UpdateRequest[ReverseProxyServiceArgs, ReverseProxyServiceState]) (infer.UpdateResponse[ReverseProxyServiceState], error) {
	p.GetLogger(ctx).Debugf("Update:ReverseProxyService[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[ReverseProxyServiceState]{
			Output: dryRunState(req.Inputs, req.State.ProxyCluster, req.State.Status),
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[ReverseProxyServiceState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	svc, err := client.ReverseProxyServices.Update(ctx, req.ID, buildServiceRequest(req.Inputs))
	if err != nil {
		return infer.UpdateResponse[ReverseProxyServiceState]{}, fmt.Errorf("updating reverse proxy service failed: %w", err)
	}

	return infer.UpdateResponse[ReverseProxyServiceState]{
		Output: serviceStateFromAPI(svc),
	}, nil
}

// dryRunState builds a preview state from inputs, carrying over server-derived
// outputs when they are known (Update) and leaving them nil otherwise (Create).
func dryRunState(inputs ReverseProxyServiceArgs, proxyCluster *string, status *ReverseProxyServiceStatus) ReverseProxyServiceState {
	return ReverseProxyServiceState{
		Name:               inputs.Name,
		Domain:             inputs.Domain,
		Enabled:            inputs.Enabled,
		Mode:               inputs.Mode,
		Targets:            inputs.Targets,
		PassHostHeader:     inputs.PassHostHeader,
		RewriteRedirects:   inputs.RewriteRedirects,
		ListenPort:         inputs.ListenPort,
		Private:            inputs.Private,
		AccessGroups:       inputs.AccessGroups,
		Auth:               inputs.Auth,
		AccessRestrictions: inputs.AccessRestrictions,
		ProxyCluster:       proxyCluster,
		Status:             status,
		Terminated:         nil,
		PortAutoAssigned:   nil,
	}
}

// Delete removes a reverse proxy service from NetBird.
func (*ReverseProxyService) Delete(ctx context.Context, req infer.DeleteRequest[ReverseProxyServiceState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:ReverseProxyService[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.ReverseProxyServices.Delete(ctx, req.ID)
	if err != nil && !isNotFoundErr(err) {
		return infer.DeleteResponse{}, fmt.Errorf("deleting reverse proxy service failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between ReverseProxyServiceArgs and ReverseProxyServiceState.
func (*ReverseProxyService) Diff(ctx context.Context, req infer.DiffRequest[ReverseProxyServiceArgs, ReverseProxyServiceState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:ReverseProxyService[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Domain != req.State.Domain {
		diff["domain"] = p.PropertyDiff{InputDiff: false, Kind: p.UpdateReplace}
	}

	if req.Inputs.Enabled != req.State.Enabled {
		diff["enabled"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.Mode, req.State.Mode) {
		diff["mode"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.PassHostHeader, req.State.PassHostHeader) {
		diff["passHostHeader"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.RewriteRedirects, req.State.RewriteRedirects) {
		diff["rewriteRedirects"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.ListenPort, req.State.ListenPort) {
		diff["listenPort"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalPtr(req.Inputs.Private, req.State.Private) {
		diff["private"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalSlicePtr(req.Inputs.AccessGroups, req.State.AccessGroups) {
		diff["accessGroups"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalReverseProxyTargets(req.Inputs.Targets, req.State.Targets) {
		diff["targets"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalOptionalDeep(req.Inputs.Auth, req.State.Auth) {
		diff["auth"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !reflect.DeepEqual(req.Inputs.AccessRestrictions, req.State.AccessRestrictions) {
		diff["accessRestrictions"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:ReverseProxyService[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a reverse proxy service.
func (*ReverseProxyService) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[ReverseProxyServiceArgs], error) {
	p.GetLogger(ctx).Debugf("Check:ReverseProxyService old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[ReverseProxyServiceArgs](ctx, req.NewInputs)

	failures = append(failures, reverseProxyCheckArgs(args)...)

	return infer.CheckResponse[ReverseProxyServiceArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// reverseProxyCheckArgs validates a ReverseProxyServiceArgs and returns all failures.
//
//nolint:cyclop
func reverseProxyCheckArgs(args ReverseProxyServiceArgs) []p.CheckFailure {
	var failures []p.CheckFailure

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{Property: "name", Reason: "name must not be empty"})
	}

	if isBlank(args.Domain) {
		failures = append(failures, p.CheckFailure{Property: "domain", Reason: "domain must not be empty"})
	}

	if len(args.Targets) == 0 {
		failures = append(failures, p.CheckFailure{Property: "targets", Reason: "at least one target is required"})
	}

	for idx, target := range args.Targets {
		if target.Port < 1 || target.Port > 65535 {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("targets[%d].port", idx),
				Reason:   "port must be between 1 and 65535",
			})
		}

		if target.TargetType == ReverseProxyTargetTypeDomain || target.TargetType == ReverseProxyTargetTypeHost {
			if target.Host == nil || isBlank(*target.Host) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("targets[%d].host", idx),
					Reason:   "host must not be empty for domain or host target types",
				})
			}
		}
	}

	if args.AccessGroups != nil {
		for i, group := range *args.AccessGroups {
			if isBlank(group) {
				failures = append(failures, p.CheckFailure{
					Property: fmt.Sprintf("accessGroups[%d]", i),
					Reason:   "access group id must not be empty",
				})
			}
		}
	}

	if boolVal(args.Private) {
		if args.Mode == nil || *args.Mode != ReverseProxyServiceModeHTTP {
			failures = append(failures, p.CheckFailure{
				Property: "mode",
				Reason:   "mode must be http when private is true",
			})
		}

		if args.AccessGroups == nil || len(*args.AccessGroups) == 0 {
			failures = append(failures, p.CheckFailure{
				Property: "accessGroups",
				Reason:   "accessGroups must not be empty when private is true",
			})
		}
	}

	return failures
}

// WireDependencies explicitly defines input/output relationships.
func (*ReverseProxyService) WireDependencies(field infer.FieldSelector, args *ReverseProxyServiceArgs, state *ReverseProxyServiceState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Domain).DependsOn(field.InputField(&args.Domain))
	field.OutputField(&state.Enabled).DependsOn(field.InputField(&args.Enabled))
	field.OutputField(&state.Mode).DependsOn(field.InputField(&args.Mode))
	field.OutputField(&state.Targets).DependsOn(field.InputField(&args.Targets))
	field.OutputField(&state.PassHostHeader).DependsOn(field.InputField(&args.PassHostHeader))
	field.OutputField(&state.RewriteRedirects).DependsOn(field.InputField(&args.RewriteRedirects))
	field.OutputField(&state.ListenPort).DependsOn(field.InputField(&args.ListenPort))
	field.OutputField(&state.Private).DependsOn(field.InputField(&args.Private))
	field.OutputField(&state.AccessGroups).DependsOn(field.InputField(&args.AccessGroups))
	field.OutputField(&state.Auth).DependsOn(field.InputField(&args.Auth))
	field.OutputField(&state.AccessRestrictions).DependsOn(field.InputField(&args.AccessRestrictions))
}

// equalReverseProxyTargets compares two slices of ReverseProxyTarget by their key fields.
func equalReverseProxyTargets(targetsA, targetsB []ReverseProxyTarget) bool {
	if len(targetsA) != len(targetsB) {
		return false
	}

	for idx := range targetsA {
		if targetsA[idx].Enabled != targetsB[idx].Enabled ||
			targetsA[idx].Port != targetsB[idx].Port ||
			targetsA[idx].Protocol != targetsB[idx].Protocol ||
			targetsA[idx].TargetType != targetsB[idx].TargetType ||
			!equalPtr(targetsA[idx].Host, targetsB[idx].Host) ||
			!equalPtr(targetsA[idx].Path, targetsB[idx].Path) ||
			!equalOptionalDeep(targetsA[idx].Options, targetsB[idx].Options) {
			return false
		}
	}

	return true
}

// equalOptionalDeep compares two pointers by the values they point to,
// treating a nil pointer as equivalent to a pointer to the zero value. The
// NetBird API always echoes back a full struct for these optional fields
// (e.g. target options, auth config) even when nothing was configured,
// while a Pulumi program that never sets the field leaves it nil; without
// this normalization a raw reflect.DeepEqual would treat nil and an
// all-zero-fields struct as different forever.
func equalOptionalDeep[T any](valueA, valueB *T) bool {
	var zero T

	dereferencedA := zero
	if valueA != nil {
		dereferencedA = *valueA
	}

	dereferencedB := zero
	if valueB != nil {
		dereferencedB = *valueB
	}

	return reflect.DeepEqual(dereferencedA, dereferencedB)
}
