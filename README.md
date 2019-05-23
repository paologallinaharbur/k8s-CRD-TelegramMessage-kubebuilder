[![Known Vulnerabilities](https://snyk.io/test/github/paologallinaharbur/k8s-CRD-TelegramMessage-kubebuilder/badge.svg?targetFile=Gopkg.lock)](https://snyk.io/test/github/paologallinaharbur/k8s-CRD-TelegramMessage-kubebuilder?targetFile=Gopkg.lock)

# kubebuilder-TelegramMessage-example

That is a demo repository to add a controller managing a custom CRD to a kubernetes clusters.

It is described in [this Medium article](https://medium.com/@paolo.gallina/leverage-k8s-crd-and-kubebuilder-to-create-a-telegram-message-resource-8ce8ac329d89).

# Usage

 - `make install` to install the CRD

 - `make run` to let the manager run

 - `kubectl apply -f config/sample/test.yaml to create the first TelegramMessage

Disclaimer: you have first to create your TelegramMessage bot and retrieve your chatID 
