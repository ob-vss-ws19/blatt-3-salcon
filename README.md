## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    sudo docker run --rm --net actors --name treeservice terraform.cs.hm.edu:5043/ob-vss-ws19-blatt-3-salcon:development-docker-treeservice --bind="treeservice.actors:8090"
    ```

    Damit das funktioniert, müssen Sie folgendes erst im Tree-Service implementieren:

    -   die `main` verarbeitet Kommandozeilenflags und
    -   der Remote-Actor nutzt den Wert des Flags
    -   wenn Sie einen anderen Port als `8090` benutzen wollen,
        müssen Sie das auch im Dockerfile ändern (`EXPOSE...`)

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    sudo docker run --rm --net actors --name treecli terraform.cs.hm.edu:5043/ob-vss-ws19-blatt-3-salcon:development-docker-treecli --bind="treecli.actors:8091" --remote="treeservice.actors:8090" ARGUMENTE

    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass z.B. eine Liste aller Tree-Ids anzeigen soll.

    Zum Ausprobieren können Sie den Service dann laufen lassen. Das CLI soll ja jedes
    Mal nur einen Befehl abarbeiten und wird dann neu gestartet.

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.

   ```
-   Befehlsübersicht  treecli:

    treecli [FLAGS] COMMAND [KEY/SIZE] [VALUE]
    FLAGS
      -bind string
            Bind to address (default "localhost:8092")
      -id int
            tree id (default -1)
      -no-preserve-tree
            force deletion of tree
      -remote string
            remote host:port (default "127.0.0.1:8093")
      -token string
            tree token

      $ newtree SIZE
            Creates new tree. SIZE parameter specifies leaf size (minimum 1). Returns id and token

      $ insert KEY VALUE
            Insert an integer KEY with given string VALUE into the tree. id and token flag must be specified

      $ search KEY
            Search the tree for KEY. Returns corresponding value if found. id and token flag must be specified
      $ remove KEY
            Removes the KEY from the tree. id and token flag must be specified
      $ traverse
            gets all keys and values in the tree sorted by keys. id and token flag must be specified
      $ trees
            Gets a list of all available tree ids
      $ delete
            Deletes the tree. id and token flag must be specified, also no-preserve-tree flag must be set to true
    
    Example:
      $ treecli newtree 3
      $ treecli --id=0 --token=d57a23df insert 1 "hello world"
      $ treecli --id=0 --token=d57a23df insert 2 "welcome"
      $ treecli --id=0 --token=d57a23df search 1
      $ treecli --id=0 --token=d57a23df remove 1
      $ treecli --id=0 --token=d57a23df insert 1 "Hello HM"
      $ treecli --id=0 --token=d57a23df traverse
      

```

## Compilieren und Ausführen der Sourcen

-   Herunterladen des Repositorys:
    ```
    git clone https://github.com/ob-vss-ws19/blatt-3-salcon/
    ```
    
-   Compilieren und Starten des treeservice:
    ```
    cd blatt-3-salcon/treeservice
    go build
    ./treeservice
    ```    

-   Compilieren und Starten der treecli in einem zweiten Terminal:
    ```
    cd blatt-3-salcon/treecli
    go build
    ./treecli trees
    ```    

