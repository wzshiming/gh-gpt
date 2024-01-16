# gh-gpt

This is a simple tool to use the GitHub Copilot API as a GPT API.

## Ensure you can access the GitHub Copilot API

There are two ways to access the GitHub Copilot API.

### Login Copilot in IDE

Windows: `~\AppData\Local\github-copilot\hosts.json`
Other: `~/.config/github-copilot/hosts.json`

The file is created when you log in to [Copilot](https://github.com/features/copilot) in an [IDE](https://github.com/settings/copilot)

- JetBrains (tested)
- VS Code
- ...

### Setup Environment

I have only tested `ghu_` of [Identifiable prefixes](https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/#identifiable-prefixes)

``` bash
GH_COPILOT_TOKEN=ghu_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

## Usage

### As a CLI

``` bash
gh-gpt run "Who are you?"
# I am an artificial intelligence designed to assist with information and tasks. How can I help you today?
```

### As a server

Listening on http://127.0.0.1:8000

``` bash
gh-gpt server --address :8000
```

``` bash
curl --location 'http://127.0.0.1:8000/v1/chat/completions' \
--header 'Content-Type: application/json' \
--data '{
  "model": "gpt-4",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Who are you?"}
   ]
}' | jq -r '.choices[0].message.content'
# I am an artificial intelligence designed to assist with information and tasks. How can I help you today?
```

### As a library

``` golang
package main

import (
	"context"
	"os"

	"github.com/wzshiming/gh-gpt/pkg/run"
)

func main() {
	_ = run.RunStream(context.Background(), "Who are you?", os.Stdout)
	// I am an artificial intelligence designed to assist with information and tasks. How can I help you today?
}
```
