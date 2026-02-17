local ds = { type: 'prometheus', uid: '${DS_PROMETHEUS}' };

local legendTable = {
  legend: {
    calcs: ['max', 'lastNotNull'],
    displayMode: 'table',
    placement: 'right',
    showLegend: true,
    sortBy: 'Last',
    sortDesc: false,
  },
};

local target(expr, legendFormat, refId) =
  {
    datasource: ds,
    expr: expr,
    refId: refId,
  }
  + (if legendFormat != null then { legendFormat: legendFormat } else {});

local row(id, title, y, collapsed=false) = {
  collapsed: collapsed,
  gridPos: { h: 1, w: 24, x: 0, y: y },
  id: id,
  panels: [],
  title: title,
  type: 'row',
};

local statPanel(id, title, expr, pos, unit='short', extra={}) = {
  datasource: ds,
  fieldConfig: {
    defaults: {
      color: { mode: 'thresholds' },
      unit: unit,
      mappings: [],
      thresholds: {
        mode: 'absolute',
        steps: [
          { color: 'green', value: null },
        ],
      },
    } + std.get(extra, 'fieldConfig', {}),
    overrides: [],
  },
  gridPos: pos,
  id: id,
  options: {
    colorMode: 'value',
    graphMode: 'none',
    justifyMode: 'auto',
    orientation: 'auto',
    reduceOptions: { calcs: ['lastNotNull'], fields: '', values: false },
    textMode: 'auto',
  } + std.get(extra, 'options', {}),
  targets: [target(expr, null, 'A')],
  title: title,
  type: 'stat',
};

local timeSeriesPanel(id, title, targetsIn, pos, unit='short', extra={}) = {
  datasource: ds,
  fieldConfig: {
    defaults: {
      color: { mode: 'palette-classic' },
      custom: {
        axisBorderShow: false,
        axisCenteredZero: false,
        axisColorMode: 'text',
        axisLabel: '',
        axisPlacement: 'auto',
        barAlignment: 0,
        drawStyle: 'line',
        fillOpacity: 10,
        gradientMode: 'none',
        hideFrom: { legend: false, tooltip: false, viz: false },
        insertNulls: false,
        lineInterpolation: 'linear',
        lineWidth: 1,
        pointSize: 5,
        scaleDistribution: { type: 'linear' },
        showPoints: 'never',
        spanNulls: false,
        stacking: { group: 'A', mode: 'none' },
        thresholdsStyle: { mode: 'off' },
      },
      unit: unit,
    } + std.get(extra, 'fieldConfig', {}),
    overrides: std.get(extra, 'overrides', []),
  },
  gridPos: pos,
  id: id,
  options: legendTable {
    tooltip: { mode: 'multi', sort: 'desc' },
  } + std.get(extra, 'options', {}),
  targets: std.mapWithIndex(
    function(i, t)
      target(
        t.expr,
        if std.objectHas(t, 'legendFormat') then t.legendFormat else null,
        std.char(65 + i),
      ),
    if std.type(targetsIn) == 'array' then targetsIn else [targetsIn],
  ),
  title: title,
  type: 'timeseries',
};

{
  ds: ds,
  target: target,
  row: row,
  statPanel: statPanel,
  timeSeriesPanel: timeSeriesPanel,
  legendTable: legendTable,
}
