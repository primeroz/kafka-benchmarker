local podspec = import './podspec.libsonnet';

{
  local defaults = { id: 0, start: '%d00' % (self.id - 1), end: '%d99' % (self.id - 1) },

  pods: [
    podspec {
      id:: p.id,
      start:: p.start,
      end:: p.end,
    }.pod
    for p in [
      defaults {
        id: 1,
        start: 0,
        end: 99,
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
