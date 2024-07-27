# Galleria

Trata-se de uma galeria de arte online, basicamente um clone simples do Pinterest. Os usuários podem mandar imagens para serem expostas na galleria.

## Funcionais

- [ ] Permitir o cadastro de novos usuários.
- [ ] Autenticação via login (e-mail/username e senha).
- [ ] Opção de recuperação de senha.
- [ ] Permitir que os usuários editem seus perfis (foto de perfil, bio, etc.).
- [ ] Exibir informações de perfil de outros usuários (nome, bio, imagens enviadas, etc.).
- [ ] Permitir que os usuários cadastrados enviem imagens.
- [ ] Formulário de envio com campos para título, descrição.
- [ ] Interface para curadores revisarem e aprovarem/rejeitarem imagens enviadas.
- [ ] Exibir imagens em uma galeria pública.
- [ ] Detalhes da imagem com título, descrição, autor e tags.
- [ ] Permitir que usuários comentem e curtam as imagens.

## Requisitos Não-Funcionais

- [ ] Os dados devem ser armazenados numa base de dados PostgreSQL.
- [ ] Os usuários devem ser identificados por um JWT.
- [ ] Vamos tentar usar para aprender os serviços da AWS.

## Regras de Negócio

- [ ] Imagens devem seguir as diretrizes de conteúdo da galeria (ex.: sem conteúdo ofensivo).
- [ ] Comentários ofensivos podem ser removidos e usuários que violarem as regras podem ser banidos.
