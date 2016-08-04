# snap plugin processor - anomalydetection
snap plugin intended to process data and hightlight outliers

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
3. [Roadmap] (#roadmap)
4. [Community Support](#community-support)
5. [Contributing](#contributing)
6. [License](#license)
7. [Acknowledgements](#acknowledgements)

### System Requirements
* Plugin supports Linux/MacOS/*BSD systems

### Installation
#### Download anomalydetection plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-processor-anomalydetection
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

## Documentation

The intention of this plugin is to reduce the amount of data that needs to be transmitted without compromising the information that can be gained from potential usages of the data. 
An simple implementation via the Tukey Filter examines each window of data, and transmits the full window if a potential anomaly is detected.
However, it may be that some activity before and/or after the event could be additionally relevant to understand the potential anomaly, outside of the window of data under test, and to achieve statistical significance, 
therefore the sample size for study needs to be selected to assure adequate results.

![anomaly-detection-picture-grafana](https://raw.githubusercontent.com/intelsdi-x/snap-plugin-processor-anomalydetection/blob/master/anomaly.png)

### Examples
Example running psutil plugin, passthru processor, and writing data into an csv file.

Documentation for snap collector psutil plugin can be found [here](https://github.com/intelsdi-x/snap-plugin-collector-psutil)

In one terminal window, open the snap daemon :
```
$ snapd -t 0 -l 1
```
The option "-l 1" it is for setting the debugging log level and "-t 0" is for disabling plugin signing.

In another terminal window:

Load collector and processor plugins
```
$ snapctl plugin load $SNAP_PSUTIL_PLUGIN/build/rootfs/snap-plugin-collector-psutil
$ snapctl plugin load $SNAP/build/plugin/snap-plugin-publisher-file
$ snapctl plugin load $SNAP_ANOMALYDETECTION_PLUGIN/rootfs/plugin/snap-plugin-processor-anomalydetection
```

See available metrics for your system
```
$ snapctl metric list
```

Create a task file. For example, sample-psutil-anomaly-task.json:
Configure Factor value, "Factor": 1.5 indicates an "outlier", and "Factor": 3.0 indicates data that is "far out".

```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/psutil/load/load1": {},
                "/intel/psutil/load/load5": {},
                "/intel/psutil/load/load15": {},
                "/intel/psutil/vm/free": {},
                "/intel/psutil/vm/used": {}
            },
            "process": [
                {
                    "plugin_name": "anomalydetection",
                    "config":
                        {
                            "BufLength": 10,
                            "Factor": 1.5
                        },
                    "process": null,
                    "publish": [
                        {
                            "plugin_name": "file",
                            "config": {
                                "file": "/tmp/published"
                            }
                        }
                    ]
                }
            ]
        }
    }
}
```

## Roadmap

1. Apply power analysis concept to help determine the sample size. while many techniques are possible and may need to be explored for specific , power analysis provides a reasonable out-of-the-box starting point.  

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Spoczynski](https://github.com/sandlbn)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.