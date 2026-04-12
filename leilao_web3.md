# Leilão Descentralizado — Projeto Educacional

> Stack: Solidity + Rust + Jetpack Compose

Um sistema de leilão on-chain onde cada lance é imutável e auditável, sem intermediários. Projeto educacional que cobre smart contracts, backend de indexação e app Android nativo.

---

## Visão geral

```
[Smart Contract Solidity]  ←→  [Backend Rust]  ←→  [App Android / Jetpack Compose]
     (Ethereum)                  (ethers-rs              (WalletConnect
  lances, regras,                 + axum)                 + Retrofit)
  pagamentos, metadados
```

O contrato é a fonte da verdade. O backend indexa os eventos e expõe uma API REST. O app é a interface do usuário — conecta a wallet, exibe leilões e envia lances.

---

## Por que blockchain faz sentido aqui?

Em um leilão tradicional você confia que o organizador não vai manipular os lances. Na blockchain isso é impossível:

- Todo lance fica gravado on-chain permanentemente
- Qualquer pessoa pode auditar o histórico completo
- O smart contract executa as regras automaticamente, sem árbitro humano
- Pagamentos e devoluções são automáticos via contrato

---

## O que pode ser leiloado?

O contrato Solidity **não sabe** o que está sendo leiloado — ele só gerencia quem fez cada lance, quanto foi ofertado e quem está ganhando. O item em si fica fora da chain.

Esse projeto adota o **Cenário 2 + 3**: item físico ou digital com metadados gravados on-chain via IPFS.

### O que fica on-chain (imutável)
- Descrição do item (texto)
- Hash IPFS das fotos e documentos
- Histórico completo de lances (endereço + valor)
- Pagamento automático pro vendedor ao finalizar

### O que fica off-chain
- Entrega física ou digital do item
- Acordo entre comprador e vendedor

### Por que IPFS?

IPFS é um sistema de armazenamento descentralizado onde cada arquivo tem um hash único. Se alguém tentar alterar a foto do item, o hash muda. Ao gravar o hash no contrato, qualquer pessoa pode verificar que as fotos não foram alteradas depois que o leilão começou — resolvendo um problema real de plataformas como Mercado Livre e OLX.

```
Vendedor cria o leilão:
1. Faz upload das fotos/docs pro IPFS → recebe hash (ex: QmXyz...)
2. Faz deploy do contrato com descrição + hash IPFS + lance mínimo + prazo

Qualquer pessoa pode verificar:
1. Busca o hash IPFS e vê as fotos originais
2. Confirma no contrato que o hash não mudou desde o início
```

Isso significa que você pode leiloar **qualquer coisa** — um carro, um serviço, um imóvel, um arquivo digital — e a integridade do processo de lances e pagamento é garantida pela blockchain.

---

## Fluxo da aplicação

```
1. Vendedor cria o leilão:
   - faz upload das fotos pro IPFS
   - preenche descrição do item no app
   - define lance mínimo e prazo
   - app faz deploy do contrato com metadados + hash IPFS

2. Participantes enviam lances direto pro contrato
   - contrato só aceita se for maior que o lance atual
   - ETH fica bloqueado no contrato

3. Quando o tempo acaba, contrato automaticamente:
   - libera ETH pro vendedor
   - devolve ETH pros perdedores
   - emite evento LeilaoFinalizado com endereço do vencedor

4. Entrega acontece off-chain entre vendedor e vencedor
```

---

## Camada 1 — Smart Contract (Solidity)

### Responsabilidades
- Armazenar metadados do item on-chain (descrição + hash IPFS)
- Validar lances (maior que o atual, dentro do prazo)
- Bloquear ETH dos participantes
- Executar pagamento automático ao finalizar
- Emitir eventos pra cada ação (indexados pelo backend)

### Estrutura do contrato

```solidity
contract Leilao {
    address public vendedor;
    address public maiorLicitante;
    uint public maiorLance;
    uint public prazo;
    bool public finalizado;

    // metadados do item — imutáveis após deploy
    string public descricaoItem;
    string public hashIpfs; // ex: "QmXyz..." aponta pra fotos/docs no IPFS

    event NovoLance(address licitante, uint valor);
    event LeilaoFinalizado(address vencedor, uint valor);

    constructor(
        string memory _descricao,
        string memory _hashIpfs,
        uint _duracaoSegundos
    ) {
        vendedor = msg.sender;
        descricaoItem = _descricao;
        hashIpfs = _hashIpfs;
        prazo = block.timestamp + _duracaoSegundos;
    }

    function darLance() external payable {
        require(block.timestamp < prazo, "Leilao encerrado");
        require(msg.value > maiorLance, "Lance insuficiente");

        // devolve ETH pro licitante anterior
        if (maiorLicitante != address(0)) {
            payable(maiorLicitante).transfer(maiorLance);
        }

        maiorLicitante = msg.sender;
        maiorLance = msg.value;

        emit NovoLance(msg.sender, msg.value);
    }

    function finalizar() external {
        require(block.timestamp >= prazo, "Leilao ainda ativo");
        require(!finalizado, "Ja finalizado");

        finalizado = true;
        payable(vendedor).transfer(maiorLance);

        emit LeilaoFinalizado(maiorLicitante, maiorLance);
    }
}
```

### O que aprender aqui
- Variáveis de estado e seu custo em gas
- Modificadores (`require`, `payable`)
- Eventos e como são indexados
- Padrões de segurança: reentrancy guard, checks-effects-interactions
- Armazenar referências IPFS on-chain
- Deploy na testnet (Sepolia) via Hardhat ou Foundry

---

## Camada 2 — Backend (Rust)

### Responsabilidades
- Receber fotos do app e fazer upload pro IPFS
- Escutar eventos do contrato via WebSocket
- Indexar lances em SQLite pra consulta rápida
- Expor API REST pro app Android consumir
- Enviar alertas via Telegram quando um lance é superado

### Stack

| Crate | Função |
|---|---|
| `ethers-rs` | Conexão com Ethereum, escuta de eventos, decodificação de ABI |
| `axum` | API REST async |
| `sqlx` | Persistência em SQLite |
| `tokio` | Runtime async |
| `teloxide` | Bot Telegram pra alertas |
| `serde` | Serialização JSON |
| `reqwest` | Upload de arquivos pro IPFS (Pinata ou NFT.Storage) |

### Estrutura de módulos

```
src/
├── main.rs           — inicializa indexer + API em paralelo
├── indexer.rs        — escuta eventos do contrato via WebSocket
├── api.rs            — rotas REST (GET /leiloes, GET /lances/:id)
├── db.rs             — queries SQLite via sqlx
├── notifier.rs       — alertas Telegram
├── ipfs.rs           — upload de fotos pro IPFS via Pinata
└── contract.rs       — ABI do contrato + tipos gerados pelo ethers-rs
```

### Rotas da API

```
GET  /leiloes               — lista todos os leilões ativos
GET  /leiloes/:id           — detalhes + metadados + hash IPFS
GET  /leiloes/:id/lances    — histórico de lances
POST /leiloes/upload        — faz upload das fotos pro IPFS, retorna hash
POST /leiloes               — faz deploy do contrato com metadados + hash IPFS
```

### Exemplo de upload pro IPFS

```rust
// recebe fotos do app, faz upload pro Pinata, retorna hash
async fn upload_ipfs(fotos: Vec<Bytes>) -> String {
    let client = reqwest::Client::new();
    let res = client
        .post("https://api.pinata.cloud/pinning/pinFileToIPFS")
        .bearer_auth(PINATA_API_KEY)
        .multipart(form_com_fotos(fotos))
        .send().await.unwrap();

    res.json::<PinataResponse>().await.unwrap().ipfs_hash
    // retorna algo como "QmXyz123..."
}
```

### Exemplo de indexação de evento

```rust
// escuta evento NovoLance em tempo real
let filter = contract.event::<NovoLanceFilter>();

filter.stream().await?.for_each(|event| async {
    let lance = event.unwrap();
    db::salvar_lance(&pool, &lance).await.unwrap();
    notifier::alertar_superado(&bot, &lance).await;
}).await;
```

### O que aprender aqui
- Programação async com tokio
- Integração com blockchain via ethers-rs
- Upload e pinning de arquivos no IPFS via Pinata
- Padrão de indexação de eventos (reorg handling)
- API REST com axum
- SQLite com sqlx e migrations

---

## Camada 3 — App Android (Jetpack Compose)

### Responsabilidades
- Conectar wallet do usuário via WalletConnect
- Listar leilões ativos com fotos (carregadas do IPFS)
- Exibir histórico de lances em tempo real
- Criar leilão — faz upload das fotos, preenche metadados, faz deploy do contrato
- Enviar lance — assina tx localmente e envia pro contrato
- Receber push notification quando lance é superado

### Stack

| Lib | Função |
|---|---|
| `Jetpack Compose` | UI declarativa nativa |
| `WalletConnect` | Conexão com MetaMask / carteiras |
| `Retrofit` | Consumo da API REST do backend |
| `web3j` | Interação com contrato Ethereum |
| `Hilt` | Injeção de dependência |
| `ViewModel + Flow` | Gerenciamento de estado |

### Telas principais

```
MainActivity
├── TelaHome          — lista de leilões ativos com foto e timer
├── TelaDetalhes      — fotos (IPFS), descrição, histórico de lances
│   └── DialogLance   — input do valor + confirmação da tx
├── TelaCriar         — upload de fotos + descrição + lance mínimo + prazo
│   └── envia fotos pro backend → recebe hash IPFS → deploy do contrato
└── TelaMeusPerfil    — leilões criados e leilões que estou participando
```

### Fluxo de dar lance no app

```
1. Usuário digita valor do lance
2. App monta a transação via web3j
3. WalletConnect envia pra MetaMask pra assinar
4. Usuário confirma na MetaMask
5. Tx é enviada pra rede Ethereum
6. Backend detecta o evento e atualiza o app via polling/WebSocket
```

### O que aprender aqui
- Jetpack Compose e estado reativo com ViewModel + Flow
- WalletConnect pra autenticação Web3 em mobile
- Assinar e enviar transações em Android
- UI/UX de apps que interagem com blockchain

---

## Roadmap sugerido

### Fase 1 — Smart Contract (2-3 semanas)
- [ ] Estudar Solidity básico (variáveis, funções, eventos)
- [ ] Implementar contrato com metadados + hash IPFS
- [ ] Escrever testes com Hardhat ou Foundry
- [ ] Deploy na testnet Sepolia
- [ ] Interagir com o contrato via Etherscan e verificar imutabilidade dos metadados

### Fase 2 — Backend Rust (2-3 semanas)
- [ ] Setup do projeto com tokio + axum + sqlx
- [ ] Integrar upload de fotos pro IPFS via Pinata
- [ ] Conectar no contrato via ethers-rs WebSocket
- [ ] Indexar eventos em SQLite
- [ ] Implementar rotas da API REST
- [ ] Adicionar alertas Telegram

### Fase 3 — App Android (3-4 semanas)
- [ ] Setup do projeto Jetpack Compose
- [ ] Telas de listagem e detalhes consumindo a API
- [ ] Exibir fotos carregadas do IPFS
- [ ] Integração WalletConnect
- [ ] Tela de criar leilão com upload de fotos
- [ ] Fluxo de dar lance com assinatura de tx
- [ ] Push notifications

### Fase 4 — Polimento (1-2 semanas)
- [ ] Tratamento de erros e edge cases
- [ ] Reorg handling no indexer
- [ ] Testes end-to-end na testnet
- [ ] Deploy do backend num VPS

**Tempo total estimado: 2 a 3 meses** dedicando algumas horas por dia.

---

## Conceitos Web3 que você vai dominar

| Conceito | Onde aparece |
|---|---|
| EVM e gas | Solidity — custo de cada operação |
| Eventos e logs | ethers-rs — indexação off-chain |
| ABI | Interface entre Rust/Android e o contrato |
| Transações e assinatura | WalletConnect no app |
| Testnet vs Mainnet | Deploy e testes sem dinheiro real |
| Reorg de blocos | Backend — lidar com reorganização da chain |
| Imutabilidade on-chain | Confiança nos lances e metadados registrados |
| IPFS e content addressing | Hash das fotos gravado no contrato — prova de integridade |
| Smart contract security | Reentrancy, access control, overflow |

---

## Recursos pra começar

- **Solidity:** [docs.soliditylang.org](https://docs.soliditylang.org)
- **Hardhat (deploy e testes):** [hardhat.org](https://hardhat.org)
- **ethers-rs:** [docs.rs/ethers](https://docs.rs/ethers)
- **IPFS:** [docs.ipfs.tech](https://docs.ipfs.tech)
- **Pinata (IPFS pinning):** [pinata.cloud](https://pinata.cloud) — plano grátis suficiente pra dev
- **WalletConnect Android:** [docs.walletconnect.com](https://docs.walletconnect.com)
- **Sepolia testnet faucet:** [sepoliafaucet.com](https://sepoliafaucet.com) — ETH de teste grátis
