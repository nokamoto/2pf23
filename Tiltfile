k8s_yaml('deployments/local/apis.yaml')
k8s_resource('ke-apis', port_forwards=9000)

k8s_yaml('deployments/local/web.yaml')
k8s_resource('web', port_forwards=3000)

# https://github.com/tilt-dev/tilt-extensions/tree/master/helm_resource
load('ext://helm_resource', 'helm_resource')
helm_resource(
  'postgresql', 
  'oci://registry-1.docker.io/bitnamicharts/postgresql', 
  flags = ['-f', 'deployments/local/postgresql.yaml'],
)

# https://github.com/tilt-dev/tilt-extensions/tree/master/ko
load('ext://ko', 'ko_build')
ko_build('ke-apis', './cmd/ke-apis', deps = ['cmd', 'pkg', 'internal'], ignore = ['*/*/*_test.go'])

# https://github.com/tilt-dev/tilt-extensions/tree/master/pack
load('ext://pack', 'pack')
pack('web', path = './web/2pf23console')
