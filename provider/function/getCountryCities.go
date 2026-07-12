package function

import (
	"context"
	"fmt"

	"github.com/mbrav/pulumi-netbird/provider/config"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GetCountryCities lists the cities within a given country from NetBird's geo-location database.
type GetCountryCities struct{}

// Annotate describes the function.
func (f *GetCountryCities) Annotate(a infer.Annotator) {
	a.Describe(f, "List the cities within a country from NetBird's geo-location database. Useful "+
		"for populating PostureCheck geo-location rules with specific cities.")
}

// GetCountryCitiesArgs are the inputs for GetCountryCities.
type GetCountryCitiesArgs struct {
	CountryCode string `pulumi:"countryCode"`
}

// Annotate provides field descriptions for GetCountryCitiesArgs.
func (a *GetCountryCitiesArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.CountryCode, "2-letter ISO 3166-1 alpha-2 country code to list cities for.")
}

// City is a single city entry.
type City struct {
	CityName  string `pulumi:"cityName"`
	GeonameID int    `pulumi:"geonameId"`
}

// Annotate provides field descriptions for City.
func (c *City) Annotate(ann infer.Annotator) {
	ann.Describe(&c.CityName, "Commonly used English name of the city.")
	ann.Describe(&c.GeonameID, "Integer ID of the record in the GeoNames database.")
}

// GetCountryCitiesResult is the output of GetCountryCities.
type GetCountryCitiesResult struct {
	Cities []City `pulumi:"cities"`
}

// Annotate provides field descriptions for GetCountryCitiesResult.
func (r *GetCountryCitiesResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Cities, "The list of cities in the country.")
}

// Invoke lists the cities in the given country.
func (f *GetCountryCities) Invoke(
	ctx context.Context,
	req infer.FunctionRequest[GetCountryCitiesArgs],
) (infer.FunctionResponse[GetCountryCitiesResult], error) {
	client, err := config.GetNetBirdClient(ctx)
	if err != nil {
		return infer.FunctionResponse[GetCountryCitiesResult]{}, fmt.Errorf("error getting NetBird client: %w", err)
	}

	apiCities, err := client.GeoLocation.ListCountryCities(ctx, req.Input.CountryCode)
	if err != nil {
		return infer.FunctionResponse[GetCountryCitiesResult]{}, fmt.Errorf("listing cities failed: %w", err)
	}

	cities := make([]City, 0, len(apiCities))

	for _, c := range apiCities {
		cities = append(cities, City{
			CityName:  c.CityName,
			GeonameID: c.GeonameId,
		})
	}

	return infer.FunctionResponse[GetCountryCitiesResult]{
		Output: GetCountryCitiesResult{
			Cities: cities,
		},
	}, nil
}
