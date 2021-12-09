{
  id:: error 'id is needed',
  start:: error 'start is needed',
  end:: error 'end is needed',
  messages:: '1000000',
  threads:: '2',
  fair:: false,
  memory:: '1Gi',

  pod:: {
    apiVersion: 'v1',
    kind: 'Pod',
    metadata: {
      labels: {
        app: 'kafka-benchmarker-%s' % $.id,
        job: 'kafka-benchmarker',
      },
      name: 'kafka-benchmarker-%s' % $.id,
      namespace: 'sre-test-kafka',
    },
    spec: {
      affinity: {
        podAntiAffinity: {
          preferredDuringSchedulingIgnoredDuringExecution: [
            {
              podAffinityTerm: {
                labelSelector: {
                  matchLabels: {
                    job: 'kafka-benchmarker',
                  },
                },
                topologyKey: 'kubernetes.io/hostname',
              },
              weight: 100,
            },
            {
              podAffinityTerm: {
                labelSelector: {
                  matchLabels: {
                    job: 'kafka-benchmarker',
                  },
                },
                topologyKey: 'failure-domain.beta.kubernetes.io/zone',
              },
              weight: 100,
            },
          ],
        },
      },
      containers: [
        {
          env: [
            {
              name: 'BOOTSTRAP_SERVERS',
              value: 'sre-test-kafka-bootstrap.sre-test-kafka.svc:9092',
            },
            {
              name: 'HOME',
              value: '/tmp',
            },
          ],
          image: 'quay.io/fciocchetti0/kafka-benchmarker:master',
          imagePullPolicy: 'Always',
          name: 'kafka-benchmarker-%s' % $.id,
          command: [
            '/usr/local/bin/producer',
          ],
          args: [
                  '--brokerList=$(BOOTSTRAP_SERVERS)',
                  '--topicRangeStart=%s' % $.start,
                  '--topicRangeEnd=%s' % $.end,
                  '--nMessages=%s' % $.messages,
                  '--nThreads=%s' % $.threads,
                ] +
                (if $.fair then ['--fair'] else []),
          resources: {
            limits: {
              cpu: '2',
              memory: $.memory,
            },
            requests: {
              cpu: '1',
              memory: $.memory,
            },
          },
          terminationMessagePath: '/dev/termination-log',
          terminationMessagePolicy: 'File',
        },
      ],
    },
  },
}
