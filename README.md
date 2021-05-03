# k8s-webhook

This repository provides the simple mutating/validating kubernetes webhook.

**Mutating webhook** : Add nodeselector to pod depending on namespaces. This means you can divide the nodes(nodegroups) in which the pod is running for each namespace.

**Validating webhook** : TBD

When you watch manifest files and source codes, you might understand how the kubernetes webhook works.

I hope this repo helps you create your original webhook :) 

## How to deploy

### 1. Use deploy.sh

```bash
$ git clone https://github.com/a2ush/k8s-webhook.git
$ cd k8s-webhook
$ ./deploy.sh
2021/05/03 03:39:33 [INFO] generate received request
2021/05/03 03:39:33 [INFO] received CSR
2021/05/03 03:39:33 [INFO] generating key: ecdsa-256
2021/05/03 03:39:33 [INFO] encoded CSR
certificatesigningrequest.certificates.k8s.io/k8s-webhook.default created
certificatesigningrequest.certificates.k8s.io/k8s-webhook.default approved
secret/k8s-webhook-secret created
deployment.apps/k8s-webhook created
service/k8s-webhook created
mutatingwebhookconfiguration.admissionregistration.k8s.io/k8s-mutating-webhook created
validatingwebhookconfiguration.admissionregistration.k8s.io/k8s-validating-webhook created
```

### 2. Manual deployment without `git clone` command

Create server-key.pem and server.csr file
```bash
$ cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "k8s-webhook.default.svc"
  ],
  "CN": "k8s-webhook.default.svc",
  "key": {
    "algo": "ecdsa",
    "size": 256
  }
}
EOF
```

Deploy csr
```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: k8s-webhook.default
spec:
  request: $(cat server.csr | base64 -w 0)
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF
```

Create server.crt file
```bash
$ kubectl certificate approve k8s-webhook.default  
$ kubectl get csr k8s-webhook.default -o jsonpath='{.status.certificate}' | base64 --decode > server.crt
$ kubectl create secret tls --save-config  k8s-webhook-secret --key server-key.pem --cert server.crt
```

Deploy webhook-server
```bash
$ kubectl apply -f https://raw.githubusercontent.com/a2ush/k8s-webhook/main/manifests/webhook-server.yaml
```

Deploy MutatingWebhookConfiguration / ValidatingWebhookConfiguration
```bash
$ curl -s https://raw.githubusercontent.com/a2ush/k8s-webhook/main/manifests/mutating-webhook.yaml | sed "s/YOUR_ENCODED_SERVER_CRT/$(cat server.crt | base64 -w 0)/" | kubectl apply -f -
$ curl -s https://raw.githubusercontent.com/a2ush/k8s-webhook/main/manifests/validating-webhook.yaml | sed "s/YOUR_ENCODED_SERVER_CRT/$(cat server.crt | base64 -w 0)/" | kubectl apply -f -
```

