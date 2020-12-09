
# prerequisite: have keptn installed with --use-case=continuous-delivery flag

# 1. install Istio https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#2 
# 2. install Keptn https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#4
# 3. configure Istio + Keptn  shttps://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#5
# 4. connect the Keptn CLI to the cluster https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-07/index.html?index=..%2F..index#6 

## CREATE PROJECT
keptn create project litmus --shipyard=./shipyard.yaml

## ONBOARD SERVICE
keptn onboard service helloservice --chart=./helloservice/helm/ --project=litmus

## ADD JMETER TESTS & CONFIG
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml

## ADD QUALITY GATE
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/prometheus/sli.yaml --resourceUri=prometheus/sli.yaml

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/slo.yaml --resourceUri=slo.yaml

## ADD LITMUS EXPERIMENT
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/litmus/experiment.yaml --resourceUri=litmus/experiment.yaml


## ADD PROMETHEUS
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.5/deploy/service.yaml

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/release-0.2.2/deploy/service.yaml

keptn configure monitoring prometheus --project=litmus --service=helloservice

kubectl apply -f helloservice/prometheus/blackbox-exporter.yaml

kubectl apply -f helloservice/prometheus/prometheus-server-conf-cm.yaml

kubectl delete pod -l app=prometheus-server -n monitoring

### LITMUS  Begins!! 

## install litmus operator & chaos CRDs 

kubectl apply -f litmus/litmus-operator-v1.9.1.yaml

## pull the chaos experiment CR (static) 

kubectl apply -f litmus/pod-delete-ChaosExperiment-CR.yaml 

## pull the chaos experiment RBAC (static) 

kubectl apply -f litmus/pod-delete-rbac.yaml 

# Litmus PreReq End!! 

## now we also have to add the chaos tests - ATTENTION right now this file is empty!
keptn add-resource --project=litmus --stage=chaos --service=carts --resource=helloservice/litmus/experiment.yaml --resourceUri=litmus/experiment.yaml


## first deployment event
keptn send event new-artifact --project=litmus --service=helloservice --image=jetzlstorfer/hello-server:v0.1.1 

## second deployment event (able to scale by editing the deploy-event.json)
keptn send event -f test-data/helloservice/deploy-event.json
