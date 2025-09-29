[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/whois-parser-ai/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/whois-parser-ai)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/whois-parser-ai)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/whois-parser-ai)](https://goreportcard.com/report/github.com/rgglez/whois-parser-ai)
[![GitHub release](https://img.shields.io/github/release/rgglez/whois-parser-ai.svg)](https://github.com/rgglez/whois-parser-ai/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/whois-parser-ai?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/whois-parser-ai?style=social)

# whois-parser-ai

[Go](https://go.dev/) module for parsing [WHOIS](https://en.wikipedia.org/wiki/WHOIS) output and extracting useful information using Azure OpenAI.

## Azure OpenAI setup

1. Sign in to [Azure Portal](https://portal.azure.com).
2. Create a Resource:
  * Search for "Azure OpenAI" in the Marketplace.
  * Click **Create** ans choose your subscription,
  resource group, and region (for instance, `West US 3`).
  * Set a **deployment name** (for instance, `gpt-4o-mini-whois`) and choose the model, for example `gpt-4o-mini`.
3. Wait for deployment. Once deployed, note your endpoint URL, API keys and deployment name.

## Installation

```bash
go get github.com/rgglez/whois-parser-ai/whoisparserai
```

## Usage

You can call the parser with this sintax:

```go
var azure *whoisparserai.AzureOpenAIClient

// Load credentials from environment variables
azureOpenAIKey = os.Getenv("AZURE_OPENAI_KEY")
azureOpenAIEndpoint = os.Getenv("AZURE_OPENAI_ENDPOINT")
azureOpenAIModel = os.Getenv("AZURE_OPENAI_MODEL")

// Create client
azure = whoisparserai.NewAzureOpenAIClient(azureOpenAIKey, azureOpenAIEndpoint, azureOpenAIModel)

// Parse the WHOIS output
var result map[string]interface{}
result, err = azure.ParseWhois(rawWhoisData)
```

This will fill a map with some of these fields:

```json
results = {
  'domain_name'
  'expiration_date'
  'creation_date'
  'registrar'
  'name_servers'
  'registrant_contact'
  'admin_contact'
  'tech_contact'
  'status'
}
```

Since WHOIS returns free form text without a fixed structure or fixed fields, results may be different from server to server.

## Dependencies

* [Microsoft Azure OpenAI](https://azure.microsoft.com/es-mx/pricing/details/cognitive-services/openai-service/)

## License

Copyright (c) 2025 Rodolfo González González.

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

