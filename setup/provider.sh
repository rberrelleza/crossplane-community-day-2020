set -euxo pipefail

accessKey=$1
accessSecretKey=$2
awsRegion='us-west-2'

BASE64ENCODED_AWS_ACCOUNT_CREDS=$(echo "[default]\naws_access_key_id = ${accessKey}\naws_secret_access_key = ${accessSecretKey}" | base64  | tr -d "\n")

cat > provider.yaml <<EOF
---
apiVersion: packages.crossplane.io/v1alpha1
kind: ClusterPackageInstall
metadata:
  name: provider-aws
  namespace: crossplane-system
spec:
  package: "crossplane/provider-aws:v0.11.0-rc"
---
apiVersion: v1
kind: Secret
metadata:
  name: aws-account-creds
  namespace: crossplane-system
type: Opaque
data:
  credentials: ${BASE64ENCODED_AWS_ACCOUNT_CREDS}
---
apiVersion: aws.crossplane.io/v1alpha3
kind: Provider
metadata:
  name: aws-provider
spec:
  region: ${awsRegion}
  credentialsSecretRef:
    namespace: crossplane-system
    name: aws-account-creds
    key: credentials
EOF

# apply it to the cluster:
kubectl apply -f "provider.yaml"

# delete the credentials variable
unset BASE64ENCODED_AWS_ACCOUNT_CREDS
rm provider.yaml