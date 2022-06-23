## Kubernetes

Kubernetes (K8S) é um produto Open Source utilizado para **automatizar a implantação**, o **dimensionamento** e o **gerenciamento** de aplicativos em contêiner.

Da onde veio: **Google**
- Borg
- Omega
- Kubernetes

### Pontos importantes
- Kubernetes é disponibilizado através de um conjunto de APIs
- Normalmente acessamos a API usando a CLI: kubectl
- Tudo é baseado em estado. Você configura o estado de cada objeto
- Kubernetes Master
  - Kube-apiserver
  - Kube-controller-manager
  - Kube-scheduler
- Outros Nodes:
  - Kubelet
  - Kubeproxy

### Dinâmica "superficial"
- Cluster: Conjunto de maquinas (nodes)
- Cada máquina possui uma quantidade de vCPU e Memória
- Pods: Unidade que contém os containers provisionados
- O Pod representa os processos rodando no cluster

### Deployment
- Provisiona os Pods
- ReplicaSet
- Transbordo de Pods para outros Nodes

Exemplo:
B = Backend => 3 réplicas => 3 Pods Backend
F = Frontend => 2 réplicas => 2 Pods Frontend

### Kind
https://kind.sigs.k8s.io/

```bash
kind create cluster --config=k8s/kind.yaml --name=fullcycle
```

### APP Go
```bash
docker build -t robsantossilva/hello-go .
docker run --rm -p 8080:8080 robsantossilva/hello-go
```

### Criando POD
```bash
kubectl apply -f k8s/pod.yaml
```

### Acessando POD
```bash
kubectl port-forward pod/goserver 8080:8080
```

### ReplicaSet
```bash
kubectl apply -f k8s/replicaset.yaml
```

### O "problema" do ReplicaSet
Se por algum motivo uma nova versão de uma imagem for gerada, mesmo que o replicaset seja configurado, os pods não serão atualizados com a nova versão.
Para que o replicaset suba a nova versão, os PODs que estão rodando precisam ser deletados para o replicaset criar o POD com a nova versão da imagem.

### Detalhes de um POD
```bash
kubectl describe pod pod-name
```

### Deployment
Deployment >>> ReplicaSet >>> Pod

```bash
kubectl apply -f k8s/deployment.yaml
```

### Historico
```bash
> kubectl rollout history deployment goserver
deployment.apps/goserver 
REVISION  CHANGE-CAUSE
1         <none>

> kubectl rollout undo deployment goserver --to-revision=1
```

### Services
É o mecanismo que possibilita o acesso aos pods

**ClusterIP**
```bash
kubectl apply -f k8s/service.yaml
```
```bash
> kubectl get services
NAME               TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
goserver-service   ClusterIP   10.96.73.79   <none>        8080/TCP   22s
kubernetes         ClusterIP   10.96.0.1     <none>        443/TCP    12h
```
```bash
kubectl port-forward service/goserver-service 8080:8080
```

**NodePort**
Node 1: 30000 > < 32767 ---> 30001
Node 2: 30001
Node 3: 30001
Node 4: 30001

**LoadBalancer**
```bash
type: LoadBalancer

> kubectl delete service goserver-service
> kubectl apply -f k8s/service.yaml
> kubectl get svc
NAME               TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
goserver-service   LoadBalancer   10.96.190.149   <pending>     8081:30954/TCP   94s
```

### Port / TargetPort
TargetPort é a porta ativa do container que ira receber as requisições
Port fornece a possibilidade de estabelecer um porta de entrada diferente da porta esperada pelo container, isso quando o TargetPort é informado.

```bash
port: 8081 //porta do service
targetPort: 8080 //porta ativa do container

kubectl port-forward service/goserver-service 9000:8081
> localhost:9000   >>>   8081(porta do service)  >>>  8080(porta ativa do container)
```

### Acessando API do Kubernetes
```bash
kubectl proxy --port=8001
```
http://localhost:8001/api/v1/namespaces/default/services/goserver-service

### Variáveis de ambiente
deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goserver
  labels:
    app: goserver
spec:
  selector:
    matchLabels:
      app: goserver    
  replicas: 1
  template:
    metadata:
      labels:
        app: "goserver"
    spec:
      containers:
        - name: goserver
          image: "robsantossilva/hello-go:v4"
          env:
            - name: "NAME"
              value: "Robson"
            - name: "AGE"
              value: "30"
```

```bash
kubectl port-forward service/goserver-service 9000:8081
```

### ConfigMap
configmap-env.yaml
```yaml
apiVersion: v1 
kind: ConfigMap
metadata:
  name: goserver-env
data:
  NAME: "Robson"
  AGE: "30"
```

deployment.yaml
```yaml
env:
  - name: "NAME"
    valueFrom:
      configMapKeyRef:
        name: "goserver-env"
        key: "NAME"
  - name: "AGE"
    valueFrom:
      configMapKeyRef:
        name: "goserver-env"
        key: "AGE"

OU

envFrom:
  - configMapRef:
      name: goserver-env
```

**Sempre que ouver mudanças no ConfigMap deve-se subir novamente o Deployment**

```bash
> kubectl apply -f k8s/configmap-env.yaml
> kubectl apply -f k8s/deployment.yaml
```

**Criando volume apartir de um ConfigMap**
```yaml
spec:
  containers:
    - name: goserver
      image: "robsantossilva/hello-go:v5"
      envFrom:
        - configMapRef:
            name: goserver-env
      volumeMounts:
        - mountPath: /go/myfamily
          name: config

  volumes:
    - name: config
      configMap:
        name: configmap-family
        items:
          - key: members
            path: family.txt
```
```bash
> kubectl apply -f k8s/configmap-family.yaml
> kubectl apply -f k8s/deployment.yaml
```

### Entrando dentro do Pod
```bash
kubectl exec -it pod-name -- sh
kubectl logs pod-name
```

### Secrets e variaveis de ambiente
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: goserver-secret
type: Opaque
data:
  USER: "YWRtaW4K"
  PASSWORD: "MTIzNDU2Cg=="
```
```bash
kubectl apply -f k8s/secret.yaml
```

### Health Check (Probes)
Possibilita a garantia de que um Pod esta funcionando.

**Liveness**
Validar se a aplicação esta disponivel
deployment.yaml
```yaml
containers:
  - name: goserver
    image: "robsantossilva/hello-go:v5.3"
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
      periodSeconds: 5
      failureThreshold: 3
      timeoutSeconds: 1
      successThreshold: 1
```
```bash
kubectl apply -f k8s/deployment.yaml && watch -n1 kubectl get pods
```

**Readiness**
Validar se aplicação esta pronta para receber chamadas.
Em alguns casos a aplicação demora um pouco para carregar tudo, exige um tempo de inicialização.
deployment.yaml
```yaml
containers:
  - name: goserver
    image: "robsantossilva/hello-go:v5.5"
    readinessProbe:
      httpGet:
        path: /readiness
        port: 8080
      periodSeconds: 3
      failureThreshold: 1
      initialDelaySeconds: 10
```
```bash
kubectl apply -f k8s/deployment.yaml && watch -n1 kubectl get pods
```

**Readiness** Válida se o container esta pronto, estando pronto o container fica disponível, se não estiver READY o trafego não é mais direcionado para o container
**Liveness** Válida se o container esta de pé, caso contrário ele reiniciar e tenta recriar o processo.

Subindo tudo:
```bash
kubectl delete deployment goserver \
&& kubectl apply -f k8s/configmap-env.yaml \
&& kubectl apply -f k8s/configmap-family.yaml \
&& kubectl apply -f k8s/secret.yaml \
&& kubectl apply -f k8s/deployment.yaml \
&& watch -n1 kubectl get pods
```

**startupProbe** Garante a inicialização do container deixando o caminho livre para o Readiness e Liveness


### Metrics-server
- O que devo observar e como definir a quantidade de pods que devo escalar?
- Como usar corretamente o Autoscaling?
- Quais os limites da aplicação para que faça sentido escalar?
- O que é o HPA (Horizontal Pod Autoscaling)

**Metrics-server** coleta em tempo real a quantidade de recursos sendo consumidos no momento.
Pode ser integrado com o Prometheus para extrair metricas de forma visual para tomada de decisão.

**Instalando Metrics-server**
https://github.com/kubernetes-sigs/metrics-server

```bash
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Renomear components.yaml para metrics-server.yaml
Adicionar essa linha em Deployment
```yaml
containers:
  - args:
    - --kubelet-insecure-tls
```

Aplicar
```bash
> kubectl apply -f metrics-server.yaml

> kubectl get apiservices
NAME                                   SERVICE                      AVAILABLE   AGE
v1beta1.metrics.k8s.io                 kube-system/metrics-server   True        37s
```

**Definindo recursos por pod**
```yaml
resources:
  requests:
    cpu: 100m
    memory: 20Mi
  limits:
    cpu: 500m
    memory: 25Mi
```

**Observando uso de recursos**
```bash
kubectl get pod
NAME                        READY   STATUS    RESTARTS   AGE
goserver-7c4c9fd54d-vhl8b   1/1     Running   0          47s

---

kubectl top pod goserver-7c4c9fd54d-vhl8b
NAME                        CPU(cores)   MEMORY(bytes)   
goserver-7c4c9fd54d-vhl8b   1m           6Mi
```

### HPA (Horizontal Pod Autoscaling)
```yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: goserver-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    name: goserver
    kind: Deployment
  minReplicas: 1
  maxReplicas: 30
  targetCPUUtilizationPercentage: 25
```

```bash
> kubectl apply -f k8s/hpa.yaml
> kubectl get hpa
NAME           REFERENCE             TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
goserver-hpa   Deployment/goserver   <unknown>/25%   1         30        0          29s
```

Teste de stress no service
https://github.com/fortio/fortio

```bash
watch -n1 kubectl get hpa
```

```bash
kubectl run -it fortio --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://goserver-service:8081/healthz"
```