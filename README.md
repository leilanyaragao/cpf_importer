# cpf_importer

Projeto em Go para importar, higienizar e validar dados de clientes a partir de um arquivo texto, armazenando os dados em um banco PostgreSQL.

---

## Funcionalidades

- Leitura de arquivo texto com dados de clientes
- Processamento para higienização e validação (CPF, CNPJ)
- Inserção em lote no banco PostgreSQL para otimizar a performance
- Reconexão automática ao banco em caso de falha na conexão
- Criação automática da tabela no banco se não existir

---

## Pré-requisitos

- Go 1.20+ instalado
- Docker e Docker Compose instalados 
- Acesso a terminal/linha de comando

---

### Instalação e Execução

1. Get the repository link [https://github.com/leilanyaragao/cpf_importer]
2. Clone o repositório
   ```
   git clone git@github.com:leilanyaragao/cpf_importer.git
   ```
3. Abra o projeto na sua IDE preferida

4. Para iniciar o serviço de importação abra o terminal na raiz do projeto e execute o seguinte comando:
    ```
    docker-compose up --build
    ```
5. Com o projeto rodando, você pode acessar o banco de dados PostgreSQL para verificar se os dados foram importados corretamente e realizar consultas. Execute o seguinte comando no terminal para acessar o banco de dados via psql:
    ```
    psql -h localhost -p 5432 -U user -d mydatabase
    ```
    **onde:**
    
    - -h localhost -> Conecta-se ao banco local.
    
    - -p 5432 -> Porta padrão do PostgreSQL, isso é configurado no projeto.
  
    - -U user -> Nome de usuário do banco configurado no docker-compose.yml.
  
     - d mydatabase -> Nome do banco de dados utilizado no projeto, isso é configurado no projeto.
       


    Após acessar o prompt do PostgreSQL, você pode executar consultas como:
  
      ```SELECT COUNT(*) FROM clients; #Esse comando retornará a quantidade total de registros inseridos na tabela clients.```


