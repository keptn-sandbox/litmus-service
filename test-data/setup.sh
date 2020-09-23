
# prerequisite: have keptn installed with --use-case=continuous-delivery flag

# 1. install Istio https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#2 
# 2. install Keptn https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#4
# 3. configure Istio + Keptn  shttps://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#5
# 4. connect the Keptn CLI to the cluster https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#6 

keptn create project litmus --shipyard=./shipyard.yaml

keptn onboard service carts-db --chart=./carts-db/ --project=litmus

keptn onboard service carts --chart=./carts --project=litmus

keptn send event new-artifact --project=litmus --service=carts-db --image=docker.io/mongo --tag=4.2.2

keptn send event new-artifact --project=litmus --service=carts --image=docker.io/keptnexamples/carts --tag=0.11.1

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.5/deploy/service.yaml

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/release-0.2.2/deploy/service.yaml

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=prometheus/sli.yaml --resourceUri=prometheus/sli.yaml 

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=slo.yaml --resourceUri=slo.yaml

keptn configure monitoring prometheus --project=litmus --service=carts




## now we also have to add the chaos tests - ATTENTION right now this file is empty!
keptn add-resource --project=litmus --stage=chaos --service=carts --resource=litmus/experiment.yaml --resourceUri=litmus/experiment.yaml

## make sure that Litmus is installed in the cluster!


# now test with a a new-artifact event
keptn send event new-artifact --project=litmus --service=carts --image=docker.io/keptnexamples/carts --tag=0.11.3

