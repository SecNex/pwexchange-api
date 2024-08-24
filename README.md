# SecNex - Password Exchange

**PwExchange** is a simple password exchange tool that allows you to share passwords with your friends or colleagues in a secure way. The tool uses a symmetric encryption algorithm to encrypt the password and then sends it to the recipient. The recipient can then decrypt the password using the same tool for one-time use.

## Encryption Algorithm

The tool uses `argon2` for key derivation and `AES` for encryption. Each secret is encrypted with a unique key. We use a server side key, a random client secret, and a random salt to derive the key. The key is then used to encrypt the secret.

**Only with this three keys, the secret can be decrypted.**

## Usage

### Installation

```bash
git clone https://github.com/secnex/pwexchange-api.git pwexchange
cd pwexchange

# Build the tool
docker build -t pwexchange:local .

SERVER_SECRET=$(openssl rand -hex 32)
AUTH_TOKEN=$(openssl rand -hex 32)

# Run the tool
docker run -p 3030:8080 -e SERVER_SECRET=$SERVER_SECRET -e AUTH_TOKEN=$AUTH_TOKEN pwexchange:local
```

### API

#### Create a new secret

```bash
curl -X POST http://localhost:3030/api/store/encrypt -d '{"password": "my-secret"}' -H "Authorization : Bearer $AUTH_TOKEN"
```

#### Decrypt a secret

```bash
curl -X POST http://localhost:3030/api/store/decrypt?id=00000000-0000-0000-0000-000000000000 -d '{"secret": "encryption-key"}' -H "Authorization : Bearer $AUTH_TOKEN"
```

#### List all secrets

```bash
curl -X GET http://localhost:3030/api/store -H "Authorization : Bearer $AUTH_TOKEN"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
