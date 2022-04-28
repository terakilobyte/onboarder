package globals

func GetReposForTeam(team string) map[string][]string {

	teamRepoMap := map[string]map[string][]string{
		"cet": {
			"10gen": {
				"mms-docs",
				"cloud-docs",
				"docs-charts",
				"docs-k8s-operator",
				"cloud-docs-osb",
				"docs-mongocli",
				"docs-datalake",
				"docs-kafka-connector",
				"docs-tutorials",
				"docs-mongodb-internal",
			},
			"mongodb": {
				"docs-assets",
				"docs",
				"docs-ecosystem",
				"docs-primer",
				"docs-compass",
				"docs-bi-connector",
				"docs-spark-connector",
				"docs-mongodb-shell",
				"mongodb-kubernetes-operator",
				"docs-realm",
				"docs-commandline-tools",
				"docs-worker-pool",
				"docs-tools",
				"docs-landing",
				"mongocli",
				"docs-atlas-cli",
			},
		},
		"tdbx": {
			"10gen": {
				"cloud-docs",
				"docs-tutorials",
				"docs-mongodb-internal",
			},
			"mongodb": {
				"docs-kafka-connector",
				"docs-ecosystem",
				"docs-spark-connector",
				"docs-worker-pool",
				"docs-tools",
				"docs-landing",
				"docs-node",
				"docs-java",
				"docs-java-other",
				"docs-visual-studio-extension",
				"docs-golang",
				"docs-ruby",
				"docs-php-library",
				"docs-mongoid",
			},
		},
	}
	return teamRepoMap[team]
}
