# Chat Talks - Backend (en Go)

#### Port utilisé ici: 8000 - les appel en frontend sont addressé à ce port. 

#### Repo du front: <a href="https://github.com/ExploryKod/chatTalksClient">voir ici</a>

## Pour le lancer en allant dans go_app:

En local il est trés intéressant d'utiliser ***air*** pour le live reloading: 
<a href="https://github.com/cosmtrek/air">Voir la doc ici pour installer air</a> 

Je ne l'ai pas intégrer au docker.

***Avec docker*** depuis la racine (mais il faudra le remonté à chaque changement de code) 
```
docker compose up -d --build
```

***Dev Local sans docker*** en allant dans go_app (et non depuis la racine) : 

```
go mod tidy
```

```shell
go run main/main.go
```

***Dev en version turbo*** (live reloading) avec air (aprés avoir get le package) et sur go_app: 

- Générer la config / default : `air init`
- Puis : `air` 
