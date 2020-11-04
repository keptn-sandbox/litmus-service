
keptn onboard service helloservice --chart=./helm/ --project=litmus

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=litmus/experiment.yaml --resourceUri=litmus/experiment.yaml

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=prometheus/sli.yaml --resourceUri=prometheus/sli.yaml

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=slo.yaml --resourceUri=slo.yaml

## first deployment event
keptn send event new-artifact --project=litmus --service=helloservice --image=jetzlstorfer/hello-server:v0.1.1 

## second deployment event (able to scale by editing the deploy-event.json)
keptn send event -f test-data/helloservice/deploy-event.json


### TODO
# configure prometheus
