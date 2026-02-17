local config = import 'config.libsonnet';

{
  config:: config,
  grafanaDashboards: {
    bigip: (import 'dashboards/bigip.jsonnet'),
  },
}
