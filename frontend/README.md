# 📌 Next-Marketing – Frontend para Gerenciamento de Campanhas de Marketing

**Next-Marketing** é o frontend do sistema de **automação de campanhas de marketing**, permitindo a criação, edição, visualização e gerenciamento de campanhas enviadas via **e-mail e WhatsApp**. Ele consome a API do backend **[go-marketing](https://github.com/seu-usuario/go-marketing)**.

---

## 🚀 Tecnologias Utilizadas

- **Next.js** – Framework React para SSR/SSG.
- **Tailwind CSS** – Estilização rápida e responsiva.
- **TypeScript** – Tipagem estática e segurança no código.
- **ShadCN/UI** – Componentes estilizados para melhor UX.
- **Axios** – Requisições HTTP para a API.
- **JWT (JSON Web Token)** – Autenticação segura.

---

## 🎯 Principais Funcionalidades

✅ **Autenticação** via OTP (E-mail/WhatsApp).  
✅ **Dashboard** com resumo de campanhas e métricas.  
✅ **Gerenciamento de Campanhas** (criação, edição, ativação e exclusão).  
✅ **Listas de Contatos** (importação via CSV e segmentação).  
✅ **Gerenciamento de Templates** para e-mails e WhatsApp.  
✅ **Configuração da Conta** (API Key, remetentes e permissões).

---

## ⚙️ Como Executar

### 🔹 1️⃣ Clone o Repositório

````sh
git clone https://github.com/seu-usuario/next-marketing.git
cd next-marketing


### 🔹 2️⃣ Instale as Dependências
```sh
npm install
# ou
yarn install
````

### 🔹 3️⃣ Configure as Variáveis de Ambiente

Crie um arquivo .env.local na raiz do projeto:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 🔹 4️⃣ Execute o Projeto

```sh
npm run dev
# ou
yarn dev
```

O frontend estará disponível em http://localhost:3000.

---

## 📌 Conectando com o Backend (go-marketing)

O frontend consome a API do backend [go-marketing](https://github.com/seu-usuario/go-marketing), que gerencia campanhas e envios.

Caso o backend ainda não esteja rodando, siga as instruções do repositório go-marketing para iniciá-lo.

---

## 🛠 Estrutura do Projeto

📁 /pages – Contém as páginas da aplicação.
📁 /components – Componentes reutilizáveis (botões, cards, etc.).
📁 /styles – Estilos globais e configurações do Tailwind.
📁 /types – Definições de tipos TypeScript.
📁 /public – Assets públicos (ícones, imagens, etc.).

---

🤝 Contribuição
Sinta-se à vontade para contribuir com melhorias!

1. Crie um fork do projeto.
2. Crie uma branch com sua funcionalidade:

```sh
git checkout -b minha-feature
```

3. Faça as alterações e commit:

```sh
git commit -m "Adicionando nova funcionalidade"
```

4. Envie para o repositório:

```sh
git push origin minha-feature
```

5. Crie um Pull Request.

---

📄 Licença

Este projeto está licenciado sob a MIT License.

---
