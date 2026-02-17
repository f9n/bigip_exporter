local config = import '../config.libsonnet';
local panels = import '../lib/panels.libsonnet';

local inst = 'instance=~"$instance"';
local part = 'partition=~"$partition"';
local f = inst + ',' + part;

local ts(expr, legend=null) = { expr: expr } + (if legend != null then { legendFormat: legend } else {});

// --- Stat Panels ---
local statPanels = [
  panels.statPanel(
    1,
    'VS Available',
    'count(bigip_vs_status_availability_state{' + f + '} == 1)',
    { h: 4, w: 5, x: 0, y: 0 },
  ),
  panels.statPanel(
    2,
    'Pool Available',
    'count(bigip_pool_status_availability_state{' + f + '} == 1)',
    { h: 4, w: 5, x: 5, y: 0 },
  ),
  panels.statPanel(
    3,
    'Node Available',
    'count(bigip_node_status_availability_state{' + f + '} == 1)',
    { h: 4, w: 5, x: 10, y: 0 },
  ),
  panels.statPanel(
    4,
    'Pool Members Available',
    'count(bigip_pool_member_status_availability_state{' + f + '} == 1)',
    { h: 4, w: 5, x: 15, y: 0 },
  ),
  panels.statPanel(
    5,
    'Scrape Duration',
    'bigip_total_scrape_duration{' + inst + '}',
    { h: 4, w: 4, x: 20, y: 0 },
    's',
  ),
];

// --- Virtual Servers ---
local vsRow = panels.row(10, 'Virtual Servers', 4);

local vsPanels = [
  panels.timeSeriesPanel(
    11,
    'VS Current Connections',
    [
      ts('bigip_vs_clientside_cur_conns{' + f + '}', '{{vs}}'),
    ],
    { h: 8, w: 8, x: 0, y: 5 },
  ),
  panels.timeSeriesPanel(
    12,
    'VS Traffic',
    [
      ts('rate(bigip_vs_clientside_bytes_in{' + f + '}[$__rate_interval])', '{{vs}} in'),
      ts('rate(bigip_vs_clientside_bytes_out{' + f + '}[$__rate_interval])', '{{vs}} out'),
    ],
    { h: 8, w: 8, x: 8, y: 5 },
    'Bps',
  ),
  panels.timeSeriesPanel(
    13,
    'VS Request Rate',
    [
      ts('rate(bigip_vs_tot_requests{' + f + '}[$__rate_interval])', '{{vs}}'),
    ],
    { h: 8, w: 8, x: 16, y: 5 },
    'reqps',
  ),
  panels.timeSeriesPanel(
    14,
    'VS Packets',
    [
      ts('rate(bigip_vs_clientside_pkts_in{' + f + '}[$__rate_interval])', '{{vs}} in'),
      ts('rate(bigip_vs_clientside_pkts_out{' + f + '}[$__rate_interval])', '{{vs}} out'),
    ],
    { h: 8, w: 12, x: 0, y: 13 },
    'pps',
  ),
  panels.timeSeriesPanel(
    15,
    'VS SYN Cookie',
    [
      ts('rate(bigip_vs_syncookie_accepts{' + f + '}[$__rate_interval])', '{{vs}} accepts'),
      ts('rate(bigip_vs_syncookie_rejects{' + f + '}[$__rate_interval])', '{{vs}} rejects'),
    ],
    { h: 8, w: 12, x: 12, y: 13 },
    'ops',
  ),
];

// --- Pools ---
local poolRow = panels.row(20, 'Pools', 21);

local poolPanels = [
  panels.timeSeriesPanel(
    21,
    'Pool Current Connections',
    [
      ts('bigip_pool_serverside_cur_conns{' + f + '}', '{{pool}}'),
    ],
    { h: 8, w: 8, x: 0, y: 22 },
  ),
  panels.timeSeriesPanel(
    22,
    'Pool Traffic',
    [
      ts('rate(bigip_pool_serverside_bytes_in{' + f + '}[$__rate_interval])', '{{pool}} in'),
      ts('rate(bigip_pool_serverside_bytes_out{' + f + '}[$__rate_interval])', '{{pool}} out'),
    ],
    { h: 8, w: 8, x: 8, y: 22 },
    'Bps',
  ),
  panels.timeSeriesPanel(
    23,
    'Pool Request Rate',
    [
      ts('rate(bigip_pool_tot_requests{' + f + '}[$__rate_interval])', '{{pool}}'),
    ],
    { h: 8, w: 8, x: 16, y: 22 },
    'reqps',
  ),
  panels.timeSeriesPanel(
    24,
    'Pool Active vs Total Members',
    [
      ts('bigip_pool_active_member_cnt{' + f + '}', '{{pool}} active'),
      ts('bigip_pool_member_total_cnt{' + f + '}', '{{pool}} total'),
    ],
    { h: 8, w: 8, x: 0, y: 30 },
  ),
  panels.timeSeriesPanel(
    25,
    'Pool Connection Queue Depth',
    [
      ts('bigip_pool_connq_depth{' + f + '}', '{{pool}}'),
    ],
    { h: 8, w: 8, x: 8, y: 30 },
  ),
  panels.timeSeriesPanel(
    26,
    'Pool Sessions',
    [
      ts('bigip_pool_cur_sessions{' + f + '}', '{{pool}}'),
    ],
    { h: 8, w: 8, x: 16, y: 30 },
  ),
];

// --- Pool Members ---
local poolMemberRow = panels.row(30, 'Pool Members', 38);

local poolMemberPanels = [
  panels.timeSeriesPanel(
    31,
    'Pool Member Connections',
    [
      ts('bigip_pool_member_serverside_cur_conns{' + f + '}', '{{pool}} / {{member}}'),
    ],
    { h: 8, w: 8, x: 0, y: 39 },
  ),
  panels.timeSeriesPanel(
    32,
    'Pool Member Traffic',
    [
      ts('rate(bigip_pool_member_serverside_bytes_in{' + f + '}[$__rate_interval])', '{{pool}} / {{member}} in'),
      ts('rate(bigip_pool_member_serverside_bytes_out{' + f + '}[$__rate_interval])', '{{pool}} / {{member}} out'),
    ],
    { h: 8, w: 8, x: 8, y: 39 },
    'Bps',
  ),
  panels.timeSeriesPanel(
    33,
    'Pool Member Request Rate',
    [
      ts('rate(bigip_pool_member_tot_requests{' + f + '}[$__rate_interval])', '{{pool}} / {{member}}'),
    ],
    { h: 8, w: 8, x: 16, y: 39 },
    'reqps',
  ),
];

// --- Nodes ---
local nodeRow = panels.row(40, 'Nodes', 47);

local nodePanels = [
  panels.timeSeriesPanel(
    41,
    'Node Current Connections',
    [
      ts('bigip_node_serverside_cur_conns{' + f + '}', '{{node}}'),
    ],
    { h: 8, w: 8, x: 0, y: 48 },
  ),
  panels.timeSeriesPanel(
    42,
    'Node Traffic',
    [
      ts('rate(bigip_node_serverside_bytes_in{' + f + '}[$__rate_interval])', '{{node}} in'),
      ts('rate(bigip_node_serverside_bytes_out{' + f + '}[$__rate_interval])', '{{node}} out'),
    ],
    { h: 8, w: 8, x: 8, y: 48 },
    'Bps',
  ),
  panels.timeSeriesPanel(
    43,
    'Node Request Rate',
    [
      ts('rate(bigip_node_tot_requests{' + f + '}[$__rate_interval])', '{{node}}'),
    ],
    { h: 8, w: 8, x: 16, y: 48 },
    'reqps',
  ),
];

// --- iRules ---
local ruleRow = panels.row(50, 'iRules', 56);

local rulePanels = [
  panels.timeSeriesPanel(
    51,
    'Rule Execution Rate',
    [
      ts('rate(bigip_rule_total_executions{' + f + '}[$__rate_interval])', '{{rule}} / {{event}}'),
    ],
    { h: 8, w: 8, x: 0, y: 57 },
    'ops',
  ),
  panels.timeSeriesPanel(
    52,
    'Rule Failures & Aborts',
    [
      ts('rate(bigip_rule_failures{' + f + '}[$__rate_interval])', '{{rule}} failures'),
      ts('rate(bigip_rule_aborts{' + f + '}[$__rate_interval])', '{{rule}} aborts'),
    ],
    { h: 8, w: 8, x: 8, y: 57 },
    'ops',
  ),
  panels.timeSeriesPanel(
    53,
    'Rule CPU Cycles',
    [
      ts('bigip_rule_avg_cycles{' + f + '}', '{{rule}} avg'),
      ts('bigip_rule_max_cycles{' + f + '}', '{{rule}} max'),
    ],
    { h: 8, w: 8, x: 16, y: 57 },
  ),
];

// --- Assemble Dashboard ---
local allPanels =
  statPanels
  + [vsRow] + vsPanels
  + [poolRow] + poolPanels
  + [poolMemberRow] + poolMemberPanels
  + [nodeRow] + nodePanels
  + [ruleRow] + rulePanels;

{
  __inputs: [
    {
      name: 'DS_PROMETHEUS',
      label: 'Prometheus',
      description: '',
      type: 'datasource',
      pluginId: 'prometheus',
      pluginName: 'Prometheus',
    },
  ],
  __requires: [
    { type: 'grafana', id: 'grafana', name: 'Grafana', version: '10.0.0' },
    { type: 'datasource', id: 'prometheus', name: 'Prometheus', version: '1.0.0' },
    { type: 'panel', id: 'stat', name: 'Stat', version: '' },
    { type: 'panel', id: 'timeseries', name: 'Time series', version: '' },
  ],
  annotations: { list: [] },
  editable: true,
  fiscalYearStartMonth: 0,
  graphTooltip: 1,
  id: null,
  links: [],
  panels: allPanels,
  schemaVersion: 39,
  tags: config.dashboardTags,
  templating: {
    list: [
      {
        current: {},
        hide: 0,
        includeAll: false,
        label: 'Datasource',
        multi: false,
        name: 'DS_PROMETHEUS',
        options: [],
        query: 'prometheus',
        queryValue: '',
        refresh: 1,
        regex: '',
        skipUrlSync: false,
        type: 'datasource',
      },
      {
        current: {},
        datasource: { type: 'prometheus', uid: '${DS_PROMETHEUS}' },
        definition: 'label_values(bigip_collector_scrape_status, instance)',
        hide: 0,
        includeAll: true,
        label: 'Instance',
        multi: true,
        name: 'instance',
        options: [],
        query: 'label_values(bigip_collector_scrape_status, instance)',
        refresh: 2,
        regex: '',
        skipUrlSync: false,
        sort: 1,
        type: 'query',
      },
      {
        current: {},
        datasource: { type: 'prometheus', uid: '${DS_PROMETHEUS}' },
        definition: 'label_values(bigip_vs_clientside_cur_conns{instance=~"$instance"}, partition)',
        hide: 0,
        includeAll: true,
        label: 'Partition',
        multi: true,
        name: 'partition',
        options: [],
        query: 'label_values(bigip_vs_clientside_cur_conns{instance=~"$instance"}, partition)',
        refresh: 2,
        regex: '',
        skipUrlSync: false,
        sort: 1,
        type: 'query',
      },
    ],
  },
  time: { from: 'now-1h', to: 'now' },
  timepicker: {},
  timezone: '',
  title: config.dashboardTitle,
  uid: config.dashboardUid,
  version: 0,
  weekStart: '',
  refresh: config.dashboardRefresh,
}
