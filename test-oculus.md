### Observações

O nome da conection precisa ser o nome da org

### mensagem mock pro Kafka

{"component": "stone/exampleA","team": "dlpco/all","repository": "dlpco/cobaia-devops", "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}

{"component": "stone/exampleB","team": "dlpco/all","repository": "dlpco/terraform-pager-duty",  "incidentManagement": {"type": "PAGERDUTY", "application": "DevLake", "instance": "BANKING"}, "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}

{"component": "stone/exampleC","team": "dlpco/all","repository": "dlpco/infra-devops-assets",  "incidentManagement": {"type": "OPSGENIE", "application": "devlake", "instance": "BANKING"}, "pipeline": {"type":"GITHUB_ACTIONS", "productionDeployJob": ""}}
