package app

var YamlValuesMap = map[string]string{
	"luke":   lukeValuesYaml,
	"leia":   leiaValuesYaml,
	"anakin": anakinValuesYaml,
}

var lukeValuesYaml = `
injector:
   enabled: false
server:
   affinity: ""
   ha:
      enabled: true
      raft: 
         enabled: true
         setNodeId: true
         config: |
            cluster_name = "luke"
            storage "raft" {
                path    = "/vault/data/"
                retry_join {
                  leader_api_addr = "http://luke-vault-0.luke-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://luke-vault-1.luke-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://luke-vault-2.luke-vault-internal.default.svc.cluster.local:8200"
                }
             }
            listener "tcp" {
               address = "[::]:8200"
               cluster_address = "[::]:8201"
               tls_disable = "true"
            }
            service_registration "kubernetes" {}
`

var leiaValuesYaml = `
injector:
   enabled: false
server:
   affinity: ""
   ha:
      enabled: true
      raft: 
         enabled: true
         setNodeId: true
         config: |
            cluster_name = "leia"
            storage "raft" {
                path    = "/vault/data/"
                retry_join {
                  leader_api_addr = "http://leia-vault-0.leia-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://leia-vault-1.leia-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://leia-vault-2.leia-vault-internal.default.svc.cluster.local:8200"
                }
             }
            listener "tcp" {
               address = "[::]:8200"
               cluster_address = "[::]:8201"
               tls_disable = "true"
            }
            service_registration "kubernetes" {}
`

var anakinValuesYaml = `
injector:
   enabled: false
server:
   affinity: ""
   ha:
      enabled: true
      raft: 
         enabled: true
         setNodeId: true
         config: |
            cluster_name = "anakin"
            storage "raft" {
                path    = "/vault/data/"
                retry_join {
                  leader_api_addr = "http://anakin-vault-0.anakin-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://anakin-vault-1.anakin-vault-internal.default.svc.cluster.local:8200"
                }
                retry_join {
                  leader_api_addr = "http://anakin-vault-2.anakin-vault-internal.default.svc.cluster.local:8200"
                }
             }
            listener "tcp" {
               address = "[::]:8200"
               cluster_address = "[::]:8201"
               tls_disable = "true"
            }
            service_registration "kubernetes" {}
`
