#!/bin/sh

cat <<EOF | cfssl genkey - | cfssljson -bare server
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

cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: k8s-webhook.default
spec:
  signerName: beta.eks.amazonaws.com/app-serving # see https://docs.aws.amazon.com/eks/latest/userguide/cert-signing.html
  request: $(cat server.csr | base64 -w 0)
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

kubectl certificate approve k8s-webhook.default  
kubectl get csr k8s-webhook.default -o jsonpath='{.status.certificate}' | base64 --decode > server.crt
kubectl create secret tls --save-config  k8s-webhook-secret --key server-key.pem --cert server.crt

kubectl apply -f manifests/webhook-server.yaml

sed -i -e "s/YOUR_ENCODED_SERVER_CRT/$(cat server.crt | base64 -w 0)/" manifests/mutating-webhook.yaml
sed -i -e "s/YOUR_ENCODED_SERVER_CRT/$(cat server.crt | base64 -w 0)/" manifests/validating-webhook.yaml

kubectl apply -f manifests/mutating-webhook.yaml
kubectl apply -f manifests/validating-webhook.yaml
