//nolint:testpackage // exercises the unexported fromAPINetworkAddresses directly
package resource

import (
	"testing"

	nbapi "github.com/netbirdio/netbird/shared/management/http/api"
)

// TestFromAPINetworkAddresses covers the peer interface list netbird v0.81.0
// added (Peer.network_addresses). The field is optional on the wire, so an
// absent list must stay absent rather than becoming an empty one — otherwise a
// peer on an older management server would show networkAddresses: [] instead of
// nothing.
func TestFromAPINetworkAddresses(t *testing.T) {
	t.Parallel()

	t.Run("nil stays nil", func(t *testing.T) {
		t.Parallel()

		if got := fromAPINetworkAddresses(nil); got != nil {
			t.Fatalf("expected nil, got %v", *got)
		}
	})

	t.Run("empty list stays empty, not nil", func(t *testing.T) {
		t.Parallel()

		got := fromAPINetworkAddresses(&[]nbapi.NetworkAddress{})
		if got == nil {
			t.Fatal("expected a non-nil pointer for an empty list")
		}

		if len(*got) != 0 {
			t.Fatalf("expected an empty list, got %v", *got)
		}
	})

	t.Run("fields map in order", func(t *testing.T) {
		t.Parallel()

		got := fromAPINetworkAddresses(&[]nbapi.NetworkAddress{
			{Mac: "aa:bb:cc:dd:ee:ff", NetIp: "10.10.1.5/24"},
			{Mac: "11:22:33:44:55:66", NetIp: "192.168.1.9/24"},
		})
		if got == nil {
			t.Fatal("expected a non-nil result")
		}

		if len(*got) != 2 {
			t.Fatalf("expected 2 addresses, got %d", len(*got))
		}

		if (*got)[0].Mac != "aa:bb:cc:dd:ee:ff" || (*got)[0].NetIP != "10.10.1.5/24" {
			t.Fatalf("first address mapped wrong: %+v", (*got)[0])
		}

		if (*got)[1].Mac != "11:22:33:44:55:66" || (*got)[1].NetIP != "192.168.1.9/24" {
			t.Fatalf("second address mapped wrong: %+v", (*got)[1])
		}
	})
}
