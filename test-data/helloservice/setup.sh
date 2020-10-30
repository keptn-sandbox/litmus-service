
keptn onboard service helloservice --chart=./helm/ --project=litmus

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=litmus --stage=chaos --service=helloservice --resource=jmeter/jmeter.conf.yaml --resourceUri=jmeter/jmeter.conf.yaml

keptn add-resource --project=litmus --stage=chaos --service=carts --resource=litmus/experiment.yaml --resourceUri=litmus/experiment.yaml


### TODO
# add SLI
# add SLO
# configure prometheus
