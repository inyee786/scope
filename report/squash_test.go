package report

import (
	"net"
	"reflect"
	"testing"
)

var (
	_, netdot1, _ = net.ParseCIDR("192.168.1.0/24")
	_, netdot2, _ = net.ParseCIDR("192.168.2.0/24")
)

func reportToSquash() Report {
	return Report{
		Endpoint: Topology{
			Adjacency: Adjacency{
				"hostA|;192.168.1.1;12345": []string{";192.168.1.2;80"},
				"hostA|;192.168.1.1;8888":  []string{";1.2.3.4;22", ";1.2.3.4;23"},
				"hostB|;192.168.1.2;80":    []string{";192.168.1.1;12345"},
				"hostZ|;192.168.2.2;80":    []string{";192.168.1.1;12345"},
			},
			EdgeMetadatas: EdgeMetadatas{
				";192.168.1.1;12345|;192.168.1.2;80": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  12,
					BytesIngress: 0,
				},
				";192.168.1.1;8888|;1.2.3.4;22": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  200,
					BytesIngress: 0,
				},
				";192.168.1.1;8888|;1.2.3.4;23": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  200,
					BytesIngress: 0,
				},
				";192.168.1.2;80|;192.168.1.1;12345": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
				";192.168.2.2;80|;192.168.1.1;12345": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
			},
			NodeMetadatas: NodeMetadatas{
				";192.168.1.1;12345": NodeMetadata{
					"pid":    "23128",
					"name":   "curl",
					"domain": "node-a.local",
				},
				";192.168.1.1;8888": NodeMetadata{
					"pid":    "55100",
					"name":   "ssh",
					"domain": "node-a.local",
				},
				";192.168.1.2;80": NodeMetadata{
					"pid":    "215",
					"name":   "apache",
					"domain": "node-b.local",
				},
				";192.168.2.2;80": NodeMetadata{
					"pid":    "213",
					"name":   "apache",
					"domain": "node-z.local",
				},
			},
		},

		Address: Topology{
			Adjacency: Adjacency{
				"hostA|;192.168.1.1": []string{";192.168.1.2", ";1.2.3.4"},
				"hostB|;192.168.1.2": []string{";192.168.1.1"},
				"hostZ|;192.168.2.2": []string{";192.168.1.1"},
			},
			EdgeMetadatas: EdgeMetadatas{
				";192.168.1.1|;192.168.1.2": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  12,
					BytesIngress: 0,
				},
				";192.168.1.1|;1.2.3.4": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  200,
					BytesIngress: 0,
				},
				";192.168.1.2|;192.168.1.1": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
				";192.168.2.2|;192.168.1.1": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
			},
			NodeMetadatas: NodeMetadatas{
				";192.168.1.1": NodeMetadata{
					"name": "host-a",
				},
				";192.168.1.2": NodeMetadata{
					"name": "host-b",
				},
				";192.168.2.2": NodeMetadata{
					"name": "host-z",
				},
			},
		},

		Host: Topology{
			Adjacency:     Adjacency{},
			EdgeMetadatas: EdgeMetadatas{},
			NodeMetadatas: NodeMetadatas{
				MakeHostNodeID("hostA"): NodeMetadata{
					"host_name":      "node-a.local",
					"os":             "Linux",
					"local_networks": netdot1.String(),
				},
				MakeHostNodeID("hostB"): NodeMetadata{
					"host_name":      "node-b.local",
					"os":             "Linux",
					"local_networks": netdot1.String(),
				},
				MakeHostNodeID("hostZ"): NodeMetadata{
					"host_name":      "node-z.local",
					"os":             "Linux",
					"local_networks": netdot2.String(),
				},
			},
		},
	}
}

func TestSquashTopology(t *testing.T) {
	// Tests just a topology
	want := Topology{
		Adjacency: Adjacency{
			"hostA|;192.168.1.1;12345": []string{";192.168.1.2;80"},
			"hostA|;192.168.1.1;8888":  []string{"theinternet"},
			"hostB|;192.168.1.2;80":    []string{";192.168.1.1;12345"},
			"hostZ|;192.168.2.2;80":    []string{";192.168.1.1;12345"},
		},
		EdgeMetadatas: EdgeMetadatas{
			";192.168.1.1;12345|;192.168.1.2;80": EdgeMetadata{
				WithBytes:    true,
				BytesEgress:  12,
				BytesIngress: 0,
			},
			";192.168.1.1;8888|theinternet": EdgeMetadata{
				WithBytes:    true,
				BytesEgress:  2 * 200,
				BytesIngress: 2 * 0,
			},
			";192.168.1.2;80|;192.168.1.1;12345": EdgeMetadata{
				WithBytes:    true,
				BytesEgress:  0,
				BytesIngress: 12,
			},
			";192.168.2.2;80|;192.168.1.1;12345": EdgeMetadata{
				WithBytes:    true,
				BytesEgress:  0,
				BytesIngress: 12,
			},
		},
		NodeMetadatas: reportToSquash().Endpoint.NodeMetadatas,
	}

	have := Squash(reportToSquash().Endpoint, EndpointIDAddresser, reportToSquash().LocalNetworks())
	if !reflect.DeepEqual(want, have) {
		t.Errorf("want\n\t%#v, have\n\t%#v", want, have)
	}
}

func TestSquashReport(t *testing.T) {
	// Tests a full report squash.
	want := Report{
		Endpoint: Topology{
			Adjacency: Adjacency{
				"hostA|;192.168.1.1;12345": []string{";192.168.1.2;80"},
				"hostA|;192.168.1.1;8888":  []string{"theinternet"},
				"hostB|;192.168.1.2;80":    []string{";192.168.1.1;12345"},
				"hostZ|;192.168.2.2;80":    []string{";192.168.1.1;12345"},
			},
			EdgeMetadatas: EdgeMetadatas{
				";192.168.1.1;12345|;192.168.1.2;80": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  12,
					BytesIngress: 0,
				},
				";192.168.1.1;8888|theinternet": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  2 * 200,
					BytesIngress: 2 * 0,
				},
				";192.168.1.2;80|;192.168.1.1;12345": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
				";192.168.2.2;80|;192.168.1.1;12345": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
			},
			NodeMetadatas: reportToSquash().Endpoint.NodeMetadatas,
		},
		Address: Topology{
			Adjacency: Adjacency{
				"hostA|;192.168.1.1": []string{";192.168.1.2", "theinternet"},
				"hostB|;192.168.1.2": []string{";192.168.1.1"},
				"hostZ|;192.168.2.2": []string{";192.168.1.1"},
			},
			EdgeMetadatas: EdgeMetadatas{
				";192.168.1.1|;192.168.1.2": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  12,
					BytesIngress: 0,
				},
				";192.168.1.1|theinternet": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  200,
					BytesIngress: 0,
				},
				";192.168.1.2|;192.168.1.1": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
				";192.168.2.2|;192.168.1.1": EdgeMetadata{
					WithBytes:    true,
					BytesEgress:  0,
					BytesIngress: 12,
				},
			},
			NodeMetadatas: NodeMetadatas{
				";192.168.1.1": NodeMetadata{
					"name": "host-a",
				},
				";192.168.1.2": NodeMetadata{
					"name": "host-b",
				},
				";192.168.2.2": NodeMetadata{
					"name": "host-z",
				},
			},
		},
		Host: reportToSquash().Host,
	}

	have := reportToSquash().SquashRemote()
	if !reflect.DeepEqual(want, have) {
		t.Error(diff(want, have))
	}
}
