// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package resource

import (
	"context"
	"reflect"

	"errors"
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A NetBird network.
type DNS struct {
	pulumi.CustomResourceState

	// Description of the nameserver group
	Description pulumi.StringOutput `pulumi:"description"`
	// Domains Match domain list. It should be empty only if primary is true.
	Domains pulumi.StringArrayOutput `pulumi:"domains"`
	// Enabled Nameserver group status
	Enabled pulumi.BoolOutput `pulumi:"enabled"`
	// Groups Distribution group IDs that defines group of peers that will use this nameserver group
	Groups pulumi.StringArrayOutput `pulumi:"groups"`
	// Name of nameserver group name
	Name pulumi.StringOutput `pulumi:"name"`
	// Nameservers Nameserver list
	Nameservers NameserverArrayOutput `pulumi:"nameservers"`
	// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
	Primary pulumi.BoolOutput `pulumi:"primary"`
	// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
	Search_domains_enabled pulumi.BoolOutput `pulumi:"search_domains_enabled"`
}

// NewDNS registers a new resource with the given unique name, arguments, and options.
func NewDNS(ctx *pulumi.Context,
	name string, args *DNSArgs, opts ...pulumi.ResourceOption) (*DNS, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Description == nil {
		return nil, errors.New("invalid value for required argument 'Description'")
	}
	if args.Domains == nil {
		return nil, errors.New("invalid value for required argument 'Domains'")
	}
	if args.Enabled == nil {
		return nil, errors.New("invalid value for required argument 'Enabled'")
	}
	if args.Groups == nil {
		return nil, errors.New("invalid value for required argument 'Groups'")
	}
	if args.Name == nil {
		return nil, errors.New("invalid value for required argument 'Name'")
	}
	if args.Nameservers == nil {
		return nil, errors.New("invalid value for required argument 'Nameservers'")
	}
	if args.Primary == nil {
		return nil, errors.New("invalid value for required argument 'Primary'")
	}
	if args.Search_domains_enabled == nil {
		return nil, errors.New("invalid value for required argument 'Search_domains_enabled'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource DNS
	err := ctx.RegisterResource("netbird:resource:DNS", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetDNS gets an existing DNS resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetDNS(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *DNSState, opts ...pulumi.ResourceOption) (*DNS, error) {
	var resource DNS
	err := ctx.ReadResource("netbird:resource:DNS", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering DNS resources.
type dnsState struct {
}

type DNSState struct {
}

func (DNSState) ElementType() reflect.Type {
	return reflect.TypeOf((*dnsState)(nil)).Elem()
}

type dnsArgs struct {
	// Description of the nameserver group
	Description string `pulumi:"description"`
	// Domains Match domain list. It should be empty only if primary is true.
	Domains []string `pulumi:"domains"`
	// Enabled Nameserver group status
	Enabled bool `pulumi:"enabled"`
	// Groups Distribution group IDs that defines group of peers that will use this nameserver group
	Groups []string `pulumi:"groups"`
	// Name of nameserver group name
	Name string `pulumi:"name"`
	// Nameservers Nameserver list
	Nameservers []Nameserver `pulumi:"nameservers"`
	// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
	Primary bool `pulumi:"primary"`
	// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
	Search_domains_enabled bool `pulumi:"search_domains_enabled"`
}

// The set of arguments for constructing a DNS resource.
type DNSArgs struct {
	// Description of the nameserver group
	Description pulumi.StringInput
	// Domains Match domain list. It should be empty only if primary is true.
	Domains pulumi.StringArrayInput
	// Enabled Nameserver group status
	Enabled pulumi.BoolInput
	// Groups Distribution group IDs that defines group of peers that will use this nameserver group
	Groups pulumi.StringArrayInput
	// Name of nameserver group name
	Name pulumi.StringInput
	// Nameservers Nameserver list
	Nameservers NameserverArrayInput
	// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
	Primary pulumi.BoolInput
	// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
	Search_domains_enabled pulumi.BoolInput
}

func (DNSArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*dnsArgs)(nil)).Elem()
}

type DNSInput interface {
	pulumi.Input

	ToDNSOutput() DNSOutput
	ToDNSOutputWithContext(ctx context.Context) DNSOutput
}

func (*DNS) ElementType() reflect.Type {
	return reflect.TypeOf((**DNS)(nil)).Elem()
}

func (i *DNS) ToDNSOutput() DNSOutput {
	return i.ToDNSOutputWithContext(context.Background())
}

func (i *DNS) ToDNSOutputWithContext(ctx context.Context) DNSOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DNSOutput)
}

// DNSArrayInput is an input type that accepts DNSArray and DNSArrayOutput values.
// You can construct a concrete instance of `DNSArrayInput` via:
//
//	DNSArray{ DNSArgs{...} }
type DNSArrayInput interface {
	pulumi.Input

	ToDNSArrayOutput() DNSArrayOutput
	ToDNSArrayOutputWithContext(context.Context) DNSArrayOutput
}

type DNSArray []DNSInput

func (DNSArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*DNS)(nil)).Elem()
}

func (i DNSArray) ToDNSArrayOutput() DNSArrayOutput {
	return i.ToDNSArrayOutputWithContext(context.Background())
}

func (i DNSArray) ToDNSArrayOutputWithContext(ctx context.Context) DNSArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DNSArrayOutput)
}

// DNSMapInput is an input type that accepts DNSMap and DNSMapOutput values.
// You can construct a concrete instance of `DNSMapInput` via:
//
//	DNSMap{ "key": DNSArgs{...} }
type DNSMapInput interface {
	pulumi.Input

	ToDNSMapOutput() DNSMapOutput
	ToDNSMapOutputWithContext(context.Context) DNSMapOutput
}

type DNSMap map[string]DNSInput

func (DNSMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*DNS)(nil)).Elem()
}

func (i DNSMap) ToDNSMapOutput() DNSMapOutput {
	return i.ToDNSMapOutputWithContext(context.Background())
}

func (i DNSMap) ToDNSMapOutputWithContext(ctx context.Context) DNSMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DNSMapOutput)
}

type DNSOutput struct{ *pulumi.OutputState }

func (DNSOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**DNS)(nil)).Elem()
}

func (o DNSOutput) ToDNSOutput() DNSOutput {
	return o
}

func (o DNSOutput) ToDNSOutputWithContext(ctx context.Context) DNSOutput {
	return o
}

// Description of the nameserver group
func (o DNSOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v *DNS) pulumi.StringOutput { return v.Description }).(pulumi.StringOutput)
}

// Domains Match domain list. It should be empty only if primary is true.
func (o DNSOutput) Domains() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *DNS) pulumi.StringArrayOutput { return v.Domains }).(pulumi.StringArrayOutput)
}

// Enabled Nameserver group status
func (o DNSOutput) Enabled() pulumi.BoolOutput {
	return o.ApplyT(func(v *DNS) pulumi.BoolOutput { return v.Enabled }).(pulumi.BoolOutput)
}

// Groups Distribution group IDs that defines group of peers that will use this nameserver group
func (o DNSOutput) Groups() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *DNS) pulumi.StringArrayOutput { return v.Groups }).(pulumi.StringArrayOutput)
}

// Name of nameserver group name
func (o DNSOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *DNS) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Nameservers Nameserver list
func (o DNSOutput) Nameservers() NameserverArrayOutput {
	return o.ApplyT(func(v *DNS) NameserverArrayOutput { return v.Nameservers }).(NameserverArrayOutput)
}

// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
func (o DNSOutput) Primary() pulumi.BoolOutput {
	return o.ApplyT(func(v *DNS) pulumi.BoolOutput { return v.Primary }).(pulumi.BoolOutput)
}

// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
func (o DNSOutput) Search_domains_enabled() pulumi.BoolOutput {
	return o.ApplyT(func(v *DNS) pulumi.BoolOutput { return v.Search_domains_enabled }).(pulumi.BoolOutput)
}

type DNSArrayOutput struct{ *pulumi.OutputState }

func (DNSArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*DNS)(nil)).Elem()
}

func (o DNSArrayOutput) ToDNSArrayOutput() DNSArrayOutput {
	return o
}

func (o DNSArrayOutput) ToDNSArrayOutputWithContext(ctx context.Context) DNSArrayOutput {
	return o
}

func (o DNSArrayOutput) Index(i pulumi.IntInput) DNSOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *DNS {
		return vs[0].([]*DNS)[vs[1].(int)]
	}).(DNSOutput)
}

type DNSMapOutput struct{ *pulumi.OutputState }

func (DNSMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*DNS)(nil)).Elem()
}

func (o DNSMapOutput) ToDNSMapOutput() DNSMapOutput {
	return o
}

func (o DNSMapOutput) ToDNSMapOutputWithContext(ctx context.Context) DNSMapOutput {
	return o
}

func (o DNSMapOutput) MapIndex(k pulumi.StringInput) DNSOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *DNS {
		return vs[0].(map[string]*DNS)[vs[1].(string)]
	}).(DNSOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*DNSInput)(nil)).Elem(), &DNS{})
	pulumi.RegisterInputType(reflect.TypeOf((*DNSArrayInput)(nil)).Elem(), DNSArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*DNSMapInput)(nil)).Elem(), DNSMap{})
	pulumi.RegisterOutputType(DNSOutput{})
	pulumi.RegisterOutputType(DNSArrayOutput{})
	pulumi.RegisterOutputType(DNSMapOutput{})
}
