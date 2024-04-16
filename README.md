# Chat Talks - Backend (en Go)

#### Repo du front: <a href="https://github.com/ExploryKod/chatTalksClient">voir ici</a>

Site en ligne (frontend) : https://chat-talks-client.vercel.app 
API en ligne : https://go-chat-docker.onrender.com/

## Installations

### Pour le lancer en utilisant les scripts: 

Dans un terminal bash : 

1. Monter un conteneur docker puis lancer l'app
```
 ./enable_run.sh start
```

Au cas où l'app ne se lance pas, aller manuellement le faire : 

```
 cd gorillachat 
```

```
 docker compose up -d --build
```

```
 docker exec -it go-api sh -c "go run ."
```


Aller ensuite consulter la BDD sur: localhost:8080
- Serveur: database
- Utilisateur: root
- Mot de passe: password
- Nom de le base de donnée: chabdd

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
