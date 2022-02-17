
local wave = import 'wave.jsonnet';

{
  identitiesAzure: wave.identities.Azure(name='foo-identity', resourceID='foo-resource-id', clientID='foo-client-id'),

  objectsServiceAccount: wave.objects.ServiceAccount(name='foo-service-account'),

  objectsServiceEmpty: wave.objects.Service(name='foo-service'),
  objectsServicePorts: wave.objects.Service(name='foo-service') +
    wave.network.Port(name='same-port', port=8080) +
    wave.network.Port(name='different-port', port=8081, targetPort=8082),

  objectsHeadlessServiceEmpty: wave.objects.HeadlessService(name='foo-service'),
  objectsHeadlessServicePorts: wave.objects.HeadlessService(name='foo-service') +
    wave.network.Port(name='same-port', port=8080) +
    wave.network.Port(name='different-port', port=8081, targetPort=8082),

  objectsDeploymentEmpty: wave.objects.Deployment(name='foo-deployment'),
  objectsDeploymentSingleContainer: wave.objects.Deployment(name='foo-deployment') +
    wave.spec.DeploymentContainer(
      wave.objects.Container('foo-container', wave.env.Version('eu.gcr.io/foo')),
    ),
  objectsDeploymentMultipleContainers: wave.objects.Deployment(name='foo-deployment') +
    wave.spec.DeploymentContainer(
      wave.objects.Container('foo-container', wave.env.Version('eu.gcr.io/foo')),
    ) +
    wave.spec.DeploymentContainer(
      wave.objects.Container('bar-container', 'eu.gcr.io/bar'),
    ),
  objectsDeploymentFullContainer: wave.objects.Deployment(name='foo-deployment') +
    wave.spec.DeploymentContainer(
      wave.objects.Container('foo-container', wave.env.Version('eu.gcr.io/foo')) +
      wave.network.ContainerPort(name='foo-port', port=8080) +
      wave.network.ContainerPort(name='bar-port', port=8081) +
      wave.features.DownwardAPI() +
      wave.features.healthchecks.HTTP() +
      wave.resources.Request(memory='256Mi'),
    ),

  objectsDeploymentAzureBind: wave.objects.Deployment(name='foo-deployment') +
    wave.identities.AzureBind(name='foo-identity'),
}
