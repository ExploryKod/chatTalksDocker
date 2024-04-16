# Chat Talks - Backend (en Go)

#### Repo du front: <a href="https://github.com/ExploryKod/chatTalksClient">voir ici</a>

Site en ligne (frontend) : https://chat-talks-client.vercel.app <br/>

API en ligne : https://go-chat-docker.onrender.com/

Repo de l'API en ligne : https://github.com/ExploryKod/go-chat-docker

## Installations

L'app se situe dans le dossier `gorillachat`

### BDD

1. Dump <br/>

Le dump de la bdd se trouve dans migrations/chatbdd.sql. 
Importez ce dump dans votre bdd custom ou dans celle créer par docker si elle ne se charge pas automatiquement.

2. Vérifiez que vous avez bien une BDD fonctionnelle en local ou utilisez :<br/>
- Celle qui sera créé par docker via notre configuration
- Ou la BDD de l'API en ligne (elle ne sera pas fonctionnelle éternellement).

Passé le jour de l'évaluation, pour utiliser la BDD en ligne, il faudra nous demander le contenu du .env.

3. Variables d'environnements

Remplissez le .env avec vos propre variables de bdd et mettez à jour la configuration dans gorillachat/main.go

Exemple de bdd:
- Serveur: database
- Utilisateur: root
- Mot de passe: password
- Nom de la base de donnée: chabdd

Créer ce fichier `.env` à la racine du projet :

```
MARIADB_ROOT_PASSWORD=password
MARIADB_DATABASE=chatbdd
# GOOS=darwin
# GOARCH=arm64
GOOS=linux
GOARCH=amd64
PORT=8000

MYSQL_ADDON_HOST=database:3306
MYSQL_ADDON_DB=chatbdd
MYSQL_ADDON_USER=root
MYSQL_ADDON_PASSWORD=password
```

Ici c'est la base de donnée créer dans les conteneurs docker: pour utiliser la vôtre, remplacer les variables.
La connexion à la BDD se configure dans gorillachat/main.go 

Puis il faut lancer docker et les commandes docker.

### Pour le lancer en utilisant les scripts: 

Dans un terminal bash : 

#### Monter un conteneur docker puis lancer l'app
```
 ./enable_run.sh start
```

#### Au cas où l'app ne se lance pas, aller manuellement le faire : 

```
 cd gorillachat 
```

```
 docker compose up -d --build
```

```
 docker exec -it go-api sh -c "go run ."
```

***Dev Local sans docker*** en allant dans gorillachat (et non depuis la racine) : 

```
go mod tidy
```

```sh
go run .
```

***Dev en version turbo*** (live reloading) avec air (aprés avoir get le package) et sur gorillachat: 

- Générer la config / default : `air init`
- Puis : `air` 
