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
			TargetID:     "peer-1",
			Enabled:      true,
			Host:         nil,
			Port:         8080,
			Protocol:     ReverseProxyTargetProtocolHTTP,
			TargetType:   ReverseProxyTargetTypePeer,
			Path:         nil,
			Options:      nil,
			AccessAction: nil,
		},
	}

	// Simulates what fromAPITargetOptions produces after a round trip through
	// the mock/real API: a non-nil pointer to an all-zero-fields struct.
	stateTargets := []ReverseProxyTarget{
		{
			TargetID:     "peer-1",
			Enabled:      true,
			Host:         nil,
			Port:         8080,
			Protocol:     ReverseProxyTargetProtocolHTTP,
			TargetType:   ReverseProxyTargetTypePeer,
			Path:         nil,
			Options:      &ReverseProxyTargetOptions{}, //nolint:exhaustruct,exhaustruct_v5
			AccessAction: nil,
		},
	}

	if !equalReverseProxyTargets(inputTargets, stateTargets) {
		t.Fatal("expected targets with nil vs zero-value options to compare equal")
	}
}

// TestEqualReverseProxyTargetsNormalizesAccessAction pins the nil-vs-default
// handling for the per-target access action added in netbird v0.81.0. A program
// that omits accessAction leaves it nil, while the API reports the effective
// value; without normalizing nil to "inherit" that pair would diff on every
// preview forever — the same class of bug fixed for options in 0.5.4 and for the
// identity headers in 0.5.7.
func TestEqualReverseProxyTargetsNormalizesAccessAction(t *testing.T) {
	t.Parallel()

	target := func(action *ReverseProxyTargetAccessAction) []ReverseProxyTarget {
		return []ReverseProxyTarget{
			{
				TargetID:     "peer-1",
				Enabled:      true,
				Host:         nil,
				Port:         8080,
				Protocol:     ReverseProxyTargetProtocolHTTP,
				TargetType:   ReverseProxyTargetTypePeer,
				Path:         nil,
				Options:      nil,
				AccessAction: action,
			},
		}
	}

	inherit := ReverseProxyTargetAccessActionInherit
	block := ReverseProxyTargetAccessActionBlock

	t.Run("nil equals an API-reported inherit", func(t *testing.T) {
		t.Parallel()

		if !equalReverseProxyTargets(target(nil), target(&inherit)) {
			t.Fatal("expected a nil accessAction to equal an explicit \"inherit\"")
		}
	})

	t.Run("nil differs from a real action", func(t *testing.T) {
		t.Parallel()

		if equalReverseProxyTargets(target(nil), target(&block)) {
			t.Fatal("expected a nil accessAction to differ from \"block\"")
		}
	})

	t.Run("differing actions are detected", func(t *testing.T) {
		t.Parallel()

		if equalReverseProxyTargets(target(&inherit), target(&block)) {
			t.Fatal("expected \"inherit\" and \"block\" to differ")
		}
	})
}
