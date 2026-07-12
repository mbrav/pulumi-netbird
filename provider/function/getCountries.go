package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetCountries lists all countries known to the NetBird geo-location database.
type GetCountries struct{}

// Annotate describes the function.
func (f *GetCountries) Annotate(a infer.Annotator) {
	a.Describe(f, "List all countries known to NetBird's geo-location database. Useful for "+
		"populating PostureCheck geo-location rules.")
}

// GetCountriesArgs are the inputs for GetCountries (none).
type GetCountriesArgs struct{}

// Country is a single country entry.
type Country struct {
	CountryCode string `pulumi:"countryCode"`
	CountryName string `pulumi:"countryName"`
}

// Annotate provides field descriptions for Country.
func (c *Country) Annotate(ann infer.Annotator) {
	ann.Describe(&c.CountryCode, "2-letter ISO 3166-1 alpha-2 country code.")
	ann.Describe(&c.CountryName, "Commonly used English name of the country.")
}

// GetCountriesResult is the output of GetCountries.
type GetCountriesResult struct {
	Countries []Country `pulumi:"countries"`
}

// Annotate provides field descriptions for GetCountriesResult.
func (r *GetCountriesResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Countries, "The list of countries.")
}

// Invoke lists all countries.
func (f *GetCountries) Invoke(
	ctx context.Context,
	_ infer.FunctionRequest[GetCountriesArgs],
) (infer.FunctionResponse[GetCountriesResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[GetCountriesResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiCountries, err := client.GeoLocation.ListCountries(ctx)
	if err != nil {
		return infer.FunctionResponse[GetCountriesResult]{}, fmt.Errorf("listing countries failed: %w", err)
	}

	countries := make([]Country, 0, len(apiCountries))

	for _, c := range apiCountries {
		countries = append(countries, Country{
			CountryCode: c.CountryCode,
			CountryName: c.CountryName,
		})
	}

	return infer.FunctionResponse[GetCountriesResult]{
		Output: GetCountriesResult{
			Countries: countries,
		},
	}, nil
}
