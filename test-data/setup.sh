

#######################################################################
# THIS FILE IS CURRENTLY NOT MAINTAINED
# PLEASE FOLLOW THE TUTORIAL ON HTTPS://TUTORIALS.KEPTN.SH TO USE THIS
#######################################################################



# prerequisite: have keptn installed with --use-case=continuous-delivery flag

# 1. install Istio https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-08/index.html?index=..%2F..index#2 
# 2. install Keptn https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-08/index.html?index=..%2F..index#4
# 3. configure Istio + Keptn  shttps://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-08/index.html?index=..%2F..index#5
# 4. connect the Keptn CLI to the cluster https://tutorials.keptn.sh/tutorials/keptn-full-tour-prometheus-08/index.html?index=..%2F..index#6 


#####################################################################
# make sure you are executing those commands in test-data folder!!!
#####################################################################

# 5. LITMUS Demo Setup Pre-Req 

## install litmus operator & chaos CRDs 
kubectl apply -f litmus/litmus-operator-v1.13.2.yaml

# wait for operator to start
sleep 10

## pull the chaos experiment CR (static) 
kubectl apply -f litmus/pod-delete-ChaosExperiment-CR.yaml 
## pull the chaos experiment RBAC (static) 
kubectl apply -f litmus/pod-delete-rbac.yaml 


# 6. Add Prometheus and Prometheus-SLI-SErvice
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.4.0/deploy/service.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/release-0.3.0/deploy/service.yaml


# 7. Install this service (litmus-service)
# kubectl apply -f ../deploy/service.yaml

# 8. Setup project and service in Keptn

## CREATE PROJECT
keptn create project litmus --shipyard=./shipyard.yaml

## ONBOARD SERVICE
keptn onboard service helloservice --chart=./helloservice/helm/ --project=litmus

## ADD JMETER TESTS & CONFIG
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=./jmeter/load.jmx --resourceUri=jmeter/load.jmx
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=./jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml

## ADD QUALITY GATE
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=./prometheus/sli.yaml --resourceUri=prometheus/sli.yaml
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=helloservice/slo.yaml --resourceUri=slo.yaml

## ADD LITMUS EXPERIMENT
keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=./litmus/experiment.yaml --resourceUri=litmus/experiment.yaml

# 9. Configure Prometheus for this project
keptn configure monitoring prometheus --project=litmus --service=helloservice

## Install blackbox exporter and change configuration of prometheus
kubectl apply -f ./prometheus/blackbox-exporter.yaml
kubectl apply -f ./prometheus/prometheus-server-conf-cm.yaml -n monitoring

## Restart prometheus
kubectl delete pod -l app=prometheus-server -n monitoring


# 10. Deploy hello-service in version 0.1.1
keptn trigger delivery --project=litmus --service=helloservice --image=jetzlstorfer/hello-server:v0.1.1

# 11. Second deployment event: Scale hello-service (see deploy-event.json)
keptn send event -f helloservice/deploy-event.json
