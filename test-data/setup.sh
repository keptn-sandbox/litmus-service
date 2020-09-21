
keptn create project litmus --shipyard=./shipyard.yaml

keptn onboard service carts-db --chart=./carts-db/ --project=litmus

keptn onboard service carts --chart=./carts --project=litmus

keptn send event new-artifact --project=litmus --service=carts-db --image=docker.io/mongo --tag=4.2.2

keptn send event new-artifact --project=litmus --service=carts --image=docker.io/keptnexamples/carts --tag=0.11.1

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=slo.yaml --resourceUri=slo.yaml

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.5/deploy/service.yaml

keptn configure monitoring prometheus --project=litmus --service=carts

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/0.2.2/deploy/service.yaml

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=prometheus/sli.yaml --resourceUri=prometheus/sli.yaml 


## now we also have to add the chaos tests - ATTENTION right now this file is empty!
keptn add-resource --project=litmus --stage=chaos --service=carts --resource=litmus/experiment.yaml --resourceUri=litmus/experiment.yaml


keptn send event new-artifact --project=litmus --service=carts --image=docker.io/keptnexamples/carts --tag=0.11.1

