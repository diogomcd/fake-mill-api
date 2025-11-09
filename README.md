# Fake Mill API

API REST de alta performance para geração de dados fake brasileiros, desenvolvida em Go com foco em simplicidade e velocidade.

<div align="center">

[![License](https://img.shields.io/badge/license-GPL%20v3-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org)

<a href="https://fakemill.com" target="_blank">Site</a> • <a href="https://fakemill.com/docs" target="_blank">Documentação</a> • <a href="https://hub.docker.com/r/diogomcd/fake-mill-api" target="_blank">Docker Hub</a>

</div>

## 📋 Índice

- [Sobre](#-sobre)
- [Características](#-características)
- [Endpoints Disponíveis](#-endpoints-disponíveis)
- [Instalação](#-instalação)
- [Docker](#-docker)
- [Uso](#-uso)
- [Desenvolvimento](#️-desenvolvimento)
- [Configuração](#-configuração)
- [Contribuindo](#-contribuindo)
- [Licença](#-licença)

## 📋 Sobre

A **Fake Mill API** é uma solução open source para geração de dados fictícios brasileiros, ideal para testes, desenvolvimento e prototipação. Todos os dados gerados são válidos e seguem os padrões brasileiros, incluindo CPF, CNPJ, RG, telefones, endereços, dados bancários e muito mais.

## ✨ Características

- 🚀 **Alta Performance**: Desenvolvida com Go e Fiber framework
- 🇧🇷 **Dados Brasileiros**: Todos os dados seguem padrões e validações brasileiras
- 📝 **Documentação Completa**: Swagger/OpenAPI integrado
- 🐳 **Docker Ready**: Imagem Docker disponível no Docker Hub
- ✅ **Validação**: Endpoints para validação de documentos e telefones
- 🔒 **Seguro**: Rate limiting e CORS configuráveis
- 🧪 **Testado**: Cobertura de testes unitários e E2E
- 📦 **Sempre Open Source**: Licenciado sob GPL v3

## 🎯 Endpoints Disponíveis

| Categoria | Método | Endpoint | Descrição |
|-----------|--------|----------|-----------|
| **Pessoa** | GET | `/api/v1/person` | Gera dados completos de uma pessoa |
| **Documentos** | GET | `/api/v1/cpf` | Gera CPF válido |
| | GET | `/api/v1/cnpj` | Gera CNPJ válido |
| | GET | `/api/v1/rg` | Gera RG válido |
| **Contato** | GET | `/api/v1/email` | Gera endereço de email |
| | GET | `/api/v1/phone` | Gera número de telefone brasileiro |
| **Financeiro** | GET | `/api/v1/bank-account` | Gera dados de conta bancária |
| | GET | `/api/v1/credit-card` | Gera dados de cartão de crédito |
| **Endereço** | GET | `/api/v1/address` | Gera endereço completo brasileiro |
| | GET | `/api/v1/zipcode` | Gera CEP válido |
| **Empresa** | GET | `/api/v1/company` | Gera dados de empresa |
| **Validação** | GET | `/api/v1/validate/cpf/:cpf` | Valida CPF |
| | GET | `/api/v1/validate/cnpj/:cnpj` | Valida CNPJ |
| | GET | `/api/v1/validate/rg/:rg` | Valida RG |
| | GET | `/api/v1/validate/phone` | Valida telefone |
| **Health Check** | GET | `/api/health` | Status da API |
| **Documentação** | GET | `/api/docs` | Documentação Swagger |

## 🚀 Instalação

### Pré-requisitos

- Go 1.24 ou superior
- Make (opcional, para usar os comandos do Makefile)

### Instalação Local

1. Clone o repositório:
```bash
git clone https://github.com/diogomcd/fake-mill-api.git
cd fake-mill-api
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente (opcional):
```bash
cp .env.example .env
```

4. Execute a aplicação:
```bash
go run cmd/api/main.go
```

Ou usando o Makefile:
```bash
make build
./api
```

A API estará disponível em `http://localhost:8080`

## 🐳 Docker

### Usando Docker Compose (Recomendado)

A forma mais simples de executar a API é usando Docker Compose:

```bash
docker-compose up -d
```

A API estará disponível em `http://localhost:8080`

Para parar o container:

```bash
docker-compose down
```

### Usando a imagem do Docker Hub

```bash
docker pull diogomcd/fake-mill-api:latest
docker run -p 8080:8080 diogomcd/fake-mill-api:latest
```

### Build local

```bash
docker build -t fake-mill-api .
docker run -p 8080:8080 fake-mill-api
```

## 📖 Uso

### Exemplo: Gerar uma pessoa completa

```bash
curl http://localhost:8080/api/v1/person
```

### Exemplo: Gerar múltiplos CPFs

```bash
curl "http://localhost:8080/api/v1/cpf?quantity=5"
```

### Exemplo: Gerar endereço de um estado específico

```bash
curl "http://localhost:8080/api/v1/address?state=SP&quantity=3"
```

### Exemplo: Validar CPF

```bash
curl http://localhost:8080/api/v1/validate/cpf/12345678909
```

## 🛠️ Desenvolvimento

### Comandos Disponíveis

```bash
make help          # Lista todos os comandos disponíveis
make test          # Executa todos os testes
make test-race     # Executa testes com race detector
make test-short    # Executa testes rápidos (desenvolvimento)
make coverage      # Gera relatório de cobertura em HTML
make fmt           # Formata o código
make vet           # Executa go vet
make build         # Compila a aplicação
make docs          # Gera documentação OpenAPI (JSON e YAML)
make ci-pr         # Replica verificações do CI para PRs
make ci-main       # Replica verificações do CI para main
```

## 🔧 Configuração

A aplicação pode ser configurada através de variáveis de ambiente:

- `HOST` - Host do servidor (padrão: `0.0.0.0`)
- `PORT` - Porta do servidor (padrão: `8080`)
- `ENV` - Ambiente (development/production)
- `LOG_LEVEL` - Nível de log (debug/info/warn/error)
- `RATE_LIMIT_ENABLED` - Habilita rate limiting (true/false)
- `RATE_LIMIT_LIMIT` - Limite de requisições por janela
- `RATE_LIMIT_WINDOW` - Janela de tempo em segundos

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para:

1. Fazer um Fork do projeto
2. Criar uma branch para sua feature (`git checkout -b feature/MinhaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona MinhaFeature'`)
4. Push para a branch (`git push origin feature/MinhaFeature`)
5. Abrir um Pull Request

### Padrões de Código

- Utilize `camelCase` para nomes de funções, classes e variáveis
- Siga as convenções Go
- Adicione testes para novas funcionalidades
- Mantenha a cobertura de testes alta
- Documente funções públicas

## 📝 Licença

Este projeto está licenciado sob a Licença GPL v3 - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

Desenvolvido com ❤️ usando Go
