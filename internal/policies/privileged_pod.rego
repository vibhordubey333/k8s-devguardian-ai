package devguardian.k8s

deny[reason]{
input.kind == "Pod"
input.spec.containers[_].securityContext.privileged == true
reason := sprintf("Privileged container found in pod %s", [input.metadata.name])
}