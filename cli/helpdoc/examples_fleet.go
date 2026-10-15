package helpdoc

// Curated help examples: fleet and SSH (nodes, ssh, cluster, ping,
// clone-fix-remediate). See examples.go — one registry, single source.
func registerFleetExamples(reg map[string]HelpExample) {
	reg["nodes"] = HelpExample{
		WhatItDoes: "Unified fleet view: aggregates nodes across SSH, the Cluster DB, and server-client networks with live reachability probing and machine ping.",
		Examples: []Example{
			{
				Invocation: "gitmap nodes",
				ExpectedOutput: `node   role     os       reachable  rtt
u1     worker   ubuntu   yes        12ms
<mchn> local    linux    yes        0ms`,
			},
		},
	}
	reg["node"] = reg["nodes"]

	reg["ssh"] = HelpExample{
		WhatItDoes: "SSH key pair generation, credential management, host enrollment, remote execution, and cluster connectivity. Use ssh-join (sj) to enroll a new machine.",
		Examples: []Example{
			{
				Invocation:     `gitmap ssh u1 "gitmap pa --status"`,
				ExpectedOutput: "remote pull-all status table from node u1",
			},
		},
	}

	reg["ssh-join"] = HelpExample{
		WhatItDoes: "Enrolls a machine into the SSH fleet by user@ip or IP: registers the host alias, authorizes your public key, and makes the host instantly recallable as a short alias.",
		Examples: []Example{
			{
				Invocation: "gitmap ssh-join deploy@192.168.1.20 --alias u2",
				ExpectedOutput: `✓ host u2 (192.168.1.20) enrolled
✓ public key authorized — passwordless login ready`,
			},
		},
	}
	reg["sj"] = reg["ssh-join"]

	reg["cluster"] = HelpExample{
		WhatItDoes: "Comprehensive cluster management: multi-node topology, remote execution, script deployment, and Ubuntu provisioning across the fleet.",
		Examples: []Example{
			{
				Invocation: "gitmap cluster status",
				ExpectedOutput: `cluster: 3 nodes — 3 reachable
u1  worker   12ms   heartbeat <date>
u2  worker   18ms   heartbeat <date>`,
			},
		},
	}
	reg["cluster-nodes"] = HelpExample{
		WhatItDoes: "Lists registered nodes in the cluster database with display ID, role, OS, and reachability status.",
		Examples: []Example{
			{
				Invocation: "gitmap cluster-nodes",
				ExpectedOutput: `id  alias  role     os      reachable
1   u1     worker   ubuntu  yes
2   u2     worker   ubuntu  yes`,
			},
		},
	}
	reg["cluster-add"] = HelpExample{
		WhatItDoes: "Enrolls a new node into the cluster topology and SSH host registry with alias registration and key authorization.",
		Examples: []Example{
			{
				Invocation: "gitmap cluster-add deploy@192.168.1.21 --alias u3",
				ExpectedOutput: `✓ node u3 enrolled (id 3)
✓ SSH key authorized`,
			},
		},
	}
	reg["cluster-status"] = HelpExample{
		WhatItDoes: "Inspects cluster connectivity, node health, heartbeat records, and ping round-trip latency.",
		Examples: []Example{
			{
				Invocation: "gitmap cluster-status",
				ExpectedOutput: `u1  reachable  rtt 12ms  heartbeat <date>
u2  reachable  rtt 18ms  heartbeat <date>`,
			},
		},
	}
	reg["cluster-remove"] = HelpExample{
		WhatItDoes: "Removes and unregisters a node from the cluster topology and SSH host registry.",
		Examples: []Example{
			{
				Invocation:     "gitmap cluster-remove u3",
				ExpectedOutput: "✓ node u3 removed from cluster topology and SSH registry",
			},
		},
	}

	reg["ping"] = HelpExample{
		WhatItDoes: "Runs machine ping and network reachability probes concurrently against all fleet nodes with a tabular terminal display.",
		Examples: []Example{
			{
				Invocation: "gitmap ping",
				ExpectedOutput: `node  reachable  rtt
u1    yes        12ms
u2    yes        18ms`,
			},
		},
	}
	reg["nodes-ping"] = reg["ping"]

	reg["cfr"] = HelpExample{
		WhatItDoes: "Clone, remediate, and synchronize a repository locally or across fleet nodes: clones the repo, fixes remotes and tracking, and syncs state.",
		Examples: []Example{
			{
				Invocation: "gitmap cfr alimtvnetwork/story-telling-v1",
				ExpectedOutput: `✓ cloned alimtvnetwork/story-telling-v1
✓ remotes fixed   ✓ tracking synced`,
			},
		},
	}
	reg["nodes-cfr"] = HelpExample{
		WhatItDoes: "Clone, remediate, and synchronize repositories across the multi-node SSH cluster and the local host, with pre-flight version discovery.",
		Examples: []Example{
			{
				Invocation: "gitmap nodes-cfr alimtvnetwork/gitmap-v28",
				ExpectedOutput: `u1  ✓ cloned+remediated
u2  ✓ cloned+remediated
local ✓ cloned+remediated`,
			},
		},
	}
}
