Zonecompare is a simple tool to read and compare two DNS zonefiles for the same domain, and output the records not found, and the differents, with extensive options to cover many different use-cases, and to configure the output(text or json). 

The tool was designed and implemented out of frustration of not being able to find similar tools, and later on, published under [GPLv3](https://github.com/VintageOps/dns-zone-compare/blob/master/LICENSE.md).

Below is a more thorough description of the tool and the different options, and detailed examples highlighting the tool's usage.

# Usage
---------

```
NAME:
   zonecompare - compare two dns zone files

USAGE:
   zonecompare [global options] command [command options] <path_zonefile1>|<address|name:port> <path_zonefile2>|<address|name:port>

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

# Detailed Documentation on its options
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
