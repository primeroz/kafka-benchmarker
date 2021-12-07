local podspec = import './podspec.libsonnet';

{
  local topicsPerPod = std.parseInt($.totalTopics) / 5,
  local defaults = { id: 0, start: '%d' % ((topicsPerPod * self.id) - topicsPerPod), end: '%d' % ((topicsPerPod * self.id) - 1) },
  totalTopics:: std.extVar('topics'),
  fair:: false,

  pods: [
    podspec {
      id:: p.id,
      start:: p.start,
      end:: p.end,
      fair:: $.fair,
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
