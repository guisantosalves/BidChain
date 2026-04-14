# Arquitetura de Contratos — BidChain

Um contrato `AuctionFactory` fica fixo na blockchain. Ele é deployado **uma única vez**
e fica responsável por criar e registrar todos os leilões.

Deploy (uma vez só)
└── AuctionFactory → endereço fixo na Sepolia

Criar leilão (chamado pelo app toda vez)
└── factory.createAuction("Guitarra", "QmHash", 1 days)
├── faz deploy do Auction internamente
├── salva o endereço na lista auctions[]
└── emite evento AuctionCreated

## Por que isso é melhor

Backend indexa por -> vários contratos -> um único evento `AuctionCreated` |

## Como o app vai usar

1. App chama `factory.createAuction()` → leilão criado on-chain
2. Factory emite `AuctionCreated` com o endereço do novo leilão
3. Backend detecta o evento e indexa o leilão novo
4. App lista leilões chamando `factory.getAuctions()`

## O script de deploy

Usado **uma vez só** pra subir a factory na Sepolia. Depois disso nunca mais.

Salva como docs/architecture.md no projeto.

## alchemy
O Alchemy não sabe o que é o seu contrato — ele só pega o que o forge manda e repassa pra rede.

forge script monta a transação de deploy
        ↓
  assina com sua private key  ← prova que é você
        ↓
  envia pelo Alchemy pra Sepolia
        ↓
  rede registra o contrato com seu endereço como "deployer"


## caso fosse fazer b2b deveria ter controle de quem acessa o contrato
import "@openzeppelin/contracts/access/Ownable.sol";

  contract AuctionFactory is Ownable {
      constructor() Ownable(msg.sender) {}

      function createAuction(...) external onlyOwner returns (address) {
          ...
      }
  }

  ┌───────────────┬───────────────────────────────────────────────┐
  │    Padrão     │                   O que faz                   │
  ├───────────────┼───────────────────────────────────────────────┤
  │ Ownable       │ um único dono                                 │
  ├───────────────┼───────────────────────────────────────────────┤
  │ AccessControl │ múltiplos papéis (admin, moderador, usuário)  │
  ├───────────────┼───────────────────────────────────────────────┤
  │ Pausable      │ dono pode pausar o contrato em emergência     │
  ├───────────────┼───────────────────────────────────────────────┤
  │ Upgradeable   │ permite atualizar o contrato depois do deploy │
  └───────────────┴───────────────────────────────────────────────┘