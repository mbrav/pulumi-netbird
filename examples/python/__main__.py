from pulumi_netbird import resource


net = resource.Network("net-test", name="Test Network", description="Test network")
