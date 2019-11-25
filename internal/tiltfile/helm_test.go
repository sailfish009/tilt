package tiltfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmSetArgs(t *testing.T) {
	f := newFixture(t)
	defer f.TearDown()

	f.setupHelm()

	f.file("Tiltfile", `
yml = helm('./helm', name='rose-quartz', namespace='garnet', set=[
  'ingress.enabled=true',
  'service.externalPort=1234',
  'service.internalPort=5678'
])
k8s_yaml(yml)
`)

	f.load()

	m := f.assertNextManifestUnresourced("garnet",
		// A service and ingress with the same name
		"rose-quartz-helloworld-chart",
		"rose-quartz-helloworld-chart")
	yaml := m.K8sTarget().YAML

	// Set on the service
	assert.Contains(t, yaml, "port: 1234")
	assert.Contains(t, yaml, "targetPort: 5678")

	// Set on the ingress
	assert.Contains(t, yaml, "serviceName: rose-quartz-helloworld-chart")
	assert.Contains(t, yaml, "servicePort: 1234")
}

const exampleHelmV2VersionOutput = `Client: v2.12.3+geecf22f`
const exampleHelmV3VersionOutput = `v3.0.0`

func TestParseHelmV2Version(t *testing.T) {
	expected := helmV2
	actual := parseVersion(exampleHelmV2VersionOutput)

	assert.Equal(t, expected, actual)
}

func TestParseHelmV3Version(t *testing.T) {
	expected := helmV3
	actual := parseVersion(exampleHelmV3VersionOutput)

	assert.Equal(t, expected, actual)
}

func TestHelmUnknownVersion(t *testing.T) {
	expected := unknownHelmVersion
	actual := parseVersion("v4.1.2")

	assert.Equal(t, expected, actual)
}
