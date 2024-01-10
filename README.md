# gh-gpt

This is a simple tool to use the GitHub Copilot API as a GPT API.

## Ensure you have a GitHub Copilot API key

The `~/.config/github-copilot/hosts.json` file is created when you log in to Copilot in an editor (VS Code, IDEA, etc.).

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
