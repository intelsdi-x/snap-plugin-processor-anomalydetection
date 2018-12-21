
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**


# Snap plugin processor - anomalydetection
Snap plugin intended to process data and hightlight outliers

[![Build Status](https://travis-ci.org/intelsdi-x/snap-plugin-processor-anomalydetection.svg?branch=master)](https://travis-ci.org/intelsdi-x/snap-plugin-processor-anomalydetection)
[![Go Report Card](https://goreportcard.com/badge/intelsdi-x/snap-plugin-processor-anomalydetection)](https://goreportcard.com/report/intelsdi-x/snap-plugin-processor-anomalydetection)

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
3. [Roadmap](#roadmap)
4. [Community Support](#community-support)
5. [Contributing](#contributing)
6. [License](#license)
7. [Acknowledgements](#acknowledgements)

### System Requirements
* Plugin supports Linux/MacOS/*BSD systems

### Installation
#### Download anomalydetection plugin binary:
You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection/releases) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).

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
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started)

## Documentation

The intention of this plugin is to reduce the amount of data that needs to be transmitted without compromising the information that can be gained from potential usages of the data. 
An simple implementation via the Tukey Filter examines each window of data, and transmits the full window if a potential anomaly is detected.
However, it may be that some activity before and/or after the event could be additionally relevant to understand the potential anomaly, outside of the window of data under test, and to achieve statistical significance, 
therefore the sample size for study needs to be selected to assure adequate results.

![anomaly-detection-picture-grafana](https://raw.githubusercontent.com/intelsdi-x/snap-plugin-processor-anomalydetection/master/anomaly.png)

### Examples
Example running psutil plugin, anomalydetection processor, and writing data into a file.

Documentation for Snap collector psutil plugin can be found [here](https://github.com/intelsdi-x/snap-plugin-collector-psutil)

In one terminal window, open the Snap daemon :
```
$ snapteld -t 0 -l 1
```
The option "-l 1" it is for setting the debugging log level and "-t 0" is for disabling plugin signing.

In another terminal window:

Download and load collector, processor and publisher plugins
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-processor-anomalydetection/latest/linux/x86_64/snap-plugin-processor-anomalydetection
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-processor-anomalydetection
```

See available metrics for your system
```
$ snaptel metric list
```

Create a task file. For example, psutil-anomalydetection-file.json:
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
          "config": {
            "BufLength": 10,
            "Factor": 1.5
          },
          "publish": [
            {
              "plugin_name": "file",
              "config": {
                "file": "/tmp/published_anomalydetection.log"
              }
            }
          ]
        }
      ]
    }
  }
}
```

Start task:
```
$ snaptel task create -t psutil-anomalydetection-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
snaptel task watch 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

This data is published to a file `/tmp/published` per task specification

Stop task:
```
$ snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

## Roadmap

1. Apply power analysis concept to help determine the sample size. while many techniques are possible and may need to be explored for specific , power analysis provides a reasonable out-of-the-box starting point.  

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-processor-anomalydetection/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Spoczynski](https://github.com/sandlbn)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
