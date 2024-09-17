
# Utilisation du CLI avec Cobra

## Build le projet pour générer un fichier `.exe`
```bash
go build -o bs3/bs3.exe
```

## Créer une commande avec le CLI Cobra
```bash
cobra-cli add "NOM DE COMMANDE"
```

## Exécuter une commande
```bash
go run main.go "NOM DE COMMANDE" "ARGS"
```
```bash
bs3/bs3.exe "NOM DE COMMANDE" "ARGS"
```

## Pour utiliser le prefix "bs3" dans bash

Ajouter : `[chemin vers]\my-cli-s3\bs3` à la variable d'environnement `PATH` de Windows.

![Exemple Bash](./exemple-cli.png)

## Commandes disponibles

- **Créer un bucket** :  
  ```bash
  bs3 create-bucket <bucket-name>
  ```

- **Uploader un fichier** :  
  ```bash
  bs3 upload-file <bucket-name> <file-path>
  ```

- **Télécharger un fichier** :  
  ```bash
  bs3 download-file <bucket-name> <file-name> <destination-path>
  ```

- **Lister les objets** :  
  ```bash
  bs3 list-object <bucket-name>
  ```

- **Uploader un objet** :  
  ```bash
  bs3 upload-object <bucket-name> <object-name> <file-path>
  ```
