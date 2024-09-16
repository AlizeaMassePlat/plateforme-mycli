Build le projet pour générer un fichier exe
go build     

Créer une commande avec le cli cobra 
cobra-cli add "NOM DE COMMANDE"

Executer une commande 
./plateforme-mycli.exe "NOM DE COMMANDE" "ARGS" 

Pour utiliser le prefix "bs3" dans bash : 

Ajouter : [chemin vers]\my-cli-s3\bs3 à la variable d'environnement PATH de windows 

![Exemple Bash](./exemple-cli.png)

Créer un bucket : bs3 create-bucket <bucket-name>
Uploader un fichier : bs3 upload-file <bucket-name> <file-path>
Télécharger un fichier : bs3 download-file <bucket-name> <file-name> <destination-path>
Lister les objets : bs3 list-object <bucket-name>
Uploader un objet : bs3 upload-object <bucket-name> <object-name> <file-path>
