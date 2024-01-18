## Desgining MongoKube Custom Resource
The following attributes can be defined by user in manifest file for mongokube;
- *mongoExpressImage*: This defines the name of mongo express container image that user wants to use.
- *mongoDbImage*: This defines the name of mongo db container image that user wants to use.
- *dbUsername*: This defines the db username user wants to use.
- *dbPassword*: This defines the db password user wants to use.

According to the above attributes, CustomResourceDefinition(CRD) is created for MongoKube custom resource.

To create a separate namespace for mongokube, run;
```
kubectl create namespace mongokube-ns
```

To create the CRD for monogkube, run;
```
kubectl create -f home/$(whoami)/mongokube-deployer/manifests/mongokube-crd.yaml 
```
