## Plugins

```shell
bin/plugins/
├── official/                    # namespace
│   ├── data-processing/         # type
│   │   ├── anonymizer/          # name
│   │   │   └── v1.0.0/          # version
│   │   │       └── darwin_arm64/ # os_arch
│   │   │           └── plugin   # binary
│   │   └── converter/
│   │       └── v1.0.0/
│   │           └── linux_amd64/
│   │               └── plugin
│   └── security/
│       └── encryptor/
│           └── v2.1.0/
│               └── darwin_arm64/
│                   └── plugin
└── third-party/
    └── analytics/
        └── metrics/
            └── v0.5.0/
                └── linux_amd64/
                    └── plugin
```