package resource

import (
	"context"
	"fmt"

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
	Name             string                   `pulumi:"name"`
	Domain           string                   `pulumi:"domain"`
	Enabled          bool                     `pulumi:"enabled"`
	Mode             *ReverseProxyServiceMode `pulumi:"mode,optional"`
	Targets          []ReverseProxyTarget     `pulumi:"targets"`
	PassHostHeader   *bool                    `pulumi:"passHostHeader,optional"`
	RewriteRedirects *bool                    `pulumi:"rewriteRedirects,optional"`
	ListenPort       *int                     `pulumi:"listenPort,optional"`
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
}

// ReverseProxyServiceState represents the output state of a reverse proxy service resource.
type ReverseProxyServiceState struct {
	Name             string                   `pulumi:"name"`
	Domain           string                   `pulumi:"domain"`
	Enabled          bool                     `pulumi:"enabled"`
	Mode             *ReverseProxyServiceMode `pulumi:"mode,optional"`
	Targets          []ReverseProxyTarget     `pulumi:"targets"`
	PassHostHeader   *bool                    `pulumi:"passHostHeader,optional"`
	RewriteRedirects *bool                    `pulumi:"rewriteRedirects,optional"`
	ListenPort       *int                     `pulumi:"listenPort,optional"`
	ProxyCluster     *string                  `pulumi:"proxyCluster,optional"`
	Status           *string                  `pulumi:"status,optional"`
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
	annotator.Describe(&r.ProxyCluster, "The proxy cluster handling this service (derived from domain).")
	annotator.Describe(&r.Status, "Current status of the service.")
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
		{Name: "domain", Value: ReverseProxyTargetTypeDomain, Description: "Domain-based target."},
		{Name: "host", Value: ReverseProxyTargetTypeHost, Description: "Host-based target."},
		{Name: "peer", Value: ReverseProxyTargetTypePeer, Description: "Peer-based target."},
		{Name: "subnet", Value: ReverseProxyTargetTypeSubnet, Description: "Subnet-based target."},
	}
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
			Options:    nil,
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
		}
	}

	return result
}

// serviceStateFromAPI builds a ReverseProxyServiceState from an API Service response.
func serviceStateFromAPI(svc *nbapi.Service) ReverseProxyServiceState {
	var mode *ReverseProxyServiceMode

	if svc.Mode != nil {
		m := ReverseProxyServiceMode(*svc.Mode)
		mode = &m
	}

	var status *string

	statusVal := string(svc.Meta.Status)
	status = &statusVal

	return ReverseProxyServiceState{
		Name:             svc.Name,
		Domain:           svc.Domain,
		Enabled:          svc.Enabled,
		Mode:             mode,
		Targets:          fromAPITargets(svc.Targets),
		PassHostHeader:   svc.PassHostHeader,
		RewriteRedirects: svc.RewriteRedirects,
		ListenPort:       svc.ListenPort,
		ProxyCluster:     svc.ProxyCluster,
		Status:           status,
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
		Auth:               nil,
		AccessRestrictions: nil,
	}
}

// Create creates a new NetBird reverse proxy service.
func (*ReverseProxyService) Create(ctx context.Context, req infer.CreateRequest[ReverseProxyServiceArgs]) (infer.CreateResponse[ReverseProxyServiceState], error) {
	p.GetLogger(ctx).Debugf("Create:ReverseProxyService name=%s, domain=%s", req.Inputs.Name, req.Inputs.Domain)

	if req.DryRun {
		return infer.CreateResponse[ReverseProxyServiceState]{
			ID: "preview",
			Output: ReverseProxyServiceState{
				Name:             req.Inputs.Name,
				Domain:           req.Inputs.Domain,
				Enabled:          req.Inputs.Enabled,
				Mode:             req.Inputs.Mode,
				Targets:          req.Inputs.Targets,
				PassHostHeader:   req.Inputs.PassHostHeader,
				RewriteRedirects: req.Inputs.RewriteRedirects,
				ListenPort:       req.Inputs.ListenPort,
				ProxyCluster:     nil,
				Status:           nil,
			},
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
		return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{}, fmt.Errorf("reading reverse proxy service failed: %w", err)
	}

	state := serviceStateFromAPI(svc)

	return infer.ReadResponse[ReverseProxyServiceArgs, ReverseProxyServiceState]{
		ID: req.ID,
		Inputs: ReverseProxyServiceArgs{
			Name:             state.Name,
			Domain:           state.Domain,
			Enabled:          state.Enabled,
			Mode:             state.Mode,
			Targets:          state.Targets,
			PassHostHeader:   state.PassHostHeader,
			RewriteRedirects: state.RewriteRedirects,
			ListenPort:       state.ListenPort,
		},
		State: state,
	}, nil
}

// Update updates a reverse proxy service in NetBird.
func (*ReverseProxyService) Update(ctx context.Context, req infer.UpdateRequest[ReverseProxyServiceArgs, ReverseProxyServiceState]) (infer.UpdateResponse[ReverseProxyServiceState], error) {
	p.GetLogger(ctx).Debugf("Update:ReverseProxyService[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[ReverseProxyServiceState]{
			Output: ReverseProxyServiceState{
				Name:             req.Inputs.Name,
				Domain:           req.Inputs.Domain,
				Enabled:          req.Inputs.Enabled,
				Mode:             req.Inputs.Mode,
				Targets:          req.Inputs.Targets,
				PassHostHeader:   req.Inputs.PassHostHeader,
				RewriteRedirects: req.Inputs.RewriteRedirects,
				ListenPort:       req.Inputs.ListenPort,
				ProxyCluster:     req.State.ProxyCluster,
				Status:           req.State.Status,
			},
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

// Delete removes a reverse proxy service from NetBird.
func (*ReverseProxyService) Delete(ctx context.Context, req infer.DeleteRequest[ReverseProxyServiceState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:ReverseProxyService[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.ReverseProxyServices.Delete(ctx, req.ID)
	if err != nil {
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

	if !equalReverseProxyTargets(req.Inputs.Targets, req.State.Targets) {
		diff["targets"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
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

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if isBlank(args.Domain) {
		failures = append(failures, p.CheckFailure{
			Property: "domain",
			Reason:   "domain must not be empty",
		})
	}

	if len(args.Targets) == 0 {
		failures = append(failures, p.CheckFailure{
			Property: "targets",
			Reason:   "at least one target is required",
		})
	}

	for i, target := range args.Targets {
		if target.Port < 1 || target.Port > 65535 {
			failures = append(failures, p.CheckFailure{
				Property: fmt.Sprintf("targets[%d].port", i),
				Reason:   "port must be between 1 and 65535",
			})
		}
	}

	return infer.CheckResponse[ReverseProxyServiceArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
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
}
