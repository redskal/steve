<p align="center"><a href="https://github.com/redskal/steve"><img alt="Steve Logo" src="assets/logo.png" width="50%" /></a></p>

[![Go](https://goreportcard.com/badge/github.com/redskal/steve)](https://goreportcard.com/report/github.com/redskal/steve) 

#### Overview
Steve is a fairly simple Azure attack surface discovery and mapping tool with several functions. It's designed so most of its output can be piped into the next command, making things pretty easy to execute.

###### Mine
`Mine` allows you to discover Azure resources through CNAME records of a given list of domains. You can pipe these in from your favourite subdomain hunting tool, and it will output a list of Azure resource domains it's found hits for.

###### Craft
`Craft` can be fed a list of azure resource names and it will output either permutations of all the resource names, or a list of masks in hashcat format. It can be useful for finding other resources within a resource group, or going insane if you choose the hashcat mask option.

###### Forage
`Forage` will search for Azure resources from a given list of resource names, checking for hits against a list of known Azure domains. It's "inspired" by [Microburst's Invoke-EnumerateAzureSubdomains](https://github.com/NetSPI/MicroBurst/blob/master/Misc/Invoke-EnumerateAzureSubDomains.ps1).

You may get false positives from `forage`. I believe this has something to do with the default list of resolvers, but need to investigate further. For now, modify your config file and set your own trusted resolvers.

#### Install
Requirements:
- Go version 1.21+

To install just run the following command:

```bash
go install github.com/redskal/steve@latest
```

#### Usage
It should be pretty simple to use:
```
  ⇱_ [◨_◧] /

  Steve - by @sam_phisher

  Steve can mine for Azure resources from a list of domains given to him,  
  craft a list of likely resource names or hashcat-style masks for resource
  discovery, and test a list of resource names for presence across known   
  Azure domains.

  Configure your DNS resolvers in:
    Linux: ~/.Steve/config.yml
    Windows: %userprofile%\Steve\config.yml

Usage:
  steve [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  craft       Craft all permutations or hashcat-style masks of a given list of resource names.
  forage      Checks a list of resource names for existence within Azure.
  help        Help about any command
  mine        Given a list of domains, mine will check for Azure resources.

Flags:
  -h, --help   help for steve

Use "steve [command] --help" for more information about a command.
```

###### Mine
```
Mine will read a list of domains from stdin and check each of them
for CNAMEs pointing to known Azure domains, indicating an Azure resource.
Optionally, provide -c to check for potential subdomain takeovers.

Example:
	# identify Azure resources from CNAMEs
	$ subfinder -d microsoft.com | steve mine

	# identify Azure resources with potential subdomain takeovers
	$ subfinder -d microsoft.com | steve mine -c
```

###### Craft
```
Craft allows you to generate all permutations or hashcat-style masks of
a given list of Azure resource names.

Use this for further brute-force discovery with your favourite tools, or
with hashcat piped into dnsx.

Example:
	$ cat resource-names.txt | steve craft | steve forage

	or (if you've got time and several 4090s)

	$ for line in $(cat resource-names.txt | steve craft -m);
	  do
	    hashcat -m 0 -a 3 --stdout $line | steve forage;
	  done;
```

###### Forage
```
Feed Forage a list of resource names on stdin, and it will enumerate known
Azure domains for the existence of those resources.

Example:
	$ cat resource-names.txt | steve forage

	If you want to enumerate permutations, add craft in:
	$ cat resource-names.txt | steve craft | steve forage
```

###### Combining it all
```
# Windows
PS> subfinder -d example.com | steve mine | ForEach-Object { $_.Split(".")[0]} | steve craft | steve forage

# *nix
$ subfinder -d example.com | steve mine | cut -f1 -d"." | steve craft | steve forage
```

#### To-dos
List of items to add or improve:
- Improve the hashcat mask processing to reduce overheads when using Hashcat.

#### License
Released under the [#YOLO Public License](https://github.com/YOLOSecFW/YoloSec-Framework/blob/master/YOLO%20Public%20License)