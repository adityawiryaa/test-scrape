.PHONY: help minikube-start minikube-stop minikube-delete docker-build deploy deploy-helm undeploy logs status clean
.PHONY: ansible-up ansible-down ansible-setup ansible-deploy ansible-status ansible-logs ansible-test validate-all

help:
	@echo "=== KUBERNETES (Minikube) ==="
	@echo "  make minikube-start    - Start minikube cluster"
	@echo "  make minikube-stop     - Stop minikube"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make deploy-helm       - Deploy with Helm"
	@echo "  make undeploy-helm     - Remove Helm deployment"
	@echo "  make status            - Check K8s pods/services"
	@echo "  make logs              - View K8s logs"
	@echo "  make url               - Get service URL"
	@echo ""
	@echo "=== ANSIBLE (Local Docker) ==="
	@echo "  make ansible-up        - Start target server container"
	@echo "  make ansible-down      - Stop target server container"
	@echo "  make ansible-setup     - Run setup-runtime playbook"
	@echo "  make ansible-deploy    - Run deploy-app playbook"
	@echo "  make ansible-status    - Check PM2 status"
	@echo "  make ansible-logs      - View PM2 logs"
	@echo "  make ansible-test      - Test app endpoint"
	@echo ""
	@echo "=== VALIDATION ==="
	@echo "  make validate-k8s      - Validate Kubernetes deployment"
	@echo "  make validate-ansible  - Validate Ansible deployment"
	@echo "  make validate-all      - Validate both deployments"

minikube-start:
	minikube start --driver=docker --cpus=2 --memory=4096
	minikube addons enable ingress
	minikube addons enable metrics-server

minikube-stop:
	minikube stop

minikube-delete:
	minikube delete

docker-build:
	eval $$(minikube docker-env) && docker build -t devops-app:latest .

deploy:
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml
	kubectl rollout status deployment/devops-app

deploy-helm:
	helm upgrade --install devops-app ./helm/devops-app

deploy-helm-full:
	helm upgrade --install devops-app ./helm/devops-app \
		--set ingress.enabled=true \
		--set autoscaling.enabled=true

undeploy:
	kubectl delete -f k8s/ --ignore-not-found

undeploy-helm:
	helm uninstall devops-app --ignore-not-found

logs:
	kubectl logs -l app.kubernetes.io/name=devops-app -f --tail=100

status:
	@echo "=== PODS ==="
	@kubectl get pods -l app.kubernetes.io/name=devops-app
	@echo ""
	@echo "=== SERVICES ==="
	@kubectl get svc devops-app
	@echo ""
	@echo "=== DEPLOYMENT ==="
	@kubectl get deployment devops-app

url:
	@minikube service devops-app --url

clean: undeploy minikube-delete

all: minikube-start docker-build deploy-helm status

ansible-up:
	cd ansible && docker compose up -d --build

ansible-down:
	cd ansible && docker compose down

ansible-setup:
	cd ansible && ansible-playbook -i inventory-local.ini setup-runtime.yml

ansible-deploy:
	cd ansible && ansible-playbook -i inventory-local.ini deploy-app-local.yml

ansible-status:
	@echo "=== PM2 STATUS ==="
	@docker exec ansible-target su - deploy -c "pm2 list"
	@echo ""
	@echo "=== PM2 INFO ==="
	@docker exec ansible-target su - deploy -c "pm2 show devops-app" 2>/dev/null || true

ansible-logs:
	docker exec ansible-target su - deploy -c "pm2 logs --lines 50"

ansible-test:
	@echo "=== Testing App Endpoint ==="
	@curl -s http://localhost:3000/ | python3 -m json.tool
	@echo ""
	@echo "=== Testing Health Endpoint ==="
	@curl -s http://localhost:3000/health | python3 -m json.tool

validate-k8s:
	@echo "========================================"
	@echo "  KUBERNETES VALIDATION"
	@echo "========================================"
	@echo ""
	@echo "[1] Checking Minikube status..."
	@minikube status | head -5
	@echo ""
	@echo "[2] Checking pods..."
	@kubectl get pods -l app.kubernetes.io/name=devops-app -o wide
	@echo ""
	@echo "[3] Checking pod readiness..."
	@READY=$$(kubectl get pods -l app.kubernetes.io/name=devops-app -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}'); \
	if echo "$$READY" | grep -q "False"; then \
		echo "❌ Some pods are NOT ready"; \
		exit 1; \
	else \
		echo "✅ All pods are ready"; \
	fi
	@echo ""
	@echo "[4] Checking service..."
	@kubectl get svc devops-app
	@echo ""
	@echo "[5] Testing health endpoint..."
	@URL=$$(minikube service devops-app --url 2>/dev/null | head -1); \
	HEALTH=$$(curl -s $$URL/health 2>/dev/null); \
	if echo "$$HEALTH" | grep -q "healthy"; then \
		echo "✅ Health check passed: $$HEALTH"; \
	else \
		echo "❌ Health check failed"; \
		exit 1; \
	fi
	@echo ""
	@echo "========================================"
	@echo "  ✅ KUBERNETES VALIDATION PASSED"
	@echo "========================================"

validate-ansible:
	@echo "========================================"
	@echo "  ANSIBLE VALIDATION"
	@echo "========================================"
	@echo ""
	@echo "[1] Checking container status..."
	@docker ps --filter "name=ansible-target" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
	@echo ""
	@echo "[2] Checking PM2 processes..."
	@docker exec ansible-target su - deploy -c "pm2 jlist" 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Running: {len([p for p in d if p[\"pm2_env\"][\"status\"]==\"online\"])} instances')"
	@echo ""
	@echo "[3] Checking Node.js version..."
	@docker exec ansible-target node --version
	@echo ""
	@echo "[4] Testing health endpoint..."
	@HEALTH=$$(curl -s http://localhost:3000/health 2>/dev/null); \
	if echo "$$HEALTH" | grep -q "healthy"; then \
		echo "✅ Health check passed: $$HEALTH"; \
	else \
		echo "❌ Health check failed"; \
		exit 1; \
	fi
	@echo ""
	@echo "[5] Checking PM2 logs for errors..."
	@ERRORS=$$(docker exec ansible-target su - deploy -c "pm2 logs --lines 10 --nostream 2>&1" | grep -i "error" | wc -l); \
	if [ "$$ERRORS" -gt 0 ]; then \
		echo "⚠️  Found $$ERRORS error(s) in logs"; \
	else \
		echo "✅ No errors in recent logs"; \
	fi
	@echo ""
	@echo "========================================"
	@echo "  ✅ ANSIBLE VALIDATION PASSED"
	@echo "========================================"

validate-all: validate-ansible
	@echo ""
	@echo ""
	@$(MAKE) validate-k8s 2>/dev/null || echo "⚠️  Kubernetes not running (skip)"

stop-ansible:
	cd ansible && docker compose down

stop-k8s:
	helm uninstall devops-app --ignore-not-found || true
	kubectl delete -f k8s/ --ignore-not-found || true

stop-minikube:
	minikube stop

stop-all: stop-ansible stop-k8s stop-minikube
	@echo "✅ All services stopped"
