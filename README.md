### Kubernetes

Kubernetes (K8S) é um produto Open Source utilizado para **automatizar a implantação**, o **dimensionamento** e o **gerenciamento** de aplicativos em contêiner.

Da onde veio: **Google**
- Borg
- Omega
- Kubernetes

#### Pontos importantes
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

#### Dinâmica "superficial"
- Cluster: Conjunto de maquinas (nodes)
- Cada máquina possui uma quantidade de vCPU e Memória
- Pods: Unidade que contém os containers provisionados
- O Pod representa os processos rodando no cluster

#### Deployment
- Provisiona os Pods
- ReplicaSet
- Transbordo de Pods para outros Nodes

Exemplo:
B = Backend => 3 réplicas => 3 Pods Backend
F = Frontend => 2 réplicas => 2 Pods Frontend

#### Kind
https://kind.sigs.k8s.io/

```bash
kind create cluster --config=k8s/kind.yaml --name=fullcycle
```

#### APP Go
```bash
docker build -t robsantossilva/hello-go .
docker run --rm -p 8080:8080 robsantossilva/hello-go
```

#### Criando POD
```bash
kubectl apply -f k8s/pod.yaml
```

#### Acessando POD
```bash
kubectl port-forward pod/goserver 8080:8080
```

#### ReplicaSet
```bash
kubectl apply -f k8s/replicaset.yaml
```

#### O "problema" do ReplicaSet
Se por algum motivo uma nova versão de uma imagem for gerada, mesmo que o replicaset seja configurado, os pods não serão atualizados com a nova versão.
Para que o replicaset suba a nova versão, os PODs que estão rodando precisam ser deletados para o replicaset criar o POD com a nova versão da imagem.

#### Detalhes de um POD
```bash
kubectl describe pod pod-name
```

#### Deployment
Deployment >>> ReplicaSet >>> Pod

```bash
kubectl apply -f k8s/deployment.yaml
```

#### Historico
```bash
> kubectl rollout history deployment goserver
deployment.apps/goserver 
REVISION  CHANGE-CAUSE
1         <none>

> kubectl rollout undo deployment goserver --to-revision=1
```

#### Services
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

#### Port / TargetPort
TargetPort é a porta ativa do container que ira receber as requisições
Port fornece a possibilidade de estabelecer um porta de entrada diferente da porta esperada pelo container, isso quando o TargetPort é informado.

```bash
port: 8081 //porta do service
targetPort: 8080 //porta ativa do container

kubectl port-forward service/goserver-service 9000:8081
> localhost:9000   >>>   8081(porta do service)  >>>  8080(porta ativa do container)
```

#### Acessando API do Kubernetes
```bash
kubectl proxy --port=8001
```
http://localhost:8001/api/v1/namespaces/default/services/goserver-service

#### Variáveis de ambiente
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

#### ConfigMap
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

#### Entrando dentro do Pod
```bash
kubectl exec -it pod-name -- sh
kubectl logs pod-name
```