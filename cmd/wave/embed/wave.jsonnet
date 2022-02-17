{
  sentry:: std.native('sentry'),

  objects: {
    Deployment: function(name) {
      deployment: {
        apiVersion: 'apps/v1',
        kind: 'Deployment',
        metadata: {name: name},
        spec: {
          replicas: 1,
          revisionHistoryLimit: 10,
          strategy: {
            rollingUpdate: {maxUnavailable: 0},
          },
          selector: {
            matchLabels: {app: name},
          },
          template: {
            metadata: {
              labels: {app: name},
            },
            spec: {
              containers: [],
            },
          },
        },
      },
    },

    StatefulSet: function(name) {
      statefulset: {
        apiVersion: 'apps/v1',
        kind: 'StatefulSet',
        metadata: {name: name},
        spec: {
          selector: {
            matchLabels: {app: name},
          },
          serviceName: name,
          updateStrategy: {
            type: 'RollingUpdate',
          },
          template: {
            metadata: {
              labels: {app: name},
            },
            spec: {
              containers: [],
              volumes: [],
            },
          },
        },
      },
    },

    Container: function(name, image) {
      name: name,
      image: image,
      ports: [],
      env: [
        {name: 'VERSION', value: std.extVar('version')},
      ],
      volumeMounts: [],
    },

    ServiceAccount: function(name) {
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      metadata: {name: name},
    },

    Service: function(name) {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        name: name,
      },
      spec: {
        selector: {app: name},
        ports: [],
      },
    },

    HeadlessService: function(name) {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {
        name: name,
      },
      spec: {
        selector: {app: name},
        ports: [],
        clusterIP: 'None'
      },
    },

    ExternalService: function(name, ip) {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: {name: name},
      spec: {
        selector: {app: name},
        ports: [],
        type: 'LoadBalancer',
        loadBalancerIP: ip,
        externalTrafficPolicy: 'Local',
      },
    },
  },

  network: {
    ContainerPort: function(name, port) {
      ports+: [
        {
          name: name,
          containerPort: port,
        },
      ],
    },

    Port: function(name, port, targetPort='same')
      if targetPort == 'same' then {
        spec+: {
          ports+: [
            {
              name: name,
              port: port,
              targetPort: port,
            },
          ],
        },
      } else {
        spec+: {
          ports+: [
            {
              name: name,
              port: port,
              targetPort: targetPort,
            },
          ],
        },
      },
  },

  env: {
    Version: function(name) name + ':' + std.extVar('version'),
  },

  identities: {
    Azure: function(name, resourceID, clientID) {
      identity: {
        apiVersion: 'aadpodidentity.k8s.io/v1',
        kind: 'AzureIdentity',
        metadata: {
          name: name,
        },
        spec: {
          type: 0,
          resourceID: resourceID,
          clientID: clientID,
        },
      },

      binding: {
        apiVersion: 'aadpodidentity.k8s.io/v1',
        kind: 'AzureIdentityBinding',
        metadata: {
          name: name,
        },
        spec: {
          azureIdentity: name,
          selector: name,
        },
      },
    },

    AzureBind: function(name) {
      deployment+: {
        spec+: {
          template+: {
            metadata+: {
              labels+: {
                aadpodidbinding: name,
              },
            },
          },
        },
      },
    },
  },

  features: {
    DownwardAPI: function() {
      env+: [
        {
          name: 'K8S_POD_IP',
          valueFrom: {
            fieldRef: {fieldPath: 'status.podIP'},
          },
        },
      ],
    },

    healthchecks: {
      HTTP: function(port=8080) {
        livenessProbe: {
          httpGet: {path: '/health', port: port},
          timeoutSeconds: 5,
          initialDelaySeconds: 10,
        },
        readinessProbe: {
          httpGet: {path: '/health', port: port},
          timeoutSeconds: 5,
          initialDelaySeconds: 10,
        },
      },
    },

    CustomSelector: function(selector) {
      spec+: {
        selector+: {
          app: selector,
        },
      },
    },
  },

  spec: {
    DeploymentContainer: function(container) {
      deployment+: {
        spec+: {
          template+: {
            spec+: {
              containers+: [container],
            },
          },
        },
      },
    },

    StatefulSetContainer: function(container) {
      statefulset+: {
        spec+: {
          template+: {
            spec+: {
              containers+: [container],
            },
          },
        },
      },
    },
  },

  resources: {
    Request: function(memory, cpu='2m') {
      resources+: {
        requests: {cpu: cpu, memory: memory},
        limits: {memory: memory},
      },
    },
  },

  volumes: {
    Google: function(name, disk) {
      statefulset+: {
        spec+: {
          template+: {
            spec+: {
              volumes+: [
                {
                  name: name,
                  gcePersistentDisk: {
                    pdName: disk,
                    fsType: 'ext4',
                  },
                },
              ],
            },
          },
        },
      },
    },

    Mount: function(name, path) {
      volumeMounts+: [
        {
          name: name,
          mountPath: path,
        },
      ],
    },
  },
}