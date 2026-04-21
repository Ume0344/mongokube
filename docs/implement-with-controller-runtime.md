## Install kubebuilder

`./scripts/install-kubebuilder.sh && kubebuilder --version`

## Initialize project
`cd mongokube-operator && kubebuilder init --domain parrot.com --repo github.com/ume0344/mongokube-operator`

## Create mongokube-operator API
`kubebuilder create api --group mongokube.wrd --version alpha1 --kind MongoKube`

`make manifests`
