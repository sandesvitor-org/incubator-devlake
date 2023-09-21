### Observações

O nome da conection precisa ser o nome da org

### API KEY OPSGENIE

US: d806cd19-1b31-4976-b43f-61a2615c9396
EU: 550be78f-46fd-4dc8-aff2-8ea1c320746d
### mensagem mock pro Kafka

{"component": "stone/exampleA","team": "dlpco/all","repository": "dlpco/cobaia-devops", "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}

{"component": "stone/exampleB","team": "dlpco/all","repository": "dlpco/terraform-pager-duty",  "incidentManagement": {"type": "PAGERDUTY", "application": "DevLake", "instance": "BANKING"}, "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}

{"component": "stone/exampleC","team": "dlpco/all","repository": "dlpco/infra-devops-assets",  "incidentManagement": {"type": "OPSGENIE", "application": "devlake", "instance": "BANKING"}, "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}
