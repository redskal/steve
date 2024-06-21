/*
Released under YOLO licence. Idgaf what you do.
*/
package azrecon

import (
	"fmt"
	"regexp"

	"github.com/miekg/dns"
)

var (
	// hashtable for CheckResourceExists
	resourceDomains = map[string]string{
		".onmicrosoft.com":             "Microsoft Hosted Domain",
		".scm.azurewebsites.net":       "App Services (Management)",
		".azurewebsites.net":           "App Services",
		".p.azurewebsites.net":         "App Services (p)",
		".cloudapp.net":                "App Services (Cloudapp)",
		".file.core.windows.net":       "Storage Accounts (Files)",
		".blob.core.windows.net":       "Storage Accounts (Blobs)",
		".queue.core.windows.net":      "Storage Accounts (Queues)",
		".table.core.windows.net":      "Storage Accounts (Tables)",
		".mail.protection.outlook.com": "Email",
		".sharepoint.com":              "SharePoint",
		".redis.cache.windows.net":     "Databases (Redis)",
		".documents.azure.com":         "Databases (Cosmos DB)",
		".database.windows.net":        "Databases (MSSQL)",
		".vault.azure.net":             "Key Vaults",
		".azureedge.net":               "CDN",
		".search.windows.net":          "Search Appliance",
		".azure-api.net":               "API Services",
		".azurecr.io":                  "Azure Container Registry",
		".trafficmanager.net":          "Traffic Manger (Load Balancer)",
	}

	// hashtable for checkAzureMatch
	domainPatterns = map[string]string{
		"Microsoft Hosted Domain":        `onmicrosoft.com.$`,
		"App Services (Management)":      `scm.azurewebsites.net.$`,
		"App Services":                   `azurewebsites.net.$`,
		"App Services (p)":               `p.azurewebsites.net.$`,
		"App Services (Cloudapp)":        `cloudapp.net.$`,
		"Storage Accounts (Files)":       `file.core.windows.net.$`,
		"Storage Accounts (Blobs)":       `blob.core.windows.net.$`,
		"Storage Accounts (Queues)":      `queue.core.windows.net.$`,
		"Storage Accounts (Tables)":      `table.core.windows.net.$`,
		"Email":                          `mail.protection.outlook.com.$`,
		"SharePoint":                     `sharepoint.com.$`,
		"Databases (Redis)":              `redis.cache.windows.net.$`,
		"Databases (Cosmos DB)":          `documents.azure.com.$`,
		"Databases (MSSQL)":              `database.windows.net.$`,
		"Key Vaults":                     `vault.azure.net.$`,
		"CDN":                            `azureedge.net.$`,
		"Search Appliance":               `search.windows.net.$`,
		"API Services":                   `azure-api.net.$`,
		"Azure Container Registry":       `azurecr.io.$`,
		"Traffic Manger (Load Balancer)": `trafficmanager.net.$`,
	}

	// some hardcoded resolvers in case there's a need.
	// TODO: may want to reduce this - getting false positives.
	Resolvers = []string{
		"1.1.1.1:53", // cloudflare
		"1.0.0.1:53",
		"8.8.8.8:53", // google
		"8.8.4.4:53",
		"76.76.2.0:53", // control d
		"76.76.10.0:53",
		"9.9.9.9:53", // quad9
		"149.112.112.112:53",
		"208.67.222.222:53", // opendns home
		"208.67.220.220:53",
		"185.228.168.9:53", // cleanbrowsing
		"185.228.169.9:53",
		"76.76.19.19:53", // alternate dns
		"76.223.122.150:53",
		"94.140.14.14:53", // adguard dns
		"94.140.15.15:53",
		"8.26.56.26:53", // commodo secure dns
		"8.20.247.20:53",
		"149.126.75.9:53", // Incapsula
	}
)

type Resource struct {
	Domain string
	Type   string
}

type Domain struct {
	Domain string
	Cnames []Cname
}

type Cname struct {
	Type     string
	Cname    string
	Takeover bool
}

// CheckResourceExists checks for DNS records against a given resource name
// concatenated with a number of Azure domains.
func CheckResourceExists(resource, resolver string) ([]Resource, error) {
	var ret []Resource
	for k, v := range resourceDomains {
		domain := fmt.Sprintf("%s%s", resource, k)
		var msg dns.Msg
		fqdn := dns.Fqdn(domain)
		msg.SetQuestion(fqdn, dns.TypeA)
		in, err := dns.Exchange(&msg, resolver)
		if err != nil {
			continue
		}
		if len(in.Answer) > 0 {
			ret = append(ret, Resource{
				Domain: domain,
				Type:   v,
			})
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("sorry, you're not a winner")
	}
	return ret, nil
}

// CheckAzureCnames checks a domain for CNAME DNS records that match
// known Azure domains.
func CheckAzureCnames(domain, resolver string, checkTakeover bool) (Domain, error) {
	ret := Domain{
		Domain: domain,
	}
	var msg dns.Msg
	fqdn := dns.Fqdn(domain)
	msg.SetQuestion(fqdn, dns.TypeCNAME)
	in, err := dns.Exchange(&msg, resolver)
	if err != nil {
		return ret, err
	}
	if len(in.Answer) < 1 {
		return ret, fmt.Errorf("no answers for that domain")
	}
	for _, answer := range in.Answer {
		if cname, ok := answer.(*dns.CNAME); ok {
			if rsrcType, err := checkAzureMatch(cname.Target); err == nil {
				temp := Cname{
					Type:  rsrcType,
					Cname: cname.Target[:len(cname.Target)-1], // the index loses the trailing '.'
				}
				if checkTakeover {
					temp.Takeover, _ = CheckTakeover(domain, resolver)
				}
				ret.Cnames = append(ret.Cnames, temp)
			}
		}
	}
	if len(ret.Cnames) == 0 {
		return ret, fmt.Errorf("no CNAME records matched Azure hostnames")
	}
	return ret, nil
}

// CheckTakeover will check a domain for Rcodes equating to
// NXDOMAIN. Use this when a CNAME relating to Azure has been
// identified.
func CheckTakeover(domain, resolver string) (bool, error) {
	var msg dns.Msg
	fqdn := dns.Fqdn(domain)
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, resolver)
	if err != nil {
		return false, err
	}
	if dns.RcodeToString[in.Rcode] == "NXDOMAIN" {
		return true, nil
	}
	return false, nil
}

// checkAzureMatch matches the given domain to a set of regular
// expressions which includes known Azure domains.
func checkAzureMatch(domain string) (string, error) {
	for k, v := range domainPatterns {
		match, err := regexp.MatchString(v, domain)
		if err != nil {
			return "", err
		}
		if match {
			return k, nil
		}
	}
	return "", fmt.Errorf("no matches found")
}
