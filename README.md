Zonecompare is a simple tool to read and compare two DNS zonefiles for the same domain, and output the records not found, and the differents, with extensive options to cover many different use-cases, and to configure the output(text or json). 

The tool was designed and implemented out of frustration of not being able to find similar tools, and later on, published under [GPLv3](https://github.com/VintageOps/dns-zone-compare/blob/master/LICENSE.md).

Below is a more thorough description of the tool and the different options, and detailed examples highlighting the tool's usage.

# Installation

You could either clone and build this locally as detailed on [Option 1](#Option-1:-Install-and-Use-Command-Line-with-git) below, or Build in a Go environment ([Option 2](#Option-2:-Install-by-building-with-Go-(in-Go-environment))), or simply make use of the Zonecompare package directly in your Go Code by importing the appropriate Library ([Option 3](#Option-3:-Import-zonecompare-library-in-your-Go-Code))

## Option 1: Install and Use Command Line with git

```bash
$ git clone https://github.com/VintageOps/dns-zone-compare.git
$ cd dns-zone-compare
$ go build -o zonecompare
```

## Option 2: Install by building with Go (in Go environment)
```bash
$ go get -u github.com/VintageOps/dns-zone-compare
```

## Option 3: Import zonecompare library in your Go Code
```go
import "github.com/VintageOps/dns-zone-compare/pkg/zonecompare"
```


# Usage
---------

## Command Line Usage

```
NAME:
   zonecompare - compare two dns zone files

USAGE:
   zonecompare [options] <path_zonefile1>|<address|name:port> <path_zonefile2>|<address|name:port>

DESCRIPTION:
   zonecompare reads or transfer two DNS zone files and by default, output the differences.
   It can also output the similarities and has an extensive set of options to customize the comparison.
   The default output is a timestamped text highlighting the differences (or similarities with <--showfound>/<-f>),
   But with the appropriate <--json>/<-j> option, it can print the same in a comprehensive json format.
   The mandatory zonefiles argument can be specified either using their file path (e.g. tmp/zonefile1 tmp/zonefile2
   Or using <address|name>:<port> format (e.g. localhost:8053 192.168.0.1:53), in which case, for a DNS zone transfer(axfr)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --domain value                                         domain to compare (e.g. example.com).Required when arguments are <ip|name>:<port>, and on zonefiles with no $ORIGIN
   --ignorettl, -t                                        Force TTL value to 604800 in both zones (default: false)
   --showfound, -f                                        Report on found records (default: false)
   --skipnotfound, -n                                     Skip not found records (default: false)
   --strict, -s                                           Consider the different order of the same record a difference (default: false)
   --ignore value, -i value [ --ignore value, -i value ]  Ignore <value> type records
   --deep value, -d value [ --deep value, -d value ]      Inspect <value> type records by merging, then splitting and sorting the content
   --deepAll, --da                                        Inspect all type records by merging, then splitting and sorting the content (default: false)
   --json, -j                                             output in json format (default: false)
   --text, -x                                             Forcing timestamped text in output, useful only to produce both json and text output (default: false)
   --labelorigin value, --lo value                        label of the origin zone, default origin filename|server:port
   --labeldestination value, --ld value                   label of the destination zone, default destination filename|server:port
   --help, -h         
```

## Importing as a Library
```go
// Import zonecompare
import "github.com/VintageOps/dns-zone-compare/pkg/zonecompare"
    ...

// Set any of the appropriate Options
// You can see examples in the example section of this README
    options := zonecompare.Opts{
    Domain:           string, //domain to compare (e.g. example.com).Required when arguments are <ip|name>:<port>, and on zonefiles with no $ORIGIN
    Origin:           string, //First zonefiles|nameserver:port (the one we are comparing to Destination)
    IgnoreTTL:        bool, //Force TTL value to 604800 in both zones (default: false)
    Ignore:           []string, //Ignore <value> type records
    Deep:             []string, //Inspect <value> type records by merging, then splitting and sorting the content
    DeepAll:          bool, //Inspect all type records by merging, then splitting and sorting the content (default: false)
    Found:            bool, //Report on found records (default: false)
    Notfound:         bool, //Skip not found records (default: false)
    Strict:           bool, //Consider the different order of the same record a difference (default: false)
    Json:             bool, //output in json format (default: false)
    Text:             bool, //Forcing timestamped text in output, useful only to produce both json and text output (default: false)
    Destination:      string, //Second zonefiles|nameserver:port (the one we are comparing to Origin)
    Labelorigin:      string, //label of the origin zone, default origin filename|server:port
    Labeldestination: string, //label of the destination zone, default destination filename|server:port
}
// Call ZoneCompare
zonecompare.ZoneCompare(options)
```

# Detailed Documentation on Options
-----

## Example DNS Zonefiles

To make this section more readable, we have uploaded two zonefiles under the [examples folder](https://github.com/VintageOps/dns-zone-compare/blob/master/examples) , [zone1](https://github.com/VintageOps/dns-zone-compare/blob/master/examples/zone1) and [zone2](https://github.com/VintageOps/dns-zone-compare/blob/master/examples/zone1) that displayed a certain number of similarities and differences.

Whenever there's a need in this documentation to showcase an option, we will be using these two files

---

### `--domain value` (e.g. example.com)

This option sets the domain name for all the records in the two zone, it is required when zonefiles argument is provided as <ip|name>:<port> or when it is provided as static path to the files, but `$ORIGIN` is not specified on the zonefiles

---

### `--ignorettl, -t`

This option forces TTLs value to be aligned to the same value (604800) for all records in both zones during the comparison.

---

### `--showfound, -f`

By default, only the differences are returned, with this option, the similarities are also returned

e.g. 

- without `--showfound`

```
$ ./zonecompare --domain example.com examples/zone1 examples/zone2 
2023/06/21 19:13:48 examples/zone1 different [additional-a-3.example.com. 3600 IN A  192.0.2.17]
2023/06/21 19:13:48 examples/zone2 different [additional-a-3.example.com. 3600 IN A  192.0.2.42]
2023/06/21 19:13:48 examples/zone1 different [additional-txt-4.example.com. 3600 IN TXT  Additional TXT record for Example.com additional-txt-4.example.com. 3600 IN TXT  Welcome to Example.com]
2023/06/21 19:13:48 examples/zone2 different [additional-txt-4.example.com. 3600 IN TXT  Welcome to the new Example.com]
2023/06/21 19:13:48 examples/zone1 different [example.com. 3600 IN TXT  Example.com welcomes you]
2023/06/21 19:13:48 examples/zone2 different [example.com. 3600 IN TXT  Welcome to Example.com]
2023/06/21 19:13:48 examples/zone1 different [mail.example.com. 3600 IN A  192.0.2.4]
2023/06/21 19:13:48 examples/zone2 different [mail.example.com. 3600 IN A  192.0.2.10]
...
```

- with `--showfound`

```
$ ./zonecompare --domain example.com --showfound examples/zone1 examples/zone2
2023/06/21 19:16:01 examples/zone1 different [additional-a-4.example.com. 3600 IN A  192.0.2.18]
2023/06/21 19:16:01 examples/zone2 different [additional-a-4.example.com. 3600 IN A  192.0.2.43]
2023/06/21 19:16:01 examples/zone1 different [additional-txt-2.example.com. 3600 IN TXT  Additional text for Example.com additional-txt-2.example.com. 3600 IN TXT  Another Example.com TXT record]
2023/06/21 19:16:01 examples/zone2 different [additional-txt-2.example.com. 3600 IN TXT  Additional text record for Example.com]
2023/06/21 19:16:01 examples/zone1 different [additional-a-1.example.com. 3600 IN A  192.0.2.15]
2023/06/21 19:16:01 examples/zone2 different [additional-a-1.example.com. 3600 IN A  192.0.2.40]
2023/06/21 19:16:01 examples/zone1 different [additional-a-2.example.com. 3600 IN A  192.0.2.16]
2023/06/21 19:16:01 examples/zone2 different [additional-a-2.example.com. 3600 IN A  192.0.2.41]
2023/06/21 19:16:01 examples/zone1 different [additional-txt-1.example.com. 3600 IN TXT  Example.com TXT record additional-txt-1.example.com. 3600 IN TXT  Welcome to Example.com]
2023/06/21 19:16:01 examples/zone2 different [additional-txt-1.example.com. 3600 IN TXT  Hello from Example.com]
2023/06/21 19:16:01 examples/zone1 found example.com. 3600 IN MX  10 mail.example.com. 20 backup-mail.example.com.
2023/06/21 19:16:01 examples/zone1 found example.com. 3600 IN A  192.0.2.1
2023/06/21 19:16:01 examples/zone1 found example.com. 3600 IN SPF  v=spf1 mx -all
...
```

---

### `--skipnotfound, -n`

This option enables to skip the records not found on the output, and reports only on the differences (and the found, if --showfound was requested)

---

### `--strict, -s`

This option is for comparing the entries in strict order of appearances, for example, the host3 entries are the same, but with a different order.
If we want to highlight the facts that the order is different, then we could add ``--strict``

```bash
$ ./zonecompare --domain example.com --strict --json examples/zone1 examples/zone2 | jq '. | with_entries(select(.key | startswith("host3")))'
{
  "host3.example.com.": {
    "A": [
      {
        "differences": [
          "Wrong Order"
        ],
        "originalRecords": {
          "examples/zone1": [
            "host3.example.com. 3600 IN A  192.0.2.51",
            "host3.example.com. 3600 IN A  192.0.2.50"
          ],
          "examples/zone2": [
            "host3.example.com. 3600 IN A  192.0.2.50",
            "host3.example.com. 3600 IN A  192.0.2.51"
          ]
        },
        "status": "different"
      }
    ]
  }
}
```
---

### `--ignore value, -i value [ --ignore value, -i value ]`

This option simply requests to ignore some value record type.
As an example, if we want not to check/report on SOA difference, we could add `--ignore SOA`

---

### `--deep value, -d value [ --deep value, -d value ] `

The deep enables to inspect <value> type records by merging, then splitting and sorting the content.
It is a very useful options for records type like TXT, which may have been splitted on multiple RRs.

As an example, considering the following two records on zonefiles we are comparing

(zone1)
```bash
host1-cpus         IN      TXT     "cpu1 cpu2 cpu3"
host1-cpus         IN      TXT     "cpu6 cpu4 cpu5"
```

(zone2)
```bash
host1-cpus         IN      TXT     "cpu6 cpu5 cpu3 cpu4 cpu1 cpu2"
```

The two contains the same information, but in a different order, without `--deep TXT`, these are reported different:

```bash
$ ./zonecompare --domain example.com --json --showfound examples/zone1 examples/zone2 | jq '. | with_entries(select(.key | startswith("host1-cpu")))'
{
  "host1-cpus.example.com.": {
    "TXT": [
      {
        "differences": {
          "examples/zone1": [
            "host1-cpus.example.com. 3600 IN TXT  cpu1 cpu2 cpu3",
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu4 cpu5"
          ],
          "examples/zone2": [
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu5 cpu3 cpu4 cpu1 cpu2"
          ]
        },
        "originalRecords": {
          "examples/zone1": [
            "host1-cpus.example.com. 3600 IN TXT  cpu1 cpu2 cpu3",
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu4 cpu5"
          ],
          "examples/zone2": [
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu5 cpu3 cpu4 cpu1 cpu2"
          ]
        },
        "status": "different"
      }
    ]
  }
}
```

With `--deep TXT`, TXT records matching for the same entry, will be compared on their values, and this will match

```bash
$ ./zonecompare --domain example.com --json --showfound --deep TXT examples/zone1 examples/zone2 | jq '. | with_entries(select(.key | startswith("host1-cpu")))'
{
  "host1-cpus.example.com.": {
    "TXT": [
      {
        "originalRecords": {
          "examples/zone1": [
            "host1-cpus.example.com. 3600 IN TXT  cpu1 cpu2 cpu3",
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu4 cpu5"
          ],
          "examples/zone2": [
            "host1-cpus.example.com. 3600 IN TXT  cpu6 cpu5 cpu3 cpu4 cpu1 cpu2"
          ]
        },
        "status": "found"
      }
    ]
  }
}
```

---
### `--deepAll, --da`

The deepAll does the same as deep, but with all Record Type

---
### `--json, -j`

This option return the output in a Json format , that could be used for any extra processing/reporting.
The format of the json output is the following (Comments inline for more explanation):

#### In case of records that have differences between the zones specified 
```json
	"<resource_record>": {              // The Resource record 
	  "<record_type>": [                // The Record Type
	    {
	      "differences": {              // If Record is different
	        "<label_zone>": [           // Where the difference is spotted (label)
	          "<full_record(s)>",       // Which record(s) is different
	        ]
	      },
          "originalRecords":{               // The Original records
		    "<label_first_zone>": [
			  "<full_record(s)>",
		    ],
		    "<label_first_zone>": [
			  "<full_record(s)>",
		    ]
          },
          "repeats": {                     // When using `--deep, -d`, if the information was repeated (see examples below) 
            "<label_zone_repeated>": [
              "<full_record(s)>"
            ]
          },
          "status": "different",           // The Status (different)
	    }
	  ]
	}
```

#### In case of records that are either not found or found (when using --showfound/-f)

```json
  "<resource_record>": {                // The Resource record
    "<record_type>": [                  // The Record Type
      {
        "originalRecords": {            // The Original records on the first zone
          "<label_first_zone>": [
            "<full_record(s)>",
          ]
        },
        "status": "notfound|found"      // The Status (notfound or found)
      }
    ]
  },
```

As example, this is a sample output for two records that are different and one notfound from our examples zones

```json
  "host1.example.com.": {
    "A": [
      {
        "differences": {
          "examples/zone1": [
            "host1.example.com. 3600 IN A  192.0.2.10",
            "host1.example.com. 3600 IN A  192.0.2.11",
            "host1.example.com. 3600 IN A  192.0.2.12",
            "host1.example.com. 3600 IN A  192.0.2.13",
            "host1.example.com. 3600 IN A  192.0.2.14"
          ],
          "examples/zone2": [
            "host1.example.com. 3600 IN A  192.0.2.20",
            "host1.example.com. 3600 IN A  192.0.2.21",
            "host1.example.com. 3600 IN A  192.0.2.22",
            "host1.example.com. 3600 IN A  192.0.2.23",
            "host1.example.com. 3600 IN A  192.0.2.24"
          ]
        },
        "originalRecords": {
          "examples/zone1": [
            "host1.example.com. 3600 IN A  192.0.2.10",
            "host1.example.com. 3600 IN A  192.0.2.11",
            "host1.example.com. 3600 IN A  192.0.2.12",
            "host1.example.com. 3600 IN A  192.0.2.13",
            "host1.example.com. 3600 IN A  192.0.2.14"
          ],
          "examples/zone2": [
            "host1.example.com. 3600 IN A  192.0.2.20",
            "host1.example.com. 3600 IN A  192.0.2.21",
            "host1.example.com. 3600 IN A  192.0.2.22",
            "host1.example.com. 3600 IN A  192.0.2.23",
            "host1.example.com. 3600 IN A  192.0.2.24"
          ]
        },
        "status": "different"
      }
    ]
  },
  "host7.example.com.": {
    "A": [
      {
        "originalRecords": {
          "examples/zone1": [
            "host7.example.com. 3600 IN A  192.168.10.100"
          ]
        },
        "status": "notfound"
      }
    ]
  },
  "host8.example.com.": {
    "A": [
      {
        "originalRecords": {
          "examples/zone1": [
            "host8.example.com. 3600 IN A  192.168.10.100"
          ]
        },
        "status": "notfound"
      }
    ]
  },
```

---

### `--text, -x`

This option is used to forced the timestamped text in output, and is only useful when we want to produce both json and text output at the same time, as text only output is the default option.

---

### ```--labelorigin value, --lo value``` and ```--labeldestination value, --ld value```

The default label on the output for the origin and target zonefile is the <filename|server:port> provided.
This option enables us to modify and customized this at will, for all output types.

For example, running the following without specifying a label, use the "example/zone1" and "example/zone2" as label

#### - With no label (sample of both default and json output)
```bash
$ ./zonecompare --domain example.com examples/zone1 examples/zone2 
2023/06/27 20:12:31 examples/zone1 different [additional-txt-2.example.com. 3600 IN TXT  Additional text for Example.com additional-txt-2.example.com. 3600 IN TXT  Another Example.com TXT record]
2023/06/27 20:12:31 examples/zone2 different [additional-txt-2.example.com. 3600 IN TXT  Additional text record for Example.com]
2023/06/27 20:12:31 examples/zone1 notfound host7.example.com. 3600 IN A  192.168.10.100
2023/06/27 20:12:31 examples/zone1 different [additional-a-1.example.com. 3600 IN A  192.0.2.15]
2023/06/27 20:12:31 examples/zone2 different [additional-a-1.example.com. 3600 IN A  192.0.2.40]
2023/06/27 20:12:31 examples/zone1 different [additional-a-2.example.com. 3600 IN A  192.0.2.16]
2023/06/27 20:12:31 examples/zone2 different [additional-a-2.example.com. 3600 IN A  192.0.2.41]
...
```

```json
$ ./zonecompare --domain example.com --json examples/zone1 examples/zone2
{
  "additional-a-1.example.com.": {
    "A": [
      {
        "differences": {
          "examples/zone1": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.15"
          ],
          "examples/zone2": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.40"
          ]
        },
        "originalRecords": {
          "examples/zone1": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.15"
          ],
          "examples/zone2": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.40"
          ]
        },
        "status": "different"
      }
    ]
  },
...
```

#### - With Custom label (sample of both default and json output)

```bash
 ./zonecompare --domain example.com --labelorigin zone1Label --labeldestination zone2Label examples/zone1 examples/zone2  
2023/06/27 20:16:41 zone1Label different [additional-txt-1.example.com. 3600 IN TXT  Example.com TXT record additional-txt-1.example.com. 3600 IN TXT  Welcome to Example.com]
2023/06/27 20:16:41 zone2Label different [additional-txt-1.example.com. 3600 IN TXT  Hello from Example.com]
2023/06/27 20:16:41 zone1Label different [mail.example.com. 3600 IN A  192.0.2.4]
2023/06/27 20:16:41 zone2Label different [mail.example.com. 3600 IN A  192.0.2.10]
2023/06/27 20:16:41 zone1Label different [additional-a-1.example.com. 3600 IN A  192.0.2.15]
2023/06/27 20:16:41 zone2Label different [additional-a-1.example.com. 3600 IN A  192.0.2.40]
2023/06/27 20:16:41 zone1Label different [additional-a-4.example.com. 3600 IN A  192.0.2.18]
2023/06/27 20:16:41 zone2Label different [additional-a-4.example.com. 3600 IN A  192.0.2.43]
```

```json
$ ./zonecompare --domain example.com --json --labelorigin zone1Label --labeldestination zone2Label examples/zone1 examples/zone2 
{
  "additional-a-1.example.com.": {
    "A": [
      {
        "differences": {
          "zone1Label": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.15"
          ],
          "zone2Label": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.40"
          ]
        },
        "originalRecords": {
          "zone1Label": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.15"
          ],
          "zone2Label": [
            "additional-a-1.example.com. 3600 IN A  192.0.2.40"
          ]
        },
        "status": "different"
      }
    ]
  },
```

---

# TODO:

- Review deepAll
- Unify json and text stream in one
- Add Contextualization
- Separate `zonecompare` according to SOLID and adjust function to receive test
- Create file to store helper functions
- Point Tests to run on CI
- Go Documentation
- Make loadZone Public
