FROM argoproj/argocd:v1.5.4
ARG FILE2CONFIGMAP_VERSION=v0.0.2

ADD https://github.com/yix/file2configmap/releases/download/${FILE2CONFIGMAP_VERSION}/file2configmap /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/file2configmap
USER argocd
