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