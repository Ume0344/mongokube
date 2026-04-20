#!/usr/bin/env bash

install_kubebuilder() {
	if command -v kubebuilder >/dev/null 2>&1; then
		echo "kubebuilder is already installed: $(command -v kubebuilder)"
		return 0
	fi

	curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
	chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
}

install_kubebuilder
