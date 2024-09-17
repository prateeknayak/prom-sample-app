.PHONY: k8s-local
k8s-local:
	kind create cluster --config k8s-local/cluster.yaml

.PHONY: create-ns
create-ns:
	kubectl apply -f k8s-local/manifests/namespace.yaml

.PHONY: destroy-ns
destroy-ns:
	kubectl delete -f k8s-local/manifests/namespace.yaml

.PHONY: deploy-op-chart
deploy-op-chart: create-ns
	helm upgrade -n o11y --install kube-prometheus-stack prometheus-community/kube-prometheus-stack -f k8s-local/helm/values.yaml

.PHONY: delete-op-chart
delete-op-chart:
	helm uninstall kube-prometheus-stack


.PHONY: deploy-prom-chart
deploy-prom-chart: create-ns
	helm upgrade -n o11y --install prometheus prometheus-community/prometheus -f k8s-local/helm/prom-values.yaml

.PHONY: delete-prom-chart
delete-prom-chart:
	helm uninstall prometheus

.PHONY: deploy-minio
deploy-minio:
	helm install -n o11y minio bitnami/minio --set persistence.enabled=false

.PHONY: delete-minio
	helm uninstall minio

.PHONY: create-thanos-bucket
create-thanos-bucket:
	kubectl run --namespace o11y minio-client --rm --tty -i --restart='Never' --env MINIO_SERVER_ROOT_USER=$(shell kubectl get secret --namespace o11y minio -o jsonpath="{.data.root-user}" | base64 -d) --env MINIO_SERVER_ROOT_PASSWORD=$(shell kubectl get secret --namespace o11y minio -o jsonpath="{.data.root-password}" | base64 -d) --env MINIO_SERVER_HOST=minio --image docker.io/bitnami/minio-client -- mc mb -p minio/thanos

.PHONY: deploy-thanos
create-thanos-secret:
	kubectl create secret generic thanos-objstore --from-file=k8s-local/helm/objstore.yml

.PHONY: deploy-thanos-component
deploy-thanos-component:
	helm upgrade -n o11y --install thanos bitnami/thanos --values k8s-local/helm/thanos-values.yaml

.PHONY: delete-thanos-component
delete-thanos-component:
	helm uninstall -n o11y thanos

.PHONY: deploy-grafana-chart
deploy-grafana-chart:
	helm install grafana grafana/grafana --namespace o11y
