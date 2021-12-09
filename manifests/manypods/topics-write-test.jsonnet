local podspec = import './podspec.libsonnet';

local stringToBool(s) =
  if s == 'true' then true
  else if s == 'false' then false
  else error 'invalid boolean: ' + std.manifestJson(s);

{
  local topicsPerPod = std.parseInt($.writeTopics) / 5,
  local defaults = { id: 0, start: '%d' % ((topicsPerPod * self.id) - topicsPerPod), end: '%d' % ((topicsPerPod * self.id) - 1) },
  writeTopics:: std.extVar('writeTopics'),
  fair:: std.extVar('fair'),
  memory:: '3Gi',

  pods: [
    podspec {
      id:: p.id,
      start:: p.start,
      end:: p.end,
      fair:: stringToBool($.fair),
      memory:: $.memory,
    }.pod
    for p in [
      defaults {
        id: 1,
      },
      defaults {
        id: 2,
      },
      defaults {
        id: 3,
      },
      defaults {
        id: 4,
      },
      defaults {
        id: 5,
      },
    ]
  ],

}
