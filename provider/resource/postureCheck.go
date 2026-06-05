package resource

import (
	"context"
	"fmt"
	"slices"

	"github.com/mbrav/pulumi-netbird/provider/config"
	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// PostureCheck represents a NetBird posture check resource.
type PostureCheck struct{}

// Annotate adds a description to the PostureCheck resource type.
func (pc *PostureCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc, "A NetBird posture check used to validate peer properties before granting policy access.")
}

// PostureCheckArgs defines input fields for creating or updating a posture check.
type PostureCheckArgs struct {
	Name        string              `pulumi:"name"`
	Description *string             `pulumi:"description,optional"`
	Checks      PostureChecksConfig `pulumi:"checks"`
}

// Annotate provides documentation for PostureCheckArgs fields.
func (pc *PostureCheckArgs) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Name, "Posture check unique name identifier.")
	annotator.Describe(&pc.Description, "Posture check friendly description.")
	annotator.Describe(&pc.Checks, "List of checks to perform against peer properties.")
}

// PostureCheckState represents the output state of a posture check resource.
type PostureCheckState struct {
	Name        string              `pulumi:"name"`
	Description *string             `pulumi:"description,optional"`
	Checks      PostureChecksConfig `pulumi:"checks"`
}

// Annotate provides documentation for PostureCheckState fields.
func (pc *PostureCheckState) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Name, "Posture check unique name identifier.")
	annotator.Describe(&pc.Description, "Posture check friendly description.")
	annotator.Describe(&pc.Checks, "List of checks to perform against peer properties.")
}

// PostureChecksConfig groups all available posture checks.
type PostureChecksConfig struct {
	GeoLocation  *PostureGeoLocationCheck      `pulumi:"geoLocationCheck,optional"`
	NbVersion    *PostureMinVersionCheck       `pulumi:"nbVersionCheck,optional"`
	OsVersion    *PostureOSVersionCheck        `pulumi:"osVersionCheck,optional"`
	NetworkRange *PosturePeerNetworkRangeCheck `pulumi:"peerNetworkRangeCheck,optional"`
	Process      *PostureProcessCheck          `pulumi:"processCheck,optional"`
}

// Annotate provides documentation for PostureChecksConfig fields.
func (pc *PostureChecksConfig) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.GeoLocation, "Posture check for geo location.")
	annotator.Describe(&pc.NbVersion, "Posture check for the minimum NetBird client version.")
	annotator.Describe(&pc.OsVersion, "Posture check for the minimum operating system version.")
	annotator.Describe(&pc.NetworkRange, "Posture check based on peer local network addresses.")
	annotator.Describe(&pc.Process, "Posture check for required binaries running on the peer.")
}

// PostureGeoLocationCheck defines a geo location posture check.
type PostureGeoLocationCheck struct {
	Action    PostureGeoLocationAction `pulumi:"action"`
	Locations []PostureLocation        `pulumi:"locations"`
}

// Annotate provides documentation for PostureGeoLocationCheck fields.
func (pc *PostureGeoLocationCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Action, "Action to take upon geo location match (allow or deny).")
	annotator.Describe(&pc.Locations, "List of geo locations to which the check applies.")
}

// PostureLocation defines a geographic location used in posture checks.
type PostureLocation struct {
	CountryCode string  `pulumi:"countryCode"`
	CityName    *string `pulumi:"cityName,optional"`
}

// Annotate provides documentation for PostureLocation fields.
func (pl *PostureLocation) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pl.CountryCode, "2-letter ISO 3166-1 alpha-2 country code.")
	annotator.Describe(&pl.CityName, "Commonly used English name of the city.")
}

// PostureGeoLocationAction defines the action for a geo location posture check.
type PostureGeoLocationAction string

const (
	// PostureGeoLocationActionAllow permits peers matching the specified locations.
	PostureGeoLocationActionAllow PostureGeoLocationAction = PostureGeoLocationAction(nbapi.GeoLocationCheckActionAllow)
	// PostureGeoLocationActionDeny rejects peers matching the specified locations.
	PostureGeoLocationActionDeny PostureGeoLocationAction = PostureGeoLocationAction(nbapi.GeoLocationCheckActionDeny)
)

// Values returns the valid enum values for PostureGeoLocationAction.
func (PostureGeoLocationAction) Values() []infer.EnumValue[PostureGeoLocationAction] {
	return []infer.EnumValue[PostureGeoLocationAction]{
		{Name: "allow", Value: PostureGeoLocationActionAllow, Description: "Allow peers from the specified locations."},
		{Name: "deny", Value: PostureGeoLocationActionDeny, Description: "Deny peers from the specified locations."},
	}
}

// PostureMinVersionCheck defines a minimum version posture check.
type PostureMinVersionCheck struct {
	MinVersion string `pulumi:"minVersion"`
}

// Annotate provides documentation for PostureMinVersionCheck fields.
func (pc *PostureMinVersionCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.MinVersion, "Minimum acceptable version string.")
}

// PostureMinKernelVersionCheck defines a minimum kernel version posture check.
type PostureMinKernelVersionCheck struct {
	MinKernelVersion string `pulumi:"minKernelVersion"`
}

// Annotate provides documentation for PostureMinKernelVersionCheck fields.
func (pc *PostureMinKernelVersionCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.MinKernelVersion, "Minimum acceptable kernel version string.")
}

// PostureOSVersionCheck defines minimum OS version requirements per platform.
type PostureOSVersionCheck struct {
	Android *PostureMinVersionCheck       `pulumi:"android,optional"`
	Darwin  *PostureMinVersionCheck       `pulumi:"darwin,optional"`
	Ios     *PostureMinVersionCheck       `pulumi:"ios,optional"`
	Linux   *PostureMinKernelVersionCheck `pulumi:"linux,optional"`
	Windows *PostureMinKernelVersionCheck `pulumi:"windows,optional"`
}

// Annotate provides documentation for PostureOSVersionCheck fields.
func (pc *PostureOSVersionCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Android, "Minimum version check for Android.")
	annotator.Describe(&pc.Darwin, "Minimum version check for macOS.")
	annotator.Describe(&pc.Ios, "Minimum version check for iOS.")
	annotator.Describe(&pc.Linux, "Minimum kernel version check for Linux.")
	annotator.Describe(&pc.Windows, "Minimum kernel version check for Windows.")
}

// PosturePeerNetworkRangeCheck defines a posture check based on peer network ranges.
type PosturePeerNetworkRangeCheck struct {
	Action PosturePeerNetworkRangeAction `pulumi:"action"`
	Ranges []string                      `pulumi:"ranges"`
}

// Annotate provides documentation for PosturePeerNetworkRangeCheck fields.
func (pc *PosturePeerNetworkRangeCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Action, "Action to take when the peer's network range matches (allow or deny).")
	annotator.Describe(&pc.Ranges, "List of CIDR network ranges to match against.")
}

// PosturePeerNetworkRangeAction defines the action for a network range posture check.
type PosturePeerNetworkRangeAction string

const (
	// PosturePeerNetworkRangeActionAllow permits peers whose local network matches.
	PosturePeerNetworkRangeActionAllow PosturePeerNetworkRangeAction = PosturePeerNetworkRangeAction(nbapi.PeerNetworkRangeCheckActionAllow)
	// PosturePeerNetworkRangeActionDeny rejects peers whose local network matches.
	PosturePeerNetworkRangeActionDeny PosturePeerNetworkRangeAction = PosturePeerNetworkRangeAction(nbapi.PeerNetworkRangeCheckActionDeny)
)

// Values returns the valid enum values for PosturePeerNetworkRangeAction.
func (PosturePeerNetworkRangeAction) Values() []infer.EnumValue[PosturePeerNetworkRangeAction] {
	return []infer.EnumValue[PosturePeerNetworkRangeAction]{
		{Name: "allow", Value: PosturePeerNetworkRangeActionAllow, Description: "Allow peers whose local network matches."},
		{Name: "deny", Value: PosturePeerNetworkRangeActionDeny, Description: "Deny peers whose local network matches."},
	}
}

// PostureProcessCheck defines a posture check for required running processes.
type PostureProcessCheck struct {
	Processes []PostureProcess `pulumi:"processes"`
}

// Annotate provides documentation for PostureProcessCheck fields.
func (pc *PostureProcessCheck) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pc.Processes, "List of processes that must be running on the peer.")
}

// PostureProcess defines a single process required by a posture check.
type PostureProcess struct {
	LinuxPath   *string `pulumi:"linuxPath,optional"`
	MacPath     *string `pulumi:"macPath,optional"`
	WindowsPath *string `pulumi:"windowsPath,optional"`
}

// Annotate provides documentation for PostureProcess fields.
func (pp *PostureProcess) Annotate(annotator infer.Annotator) {
	annotator.Describe(&pp.LinuxPath, "Path to the process executable on Linux.")
	annotator.Describe(&pp.MacPath, "Path to the process executable on macOS.")
	annotator.Describe(&pp.WindowsPath, "Path to the process executable on Windows.")
}

// toAPIOSVersionCheck converts a PostureOSVersionCheck to an api.OSVersionCheck.
func toAPIOSVersionCheck(check *PostureOSVersionCheck) *nbapi.OSVersionCheck {
	result := nbapi.OSVersionCheck{
		Android: nil,
		Darwin:  nil,
		Ios:     nil,
		Linux:   nil,
		Windows: nil,
	}

	if check.Android != nil {
		mv := nbapi.MinVersionCheck{MinVersion: check.Android.MinVersion}
		result.Android = &mv
	}

	if check.Darwin != nil {
		mv := nbapi.MinVersionCheck{MinVersion: check.Darwin.MinVersion}
		result.Darwin = &mv
	}

	if check.Ios != nil {
		mv := nbapi.MinVersionCheck{MinVersion: check.Ios.MinVersion}
		result.Ios = &mv
	}

	if check.Linux != nil {
		mk := nbapi.MinKernelVersionCheck{MinKernelVersion: check.Linux.MinKernelVersion}
		result.Linux = &mk
	}

	if check.Windows != nil {
		mk := nbapi.MinKernelVersionCheck{MinKernelVersion: check.Windows.MinKernelVersion}
		result.Windows = &mk
	}

	return &result
}

// fromAPIOSVersionCheck converts an api.OSVersionCheck to a PostureOSVersionCheck.
func fromAPIOSVersionCheck(check *nbapi.OSVersionCheck) *PostureOSVersionCheck {
	result := PostureOSVersionCheck{
		Android: nil,
		Darwin:  nil,
		Ios:     nil,
		Linux:   nil,
		Windows: nil,
	}

	if check.Android != nil {
		result.Android = &PostureMinVersionCheck{MinVersion: check.Android.MinVersion}
	}

	if check.Darwin != nil {
		result.Darwin = &PostureMinVersionCheck{MinVersion: check.Darwin.MinVersion}
	}

	if check.Ios != nil {
		result.Ios = &PostureMinVersionCheck{MinVersion: check.Ios.MinVersion}
	}

	if check.Linux != nil {
		result.Linux = &PostureMinKernelVersionCheck{MinKernelVersion: check.Linux.MinKernelVersion}
	}

	if check.Windows != nil {
		result.Windows = &PostureMinKernelVersionCheck{MinKernelVersion: check.Windows.MinKernelVersion}
	}

	return &result
}

// toAPIChecks converts a PostureChecksConfig to an api.Checks.
func toAPIChecks(checks PostureChecksConfig) nbapi.Checks {
	result := nbapi.Checks{
		GeoLocationCheck:      nil,
		NbVersionCheck:        nil,
		OsVersionCheck:        nil,
		PeerNetworkRangeCheck: nil,
		ProcessCheck:          nil,
	}

	if checks.GeoLocation != nil {
		locations := make([]nbapi.Location, len(checks.GeoLocation.Locations))
		for idx, loc := range checks.GeoLocation.Locations {
			locations[idx] = nbapi.Location{
				CountryCode: loc.CountryCode,
				CityName:    loc.CityName,
			}
		}

		result.GeoLocationCheck = &nbapi.GeoLocationCheck{
			Action:    nbapi.GeoLocationCheckAction(checks.GeoLocation.Action),
			Locations: locations,
		}
	}

	if checks.NbVersion != nil {
		nb := nbapi.MinVersionCheck{MinVersion: checks.NbVersion.MinVersion}
		result.NbVersionCheck = &nb
	}

	if checks.OsVersion != nil {
		result.OsVersionCheck = toAPIOSVersionCheck(checks.OsVersion)
	}

	if checks.NetworkRange != nil {
		result.PeerNetworkRangeCheck = &nbapi.PeerNetworkRangeCheck{
			Action: nbapi.PeerNetworkRangeCheckAction(checks.NetworkRange.Action),
			Ranges: checks.NetworkRange.Ranges,
		}
	}

	if checks.Process != nil {
		processes := make([]nbapi.Process, len(checks.Process.Processes))
		for idx, proc := range checks.Process.Processes {
			processes[idx] = nbapi.Process{
				LinuxPath:   proc.LinuxPath,
				MacPath:     proc.MacPath,
				WindowsPath: proc.WindowsPath,
			}
		}

		result.ProcessCheck = &nbapi.ProcessCheck{Processes: processes}
	}

	return result
}

// fromAPIChecks converts an api.Checks to a PostureChecksConfig.
func fromAPIChecks(checks nbapi.Checks) PostureChecksConfig {
	result := PostureChecksConfig{
		GeoLocation:  nil,
		NbVersion:    nil,
		OsVersion:    nil,
		NetworkRange: nil,
		Process:      nil,
	}

	if checks.GeoLocationCheck != nil {
		locations := make([]PostureLocation, len(checks.GeoLocationCheck.Locations))
		for idx, loc := range checks.GeoLocationCheck.Locations {
			locations[idx] = PostureLocation{
				CountryCode: loc.CountryCode,
				CityName:    loc.CityName,
			}
		}

		result.GeoLocation = &PostureGeoLocationCheck{
			Action:    PostureGeoLocationAction(checks.GeoLocationCheck.Action),
			Locations: locations,
		}
	}

	if checks.NbVersionCheck != nil {
		result.NbVersion = &PostureMinVersionCheck{MinVersion: checks.NbVersionCheck.MinVersion}
	}

	if checks.OsVersionCheck != nil {
		result.OsVersion = fromAPIOSVersionCheck(checks.OsVersionCheck)
	}

	if checks.PeerNetworkRangeCheck != nil {
		result.NetworkRange = &PosturePeerNetworkRangeCheck{
			Action: PosturePeerNetworkRangeAction(checks.PeerNetworkRangeCheck.Action),
			Ranges: checks.PeerNetworkRangeCheck.Ranges,
		}
	}

	if checks.ProcessCheck != nil {
		processes := make([]PostureProcess, len(checks.ProcessCheck.Processes))
		for idx, proc := range checks.ProcessCheck.Processes {
			processes[idx] = PostureProcess{
				LinuxPath:   proc.LinuxPath,
				MacPath:     proc.MacPath,
				WindowsPath: proc.WindowsPath,
			}
		}

		result.Process = &PostureProcessCheck{Processes: processes}
	}

	return result
}

// postureCheckStateFromAPI builds a PostureCheckState from an API PostureCheck response.
func postureCheckStateFromAPI(apiCheck *nbapi.PostureCheck) PostureCheckState {
	return PostureCheckState{
		Name:        apiCheck.Name,
		Description: apiCheck.Description,
		Checks:      fromAPIChecks(apiCheck.Checks),
	}
}

// buildPostureCheckRequest constructs an API PostureCheckUpdate from Pulumi inputs.
func buildPostureCheckRequest(args PostureCheckArgs) nbapi.PostureCheckUpdate {
	description := ""
	if args.Description != nil {
		description = *args.Description
	}

	apiChecks := toAPIChecks(args.Checks)

	return nbapi.PostureCheckUpdate{
		Name:        args.Name,
		Description: description,
		Checks:      &apiChecks,
	}
}

// Create creates a new NetBird posture check.
func (*PostureCheck) Create(ctx context.Context, req infer.CreateRequest[PostureCheckArgs]) (infer.CreateResponse[PostureCheckState], error) {
	p.GetLogger(ctx).Debugf("Create:PostureCheck name=%s", req.Inputs.Name)

	if req.DryRun {
		return infer.CreateResponse[PostureCheckState]{
			ID: "preview",
			Output: PostureCheckState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
				Checks:      req.Inputs.Checks,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.CreateResponse[PostureCheckState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	created, err := client.PostureChecks.Create(ctx, buildPostureCheckRequest(req.Inputs))
	if err != nil {
		return infer.CreateResponse[PostureCheckState]{}, fmt.Errorf("creating posture check failed: %w", err)
	}

	return infer.CreateResponse[PostureCheckState]{
		ID:     created.Id,
		Output: postureCheckStateFromAPI(created),
	}, nil
}

// Read reads a posture check from NetBird.
func (*PostureCheck) Read(ctx context.Context, req infer.ReadRequest[PostureCheckArgs, PostureCheckState]) (infer.ReadResponse[PostureCheckArgs, PostureCheckState], error) {
	p.GetLogger(ctx).Debugf("Read:PostureCheck[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.ReadResponse[PostureCheckArgs, PostureCheckState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiCheck, err := client.PostureChecks.Get(ctx, req.ID)
	if err != nil {
		return infer.ReadResponse[PostureCheckArgs, PostureCheckState]{}, fmt.Errorf("reading posture check failed: %w", err)
	}

	state := postureCheckStateFromAPI(apiCheck)

	if req.Inputs.Description == nil {
		state.Description = nil
	}

	return infer.ReadResponse[PostureCheckArgs, PostureCheckState]{
		ID:     req.ID,
		Inputs: PostureCheckArgs(state),
		State:  state,
	}, nil
}

// Update updates a posture check in NetBird.
func (*PostureCheck) Update(ctx context.Context, req infer.UpdateRequest[PostureCheckArgs, PostureCheckState]) (infer.UpdateResponse[PostureCheckState], error) {
	p.GetLogger(ctx).Debugf("Update:PostureCheck[%s]", req.ID)

	if req.DryRun {
		return infer.UpdateResponse[PostureCheckState]{
			Output: PostureCheckState{
				Name:        req.Inputs.Name,
				Description: req.Inputs.Description,
				Checks:      req.Inputs.Checks,
			},
		}, nil
	}

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.UpdateResponse[PostureCheckState]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	updated, err := client.PostureChecks.Update(ctx, req.ID, buildPostureCheckRequest(req.Inputs))
	if err != nil {
		return infer.UpdateResponse[PostureCheckState]{}, fmt.Errorf("updating posture check failed: %w", err)
	}

	return infer.UpdateResponse[PostureCheckState]{
		Output: postureCheckStateFromAPI(updated),
	}, nil
}

// Delete removes a posture check from NetBird.
func (*PostureCheck) Delete(ctx context.Context, req infer.DeleteRequest[PostureCheckState]) (infer.DeleteResponse, error) {
	p.GetLogger(ctx).Debugf("Delete:PostureCheck[%s]", req.ID)

	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	err = client.PostureChecks.Delete(ctx, req.ID)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("deleting posture check failed: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Diff detects changes between PostureCheckArgs and PostureCheckState.
func (*PostureCheck) Diff(ctx context.Context, req infer.DiffRequest[PostureCheckArgs, PostureCheckState]) (infer.DiffResponse, error) {
	p.GetLogger(ctx).Debugf("Diff:PostureCheck[%s]", req.ID)

	diff := map[string]p.PropertyDiff{}

	if req.Inputs.Name != req.State.Name {
		diff["name"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if req.Inputs.Description != nil && !equalPtr(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	if !equalChecks(req.Inputs.Checks, req.State.Checks) {
		diff["checks"] = p.PropertyDiff{InputDiff: false, Kind: p.Update}
	}

	p.GetLogger(ctx).Debugf("Diff:PostureCheck[%s] diff=%d", req.ID, len(diff))

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Check validates input fields for a posture check.
func (*PostureCheck) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[PostureCheckArgs], error) {
	p.GetLogger(ctx).Debugf("Check:PostureCheck old=%s, new=%s", req.OldInputs.GoString(), req.NewInputs.GoString())

	args, failures, err := infer.DefaultCheck[PostureCheckArgs](ctx, req.NewInputs)

	if isBlank(args.Name) {
		failures = append(failures, p.CheckFailure{
			Property: "name",
			Reason:   "name must not be empty",
		})
	}

	if args.Checks.GeoLocation != nil && len(args.Checks.GeoLocation.Locations) == 0 {
		failures = append(failures, p.CheckFailure{
			Property: "checks.geoLocationCheck.locations",
			Reason:   "at least one location is required when geoLocationCheck is set",
		})
	}

	if args.Checks.NetworkRange != nil && len(args.Checks.NetworkRange.Ranges) == 0 {
		failures = append(failures, p.CheckFailure{
			Property: "checks.peerNetworkRangeCheck.ranges",
			Reason:   "at least one range is required when peerNetworkRangeCheck is set",
		})
	}

	if args.Checks.Process != nil && len(args.Checks.Process.Processes) == 0 {
		failures = append(failures, p.CheckFailure{
			Property: "checks.processCheck.processes",
			Reason:   "at least one process is required when processCheck is set",
		})
	}

	return infer.CheckResponse[PostureCheckArgs]{
		Inputs:   args,
		Failures: failures,
	}, err
}

// WireDependencies explicitly defines input/output relationships.
func (*PostureCheck) WireDependencies(field infer.FieldSelector, args *PostureCheckArgs, state *PostureCheckState) {
	field.OutputField(&state.Name).DependsOn(field.InputField(&args.Name))
	field.OutputField(&state.Description).DependsOn(field.InputField(&args.Description))
	field.OutputField(&state.Checks).DependsOn(field.InputField(&args.Checks))
}

// equalChecks compares two PostureChecksConfig values for equality.
func equalChecks(checksA, checksB PostureChecksConfig) bool {
	return equalGeoLocationCheck(checksA.GeoLocation, checksB.GeoLocation) &&
		equalMinVersionCheck(checksA.NbVersion, checksB.NbVersion) &&
		equalOSVersionCheck(checksA.OsVersion, checksB.OsVersion) &&
		equalNetworkRangeCheck(checksA.NetworkRange, checksB.NetworkRange) &&
		equalProcessCheck(checksA.Process, checksB.Process)
}

func equalGeoLocationCheck(checkA, checkB *PostureGeoLocationCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	if checkA.Action != checkB.Action || len(checkA.Locations) != len(checkB.Locations) {
		return false
	}

	locA := slices.Clone(checkA.Locations)
	locB := slices.Clone(checkB.Locations)

	slices.SortFunc(locA, func(locationA, locationB PostureLocation) int {
		if locationA.CountryCode != locationB.CountryCode {
			if locationA.CountryCode < locationB.CountryCode {
				return -1
			}

			return 1
		}

		return 0
	})

	slices.SortFunc(locB, func(locationA, locationB PostureLocation) int {
		if locationA.CountryCode != locationB.CountryCode {
			if locationA.CountryCode < locationB.CountryCode {
				return -1
			}

			return 1
		}

		return 0
	})

	for idx := range locA {
		if locA[idx].CountryCode != locB[idx].CountryCode ||
			!equalPtr(locA[idx].CityName, locB[idx].CityName) {
			return false
		}
	}

	return true
}

func equalMinVersionCheck(checkA, checkB *PostureMinVersionCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	return checkA.MinVersion == checkB.MinVersion
}

func equalMinKernelVersionCheck(checkA, checkB *PostureMinKernelVersionCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	return checkA.MinKernelVersion == checkB.MinKernelVersion
}

func equalOSVersionCheck(checkA, checkB *PostureOSVersionCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	return equalMinVersionCheck(checkA.Android, checkB.Android) &&
		equalMinVersionCheck(checkA.Darwin, checkB.Darwin) &&
		equalMinVersionCheck(checkA.Ios, checkB.Ios) &&
		equalMinKernelVersionCheck(checkA.Linux, checkB.Linux) &&
		equalMinKernelVersionCheck(checkA.Windows, checkB.Windows)
}

func equalNetworkRangeCheck(checkA, checkB *PosturePeerNetworkRangeCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	return checkA.Action == checkB.Action && equalSlice(checkA.Ranges, checkB.Ranges)
}

func equalProcessCheck(checkA, checkB *PostureProcessCheck) bool {
	if checkA == nil && checkB == nil {
		return true
	}

	if checkA == nil || checkB == nil {
		return false
	}

	if len(checkA.Processes) != len(checkB.Processes) {
		return false
	}

	for idx := range checkA.Processes {
		if !equalPtr(checkA.Processes[idx].LinuxPath, checkB.Processes[idx].LinuxPath) ||
			!equalPtr(checkA.Processes[idx].MacPath, checkB.Processes[idx].MacPath) ||
			!equalPtr(checkA.Processes[idx].WindowsPath, checkB.Processes[idx].WindowsPath) {
			return false
		}
	}

	return true
}
