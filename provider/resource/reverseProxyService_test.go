//nolint:testpackage // exercises unexported equalOptionalDeep/equalReverseProxyTargets directly
package resource

import "testing"

// Regression test for the perpetual-diff bug fixed in v0.5.4: the NetBird API
// always echoes a full struct for optional fields like target.options and
// auth (all sub-fields nil/zero when unconfigured), while a Pulumi program
// that never sets them leaves the field nil. A raw reflect.DeepEqual treated
// nil and an all-zero-fields struct as different forever.
func TestEqualOptionalDeep(t *testing.T) {
	t.Parallel()

	t.Run("nil options equals API-echoed zero-value options", func(t *testing.T) {
		t.Parallel()

		var nilOptions *ReverseProxyTargetOptions

		zeroOptions := &ReverseProxyTargetOptions{
			CustomHeaders:      nil,
			DirectUpstream:     nil,
			PathRewrite:        nil,
			ProxyProtocol:      nil,
			RequestTimeout:     nil,
			SessionIdleTimeout: nil,
			SkipTLSVerify:      nil,
		}

		if !equalOptionalDeep(nilOptions, zeroOptions) {
			t.Fatal("expected nil options to equal an all-zero-fields options struct")
		}
	})

	t.Run("nil auth equals API-echoed zero-value auth", func(t *testing.T) {
		t.Parallel()

		var nilAuth *ReverseProxyAuth

		zeroAuth := &ReverseProxyAuth{
			BearerAuth:   nil,
			HeaderAuths:  nil,
			LinkAuth:     nil,
			PasswordAuth: nil,
			PinAuth:      nil,
		}

		if !equalOptionalDeep(nilAuth, zeroAuth) {
			t.Fatal("expected nil auth to equal an all-zero-fields auth struct")
		}
	})

	t.Run("nil vs nil is equal", func(t *testing.T) {
		t.Parallel()

		var a, b *ReverseProxyTargetOptions

		if !equalOptionalDeep(a, b) {
			t.Fatal("expected two nil pointers to be equal")
		}
	})

	t.Run("genuinely different values are still detected", func(t *testing.T) {
		t.Parallel()

		enabled := true

		configured := &ReverseProxyTargetOptions{
			CustomHeaders:      nil,
			DirectUpstream:     &enabled,
			PathRewrite:        nil,
			ProxyProtocol:      nil,
			RequestTimeout:     nil,
			SessionIdleTimeout: nil,
			SkipTLSVerify:      nil,
		}

		var nilOptions *ReverseProxyTargetOptions

		if equalOptionalDeep(nilOptions, configured) {
			t.Fatal("expected nil options to differ from options with a configured field")
		}
	})
}

func TestEqualReverseProxyTargetsIgnoresNilVsZeroOptions(t *testing.T) {
	t.Parallel()

	inputTargets := []ReverseProxyTarget{
		{
			TargetID:   "peer-1",
			Enabled:    true,
			Host:       nil,
			Port:       8080,
			Protocol:   ReverseProxyTargetProtocolHTTP,
			TargetType: ReverseProxyTargetTypePeer,
			Path:       nil,
			Options:    nil,
		},
	}

	// Simulates what fromAPITargetOptions produces after a round trip through
	// the mock/real API: a non-nil pointer to an all-zero-fields struct.
	stateTargets := []ReverseProxyTarget{
		{
			TargetID:   "peer-1",
			Enabled:    true,
			Host:       nil,
			Port:       8080,
			Protocol:   ReverseProxyTargetProtocolHTTP,
			TargetType: ReverseProxyTargetTypePeer,
			Path:       nil,
			Options:    &ReverseProxyTargetOptions{}, //nolint:exhaustruct,exhaustruct_v5
		},
	}

	if !equalReverseProxyTargets(inputTargets, stateTargets) {
		t.Fatal("expected targets with nil vs zero-value options to compare equal")
	}
}
